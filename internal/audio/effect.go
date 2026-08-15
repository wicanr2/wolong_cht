package audio

import "fmt"

// `SOUND.DAT` 的重放（`docs/re/57` §6）。
//
// 音效與音樂用同一顆晶片但不同通道：音樂的六個聲軌佔滿六組 4-operator
// 通道（0–5 × 兩個 bank），音效走剩下的 **2-operator 通道 6、7、8**。
// 原版的音效常式（`0x10828`）就是對 `0x30+ch`／`0x50+ch`… 那一組寫的，
// 也就是通道 6 的 operator 位移 `0x10`／`0x13`。

// Effect 是一筆 16 bytes 的音效記錄。
type Effect struct {
	Raw []byte
}

// Effects 把 `SOUND.DAT` 拆成記錄。
func Effects(data []byte) ([]Effect, error) {
	if len(data) == 0 || len(data)%effectSize != 0 {
		return nil, fmt.Errorf("SOUND.DAT 是 %d bytes，不是 %d 的倍數", len(data), effectSize)
	}
	out := make([]Effect, len(data)/effectSize)
	for i := range out {
		out[i] = Effect{Raw: data[i*effectSize : (i+1)*effectSize]}
	}
	return out, nil
}

// Next 是接續播放的音效編號，Delay 是接續前要等幾個 tick。
func (e Effect) Next() int    { return int(e.Raw[0x0D]) }
func (e Effect) Delay() int   { return int(e.Raw[0x0E]) }
func (e Effect) Silent() bool { return e.Raw[0x02] == 0x3F && e.Raw[0x03] == 0x3F }

// program 把一筆音效寫進晶片的 2-op 通道 ch（6/7/8）。
//
// 暫存器位移照原版：通道 6 的兩個 operator 在 `0x10`／`0x13`，
// 所以 `0x20` 族落在 `0x30+n`、`0x40` 族落在 `0x50+n`，以此類推。
func (e Effect) program(c *OPL3, bank, ch int) {
	n := uint8(ch - 6)
	o1, o2 := 0x10+n, 0x13+n
	for i, base := range []uint8{0x20, 0x40, 0x60, 0x80} {
		c.Write(bank, base+o1, e.Raw[i*2])
		c.Write(bank, base+o2, e.Raw[i*2+1])
	}
	c.Write(bank, 0xE0+o1, e.Raw[0x08])
	c.Write(bank, 0xE0+o2, e.Raw[0x09])
	c.Write(bank, 0xC0+uint8(ch), e.Raw[0x0A])
	c.Write(bank, 0xA0+uint8(ch), e.Raw[0x0B])
	c.Write(bank, 0xB0+uint8(ch), e.Raw[0x0C])
}

// RenderEffect 把一筆音效（含接續鏈）渲染成 PCM。
//
// 鏈結指到 0 就是結束——**記錄 #0 是靜音哨兵**，不是「第 0 號音效」
// （`docs/re/57` §6）。tail 是最後一段的殘響要多留幾秒。
func RenderEffect(list []Effect, id int, tail float64) (*PCM, error) {
	if id < 0 || id >= len(list) {
		return nil, fmt.Errorf("音效編號 %d 超出範圍（0–%d）", id, len(list)-1)
	}
	const ch = 6 // 三個 2-op 通道裡的第一個
	c := NewOPL3(0)
	c.Write(1, 0x05, 0x01)
	c.Write(1, 0x04, 0x3F)

	out := &PCM{Rate: NativeRate}
	// 鏈結的延遲單位是 INT 1Ch 的一格，不是音樂 tick（`docs/re/57` §5）。
	perSample := isrHz / slowTickDiv / NativeRate
	acc, ticksLeft := 0.0, 0
	cur := id
	list[cur].program(c, 0, ch)
	ticksLeft = list[cur].Delay()

	// 上限擋住壞資料造成的無限鏈。正常的鏈都在幾十個 tick 內收掉。
	limit := int(NativeRate * (tail + 10))
	for i := 0; i < limit; i++ {
		acc += perSample
		for acc >= 1 {
			acc--
			if cur < 0 {
				continue
			}
			ticksLeft--
			if ticksLeft > 0 {
				continue
			}
			next := list[cur].Next()
			if next == 0 || next >= len(list) || next == cur {
				cur = -1 // 鏈結束，讓殘響自己收尾
				continue
			}
			cur = next
			list[cur].program(c, 0, ch)
			ticksLeft = list[cur].Delay()
		}
		l, r := c.Generate()
		out.L = append(out.L, float32(l/8192))
		out.R = append(out.R, float32(r/8192))
		if cur < 0 && float64(len(out.L)) > NativeRate*tail {
			break
		}
	}
	return out, nil
}
