/* Fixed IDA native C input; original segment locators; spec/231, re/111. */
/* Original sub_18810; IDA 0x18810. */
static void mx_body_18810(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x18810:goto L_18810; default:assert(0);}
L_18810: ; /* IDA 0x18810; push    bx */
    ec_push(m,m->bx);
L_18811: ; /* IDA 0x18811; push    dx */
    ec_push(m,m->dx);
L_18812: ; /* IDA 0x18812; push    ax */
    ec_push(m,m->ax);
L_18813: ; /* IDA 0x18813; push    cx */
    ec_push(m,m->cx);
L_18814: ; /* IDA 0x18814; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_18817: ; /* IDA 0x18817; mov     bx, 8 */
    m->bx=(uint16_t)(8);
L_1881A: ; /* IDA 0x1881A; mov     al, 1Eh */
    vg_al(m,(uint8_t)(30));
L_1881C: ; /* IDA 0x1881C; mov     cx, 510h */
    m->cx=(uint16_t)(1296);
L_1881F: ; /* IDA 0x1881F; call    sub_1895D */
    mx_call(m,35165,34850,0x10000,h);
L_18822: ; /* IDA 0x18822; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_18823: ; /* IDA 0x18823; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_18824: ; /* IDA 0x18824; call    sub_101B4 */
    mx_call(m,436,34855,0x10000,h);
L_18827: ; /* IDA 0x18827; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_1882A: ; /* IDA 0x1882A; mov     bx, 0Ah */
    m->bx=(uint16_t)(10);
L_1882D: ; /* IDA 0x1882D; call    sub_1075B */
    mx_call(m,1883,34864,0x10000,h);
L_18830: ; /* IDA 0x18830; call    sub_101DB */
    mx_call(m,475,34867,0x10000,h);
L_18833: ; /* IDA 0x18833; call    loc_1222B */
    mx_call(m,8747,34870,0x10000,h);
L_18836: ; /* IDA 0x18836; call    sub_101B4 */
    mx_call(m,436,34873,0x10000,h);
L_18839: ; /* IDA 0x18839; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_1883C: ; /* IDA 0x1883C; mov     bx, 8 */
    m->bx=(uint16_t)(8);
L_1883F: ; /* IDA 0x1883F; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_18841: ; /* IDA 0x18841; mov     cx, 510h */
    m->cx=(uint16_t)(1296);
L_18844: ; /* IDA 0x18844; call    sub_1895D */
    mx_call(m,35165,34887,0x10000,h);
L_18847: ; /* IDA 0x18847; mov     cs:byte_198A5, 8 */
#if KI_INPUT_MUTATION == 12
    ; /* lost message close marker */
#else
    dl_write8(m,m->cs,39077,8);
#endif
L_1884D: ; /* IDA 0x1884D; call    sub_101DB */
    mx_call(m,475,34896,0x10000,h);
L_18850: ; /* IDA 0x18850; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_18851: ; /* IDA 0x18851; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_18852: ; /* IDA 0x18852; retn */
    return;
}
/* Original sub_121E7; IDA 0x121E7. */
static void mx_body_121E7(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x121e7:goto L_121E7; default:assert(0);}
L_121E7: ; /* IDA 0x121E7; push    ax */
    ec_push(m,m->ax);
L_121E8: ; /* IDA 0x121E8; push    bx */
    ec_push(m,m->bx);
L_121E9: ; /* IDA 0x121E9; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_121EC: ; /* IDA 0x121EC; mov     bx, 1 */
#if KI_INPUT_MUTATION == 1
    m->bx=(uint16_t)(0);
#else
    m->bx=(uint16_t)(1);
#endif
L_121EF: ; /* IDA 0x121EF; call    sub_20000 */
    mx_far(m,0x21ef,0x21f4,h);
L_121F4: ; /* IDA 0x121F4; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_121F6: ; /* IDA 0x121F6; jnz     short loc_12212 */
    if(!(m->flags&0x40)) goto L_12212;
L_121F8: ; /* IDA 0x121F8; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_121FB: ; /* IDA 0x121FB; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_121FD: ; /* IDA 0x121FD; call    sub_20000 */
    mx_far(m,0x21fd,0x2202,h);
L_12202: ; /* IDA 0x12202; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_12204: ; /* IDA 0x12204; jz      short loc_121E9 */
    if(m->flags&0x40) goto L_121E9;
L_12206: ; /* IDA 0x12206; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_12209: ; /* IDA 0x12209; call    sub_20000 */
    mx_far(m,0x2209,0x220e,h);
L_1220E: ; /* IDA 0x1220E; clc */
    m->flags&=(uint16_t)~1u;
L_1220F: ; /* IDA 0x1220F; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_12210: ; /* IDA 0x12210; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_12211: ; /* IDA 0x12211; retn */
    return;
L_12212: ; /* IDA 0x12212; stc */
    m->flags|=1;
L_12213: ; /* IDA 0x12213; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_12214: ; /* IDA 0x12214; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_12215: ; /* IDA 0x12215; retn */
    return;
}
/* Original sub_20000; IDA 0x20000. */
static void mx_body_20000(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x20000:goto L_20000; default:assert(0);}
L_20000: ; /* IDA 0x20000; cli */
    m->flags&=(uint16_t)~0x200u;
L_20001: ; /* IDA 0x20001; push    si */
    ec_push(m,m->si);
L_20002: ; /* IDA 0x20002; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_20004: ; /* IDA 0x20004; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_20006: ; /* IDA 0x20006; call    cs:funcs_20006[si] */
    mx_call(m,dl_read16(m,m->cs,(uint16_t)(m->si+14)),11,0x20000,h);
L_2000B: ; /* IDA 0x2000B; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_2000C: ; /* IDA 0x2000C; sti */
#if KI_INPUT_MUTATION == 10
    ; /* lost original STI */
