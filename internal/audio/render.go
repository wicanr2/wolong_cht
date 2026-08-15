package audio

import (
	"encoding/binary"
	"io"
	"math"
)

// PCM 是立體聲取樣，範圍 −1..1。
type PCM struct {
	Rate      float64
	L, R      []float32
	TickCount int // 實際走了幾個音樂 tick
}

// Render 把一首曲子渲染成 PCM。
//
// 晶片一律跑原生的 49,716 Hz——包絡與相位都以晶片時脈為準，
// 直接用 44,100 建晶片會讓包絡的時間軸偏掉。要別的取樣率請用 Resample。
//
// seconds 是硬上限。原版靠控制事件無限循環（`docs/re/56` §5），
// **沒有「曲子結束」這回事**，所以呼叫端一定要給一個長度。
func Render(p *Player, seconds float64) *PCM {
	n := int(NativeRate * seconds)
	out := &PCM{Rate: NativeRate, L: make([]float32, 0, n), R: make([]float32, 0, n)}
	perSample := isrHz / NativeRate
	acc := 0.0
	for i := 0; i < n; i++ {
		acc += perSample
		for acc >= 1 {
			acc--
			if p.isr() {
				out.TickCount++
			}
		}
		l, r := p.chip.Generate()
		out.L = append(out.L, float32(l/8192))
		out.R = append(out.R, float32(r/8192))
	}
	return out
}

// Resample 是線性內插的重取樣。
//
// 用途只有「49,716 → 44,100」這一步，落差不到 13%，線性內插的誤差
// 遠低於 OPL3 模型本身的近似（`opl3.go` 開頭的「模型」那半邊）。
// 要更好的品質請在外面用 ffmpeg，不要在這裡加濾波器——
// **這一層要保持看得懂**。
func Resample(in *PCM, rate float64) *PCM {
	if rate <= 0 || rate == in.Rate {
		return in
	}
	ratio := in.Rate / rate
	n := int(float64(len(in.L)) / ratio)
	out := &PCM{Rate: rate, L: make([]float32, n), R: make([]float32, n),
		TickCount: in.TickCount}
	for i := 0; i < n; i++ {
		pos := float64(i) * ratio
		j := int(pos)
		f := float32(pos - float64(j))
		k := j + 1
		if k >= len(in.L) {
			k = len(in.L) - 1
		}
		out.L[i] = in.L[j] + (in.L[k]-in.L[j])*f
		out.R[i] = in.R[j] + (in.R[k]-in.R[j])*f
	}
	return out
}

// Normalize 把峰值拉到 peak（0..1）。回傳套用的增益。
//
// 原版錄音的峰值只有滿刻度的 4.7%（`docs/playtest/25` §2），
// 直接輸出會安靜到以為沒聲音。**這是 remake 差異，要標記。**
func Normalize(p *PCM, peak float64) float64 {
	max := 0.0
	for i := range p.L {
		max = math.Max(max, math.Abs(float64(p.L[i])))
		max = math.Max(max, math.Abs(float64(p.R[i])))
	}
	if max == 0 {
		return 0
	}
	gain := peak / max
	for i := range p.L {
		p.L[i] = float32(float64(p.L[i]) * gain)
		p.R[i] = float32(float64(p.R[i]) * gain)
	}
	return gain
}

// SegmentRMS 把 PCM 切成 n 段各算 RMS。
//
// **驗「不是靜音」一定要分段**：只渲染到片頭的話後段是 0，
// 而總 RMS 看不出來（`docs/playtest/25` §2 踩過）。
func SegmentRMS(p *PCM, n int) []float64 {
	if n <= 0 || len(p.L) == 0 {
		return nil
	}
	out := make([]float64, n)
	size := len(p.L) / n
	if size == 0 {
		return nil
	}
	for s := 0; s < n; s++ {
		var sum float64
		for i := s * size; i < (s+1)*size; i++ {
			sum += float64(p.L[i])*float64(p.L[i]) + float64(p.R[i])*float64(p.R[i])
		}
		out[s] = math.Sqrt(sum / float64(size*2))
	}
	return out
}

// WriteWAV 寫 16-bit 立體聲 PCM。
func WriteWAV(w io.Writer, p *PCM) error {
	n := len(p.L)
	dataLen := n * 4
	hdr := []any{
		[]byte("RIFF"), uint32(36 + dataLen), []byte("WAVEfmt "),
		uint32(16), uint16(1), uint16(2), uint32(p.Rate),
		uint32(p.Rate) * 4, uint16(4), uint16(16),
		[]byte("data"), uint32(dataLen),
	}
	for _, f := range hdr {
		if err := binary.Write(w, binary.LittleEndian, f); err != nil {
			return err
		}
	}
	buf := make([]byte, dataLen)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint16(buf[i*4:], uint16(int16(clip(p.L[i])*32767)))
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(int16(clip(p.R[i])*32767)))
	}
	_, err := w.Write(buf)
	return err
}

func clip(v float32) float32 {
	if v > 1 {
		return 1
	}
	if v < -1 {
		return -1
	}
	return v
}
