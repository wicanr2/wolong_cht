/* Original native C nonlocal exit; re/125, spec/245. Navigation-only main loop is not generated. */
/* Original sub_10A1C; IDA 0x10A1C. */
static void nle_body_10A1C(KiMachine16 *m,const KiEconomyHooks *h) {
L_10A1C: ; /* IDA 0x10A1C; push    ax */
    ec_push(m,m->ax);
L_10A1D: ; /* IDA 0x10A1D; push    bx */
    ec_push(m,m->bx);
L_10A1E: ; /* IDA 0x10A1E; push    cx */
    ec_push(m,m->cx);
L_10A1F: ; /* IDA 0x10A1F; push    dx */
    ec_push(m,m->dx);
L_10A20: ; /* IDA 0x10A20; push    si */
    ec_push(m,m->si);
L_10A21: ; /* IDA 0x10A21; mov     cl, 10h */
#if KI_MAIN_MUTATION == 5
    vg_cl(m,15);
#else
    vg_cl(m,(uint8_t)(16));
#endif
L_10A23: ; /* IDA 0x10A23; mov     ch, 0 */
    vg_ch(m,(uint8_t)(0));
L_10A25: ; /* IDA 0x10A25; mov     si, 1964h */
    m->si=(uint16_t)(6500);
L_10A28: ; /* IDA 0x10A28; mov     dx, 3DAh */
    m->dx=(uint16_t)(986);
L_10A2B: ; /* IDA 0x10A2B; in      al, dx; Video status bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
L_10A2C: ; /* IDA 0x10A2C; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_10A2E: ; /* IDA 0x10A2E; jnz     short loc_10A2B */
    if(!(m->flags&0x40)) goto L_10A2B;
L_10A30: ; /* IDA 0x10A30; in      al, dx; Video status bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
L_10A31: ; /* IDA 0x10A31; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_10A33: ; /* IDA 0x10A33; jz      short loc_10A30 */
    if(m->flags&0x40) goto L_10A30;
L_10A35: ; /* IDA 0x10A35; in      al, dx; Video status bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
L_10A36: ; /* IDA 0x10A36; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_10A38: ; /* IDA 0x10A38; jnz     short loc_10A35 */
    if(!(m->flags&0x40)) goto L_10A35;
L_10A3A: ; /* IDA 0x10A3A; in      al, dx; Video status bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
L_10A3B: ; /* IDA 0x10A3B; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_10A3D: ; /* IDA 0x10A3D; jz      short loc_10A3A */
    if(m->flags&0x40) goto L_10A3A;
L_10A3F: ; /* IDA 0x10A3F; mov     bh, cs:[si] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->si))));
L_10A42: ; /* IDA 0x10A42; mov     dl, cs:[si+1] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->si+1))));
L_10A46: ; /* IDA 0x10A46; mov     dh, cs:[si+2] */
    vg_dh(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->si+2))));