#else
    m->flags|=0x200;
#endif
L_2000D: ; /* IDA 0x2000D; retf */
    return;
}
/* Original sub_2002E; IDA 0x2002E. */
static void mx_body_2002E(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x2002e:goto L_2002E; case 0x2006f:goto L_2006F; default:assert(0);}
L_2002E: ; /* IDA 0x2002E; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_20030: ; /* IDA 0x20030; int     33h; - MS MOUSE - RESET DRIVER AND READ STATUS */
    gl_interrupt(m,51,50);
L_20032: ; /* IDA 0x20032; and     ax, ax */
    {uint16_t result=((m->ax)&(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_20034: ; /* IDA 0x20034; jz      short nullsub_4 */
    if(m->flags&0x40) goto L_2006F;
L_20036: ; /* IDA 0x20036; push    es */
    ec_push(m,m->es);
L_20037: ; /* IDA 0x20037; push    bx */
    ec_push(m,m->bx);
L_20038: ; /* IDA 0x20038; push    cx */
    ec_push(m,m->cx);
L_20039: ; /* IDA 0x20039; push    dx */
    ec_push(m,m->dx);
L_2003A: ; /* IDA 0x2003A; mov     ax, 7 */
    m->ax=(uint16_t)(7);
L_2003D: ; /* IDA 0x2003D; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_2003F: ; /* IDA 0x2003F; mov     dx, 27Fh */
    m->dx=(uint16_t)(639);
L_20042: ; /* IDA 0x20042; int     33h; - MS MOUSE - DEFINE HORIZONTAL CURSOR RANGE */
    gl_interrupt(m,51,68);
L_20044: ; /* IDA 0x20044; mov     ax, 8 */
    m->ax=(uint16_t)(8);
L_20047: ; /* IDA 0x20047; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_20049: ; /* IDA 0x20049; mov     dx, 18Fh */
    m->dx=(uint16_t)(399);
L_2004C: ; /* IDA 0x2004C; int     33h; - MS MOUSE - DEFINE VERTICAL CURSOR RANGE */
    gl_interrupt(m,51,78);
L_2004E: ; /* IDA 0x2004E; mov     ax, 0Fh */
    m->ax=(uint16_t)(15);
L_20051: ; /* IDA 0x20051; mov     cx, 8 */
    m->cx=(uint16_t)(8);
L_20054: ; /* IDA 0x20054; mov     dx, 8 */
    m->dx=(uint16_t)(8);
L_20057: ; /* IDA 0x20057; int     33h; - MS MOUSE - DEFINE MICKEY/PIXEL RATIO */
    gl_interrupt(m,51,89);
L_20059: ; /* IDA 0x20059; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_2005B: ; /* IDA 0x2005B; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_2005D: ; /* IDA 0x2005D; mov     dx, 101h */
    m->dx=(uint16_t)(257);
L_20060: ; /* IDA 0x20060; mov     ax, 0Ch */
    m->ax=(uint16_t)(12);
L_20063: ; /* IDA 0x20063; mov     cx, 1 */
    m->cx=(uint16_t)(1);
L_20066: ; /* IDA 0x20066; int     33h; - MS MOUSE - DEFINE INTERRUPT SUBROUTINE PARAMETERS */
    gl_interrupt(m,51,104);
L_20068: ; /* IDA 0x20068; mov     ax, 0FFFFh */
    m->ax=(uint16_t)(65535);
L_2006B: ; /* IDA 0x2006B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_2006C: ; /* IDA 0x2006C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_2006D: ; /* IDA 0x2006D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_2006E: ; /* IDA 0x2006E; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_2006F: ; /* IDA 0x2006F; retn */
    if(entry!=0x2006f) {m->ip=0x6f;if(h&&h->enter) h->enter(m,0x6f,h->user);}
    return;
}
/* Original sub_20070; IDA 0x20070. */
static void mx_body_20070(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x20070:goto L_20070; default:assert(0);}
L_20070: ; /* IDA 0x20070; cmp     cs:byte_20100, 1 */
    ec_cmp(m,dl_read8(m,m->cs,256),1,8);
L_20076: ; /* IDA 0x20076; jz      short locret_20099 */
    if(m->flags&0x40) goto L_20099;
L_20078: ; /* IDA 0x20078; push    bx */
    ec_push(m,m->bx);
L_20079: ; /* IDA 0x20079; push    cx */
    ec_push(m,m->cx);
L_2007A: ; /* IDA 0x2007A; push    dx */
    ec_push(m,m->dx);
L_2007B: ; /* IDA 0x2007B; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_2007E: ; /* IDA 0x2007E; int     33h; - MS MOUSE - RETURN POSITION AND BUTTON STATUS */
    gl_interrupt(m,51,128);
L_20080: ; /* IDA 0x20080; mov     cs:word_200FC, cx */
    dl_write16(m,m->cs,252,m->cx);
L_20085: ; /* IDA 0x20085; mov     cs:word_200FE, dx */
    dl_write16(m,m->cs,254,m->dx);
L_2008A: ; /* IDA 0x2008A; mov     cs:byte_20100, 1 */
    dl_write8(m,m->cs,256,1);
L_20090: ; /* IDA 0x20090; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_20093: ; /* IDA 0x20093; call    sub_20137 */
    mx_call(m,311,150,0x20000,h);
L_20096: ; /* IDA 0x20096; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_20097: ; /* IDA 0x20097; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_20098: ; /* IDA 0x20098; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_20099: ; /* IDA 0x20099; retn */
    return;
}
/* Original sub_2009A; IDA 0x2009A. */
static void mx_body_2009A(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x2009a:goto L_2009A; default:assert(0);}
L_2009A: ; /* IDA 0x2009A; push    cx */
    ec_push(m,m->cx);
L_2009B: ; /* IDA 0x2009B; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_2009D: ; /* IDA 0x2009D; xchg    al, cs:byte_20100 */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->cs,256);vg_al(m,(uint8_t)(b));dl_write8(m,m->cs,256,a);}
L_200A2: ; /* IDA 0x200A2; dec     ax */
    m->ax=(uint16_t)(st_dec(m,m->ax,16));
