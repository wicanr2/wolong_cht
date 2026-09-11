import sys
sys.path.insert(0,"tools")
from parity_diff import read_png
BLUE=(48,65,211); RED=(211,0,0)
def pts(path,col):
    W,H,px = read_png(path)
    out=set()
    for y in range(80,208):
        for x in range(480,640):
            if px[y][x]==col: out.add((x-480,y-80))
    # 2×2 點取左上角
    return {(x,y) for (x,y) in out if (x-1,y) not in out and (x,y-1) not in out}
for name,col in (("藍",BLUE),("紅",RED)):
    a=pts(sys.argv[1],col); b=pts(sys.argv[2],col)
    print("%s：原版 %d 點、remake %d 點，共同 %d" % (name,len(a),len(b),len(a&b)))
    oa=sorted(a-b); ob=sorted(b-a)
    print("   只在原版(%d)：%s" % (len(oa), oa[:12]))
    print("   只在remake(%d)：%s" % (len(ob), ob[:12]))