L_10A4A: ; /* IDA 0x10A4A; mov     al, ch */
    vg_al(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_10A4C: ; /* IDA 0x10A4C; mov     bl, cl */
    gw_bl(m,(uint8_t)((uint8_t)m->cx));
L_10A4E: ; /* IDA 0x10A4E; call    sub_1EBDC */
    nle_call(m,60380,2641,h);
L_10A51: ; /* IDA 0x10A51; add     si, 3 */
    m->si=(uint16_t)(ec_math(m,m->si,3,16,0,0));
L_10A54: ; /* IDA 0x10A54; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_10A56: ; /* IDA 0x10A56; cmp     ch, 10h */
    ec_cmp(m,(uint8_t)(m->cx>>8),16,8);
L_10A59: ; /* IDA 0x10A59; jb      short loc_10A3F */
    if(m->flags&1) goto L_10A3F;
L_10A5B: ; /* IDA 0x10A5B; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_10A5D: ; /* IDA 0x10A5D; jge     short loc_10A23 */
    if(!ec_less(m)) goto L_10A23;
L_10A5F: ; /* IDA 0x10A5F; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_10A60: ; /* IDA 0x10A60; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10A61: ; /* IDA 0x10A61; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10A62: ; /* IDA 0x10A62; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10A63: ; /* IDA 0x10A63; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10A64: ; /* IDA 0x10A64; retn */
    return;
}
/* Original sub_11CB1; IDA 0x11CB1. */
static void nle_body_11CB1(KiMachine16 *m,const KiEconomyHooks *h) {
L_11CB1: ; /* IDA 0x11CB1; mov     ss, cs:word_19903 */
#if KI_MAIN_MUTATION == 1
    m->ss=dl_read16(m,m->cs,0x9901);
#else
    m->ss=(uint16_t)(dl_read16(m,m->cs,39171));
#endif
L_11CB6: ; /* IDA 0x11CB6; mov     sp, cs:word_19901 */
#if KI_MAIN_MUTATION == 2
    m->sp=dl_read16(m,m->cs,0x9903);
#else
    m->sp=(uint16_t)(dl_read16(m,m->cs,39169));
#endif
L_11CBB: ; /* IDA 0x11CBB; push    ax */
#if KI_MAIN_MUTATION == 3
    ; /* lost AX push on restored outer stack */
#else
    ec_push(m,m->ax);
#endif
L_11CBC: ; /* IDA 0x11CBC; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_11CBF: ; /* IDA 0x11CBF; call    sub_20000 */
    mx_far(m,0x1cbf,0x1cc4,h);
L_11CC4: ; /* IDA 0x11CC4; call    sub_10A1C */
    nle_call(m,2588,7367,h);
L_11CC7: ; /* IDA 0x11CC7; pop     ax */
#if KI_MAIN_MUTATION == 4
    m->ax=ec_pop(m)^1;
#else
    m->ax=(uint16_t)(ec_pop(m));
#endif
L_11CC8: ; /* IDA 0x11CC8; retn */
    return;
}
/* Original sub_1EBDC; IDA 0x1EBDC. */
static void nle_body_1EBDC(KiMachine16 *m,const KiEconomyHooks *h) {
L_1EBDC: ; /* IDA 0x1EBDC; push    cx */
    ec_push(m,m->cx);
L_1EBDD: ; /* IDA 0x1EBDD; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1EBDF: ; /* IDA 0x1EBDF; mov     dx, 3C8h */
    m->dx=(uint16_t)(968);
L_1EBE2: ; /* IDA 0x1EBE2; out     dx, al */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EBE3: ; /* IDA 0x1EBE3; mov     dx, 3C9h */
    m->dx=(uint16_t)(969);
L_1EBE6: ; /* IDA 0x1EBE6; mov     al, bh */
    vg_al(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_1EBE8: ; /* IDA 0x1EBE8; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EBEA: ; /* IDA 0x1EBEA; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EBEC: ; /* IDA 0x1EBEC; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EBEE: ; /* IDA 0x1EBEE; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EBF0: ; /* IDA 0x1EBF0; mul     bl */
    dl_mul(m,(uint8_t)m->bx,8);
L_1EBF2: ; /* IDA 0x1EBF2; add     ax, 80h */
#if KI_MAIN_MUTATION == 6
    m->ax=ec_math(m,m->ax,0,16,0,0);
#else
    m->ax=(uint16_t)(ec_math(m,m->ax,128,16,0,0));
#endif
L_1EBF5: ; /* IDA 0x1EBF5; shl     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,4)));
L_1EBF7: ; /* IDA 0x1EBF7; shl     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,4)));
L_1EBF9: ; /* IDA 0x1EBF9; mov     bh, ah */
    gw_bh(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1EBFB: ; /* IDA 0x1EBFB; mov     al, ch */
    vg_al(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_1EBFD: ; /* IDA 0x1EBFD; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EBFF: ; /* IDA 0x1EBFF; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EC01: ; /* IDA 0x1EC01; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EC03: ; /* IDA 0x1EC03; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EC05: ; /* IDA 0x1EC05; mul     bl */
    dl_mul(m,(uint8_t)m->bx,8);
L_1EC07: ; /* IDA 0x1EC07; add     ax, 80h */
    m->ax=(uint16_t)(ec_math(m,m->ax,128,16,0,0));
L_1EC0A: ; /* IDA 0x1EC0A; shl     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,4)));
L_1EC0C: ; /* IDA 0x1EC0C; shl     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,4)));
L_1EC0E: ; /* IDA 0x1EC0E; mov     ch, ah */
    vg_ch(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1EC10: ; /* IDA 0x1EC10; mov     al, cl */
    vg_al(m,(uint8_t)((uint8_t)m->cx));
L_1EC12: ; /* IDA 0x1EC12; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EC14: ; /* IDA 0x1EC14; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EC16: ; /* IDA 0x1EC16; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EC18: ; /* IDA 0x1EC18; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_1EC1A: ; /* IDA 0x1EC1A; mul     bl */
    dl_mul(m,(uint8_t)m->bx,8);
L_1EC1C: ; /* IDA 0x1EC1C; add     ax, 80h */
    m->ax=(uint16_t)(ec_math(m,m->ax,128,16,0,0));
L_1EC1F: ; /* IDA 0x1EC1F; shl     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,4)));
L_1EC21: ; /* IDA 0x1EC21; shl     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,4)));
L_1EC23: ; /* IDA 0x1EC23; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1EC25: ; /* IDA 0x1EC25; out     dx, al */
#if KI_MAIN_MUTATION == 7
    dl_out(m,m->dx,(uint8_t)(m->bx>>8),8);
#else
    dl_out(m,m->dx,(uint8_t)m->ax,8);
#endif
L_1EC26: ; /* IDA 0x1EC26; mov     al, ch */
    vg_al(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_1EC28: ; /* IDA 0x1EC28; out     dx, al */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EC29: ; /* IDA 0x1EC29; mov     al, bh */
    vg_al(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_1EC2B: ; /* IDA 0x1EC2B; out     dx, al */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EC2C: ; /* IDA 0x1EC2C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1EC2D: ; /* IDA 0x1EC2D; retn */
    return;
}
int ki_main_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0xa1c:nle_body_10A1C(m,h);break;
case 0x1cb1:nle_body_11CB1(m,h);break;
case 0xebdc:nle_body_1EBDC(m,h);break;
default:return 0;}
if(t==0x1cb1) {
m->ip=ec_pop(m); /* Original RET consumes the restored outer frame. */
#if KI_MAIN_MUTATION == 8
return 1; /* Incorrectly resume old native caller after virtual nonlocal RET. */
#else
wolong_main_transfer(m,h?h->user:0);
#endif
}
return 1;}
static int nle_known(uint16_t t) {switch(t) {case 0xa1c: case 0x1cb1: case 0xebdc: return 1;default:return 0;}}
void sub_10A1C(KiMachine16 *m,const KiEconomyHooks *h) {ki_main_body(m,0xa1c,h);m->ip=ec_pop(m);}
void sub_11CB1(KiMachine16 *m,const KiEconomyHooks *h) {ki_main_body(m,0x1cb1,h);m->ip=ec_pop(m);}
void sub_1EBDC(KiMachine16 *m,const KiEconomyHooks *h) {ki_main_body(m,0xebdc,h);m->ip=ec_pop(m);}