L_200A3: ; /* IDA 0x200A3; mov     cx, cs:word_200FC */
    m->cx=(uint16_t)(dl_read16(m,m->cs,252));
L_200A8: ; /* IDA 0x200A8; mov     dx, cs:word_200FE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,254));
L_200AD: ; /* IDA 0x200AD; push    ax */
    ec_push(m,m->ax);
L_200AE: ; /* IDA 0x200AE; push    cx */
    ec_push(m,m->cx);
L_200AF: ; /* IDA 0x200AF; push    dx */
    ec_push(m,m->dx);
L_200B0: ; /* IDA 0x200B0; cmp     ax, 0FFFFh */
    ec_cmp(m,m->ax,65535,16);
L_200B3: ; /* IDA 0x200B3; jz      short loc_200B8 */
    if(m->flags&0x40) goto L_200B8;
L_200B5: ; /* IDA 0x200B5; call    sub_20137 */
    mx_call(m,311,184,0x20000,h);
L_200B8: ; /* IDA 0x200B8; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_200B9: ; /* IDA 0x200B9; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_200BA: ; /* IDA 0x200BA; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_200BB: ; /* IDA 0x200BB; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_200BC: ; /* IDA 0x200BC; retn */
    return;
}
/* Original sub_200BD; IDA 0x200BD. */
static void mx_body_200BD(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x200bd:goto L_200BD; default:assert(0);}
L_200BD: ; /* IDA 0x200BD; int     33h; - MS MOUSE - */
    gl_interrupt(m,51,191);
L_200BF: ; /* IDA 0x200BF; retn */
    return;
}
/* Original sub_200C0; IDA 0x200C0. */
static void mx_body_200C0(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x200c0:goto L_200C0; default:assert(0);}
L_200C0: ; /* IDA 0x200C0; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_200C2: ; /* IDA 0x200C2; xchg    al, cs:byte_20100 */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->cs,256);vg_al(m,(uint8_t)(b));dl_write8(m,m->cs,256,a);}
L_200C7: ; /* IDA 0x200C7; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_200C9: ; /* IDA 0x200C9; jnz     short loc_200D4 */
    if(!(m->flags&0x40)) goto L_200D4;
L_200CB: ; /* IDA 0x200CB; mov     cs:word_200FC, 0FFFFh */
    dl_write16(m,m->cs,252,65535);
L_200D2: ; /* IDA 0x200D2; jmp     short loc_200E2 */
    goto L_200E2;
L_200D4: ; /* IDA 0x200D4; cmp     dx, cs:word_200FC */
    ec_cmp(m,m->dx,dl_read16(m,m->cs,252),16);
L_200D9: ; /* IDA 0x200D9; jnz     short loc_200E2 */
    if(!(m->flags&0x40)) goto L_200E2;
L_200DB: ; /* IDA 0x200DB; cmp     bx, cs:word_200FE */
    ec_cmp(m,m->bx,dl_read16(m,m->cs,254),16);
L_200E0: ; /* IDA 0x200E0; jz      short locret_200F5 */
    if(m->flags&0x40) goto L_200F5;
L_200E2: ; /* IDA 0x200E2; push    cx */
    ec_push(m,m->cx);
L_200E3: ; /* IDA 0x200E3; xchg    dx, cs:word_200FC */
    {uint16_t a=m->dx,b=dl_read16(m,m->cs,252);m->dx=(uint16_t)(b);dl_write16(m,m->cs,252,a);}
L_200E8: ; /* IDA 0x200E8; xchg    bx, cs:word_200FE */
    {uint16_t a=m->bx,b=dl_read16(m,m->cs,254);m->bx=(uint16_t)(b);dl_write16(m,m->cs,254,a);}
L_200ED: ; /* IDA 0x200ED; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_200EF: ; /* IDA 0x200EF; mov     dx, bx */
    m->dx=(uint16_t)(m->bx);
L_200F1: ; /* IDA 0x200F1; call    sub_20137 */
    mx_call(m,311,244,0x20000,h);
L_200F4: ; /* IDA 0x200F4; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_200F5: ; /* IDA 0x200F5; retn */
    return;
}
/* Original sub_20137; IDA 0x20137. */
static void mx_body_20137(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x20137:goto L_20137; default:assert(0);}
L_20137: ; /* IDA 0x20137; push    ds */
    ec_push(m,m->ds);
L_20138: ; /* IDA 0x20138; push    si */
    ec_push(m,m->si);
L_20139: ; /* IDA 0x20139; push    di */
    ec_push(m,m->di);
L_2013A: ; /* IDA 0x2013A; push    bp */
    ec_push(m,m->bp);
L_2013B: ; /* IDA 0x2013B; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_2013E: ; /* IDA 0x2013E; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_20140: ; /* IDA 0x20140; cmp     cx, 0FFFFh */
    ec_cmp(m,m->cx,65535,16);
L_20143: ; /* IDA 0x20143; jz      short loc_20153 */
    if(m->flags&0x40) goto L_20153;
L_20145: ; /* IDA 0x20145; call    sub_201C6 */
    mx_call(m,454,328,0x20000,h);
L_20148: ; /* IDA 0x20148; call    sub_20249 */
    mx_call(m,585,331,0x20000,h);
L_2014B: ; /* IDA 0x2014B; cmp     cs:byte_20100, 0 */
    ec_cmp(m,dl_read8(m,m->cs,256),0,8);
L_20151: ; /* IDA 0x20151; jz      short loc_20166 */
    if(m->flags&0x40) goto L_20166;
L_20153: ; /* IDA 0x20153; mov     cx, cs:word_200FC */
    m->cx=(uint16_t)(dl_read16(m,m->cs,252));
L_20158: ; /* IDA 0x20158; mov     dx, cs:word_200FE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,254));
L_2015D: ; /* IDA 0x2015D; call    sub_201C6 */
    mx_call(m,454,352,0x20000,h);
L_20160: ; /* IDA 0x20160; call    sub_202BD */
    mx_call(m,701,355,0x20000,h);
L_20163: ; /* IDA 0x20163; call    sub_201E4 */
    mx_call(m,484,358,0x20000,h);
L_20166: ; /* IDA 0x20166; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_20167: ; /* IDA 0x20167; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_20168: ; /* IDA 0x20168; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_20169: ; /* IDA 0x20169; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_2016A: ; /* IDA 0x2016A; retn */
    return;
}
/* Original sub_2016B; IDA 0x2016B. */
static void mx_body_2016B(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x2016b:goto L_2016B; default:assert(0);}
L_2016B: ; /* IDA 0x2016B; push    cx */
    ec_push(m,m->cx);
L_2016C: ; /* IDA 0x2016C; push    dx */
    ec_push(m,m->dx);
L_2016D: ; /* IDA 0x2016D; mov     bx, 0F6h */
    m->bx=(uint16_t)(246);
L_20170: ; /* IDA 0x20170; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_20173: ; /* IDA 0x20173; in      al, dx; EGA: graph 1 and 2 addr reg: */
    vg_al(m,(uint8_t)(wolong_input_in(m->dx)));
L_20174: ; /* IDA 0x20174; mov     [bp+6], al */
    dl_write8(m,m->ss,(uint16_t)(m->bp+6),(uint8_t)m->ax);
L_20177: ; /* IDA 0x20177; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_2017A: ; /* IDA 0x2017A; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_2017C: ; /* IDA 0x2017C; mov     al, cs:[bx+si] */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->bx+m->si))));
L_2017F: ; /* IDA 0x2017F; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_20180: ; /* IDA 0x20180; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_20181: ; /* IDA 0x20181; in      al, dx; EGA port: graphics controller data register */
    vg_al(m,(uint8_t)(wolong_input_in(m->dx)));
L_20182: ; /* IDA 0x20182; mov     [bp+si], al */
    dl_write8(m,m->ss,(uint16_t)(m->bp+m->si),(uint8_t)m->ax);
L_20184: ; /* IDA 0x20184; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_20185: ; /* IDA 0x20185; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_20186: ; /* IDA 0x20186; loop    loc_2017C */
    if(--m->cx) goto L_2017C;
L_20188: ; /* IDA 0x20188; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_2018B: ; /* IDA 0x2018B; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_2018D: ; /* IDA 0x2018D; out     dx, al; EGA: sequencer address reg */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_2018E: ; /* IDA 0x2018E; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_2018F: ; /* IDA 0x2018F; in      al, dx; EGA port: sequencer data register */
    vg_al(m,(uint8_t)(wolong_input_in(m->dx)));
L_20190: ; /* IDA 0x20190; mov     [bp+7], al */
    dl_write8(m,m->ss,(uint16_t)(m->bp+7),(uint8_t)m->ax);
L_20193: ; /* IDA 0x20193; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_20196: ; /* IDA 0x20196; mov     ax, 0FF08h */
    m->ax=(uint16_t)(65288);
L_20199: ; /* IDA 0x20199; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_2019A: ; /* IDA 0x2019A; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_2019B: ; /* IDA 0x2019B; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_2019C: ; /* IDA 0x2019C; retn */
    return;
}
/* Original sub_2019D; IDA 0x2019D. */
static void mx_body_2019D(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x2019d:goto L_2019D; default:assert(0);}
L_2019D: ; /* IDA 0x2019D; push    cx */
    ec_push(m,m->cx);
L_2019E: ; /* IDA 0x2019E; push    dx */
    ec_push(m,m->dx);
L_2019F: ; /* IDA 0x2019F; mov     bx, 0F6h */
    m->bx=(uint16_t)(246);
L_201A2: ; /* IDA 0x201A2; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_201A5: ; /* IDA 0x201A5; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_201A8: ; /* IDA 0x201A8; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_201AA: ; /* IDA 0x201AA; mov     al, cs:[bx+si] */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->bx+m->si))));
L_201AD: ; /* IDA 0x201AD; mov     ah, [bp+si] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+m->si))));
L_201AF: ; /* IDA 0x201AF; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_201B0: ; /* IDA 0x201B0; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_201B1: ; /* IDA 0x201B1; loop    loc_201AA */
    if(--m->cx) goto L_201AA;
L_201B3: ; /* IDA 0x201B3; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_201B6: ; /* IDA 0x201B6; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_201B8: ; /* IDA 0x201B8; mov     ah, [bp+7] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+7))));
L_201BB: ; /* IDA 0x201BB; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_201BC: ; /* IDA 0x201BC; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_201BF: ; /* IDA 0x201BF; mov     al, [bp+6] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+6))));
L_201C2: ; /* IDA 0x201C2; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_201C3: ; /* IDA 0x201C3; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_201C4: ; /* IDA 0x201C4; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_201C5: ; /* IDA 0x201C5; retn */
    return;
}
/* Original sub_201C6; IDA 0x201C6. */
static void mx_body_201C6(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x201c6:goto L_201C6; default:assert(0);}
L_201C6: ; /* IDA 0x201C6; mov     bx, cx */
    m->bx=(uint16_t)(m->cx);
L_201C8: ; /* IDA 0x201C8; mov     di, cx */
    m->di=(uint16_t)(m->cx);
L_201CA: ; /* IDA 0x201CA; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_201CC: ; /* IDA 0x201CC; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_201CE: ; /* IDA 0x201CE; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_201D0: ; /* IDA 0x201D0; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_201D2: ; /* IDA 0x201D2; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_201D4: ; /* IDA 0x201D4; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_201D6: ; /* IDA 0x201D6; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_201D8: ; /* IDA 0x201D8; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_201DA: ; /* IDA 0x201DA; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_201DC: ; /* IDA 0x201DC; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_201DE: ; /* IDA 0x201DE; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_201E0: ; /* IDA 0x201E0; and     cl, 7 */
    {uint16_t result=(((uint8_t)m->cx)&(7))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_201E3: ; /* IDA 0x201E3; retn */
    return;
}
/* Original sub_201E4; IDA 0x201E4. */
static void mx_body_201E4(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x201e4:goto L_201E4; default:assert(0);}
L_201E4: ; /* IDA 0x201E4; mov     si, 31Bh */
#if KI_INPUT_MUTATION == 5
    m->si=(uint16_t)(827);
#else
    m->si=(uint16_t)(795);
#endif
L_201E7: ; /* IDA 0x201E7; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_201EA: ; /* IDA 0x201EA; mov     ax, 305h */
    m->ax=(uint16_t)(773);
L_201ED: ; /* IDA 0x201ED; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_201EE: ; /* IDA 0x201EE; mov     ax, 0FF08h */
    m->ax=(uint16_t)(65288);
L_201F1: ; /* IDA 0x201F1; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_201F2: ; /* IDA 0x201F2; mov     ax, 0F00h */
    m->ax=(uint16_t)(3840);
L_201F5: ; /* IDA 0x201F5; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_201F6: ; /* IDA 0x201F6; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_201F9: ; /* IDA 0x201F9; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_201FA: ; /* IDA 0x201FA; call    sub_2020C */
    mx_call(m,524,509,0x20000,h);
L_201FD: ; /* IDA 0x201FD; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_20200: ; /* IDA 0x20200; mov     ax, 0A00h */
    m->ax=(uint16_t)(2560);
L_20203: ; /* IDA 0x20203; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_20204: ; /* IDA 0x20204; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_20207: ; /* IDA 0x20207; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_20208: ; /* IDA 0x20208; call    sub_2020C */
    mx_call(m,524,523,0x20000,h);
L_2020B: ; /* IDA 0x2020B; retn */
    return;
}
/* Original sub_2020C; IDA 0x2020C. */
static void mx_body_2020C(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x2020c:goto L_2020C; default:assert(0);}
L_2020C: ; /* IDA 0x2020C; mov     ch, 10h */
    vg_ch(m,(uint8_t)(16));
L_2020E: ; /* IDA 0x2020E; cmp     bx, 7D00h */
    ec_cmp(m,m->bx,32000,16);
L_20212: ; /* IDA 0x20212; jnb     short loc_2023B */
#if KI_INPUT_MUTATION == 6
    if(m->flags&1) goto L_2023B;
#else
    if(!(m->flags&1)) goto L_2023B;
#endif
L_20214: ; /* IDA 0x20214; mov     ax, cs:[si] */
    m->ax=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->si)));
L_20217: ; /* IDA 0x20217; mov     dh, al */
    vg_dh(m,(uint8_t)((uint8_t)m->ax));
L_20219: ; /* IDA 0x20219; xor     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_2021B: ; /* IDA 0x2021B; shr     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,5));
L_2021D: ; /* IDA 0x2021D; shr     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
L_2021F: ; /* IDA 0x2021F; mov     dh, [bx] */
#if KI_INPUT_MUTATION == 9
    ; /* lost cursor latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx))));
#endif
L_20221: ; /* IDA 0x20221; mov     [bx], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx),(uint8_t)(m->ax>>8));
L_20223: ; /* IDA 0x20223; cmp     di, 278h */
#if KI_INPUT_MUTATION == 7
    ec_cmp(m,m->di,631,16);
#else
    ec_cmp(m,m->di,632,16);
#endif
L_20227: ; /* IDA 0x20227; jnb     short loc_2023B */
    if(!(m->flags&1)) goto L_2023B;
L_20229: ; /* IDA 0x20229; mov     dh, [bx+1] */
#if KI_INPUT_MUTATION == 9
    ; /* lost cursor latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
#endif
L_2022C: ; /* IDA 0x2022C; mov     [bx+1], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+1),(uint8_t)m->ax);
L_2022F: ; /* IDA 0x2022F; cmp     di, 270h */
    ec_cmp(m,m->di,624,16);
L_20233: ; /* IDA 0x20233; jnb     short loc_2023B */
    if(!(m->flags&1)) goto L_2023B;
L_20235: ; /* IDA 0x20235; mov     dh, [bx+2] */
#if KI_INPUT_MUTATION == 9
    ; /* lost cursor latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+2))));
#endif
L_20238: ; /* IDA 0x20238; mov     [bx+2], dl */
    dl_write8(m,m->ds,(uint16_t)(m->bx+2),(uint8_t)m->dx);
L_2023B: ; /* IDA 0x2023B; add     bx, 50h ; 'P' */
    m->bx=(uint16_t)(ec_math(m,m->bx,80,16,0,0));
L_2023E: ; /* IDA 0x2023E; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_2023F: ; /* IDA 0x2023F; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_20240: ; /* IDA 0x20240; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_20242: ; /* IDA 0x20242; jnz     short loc_2020E */
    if(!(m->flags&0x40)) goto L_2020E;
L_20244: ; /* IDA 0x20244; sub     bx, 500h */
    m->bx=(uint16_t)(ec_math(m,m->bx,1280,16,0,1));
L_20248: ; /* IDA 0x20248; retn */
    return;
}
/* Original sub_20249; IDA 0x20249. */
static void mx_body_20249(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x20249:goto L_20249; default:assert(0);}
L_20249: ; /* IDA 0x20249; push    bx */
    ec_push(m,m->bx);
L_2024A: ; /* IDA 0x2024A; mov     ax, di */
    m->ax=(uint16_t)(m->di);
L_2024C: ; /* IDA 0x2024C; sub     ax, 268h */
#if KI_INPUT_MUTATION == 8
    m->ax=(uint16_t)(ec_math(m,m->ax,624,16,0,1));
#else
    m->ax=(uint16_t)(ec_math(m,m->ax,616,16,0,1));
#endif
L_2024F: ; /* IDA 0x2024F; jbe     short loc_20259 */
    if(m->flags&(1|0x40)) goto L_20259;
L_20251: ; /* IDA 0x20251; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_20253: ; /* IDA 0x20253; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_20255: ; /* IDA 0x20255; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_20257: ; /* IDA 0x20257; sub     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,1));
L_20259: ; /* IDA 0x20259; mov     si, 35Bh */
    m->si=(uint16_t)(859);
L_2025C: ; /* IDA 0x2025C; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_2025F: ; /* IDA 0x2025F; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_20262: ; /* IDA 0x20262; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_20263: ; /* IDA 0x20263; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_20266: ; /* IDA 0x20266; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_20267: ; /* IDA 0x20267; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_2026A: ; /* IDA 0x2026A; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_2026B: ; /* IDA 0x2026B; mov     ax, 0FF08h */
    m->ax=(uint16_t)(65288);
L_2026E: ; /* IDA 0x2026E; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_2026F: ; /* IDA 0x2026F; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_20272: ; /* IDA 0x20272; mov     ax, 102h */
    m->ax=(uint16_t)(258);
L_20275: ; /* IDA 0x20275; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_20276: ; /* IDA 0x20276; call    sub_202A0 */
    mx_call(m,672,633,0x20000,h);
L_20279: ; /* IDA 0x20279; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_2027C: ; /* IDA 0x2027C; mov     ax, 202h */
    m->ax=(uint16_t)(514);
L_2027F: ; /* IDA 0x2027F; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_20280: ; /* IDA 0x20280; call    sub_202A0 */
    mx_call(m,672,643,0x20000,h);
L_20283: ; /* IDA 0x20283; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_20286: ; /* IDA 0x20286; mov     ax, 402h */
    m->ax=(uint16_t)(1026);
L_20289: ; /* IDA 0x20289; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_2028A: ; /* IDA 0x2028A; call    sub_202A0 */
    mx_call(m,672,653,0x20000,h);
L_2028D: ; /* IDA 0x2028D; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_20290: ; /* IDA 0x20290; mov     ax, 802h */
    m->ax=(uint16_t)(2050);
L_20293: ; /* IDA 0x20293; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_20294: ; /* IDA 0x20294; call    sub_202A0 */
    mx_call(m,672,663,0x20000,h);
L_20297: ; /* IDA 0x20297; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_2029A: ; /* IDA 0x2029A; mov     ax, 0F02h */
    m->ax=(uint16_t)(3842);
L_2029D: ; /* IDA 0x2029D; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_2029E: ; /* IDA 0x2029E; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_2029F: ; /* IDA 0x2029F; retn */
    return;
}
/* Original sub_202A0; IDA 0x202A0. */
static void mx_body_202A0(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x202a0:goto L_202A0; default:assert(0);}
L_202A0: ; /* IDA 0x202A0; mov     ch, 10h */
    vg_ch(m,(uint8_t)(16));
L_202A2: ; /* IDA 0x202A2; mov     ax, cs:[si] */
    m->ax=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->si)));
L_202A5: ; /* IDA 0x202A5; mov     dl, cs:[si+2] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->si+2))));
L_202A9: ; /* IDA 0x202A9; mov     [bx], ax */
    dl_write16(m,m->ds,(uint16_t)(m->bx),m->ax);
L_202AB: ; /* IDA 0x202AB; mov     [bx+2], dl */
    dl_write8(m,m->ds,(uint16_t)(m->bx+2),(uint8_t)m->dx);
L_202AE: ; /* IDA 0x202AE; add     si, 3 */
    m->si=(uint16_t)(ec_math(m,m->si,3,16,0,0));
L_202B1: ; /* IDA 0x202B1; add     bx, 50h ; 'P' */
    m->bx=(uint16_t)(ec_math(m,m->bx,80,16,0,0));
L_202B4: ; /* IDA 0x202B4; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_202B6: ; /* IDA 0x202B6; jnz     short loc_202A2 */
    if(!(m->flags&0x40)) goto L_202A2;
L_202B8: ; /* IDA 0x202B8; sub     bx, 500h */
    m->bx=(uint16_t)(ec_math(m,m->bx,1280,16,0,1));
L_202BC: ; /* IDA 0x202BC; retn */
    return;
}
/* Original sub_202BD; IDA 0x202BD. */
static void mx_body_202BD(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x202bd:goto L_202BD; default:assert(0);}
L_202BD: ; /* IDA 0x202BD; push    bx */
    ec_push(m,m->bx);
L_202BE: ; /* IDA 0x202BE; mov     ax, di */
    m->ax=(uint16_t)(m->di);
L_202C0: ; /* IDA 0x202C0; sub     ax, 268h */
    m->ax=(uint16_t)(ec_math(m,m->ax,616,16,0,1));
L_202C3: ; /* IDA 0x202C3; jbe     short loc_202CD */
    if(m->flags&(1|0x40)) goto L_202CD;
L_202C5: ; /* IDA 0x202C5; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_202C7: ; /* IDA 0x202C7; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_202C9: ; /* IDA 0x202C9; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_202CB: ; /* IDA 0x202CB; sub     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,1));
L_202CD: ; /* IDA 0x202CD; mov     si, 35Bh */
    m->si=(uint16_t)(859);
L_202D0: ; /* IDA 0x202D0; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_202D3: ; /* IDA 0x202D3; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_202D6: ; /* IDA 0x202D6; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_202D7: ; /* IDA 0x202D7; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_202D9: ; /* IDA 0x202D9; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_202DA: ; /* IDA 0x202DA; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_202DB: ; /* IDA 0x202DB; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_202DD: ; /* IDA 0x202DD; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_202DE: ; /* IDA 0x202DE; call    sub_202FE */
    mx_call(m,766,737,0x20000,h);
L_202E1: ; /* IDA 0x202E1; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_202E4: ; /* IDA 0x202E4; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_202E6: ; /* IDA 0x202E6; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_202E7: ; /* IDA 0x202E7; call    sub_202FE */
    mx_call(m,766,746,0x20000,h);
L_202EA: ; /* IDA 0x202EA; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_202ED: ; /* IDA 0x202ED; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_202EF: ; /* IDA 0x202EF; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_202F0: ; /* IDA 0x202F0; call    sub_202FE */
    mx_call(m,766,755,0x20000,h);
L_202F3: ; /* IDA 0x202F3; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_202F6: ; /* IDA 0x202F6; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_202F8: ; /* IDA 0x202F8; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_202F9: ; /* IDA 0x202F9; call    sub_202FE */
    mx_call(m,766,764,0x20000,h);
L_202FC: ; /* IDA 0x202FC; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_202FD: ; /* IDA 0x202FD; retn */
    return;
}
/* Original sub_202FE; IDA 0x202FE. */
static void mx_body_202FE(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x202fe:goto L_202FE; default:assert(0);}
L_202FE: ; /* IDA 0x202FE; mov     ch, 10h */
    vg_ch(m,(uint8_t)(16));
L_20300: ; /* IDA 0x20300; mov     ax, [bx] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx)));
L_20302: ; /* IDA 0x20302; mov     dl, [bx+2] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+2))));
L_20305: ; /* IDA 0x20305; mov     cs:[si], ax */
    dl_write16(m,m->cs,(uint16_t)(m->si),m->ax);
L_20308: ; /* IDA 0x20308; mov     cs:[si+2], dl */
    dl_write8(m,m->cs,(uint16_t)(m->si+2),(uint8_t)m->dx);
L_2030C: ; /* IDA 0x2030C; add     si, 3 */
    m->si=(uint16_t)(ec_math(m,m->si,3,16,0,0));
L_2030F: ; /* IDA 0x2030F; add     bx, 50h ; 'P' */
    m->bx=(uint16_t)(ec_math(m,m->bx,80,16,0,0));
L_20312: ; /* IDA 0x20312; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_20314: ; /* IDA 0x20314; jnz     short loc_20300 */
    if(!(m->flags&0x40)) goto L_20300;
L_20316: ; /* IDA 0x20316; sub     bx, 500h */
    m->bx=(uint16_t)(ec_math(m,m->bx,1280,16,0,1));
L_2031A: ; /* IDA 0x2031A; retn */
    return;
}
/* Original loc_1222B; IDA 0x1222B. */
static void mx_body_1222B(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x1222b:goto L_1222B; default:assert(0);}
L_1222B: ; /* IDA 0x1222B; push    ax */
    ec_push(m,m->ax);
L_1222C: ; /* IDA 0x1222C; push    bx */
    ec_push(m,m->bx);
L_1222D: ; /* IDA 0x1222D; push    cx */
    ec_push(m,m->cx);
L_1222E: ; /* IDA 0x1222E; push    dx */
    ec_push(m,m->dx);
L_1222F: ; /* IDA 0x1222F; push    si */
    ec_push(m,m->si);
L_12230: ; /* IDA 0x12230; mov     cs:byte_10D2D, 0 */
    dl_write8(m,m->cs,3373,0);
L_12236: ; /* IDA 0x12236; mov     si, 0Ah */
#if KI_INPUT_MUTATION == 3
    m->si=(uint16_t)(9);
#else
    m->si=(uint16_t)(10);
#endif
L_12239: ; /* IDA 0x12239; jmp     short near ptr byte_12246+4 */
#if KI_INPUT_MUTATION == 2
    {uint16_t patch=dl_read16(m,m->cs,0x2239);assert(patch==0x0feb||patch==0x9090);}
#else
    {uint16_t patch=dl_read16(m,m->cs,0x2239);if(patch==0x0feb) goto L_1224A;assert(patch==0x9090);}
#endif
L_1223B: ; /* IDA 0x1223B; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_1223E: ; /* IDA 0x1223E; mov     bx, 0 */
    m->bx=(uint16_t)(0);
L_12241: ; /* IDA 0x12241; call    sub_20000 */
    mx_far(m,0x2241,0x2246,h);
L_12246: ; /* IDA 0x12246; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_12248: ; /* IDA 0x12248; jnz     short near ptr byte_12255+15h */
    if(!(m->flags&0x40)) goto L_1226A;
L_1224A: ; /* IDA 0x1224A; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_1224D: ; /* IDA 0x1224D; mov     bx, 1 */
    m->bx=(uint16_t)(1);
L_12250: ; /* IDA 0x12250; call    sub_20000 */
    mx_far(m,0x2250,0x2255,h);
L_12255: ; /* IDA 0x12255; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_12257: ; /* IDA 0x12257; jnz     short near ptr unk_1226A */
    if(!(m->flags&0x40)) goto L_1226A;
L_12259: ; /* IDA 0x12259; cmp     cs:byte_10D2D, 0FFh */
    m->ip=0x2259;if(h&&h->enter) h->enter(m,0x2259,h->user);
    ec_cmp(m,dl_read8(m,m->cs,3373),255,8);
L_1225F: ; /* IDA 0x1225F; jb      short loc_12239 */
    if(m->flags&1) goto L_12239;
L_12261: ; /* IDA 0x12261; mov     cs:byte_10D2D, 0 */
    dl_write8(m,m->cs,3373,0);
L_12267: ; /* IDA 0x12267; dec     si */
    m->si=(uint16_t)(st_dec(m,m->si,16));
L_12268: ; /* IDA 0x12268; jnz     short loc_12239 */
    if(!(m->flags&0x40)) goto L_12239;
L_1226A: ; /* IDA 0x1226A; mov     cs:byte_10D2D, 0 */
    dl_write8(m,m->cs,3373,0);
L_12270: ; /* IDA 0x12270; mov     cs:byte_10D2C, 0 */
    dl_write8(m,m->cs,3372,0);
L_12276: ; /* IDA 0x12276; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_12279: ; /* IDA 0x12279; xor     bx, bx */
#if KI_INPUT_MUTATION == 4
    ; /* lost pending-left cleanup */
#else
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
#endif
L_1227B: ; /* IDA 0x1227B; call    sub_20000 */
    mx_far(m,0x227b,0x2280,h);
L_12280: ; /* IDA 0x12280; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_12281: ; /* IDA 0x12281; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_12282: ; /* IDA 0x12282; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_12283: ; /* IDA 0x12283; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_12284: ; /* IDA 0x12284; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_12285: ; /* IDA 0x12285; retn */
    return;
}
/* Original (raw code); IDA 0x20101. */
static void mx_body_20101(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {
switch(entry) {case 0x20101:goto L_20101; default:assert(0);}
L_20101: ; /* IDA 0x20101; pushf */
    ec_push(m,m->flags);
L_20102: ; /* IDA 0x20102; cli */
    m->flags&=(uint16_t)~0x200u;
L_20103: ; /* IDA 0x20103; cmp     cs:byte_20100, 1 */
    ec_cmp(m,dl_read8(m,m->cs,256),1,8);
L_20109: ; /* IDA 0x20109; jnz     short loc_20134 */
    if(!(m->flags&0x40)) goto L_20134;
L_2010B: ; /* IDA 0x2010B; cmp     cx, cs:word_200FC */
    ec_cmp(m,m->cx,dl_read16(m,m->cs,252),16);
L_20110: ; /* IDA 0x20110; jnz     short loc_20119 */
    if(!(m->flags&0x40)) goto L_20119;
L_20112: ; /* IDA 0x20112; cmp     dx, cs:word_200FE */
    ec_cmp(m,m->dx,dl_read16(m,m->cs,254),16);
L_20117: ; /* IDA 0x20117; jz      short loc_20134 */
    if(m->flags&0x40) goto L_20134;
L_20119: ; /* IDA 0x20119; sub     sp, 8 */
    m->sp=(uint16_t)(ec_math(m,m->sp,8,16,0,1));
L_2011C: ; /* IDA 0x2011C; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_2011E: ; /* IDA 0x2011E; call    sub_2016B */
    mx_call(m,363,289,0x20000,h);
L_20121: ; /* IDA 0x20121; xchg    cx, cs:word_200FC */
    {uint16_t a=m->cx,b=dl_read16(m,m->cs,252);m->cx=(uint16_t)(b);dl_write16(m,m->cs,252,a);}
L_20126: ; /* IDA 0x20126; xchg    dx, cs:word_200FE */
    {uint16_t a=m->dx,b=dl_read16(m,m->cs,254);m->dx=(uint16_t)(b);dl_write16(m,m->cs,254,a);}
L_2012B: ; /* IDA 0x2012B; call    sub_20137 */
    mx_call(m,311,302,0x20000,h);
L_2012E: ; /* IDA 0x2012E; call    sub_2019D */
    mx_call(m,413,305,0x20000,h);
L_20131: ; /* IDA 0x20131; add     sp, 8 */
    m->sp=(uint16_t)(ec_math(m,m->sp,8,16,0,0));
L_20134: ; /* IDA 0x20134; sti */
    m->flags|=0x200;
L_20135: ; /* IDA 0x20135; popf */
#if KI_INPUT_MUTATION == 11
    ; /* lost callback POPF */
#else
    m->flags=ec_pop(m);
#endif
L_20136: ; /* IDA 0x20136; retf */
    return;
}
int ki_input_body(KiMachine16 *m,uint32_t loc,const KiEconomyHooks *h) {switch(loc) {
case 0x18810:mx_body_18810(m,h,loc);break;
case 0x121e7:mx_body_121E7(m,h,loc);break;
case 0x20000:mx_body_20000(m,h,loc);break;
case 0x2002e:mx_body_2002E(m,h,loc);break;
case 0x2006f:mx_body_2002E(m,h,loc);break;
case 0x20070:mx_body_20070(m,h,loc);break;
case 0x2009a:mx_body_2009A(m,h,loc);break;
case 0x200bd:mx_body_200BD(m,h,loc);break;
case 0x200c0:mx_body_200C0(m,h,loc);break;
case 0x20137:mx_body_20137(m,h,loc);break;
case 0x2016b:mx_body_2016B(m,h,loc);break;
case 0x2019d:mx_body_2019D(m,h,loc);break;
case 0x201c6:mx_body_201C6(m,h,loc);break;
case 0x201e4:mx_body_201E4(m,h,loc);break;
case 0x2020c:mx_body_2020C(m,h,loc);break;
case 0x20249:mx_body_20249(m,h,loc);break;
case 0x202a0:mx_body_202A0(m,h,loc);break;
case 0x202bd:mx_body_202BD(m,h,loc);break;
case 0x202fe:mx_body_202FE(m,h,loc);break;
case 0x1222b:mx_body_1222B(m,h,loc);break;
case 0x20101:mx_body_20101(m,h,loc);break;
default:return 0;}return 1;}
