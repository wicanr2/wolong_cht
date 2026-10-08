/* Generated from fixed IDA evidence; native C statements, no opcode fetch.
 * Original names/address/operands remain in comments; spec/227, re/107. */
static void dl_body_1030F(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_10337(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1E993(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1E9A7(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1E9C1(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1FB29(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1FBA7(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1F6DC(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1F878(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EA1F(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EA26(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EA30(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EA3A(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EA44(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EA5A(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EA70(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EA96(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EAB3(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EAD0(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1EDFE(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1F26E(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1F3CF(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1F426(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1F43B(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
static void dl_body_1F465(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);
/* original sub_1030F; IDA 0x1030F; raw control flow; evidence re/107 */
static void dl_body_1030F(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0x30f:goto L_1030F; default:assert(0);}
L_1030F: ; /* IDA 0x1030f; push    ax */
    ec_push(m,m->ax);
L_10310: ; /* IDA 0x10310; push    bx */
    ec_push(m,m->bx);
L_10311: ; /* IDA 0x10311; push    cx */
    ec_push(m,m->cx);
L_10312: ; /* IDA 0x10312; push    dx */
    ec_push(m,m->dx);
L_10313: ; /* IDA 0x10313; mov     dx, cs */
    m->dx=(uint16_t)(m->cs);
L_10315: ; /* IDA 0x10315; mov     ax, 0E16h */
    m->ax=(uint16_t)(3606);
L_10318: ; /* IDA 0x10318; mov     bx, cs */
    m->bx=(uint16_t)(m->cs);
L_1031A: ; /* IDA 0x1031a; mov     cx, cs:word_10D50 */
    m->cx=(uint16_t)(dl_read16(m,m->cs,3408));
L_1031F: ; /* IDA 0x1031f; call    sub_1E993 */
    dl_call(m,59795,802,h);
L_10322: ; /* IDA 0x10322; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_10324: ; /* IDA 0x10324; mov     dx, cs:word_10D50 */
    m->dx=(uint16_t)(dl_read16(m,m->cs,3408));
L_10329: ; /* IDA 0x10329; mov     ax, 1800h */
    m->ax=(uint16_t)(6144);
L_1032C: ; /* IDA 0x1032c; mov     cx, 2020h */
    m->cx=(uint16_t)(8224);
L_1032F: ; /* IDA 0x1032f; call    sub_1E9A7 */
    dl_call(m,59815,818,h);
L_10332: ; /* IDA 0x10332; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10333: ; /* IDA 0x10333; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10334: ; /* IDA 0x10334; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10335: ; /* IDA 0x10335; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10336: ; /* IDA 0x10336; retn */
    return;
}
/* original sub_10337; IDA 0x10337; raw control flow; evidence re/107 */
static void dl_body_10337(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0x337:goto L_10337; default:assert(0);}
L_10337: ; /* IDA 0x10337; call    sub_1E9C1 */
    dl_call(m,59841,826,h);
L_1033A: ; /* IDA 0x1033a; retn */
    return;
}
/* original sub_1E993; IDA 0x1E993; raw control flow; evidence re/107 */
static void dl_body_1E993(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xe993:goto L_1E993; default:assert(0);}
L_1E993: ; /* IDA 0x1e993; mov     word ptr cs:dword_1EAE9, ax */
    dl_write16(m,m->cs,60137,m->ax);
L_1E997: ; /* IDA 0x1e997; mov     word ptr cs:dword_1EAE9+2, dx */
    dl_write16(m,m->cs,60139,m->dx);
L_1E99C: ; /* IDA 0x1e99c; mov     cs:word_1EAED, bx */
    dl_write16(m,m->cs,60141,m->bx);
L_1E9A1: ; /* IDA 0x1e9a1; mov     cs:word_1EAEF, cx */
    dl_write16(m,m->cs,60143,m->cx);
L_1E9A6: ; /* IDA 0x1e9a6; retn */
    return;
}
/* original sub_1E9A7; IDA 0x1E9A7; raw control flow; evidence re/107 */
static void dl_body_1E9A7(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xe9a7:goto L_1E9A7; default:assert(0);}
L_1E9A7: ; /* IDA 0x1e9a7; push    bx */
    ec_push(m,m->bx);
L_1E9A8: ; /* IDA 0x1e9a8; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1E9AA: ; /* IDA 0x1e9aa; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1E9AC: ; /* IDA 0x1e9ac; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1E9AE: ; /* IDA 0x1e9ae; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1E9B0: ; /* IDA 0x1e9b0; add     bx, 0EAF1h */
#if KI_DISPLAY_MUTATION == 10
    m->bx=(uint16_t)(ec_math(m,m->bx,60146,16,0,0));
#else
    m->bx=(uint16_t)(ec_math(m,m->bx,60145,16,0,0));
#endif
L_1E9B4: ; /* IDA 0x1e9b4; mov     cs:[bx], ax */
    dl_write16(m,m->cs,(uint16_t)(m->bx),m->ax);
L_1E9B7: ; /* IDA 0x1e9b7; mov     cs:[bx+2], dx */
    dl_write16(m,m->cs,(uint16_t)(m->bx+2),m->dx);
L_1E9BB: ; /* IDA 0x1e9bb; mov     cs:[bx+4], cx */
    dl_write16(m,m->cs,(uint16_t)(m->bx+4),m->cx);
L_1E9BF: ; /* IDA 0x1e9bf; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1E9C0: ; /* IDA 0x1e9c0; retn */
    return;
}
/* original sub_1E9C1; IDA 0x1E9C1; raw control flow; evidence re/107 */
static void dl_body_1E9C1(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xe9c1:goto L_1E9C1; default:assert(0);}
L_1E9C1: ; /* IDA 0x1e9c1; push    ds */
    ec_push(m,m->ds);
L_1E9C2: ; /* IDA 0x1e9c2; push    es */
    ec_push(m,m->es);
L_1E9C3: ; /* IDA 0x1e9c3; push    ax */
    ec_push(m,m->ax);
L_1E9C4: ; /* IDA 0x1e9c4; push    bx */
    ec_push(m,m->bx);
L_1E9C5: ; /* IDA 0x1e9c5; push    cx */
    ec_push(m,m->cx);
L_1E9C6: ; /* IDA 0x1e9c6; push    dx */
    ec_push(m,m->dx);
L_1E9C7: ; /* IDA 0x1e9c7; push    si */
    ec_push(m,m->si);
L_1E9C8: ; /* IDA 0x1e9c8; push    di */
    ec_push(m,m->di);
L_1E9C9: ; /* IDA 0x1e9c9; push    bp */
    ec_push(m,m->bp);
L_1E9CA: ; /* IDA 0x1e9ca; sub     sp, 4 */
    m->sp=(uint16_t)(ec_math(m,m->sp,4,16,0,1));
L_1E9CD: ; /* IDA 0x1e9cd; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_1E9CF: ; /* IDA 0x1e9cf; mov     [bp+0], dx */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->dx);
L_1E9D2: ; /* IDA 0x1e9d2; mov     [bp+2], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->bx);
L_1E9D5: ; /* IDA 0x1e9d5; lds     si, cs:dword_1EAE9 */
    {uint16_t off=60137,seg=m->cs;uint16_t value=dl_read16(m,seg,off),bank=dl_read16(m,seg,(uint16_t)(off+2));m->si=(uint16_t)(value);m->ds=bank;}
L_1E9DA: ; /* IDA 0x1e9da; cmp     byte ptr [si], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),0,8);
L_1E9DD: ; /* IDA 0x1e9dd; jnz     short loc_1E9E5 */
    if(!(m->flags&0x40)) goto L_1E9E5;
L_1E9DF: ; /* IDA 0x1e9df; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1E9E1: ; /* IDA 0x1e9e1; jz      short loc_1E9EA */
    if(m->flags&0x40) goto L_1E9EA;
L_1E9E3: ; /* IDA 0x1e9e3; dec     al */
    vg_al(m,(uint8_t)(st_dec(m,(uint8_t)m->ax,8)));
L_1E9E5: ; /* IDA 0x1e9e5; add     si, 0Ch */
    m->si=(uint16_t)(ec_math(m,m->si,12,16,0,0));
L_1E9E8: ; /* IDA 0x1e9e8; jmp     short loc_1E9DA */
    goto L_1E9DA;
L_1E9EA: ; /* IDA 0x1e9ea; add     si, 0Ch */
#if KI_DISPLAY_MUTATION == 1
    m->si=(uint16_t)(ec_math(m,m->si,24,16,0,0));
#else
    m->si=(uint16_t)(ec_math(m,m->si,12,16,0,0));
#endif
L_1E9ED: ; /* IDA 0x1e9ed; mov     bl, [si] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
L_1E9EF: ; /* IDA 0x1e9ef; and     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)&((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_1E9F1: ; /* IDA 0x1e9f1; jz      short loc_1EA00 */
    if(m->flags&0x40) goto L_1EA00;
L_1E9F3: ; /* IDA 0x1e9f3; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1E9F5: ; /* IDA 0x1e9f5; dec     bx */
    m->bx=(uint16_t)(st_dec(m,m->bx,16));
L_1E9F6: ; /* IDA 0x1e9f6; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1E9F8: ; /* IDA 0x1e9f8; push    si */
    ec_push(m,m->si);
L_1E9F9: ; /* IDA 0x1e9f9; call    word ptr [bx-15F3h] */
#if KI_DISPLAY_MUTATION == 3
    dl_call(m,dl_read16(m,m->cs,(uint16_t)(m->bx+59917)),59901,h);
#else
    dl_call(m,dl_read16(m,m->ds,(uint16_t)(m->bx+59917)),59901,h);
#endif
L_1E9FD: ; /* IDA 0x1e9fd; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1E9FE: ; /* IDA 0x1e9fe; jmp     short loc_1E9EA */
    goto L_1E9EA;
L_1EA00: ; /* IDA 0x1ea00; add     sp, 4 */
    m->sp=(uint16_t)(ec_math(m,m->sp,4,16,0,0));
L_1EA03: ; /* IDA 0x1ea03; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1EA04: ; /* IDA 0x1ea04; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1EA05: ; /* IDA 0x1ea05; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1EA06: ; /* IDA 0x1ea06; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1EA07: ; /* IDA 0x1ea07; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1EA08: ; /* IDA 0x1ea08; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1EA09: ; /* IDA 0x1ea09; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1EA0A: ; /* IDA 0x1ea0a; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1EA0B: ; /* IDA 0x1ea0b; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1EA0C: ; /* IDA 0x1ea0c; retn */
    return;
}
/* original sub_1FB29; IDA 0x1FB29; raw control flow; evidence re/107 */
static void dl_body_1FB29(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xfb29:goto L_1FB29; default:assert(0);}
L_1FB29: ; /* IDA 0x1fb29; push    ds */
    ec_push(m,m->ds);
L_1FB2A: ; /* IDA 0x1fb2a; push    ax */
    ec_push(m,m->ax);
L_1FB2B: ; /* IDA 0x1fb2b; push    bx */
    ec_push(m,m->bx);
L_1FB2C: ; /* IDA 0x1fb2c; push    cx */
    ec_push(m,m->cx);
L_1FB2D: ; /* IDA 0x1fb2d; push    dx */
    ec_push(m,m->dx);
L_1FB2E: ; /* IDA 0x1fb2e; push    cx */
    ec_push(m,m->cx);
L_1FB2F: ; /* IDA 0x1fb2f; mov     cx, 0A0C8h */
    m->cx=(uint16_t)(41160);
L_1FB32: ; /* IDA 0x1fb32; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_1FB34: ; /* IDA 0x1fb34; mov     cx, 4 */
    m->cx=(uint16_t)(4);
L_1FB37: ; /* IDA 0x1fb37; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1FB39: ; /* IDA 0x1fb39; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1FB3B: ; /* IDA 0x1fb3b; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1FB3D: ; /* IDA 0x1fb3d; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1FB3F: ; /* IDA 0x1fb3f; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1FB41: ; /* IDA 0x1fb41; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1FB43: ; /* IDA 0x1fb43; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1FB45: ; /* IDA 0x1fb45; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1FB47: ; /* IDA 0x1fb47; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_1FB49: ; /* IDA 0x1fb49; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1FB4B: ; /* IDA 0x1fb4b; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1FB4D: ; /* IDA 0x1fb4d; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_1FB4F: ; /* IDA 0x1fb4f; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1FB52: ; /* IDA 0x1fb52; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_1FB54: ; /* IDA 0x1fb54; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FB55: ; /* IDA 0x1fb55; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1FB56: ; /* IDA 0x1fb56; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1FB58: ; /* IDA 0x1fb58; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FB59: ; /* IDA 0x1fb59; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1FB5A: ; /* IDA 0x1fb5a; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1FB5C: ; /* IDA 0x1fb5c; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FB5D: ; /* IDA 0x1fb5d; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1FB5E: ; /* IDA 0x1fb5e; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1FB60: ; /* IDA 0x1fb60; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FB61: ; /* IDA 0x1fb61; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1FB62: ; /* IDA 0x1fb62; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1FB64: ; /* IDA 0x1fb64; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FB65: ; /* IDA 0x1fb65; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1FB66: ; /* IDA 0x1fb66; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1FB68: ; /* IDA 0x1fb68; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FB69: ; /* IDA 0x1fb69; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1FB6A: ; /* IDA 0x1fb6a; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_1FB6C: ; /* IDA 0x1fb6c; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FB6D: ; /* IDA 0x1fb6d; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1FB6E: ; /* IDA 0x1fb6e; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1FB70: ; /* IDA 0x1fb70; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FB71: ; /* IDA 0x1fb71; and     cx, 7 */
    {uint16_t result=((m->cx)&(7))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_1FB74: ; /* IDA 0x1fb74; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1FB76: ; /* IDA 0x1fb76; shr     al, cl */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,(uint8_t)m->cx,8,5)));
L_1FB78: ; /* IDA 0x1fb78; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1FB79: ; /* IDA 0x1fb79; add     dx, cx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->cx,16,0,0));
L_1FB7B: ; /* IDA 0x1fb7b; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1FB7D: ; /* IDA 0x1fb7d; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1FB7F: ; /* IDA 0x1fb7f; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1FB81: ; /* IDA 0x1fb81; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1FB83: ; /* IDA 0x1fb83; and     cx, 7 */
    {uint16_t result=((m->cx)&(7))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_1FB86: ; /* IDA 0x1fb86; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_1FB88: ; /* IDA 0x1fb88; shr     ah, cl */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),(uint8_t)m->cx,8,5)));
L_1FB8A: ; /* IDA 0x1fb8a; not     ah */
    vg_ah(m,(uint8_t)(~((uint8_t)(m->ax>>8))));
L_1FB8C: ; /* IDA 0x1fb8c; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1FB8E: ; /* IDA 0x1fb8e; mov     dl, [bx] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx))));
L_1FB90: ; /* IDA 0x1fb90; mov     [bx], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx),(uint8_t)m->ax);
L_1FB92: ; /* IDA 0x1fb92; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1FB93: ; /* IDA 0x1fb93; dec     cx */
    m->cx=(uint16_t)(st_dec(m,m->cx,16));
L_1FB94: ; /* IDA 0x1fb94; jz      short loc_1FB9D */
    if(m->flags&0x40) goto L_1FB9D;
L_1FB96: ; /* IDA 0x1fb96; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1FB98: ; /* IDA 0x1fb98; mov     [bx], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx),(uint8_t)m->ax);
L_1FB9A: ; /* IDA 0x1fb9a; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1FB9B: ; /* IDA 0x1fb9b; loop    loc_1FB98 */
    if(--m->cx) goto L_1FB98;
L_1FB9D: ; /* IDA 0x1fb9d; mov     dl, [bx] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx))));
L_1FB9F: ; /* IDA 0x1fb9f; mov     [bx], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx),(uint8_t)(m->ax>>8));
L_1FBA1: ; /* IDA 0x1fba1; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1FBA2: ; /* IDA 0x1fba2; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1FBA3: ; /* IDA 0x1fba3; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1FBA4: ; /* IDA 0x1fba4; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1FBA5: ; /* IDA 0x1fba5; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1FBA6: ; /* IDA 0x1fba6; retn */
    return;
}
/* original sub_1FBA7; IDA 0x1FBA7; raw control flow; evidence re/107 */
static void dl_body_1FBA7(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xfba7:goto L_1FBA7; default:assert(0);}
L_1FBA7: ; /* IDA 0x1fba7; push    ds */
    ec_push(m,m->ds);
L_1FBA8: ; /* IDA 0x1fba8; push    ax */
    ec_push(m,m->ax);
L_1FBA9: ; /* IDA 0x1fba9; push    bx */
    ec_push(m,m->bx);
L_1FBAA: ; /* IDA 0x1fbaa; push    cx */
    ec_push(m,m->cx);
L_1FBAB: ; /* IDA 0x1fbab; push    dx */
    ec_push(m,m->dx);
L_1FBAC: ; /* IDA 0x1fbac; push    cx */
    ec_push(m,m->cx);
L_1FBAD: ; /* IDA 0x1fbad; mov     cx, 0A0C8h */
    m->cx=(uint16_t)(41160);
L_1FBB0: ; /* IDA 0x1fbb0; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_1FBB2: ; /* IDA 0x1fbb2; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1FBB4: ; /* IDA 0x1fbb4; shr     dx, 3 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,3,16,5));
L_1FBB7: ; /* IDA 0x1fbb7; shl     bx, 4 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,4,16,4));
L_1FBBA: ; /* IDA 0x1fbba; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_1FBBC: ; /* IDA 0x1fbbc; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1FBBE: ; /* IDA 0x1fbbe; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1FBC0: ; /* IDA 0x1fbc0; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_1FBC2: ; /* IDA 0x1fbc2; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1FBC5: ; /* IDA 0x1fbc5; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_1FBC7: ; /* IDA 0x1fbc7; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FBC8: ; /* IDA 0x1fbc8; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1FBC9: ; /* IDA 0x1fbc9; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1FBCB: ; /* IDA 0x1fbcb; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FBCC: ; /* IDA 0x1fbcc; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1FBCD: ; /* IDA 0x1fbcd; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1FBCF: ; /* IDA 0x1fbcf; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FBD0: ; /* IDA 0x1fbd0; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1FBD1: ; /* IDA 0x1fbd1; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1FBD3: ; /* IDA 0x1fbd3; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FBD4: ; /* IDA 0x1fbd4; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1FBD5: ; /* IDA 0x1fbd5; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1FBD7: ; /* IDA 0x1fbd7; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FBD8: ; /* IDA 0x1fbd8; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1FBD9: ; /* IDA 0x1fbd9; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1FBDB: ; /* IDA 0x1fbdb; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FBDC: ; /* IDA 0x1fbdc; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1FBDD: ; /* IDA 0x1fbdd; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_1FBDF: ; /* IDA 0x1fbdf; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FBE0: ; /* IDA 0x1fbe0; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1FBE1: ; /* IDA 0x1fbe1; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1FBE3: ; /* IDA 0x1fbe3; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1FBE4: ; /* IDA 0x1fbe4; and     cl, 7 */
    {uint16_t result=(((uint8_t)m->cx)&(7))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_1FBE7: ; /* IDA 0x1fbe7; mov     al, 80h */
    vg_al(m,(uint8_t)(128));
L_1FBE9: ; /* IDA 0x1fbe9; shr     al, cl */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,(uint8_t)m->cx,8,5)));
L_1FBEB: ; /* IDA 0x1fbeb; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1FBEC: ; /* IDA 0x1fbec; mov     ah, [bx] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx))));
L_1FBEE: ; /* IDA 0x1fbee; mov     [bx], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx),(uint8_t)m->ax);
L_1FBF0: ; /* IDA 0x1fbf0; add     bx, 50h ; 'P' */
#if KI_DISPLAY_MUTATION == 8
    m->bx=(uint16_t)(ec_math(m,m->bx,79,16,0,0));
#else
    m->bx=(uint16_t)(ec_math(m,m->bx,80,16,0,0));
#endif
L_1FBF3: ; /* IDA 0x1fbf3; loop    loc_1FBEC */
    if(--m->cx) goto L_1FBEC;
L_1FBF5: ; /* IDA 0x1fbf5; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1FBF6: ; /* IDA 0x1fbf6; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1FBF7: ; /* IDA 0x1fbf7; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1FBF8: ; /* IDA 0x1fbf8; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1FBF9: ; /* IDA 0x1fbf9; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1FBFA: ; /* IDA 0x1fbfa; retn */
    return;
}
/* original sub_1F6DC; IDA 0x1F6DC; raw control flow; evidence re/107 */
static void dl_body_1F6DC(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xf6dc:goto L_1F6DC; default:assert(0);}
L_1F6DC: ; /* IDA 0x1f6dc; push    cx */
    ec_push(m,m->cx);
L_1F6DD: ; /* IDA 0x1f6dd; push    di */
    ec_push(m,m->di);
L_1F6DE: ; /* IDA 0x1f6de; mov     di, 10h */
    m->di=(uint16_t)(16);
L_1F6E1: ; /* IDA 0x1f6e1; mov     ch, [si] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
L_1F6E3: ; /* IDA 0x1f6e3; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_1F6E4: ; /* IDA 0x1f6e4; cmp     ch, 0 */
    ec_cmp(m,(uint8_t)(m->cx>>8),0,8);
L_1F6E7: ; /* IDA 0x1f6e7; jz      short loc_1F71D */
    if(m->flags&0x40) goto L_1F71D;
L_1F6E9: ; /* IDA 0x1f6e9; mov     cl, [si] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
L_1F6EB: ; /* IDA 0x1f6eb; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_1F6EC: ; /* IDA 0x1f6ec; xchg    ax, cx */
    {uint16_t a=m->ax,b=m->cx;m->ax=(uint16_t)(b);m->cx=(uint16_t)(a);}
L_1F6ED: ; /* IDA 0x1f6ed; call    sub_1F878 */
    dl_call(m,63608,63216,h);
L_1F6F0: ; /* IDA 0x1f6f0; jnb     short loc_1F6F6 */
    if(!(m->flags&1)) goto L_1F6F6;
L_1F6F2: ; /* IDA 0x1f6f2; mov     di, 8 */
    m->di=(uint16_t)(8);
L_1F6F5: ; /* IDA 0x1f6f5; dec     si */
#if KI_DISPLAY_MUTATION == 5
    ; /* lost byte rewind */
#else
    m->si=(uint16_t)(st_dec(m,m->si,16));
#endif
L_1F6F6: ; /* IDA 0x1f6f6; xchg    ax, cx */
    {uint16_t a=m->ax,b=m->cx;m->ax=(uint16_t)(b);m->cx=(uint16_t)(a);}
L_1F6F7: ; /* IDA 0x1f6f7; test    al, 4 */
#if KI_DISPLAY_MUTATION == 6
    {uint16_t result=(((uint8_t)m->ax)&(0))&255;ec_logic(m,result,8);}
#else
    {uint16_t result=(((uint8_t)m->ax)&(4))&255;ec_logic(m,result,8);}
#endif
L_1F6F9: ; /* IDA 0x1f6f9; jz      short loc_1F710 */
    if(m->flags&0x40) goto L_1F710;
L_1F6FB: ; /* IDA 0x1f6fb; push    ax */
    ec_push(m,m->ax);
L_1F6FC: ; /* IDA 0x1f6fc; push    bx */
    ec_push(m,m->bx);
L_1F6FD: ; /* IDA 0x1f6fd; push    dx */
    ec_push(m,m->dx);
L_1F6FE: ; /* IDA 0x1f6fe; push    cx */
    ec_push(m,m->cx);
L_1F6FF: ; /* IDA 0x1f6ff; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1F700: ; /* IDA 0x1f700; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F701: ; /* IDA 0x1f701; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_1F703: ; /* IDA 0x1f703; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_1F705: ; /* IDA 0x1f705; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_1F707: ; /* IDA 0x1f707; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_1F709: ; /* IDA 0x1f709; call    loc_1F75E */
    dl_call(m,63326,63244,h);
L_1F70C: ; /* IDA 0x1f70c; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1F70D: ; /* IDA 0x1f70d; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1F70E: ; /* IDA 0x1f70e; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1F70F: ; /* IDA 0x1f70f; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1F710: ; /* IDA 0x1f710; push    ax */
    ec_push(m,m->ax);
L_1F711: ; /* IDA 0x1f711; push    bx */
    ec_push(m,m->bx);
L_1F712: ; /* IDA 0x1f712; push    dx */
    ec_push(m,m->dx);
L_1F713: ; /* IDA 0x1f713; call    loc_1F75E */
    dl_call(m,63326,63254,h);
L_1F716: ; /* IDA 0x1f716; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1F717: ; /* IDA 0x1f717; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1F718: ; /* IDA 0x1f718; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1F719: ; /* IDA 0x1f719; add     dx, di */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->di,16,0,0));
L_1F71B: ; /* IDA 0x1f71b; jmp     short loc_1F6DE */
    goto L_1F6DE;
L_1F71D: ; /* IDA 0x1f71d; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1F71E: ; /* IDA 0x1f71e; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1F71F: ; /* IDA 0x1f71f; retn */
    return;
}
/* original sub_1F878; IDA 0x1F878; raw control flow; evidence re/107 */
static void dl_body_1F878(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xf878:goto L_1F878; default:assert(0);}
L_1F878: ; /* IDA 0x1f878; cmp     ah, 80h */
    ec_cmp(m,(uint8_t)(m->ax>>8),128,8);
L_1F87B: ; /* IDA 0x1f87b; jb      short loc_1F880 */
    if(m->flags&1) goto L_1F880;
L_1F87D: ; /* IDA 0x1f87d; clc */
    m->flags&=(uint16_t)~1u;
L_1F87E: ; /* IDA 0x1f87e; jmp     short locret_1F887 */
    goto L_1F887;
L_1F880: ; /* IDA 0x1f880; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1F882: ; /* IDA 0x1f882; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1F884: ; /* IDA 0x1f884; stc */
    m->flags|=1;
L_1F885: ; /* IDA 0x1f885; jmp     short $+2 */
    goto L_1F887;
L_1F887: ; /* IDA 0x1f887; retn */
    return;
}
/* original (no IDA function); IDA 0x1EA1F; raw control flow; evidence re/107 */
static void dl_body_1EA1F(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xea1f:goto L_1EA1F; default:assert(0);}
L_1EA1F: ; /* IDA 0x1ea1f; mov     ah, [si+0Ah] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+10))));
L_1EA22: ; /* IDA 0x1ea22; call    near ptr byte_1ECFE+100h */
    dl_call(m,60926,59941,h);
L_1EA25: ; /* IDA 0x1ea25; retn */
    return;
}
/* original (no IDA function); IDA 0x1EA26; raw control flow; evidence re/107 */
static void dl_body_1EA26(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xea26:goto L_1EA26; default:assert(0);}
L_1EA26: ; /* IDA 0x1ea26; mov     ah, [si+0Ah] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+10))));
L_1EA29: ; /* IDA 0x1ea29; call    near ptr unk_1EAD0 */
    dl_call(m,60112,59948,h);
L_1EA2C: ; /* IDA 0x1ea2c; call    sub_1F020 */
    dl_call(m,61472,59951,h);
L_1EA2F: ; /* IDA 0x1ea2f; retn */
    return;
}
/* original (no IDA function); IDA 0x1EA30; raw control flow; evidence re/107 */
static void dl_body_1EA30(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xea30:goto L_1EA30; default:assert(0);}
L_1EA30: ; /* IDA 0x1ea30; mov     ah, [si+0Ah] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+10))));
L_1EA33: ; /* IDA 0x1ea33; call    near ptr unk_1EAD0 */
    dl_call(m,60112,59958,h);
L_1EA36: ; /* IDA 0x1ea36; call    sub_1F1A3 */
    dl_call(m,61859,59961,h);
L_1EA39: ; /* IDA 0x1ea39; retn */
    return;
}
/* original (no IDA function); IDA 0x1EA3A; raw control flow; evidence re/107 */
static void dl_body_1EA3A(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xea3a:goto L_1EA3A; default:assert(0);}
L_1EA3A: ; /* IDA 0x1ea3a; mov     ax, [si+0Ah] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+10)));
L_1EA3D: ; /* IDA 0x1ea3d; call    near ptr unk_1EAD0 */
    dl_call(m,60112,59968,h);
L_1EA40: ; /* IDA 0x1ea40; call    loc_1F465 */
    dl_call(m,62565,59971,h);
L_1EA43: ; /* IDA 0x1ea43; retn */
    return;
}
/* original (no IDA function); IDA 0x1EA44; raw control flow; evidence re/107 */
static void dl_body_1EA44(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xea44:goto L_1EA44; default:assert(0);}
L_1EA44: ; /* IDA 0x1ea44; mov     ah, [si+8] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+8))));
L_1EA47: ; /* IDA 0x1ea47; mov     dx, [bp+0] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1EA4A: ; /* IDA 0x1ea4a; mov     bx, [bp+2] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_1EA4D: ; /* IDA 0x1ea4d; add     dx, [si+2] */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+2)),16,0,0));
L_1EA50: ; /* IDA 0x1ea50; add     bx, [si+4] */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->si+4)),16,0,0));
L_1EA53: ; /* IDA 0x1ea53; mov     cx, [si+6] */
    m->cx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+6)));
L_1EA56: ; /* IDA 0x1ea56; call    sub_1FB29 */
    dl_call(m,64297,59993,h);
L_1EA59: ; /* IDA 0x1ea59; retn */
    return;
}
/* original (no IDA function); IDA 0x1EA5A; raw control flow; evidence re/107 */
static void dl_body_1EA5A(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xea5a:goto L_1EA5A; default:assert(0);}
L_1EA5A: ; /* IDA 0x1ea5a; mov     ah, [si+8] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+8))));
L_1EA5D: ; /* IDA 0x1ea5d; mov     dx, [bp+0] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1EA60: ; /* IDA 0x1ea60; mov     bx, [bp+2] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_1EA63: ; /* IDA 0x1ea63; add     dx, [si+2] */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+2)),16,0,0));
L_1EA66: ; /* IDA 0x1ea66; add     bx, [si+4] */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->si+4)),16,0,0));
L_1EA69: ; /* IDA 0x1ea69; mov     cx, [si+6] */
    m->cx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+6)));
L_1EA6C: ; /* IDA 0x1ea6c; call    sub_1FBA7 */
    dl_call(m,64423,60015,h);
L_1EA6F: ; /* IDA 0x1ea6f; retn */
    return;
}
/* original (no IDA function); IDA 0x1EA70; raw control flow; evidence re/107 */
static void dl_body_1EA70(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xea70:goto L_1EA70; default:assert(0);}
L_1EA70: ; /* IDA 0x1ea70; push    ds */
    ec_push(m,m->ds);
L_1EA71: ; /* IDA 0x1ea71; mov     bl, [si+0Ah] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+10))));
L_1EA74: ; /* IDA 0x1ea74; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1EA76: ; /* IDA 0x1ea76; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1EA78: ; /* IDA 0x1ea78; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1EA7A: ; /* IDA 0x1ea7a; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1EA7C: ; /* IDA 0x1ea7c; add     bx, 0EAF1h */
    m->bx=(uint16_t)(ec_math(m,m->bx,60145,16,0,0));
L_1EA80: ; /* IDA 0x1ea80; mov     ax, cs:[bx] */
    m->ax=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->bx)));
L_1EA83: ; /* IDA 0x1ea83; mov     cx, cs:[bx+4] */
    m->cx=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->bx+4)));
L_1EA87: ; /* IDA 0x1ea87; push    word ptr cs:[bx+2] */
    ec_push(m,dl_read16(m,m->cs,(uint16_t)(m->bx+2)));
L_1EA8B: ; /* IDA 0x1ea8b; push    cx */
    ec_push(m,m->cx);
L_1EA8C: ; /* IDA 0x1ea8c; call    near ptr unk_1EAD0 */
    dl_call(m,60112,60047,h);
L_1EA8F: ; /* IDA 0x1ea8f; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1EA90: ; /* IDA 0x1ea90; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1EA91: ; /* IDA 0x1ea91; call    loc_1F26E */
    dl_call(m,62062,60052,h);
L_1EA94: ; /* IDA 0x1ea94; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1EA95: ; /* IDA 0x1ea95; retn */
    return;
}
/* original (no IDA function); IDA 0x1EA96; raw control flow; evidence re/107 */
static void dl_body_1EA96(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xea96:goto L_1EA96; default:assert(0);}
L_1EA96: ; /* IDA 0x1ea96; push    ds */
    ec_push(m,m->ds);
L_1EA97: ; /* IDA 0x1ea97; mov     dx, [bp+0] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1EA9A: ; /* IDA 0x1ea9a; mov     bx, [bp+2] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_1EA9D: ; /* IDA 0x1ea9d; add     dx, [si+2] */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+2)),16,0,0));
L_1EAA0: ; /* IDA 0x1eaa0; add     bx, [si+4] */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->si+4)),16,0,0));
L_1EAA3: ; /* IDA 0x1eaa3; mov     ax, [si+8] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+8)));
L_1EAA6: ; /* IDA 0x1eaa6; mov     si, [si+6] */
    m->si=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+6)));
L_1EAA9: ; /* IDA 0x1eaa9; mov     ds, cs:word_1EAED */
    m->ds=(uint16_t)(dl_read16(m,m->cs,60141));
L_1EAAE: ; /* IDA 0x1eaae; call    sub_1F6DC */
    dl_call(m,63196,60081,h);
L_1EAB1: ; /* IDA 0x1eab1; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1EAB2: ; /* IDA 0x1eab2; retn */
    return;
}
/* original (no IDA function); IDA 0x1EAB3; raw control flow; evidence re/107 */
static void dl_body_1EAB3(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xeab3:goto L_1EAB3; default:assert(0);}
L_1EAB3: ; /* IDA 0x1eab3; push    ds */
    ec_push(m,m->ds);
L_1EAB4: ; /* IDA 0x1eab4; mov     dx, [bp+0] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1EAB7: ; /* IDA 0x1eab7; mov     bx, [bp+2] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_1EABA: ; /* IDA 0x1eaba; add     dx, [si+2] */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+2)),16,0,0));
L_1EABD: ; /* IDA 0x1eabd; add     bx, [si+4] */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->si+4)),16,0,0));
L_1EAC0: ; /* IDA 0x1eac0; mov     cx, [si+8] */
    m->cx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+8)));
L_1EAC3: ; /* IDA 0x1eac3; mov     si, [si+6] */
    m->si=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+6)));
L_1EAC6: ; /* IDA 0x1eac6; mov     ds, cs:word_1EAEF */
#if KI_DISPLAY_MUTATION == 4
    m->ds=(uint16_t)(dl_read16(m,m->cs,60141));
#else
    m->ds=(uint16_t)(dl_read16(m,m->cs,60143));
#endif
L_1EACB: ; /* IDA 0x1eacb; call    sub_1F888 */
    dl_call(m,63624,60110,h);
L_1EACE: ; /* IDA 0x1eace; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1EACF: ; /* IDA 0x1eacf; retn */
    return;
}
/* original unk_1EAD0; IDA 0x1EAD0; raw control flow; evidence re/107 */
static void dl_body_1EAD0(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xead0:goto L_1EAD0; default:assert(0);}
L_1EAD0: ; /* IDA 0x1ead0; mov     dx, [bp+0] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1EAD3: ; /* IDA 0x1ead3; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1EAD5: ; /* IDA 0x1ead5; mov     bx, [bp+2] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_1EAD8: ; /* IDA 0x1ead8; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_1EADA: ; /* IDA 0x1eada; add     dx, [si+2] */
#if KI_DISPLAY_MUTATION == 2
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+4)),16,0,0));
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+2)),16,0,0));
#endif
L_1EADD: ; /* IDA 0x1eadd; add     bx, [si+4] */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->si+4)),16,0,0));
L_1EAE0: ; /* IDA 0x1eae0; add     cx, [si+6] */
    m->cx=(uint16_t)(ec_math(m,m->cx,dl_read16(m,m->ds,(uint16_t)(m->si+6)),16,0,0));
L_1EAE3: ; /* IDA 0x1eae3; add     di, [si+8] */
    m->di=(uint16_t)(ec_math(m,m->di,dl_read16(m,m->ds,(uint16_t)(m->si+8)),16,0,0));
L_1EAE6: ; /* IDA 0x1eae6; mov     si, cx */
    m->si=(uint16_t)(m->cx);
L_1EAE8: ; /* IDA 0x1eae8; retn */
    return;
}
/* original (no IDA function); IDA 0x1EDFE; raw control flow; evidence re/107 */
static void dl_body_1EDFE(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xedfe:goto L_1EDFE; case 0xefef:goto L_1EFEF; default:assert(0);}
L_1EDFE: ; /* IDA 0x1edfe; push    ds */
    ec_push(m,m->ds);
L_1EDFF: ; /* IDA 0x1edff; push    cx */
    ec_push(m,m->cx);
L_1EE00: ; /* IDA 0x1ee00; push    bp */
    ec_push(m,m->bp);
L_1EE01: ; /* IDA 0x1ee01; mov     cx, 0A0C8h */
    m->cx=(uint16_t)(41160);
L_1EE04: ; /* IDA 0x1ee04; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_1EE06: ; /* IDA 0x1ee06; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1EE08: ; /* IDA 0x1ee08; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1EE0B: ; /* IDA 0x1ee0b; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_1EE0D: ; /* IDA 0x1ee0d; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EE0E: ; /* IDA 0x1ee0e; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1EE0F: ; /* IDA 0x1ee0f; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1EE11: ; /* IDA 0x1ee11; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EE12: ; /* IDA 0x1ee12; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1EE13: ; /* IDA 0x1ee13; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1EE15: ; /* IDA 0x1ee15; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EE16: ; /* IDA 0x1ee16; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1EE17: ; /* IDA 0x1ee17; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1EE19: ; /* IDA 0x1ee19; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EE1A: ; /* IDA 0x1ee1a; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1EE1B: ; /* IDA 0x1ee1b; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1EE1D: ; /* IDA 0x1ee1d; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EE1E: ; /* IDA 0x1ee1e; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1EE1F: ; /* IDA 0x1ee1f; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1EE21: ; /* IDA 0x1ee21; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EE22: ; /* IDA 0x1ee22; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1EE23: ; /* IDA 0x1ee23; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_1EE25: ; /* IDA 0x1ee25; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EE26: ; /* IDA 0x1ee26; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1EE27: ; /* IDA 0x1ee27; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1EE29: ; /* IDA 0x1ee29; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1EE2A: ; /* IDA 0x1ee2a; mov     dx, cx */
    m->dx=(uint16_t)(m->cx);
L_1EE2C: ; /* IDA 0x1ee2c; add     dx, 8000h */
    m->dx=(uint16_t)(ec_math(m,m->dx,32768,16,0,0));
L_1EE30: ; /* IDA 0x1ee30; add     si, 8000h */
    m->si=(uint16_t)(ec_math(m,m->si,32768,16,0,0));
L_1EE34: ; /* IDA 0x1ee34; add     bx, 8000h */
    m->bx=(uint16_t)(ec_math(m,m->bx,32768,16,0,0));
L_1EE38: ; /* IDA 0x1ee38; add     di, 8000h */
    m->di=(uint16_t)(ec_math(m,m->di,32768,16,0,0));
L_1EE3C: ; /* IDA 0x1ee3c; mov     ax, di */
    m->ax=(uint16_t)(m->di);
L_1EE3E: ; /* IDA 0x1ee3e; sub     ax, bx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->bx,16,0,1));
L_1EE40: ; /* IDA 0x1ee40; jnb     short near ptr unk_1EE48 */
    if(!(m->flags&1)) goto L_1EE48;
L_1EE42: ; /* IDA 0x1ee42; xchg    dx, si */
    {uint16_t a=m->dx,b=m->si;m->dx=(uint16_t)(b);m->si=(uint16_t)(a);}
L_1EE44: ; /* IDA 0x1ee44; xchg    bx, di */
    {uint16_t a=m->bx,b=m->di;m->bx=(uint16_t)(b);m->di=(uint16_t)(a);}
L_1EE46: ; /* IDA 0x1ee46; neg     ax */
    m->ax=(uint16_t)(ec_math(m,0,m->ax,16,0,1));
L_1EE48: ; /* IDA 0x1ee48; mov     bp, si */
    m->bp=(uint16_t)(m->si);
L_1EE4A: ; /* IDA 0x1ee4a; sub     bp, dx */
    m->bp=(uint16_t)(ec_math(m,m->bp,m->dx,16,0,1));
L_1EE4C: ; /* IDA 0x1ee4c; jnb     short near ptr unk_1EE57 */
    if(!(m->flags&1)) goto L_1EE57;
L_1EE4E: ; /* IDA 0x1ee4e; neg     bp */
    m->bp=(uint16_t)(ec_math(m,0,m->bp,16,0,1));
L_1EE50: ; /* IDA 0x1ee50; cmp     bp, ax */
    ec_cmp(m,m->bp,m->ax,16);
L_1EE52: ; /* IDA 0x1ee52; jnb     short near ptr unk_1EE64 */
    if(!(m->flags&1)) goto L_1EE64;
L_1EE54: ; /* IDA 0x1ee54; jmp     near ptr unk_1EF29 */
    goto L_1EF29;
L_1EE57: ; /* IDA 0x1ee57; cmp     bp, ax */
    ec_cmp(m,m->bp,m->ax,16);
L_1EE59: ; /* IDA 0x1ee59; jnb     short near ptr unk_1EEC6 */
    if(!(m->flags&1)) goto L_1EEC6;
L_1EE5B: ; /* IDA 0x1ee5b; jmp     near ptr unk_1EF8C */
    goto L_1EF8C;
L_1EE5E: ; /* IDA 0x1ee5e; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1EE5F: ; /* IDA 0x1ee5f; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1EE60: ; /* IDA 0x1ee60; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1EE61: ; /* IDA 0x1ee61; retn */
    return;
L_1EE64: ; /* IDA 0x1ee64; mov     cs:word_1EE62, si */
    dl_write16(m,m->cs,61026,m->si);
L_1EE69: ; /* IDA 0x1ee69; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_1EE6B: ; /* IDA 0x1ee6b; cmp     dx, 8000h */
    ec_cmp(m,m->dx,32768,16);
L_1EE6F: ; /* IDA 0x1ee6f; jb      short near ptr unk_1EEC4 */
    if(m->flags&1) goto L_1EEC4;
L_1EE71: ; /* IDA 0x1ee71; cmp     bx, 8190h */
    ec_cmp(m,m->bx,33168,16);
L_1EE75: ; /* IDA 0x1ee75; jnb     short near ptr unk_1EEC4 */
    if(!(m->flags&1)) goto L_1EEC4;
L_1EE77: ; /* IDA 0x1ee77; cmp     dx, 8280h */
    ec_cmp(m,m->dx,33408,16);
L_1EE7B: ; /* IDA 0x1ee7b; jnb     short near ptr unk_1EE83 */
    if(!(m->flags&1)) goto L_1EE83;
L_1EE7D: ; /* IDA 0x1ee7d; cmp     bx, 8000h */
    ec_cmp(m,m->bx,32768,16);
L_1EE81: ; /* IDA 0x1ee81; jnb     short near ptr unk_1EE96 */
    if(!(m->flags&1)) goto L_1EE96;
L_1EE83: ; /* IDA 0x1ee83; cmp     dx, cs:word_1EE62 */
    ec_cmp(m,m->dx,dl_read16(m,m->cs,61026),16);
L_1EE88: ; /* IDA 0x1ee88; jz      short near ptr unk_1EEC4 */
    if(m->flags&0x40) goto L_1EEC4;
L_1EE8A: ; /* IDA 0x1ee8a; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1EE8B: ; /* IDA 0x1ee8b; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_1EE8D: ; /* IDA 0x1ee8d; cmp     si, bp */
    ec_cmp(m,m->si,m->bp,16);
L_1EE8F: ; /* IDA 0x1ee8f; jb      short loc_1EE77 */
    if(m->flags&1) goto L_1EE77;
L_1EE91: ; /* IDA 0x1ee91; sub     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,1));
L_1EE93: ; /* IDA 0x1ee93; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1EE94: ; /* IDA 0x1ee94; jmp     short loc_1EE77 */
    goto L_1EE77;
L_1EE96: ; /* IDA 0x1ee96; call    near ptr unk_1EFEF */
    dl_call(m,61423,61081,h);
L_1EE99: ; /* IDA 0x1ee99; cmp     dx, 8000h */
    ec_cmp(m,m->dx,32768,16);
L_1EE9D: ; /* IDA 0x1ee9d; jb      short near ptr unk_1EEC4 */
    if(m->flags&1) goto L_1EEC4;
L_1EE9F: ; /* IDA 0x1ee9f; cmp     bx, 8190h */
    ec_cmp(m,m->bx,33168,16);
L_1EEA3: ; /* IDA 0x1eea3; jnb     short near ptr unk_1EEC4 */
    if(!(m->flags&1)) goto L_1EEC4;
L_1EEA5: ; /* IDA 0x1eea5; mov     cl, [di] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
L_1EEA7: ; /* IDA 0x1eea7; mov     [di], ch */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)(m->cx>>8));
L_1EEA9: ; /* IDA 0x1eea9; cmp     dx, cs:word_1EE62 */
    ec_cmp(m,m->dx,dl_read16(m,m->cs,61026),16);
L_1EEAE: ; /* IDA 0x1eeae; jz      short near ptr unk_1EEC4 */
    if(m->flags&0x40) goto L_1EEC4;
L_1EEB0: ; /* IDA 0x1eeb0; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1EEB1: ; /* IDA 0x1eeb1; rol     ch, 1 */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),1,8,0)));
L_1EEB3: ; /* IDA 0x1eeb3; jnb     short near ptr unk_1EEB6 */
    if(!(m->flags&1)) goto L_1EEB6;
L_1EEB5: ; /* IDA 0x1eeb5; dec     di */
    m->di=(uint16_t)(st_dec(m,m->di,16));
L_1EEB6: ; /* IDA 0x1eeb6; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_1EEB8: ; /* IDA 0x1eeb8; cmp     si, bp */
    ec_cmp(m,m->si,m->bp,16);
L_1EEBA: ; /* IDA 0x1eeba; jb      short loc_1EE99 */
    if(m->flags&1) goto L_1EE99;
L_1EEBC: ; /* IDA 0x1eebc; sub     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,1));
L_1EEBE: ; /* IDA 0x1eebe; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1EEBF: ; /* IDA 0x1eebf; add     di, 50h ; 'P' */
    m->di=(uint16_t)(ec_math(m,m->di,80,16,0,0));
L_1EEC2: ; /* IDA 0x1eec2; jmp     short loc_1EE99 */
    goto L_1EE99;
L_1EEC4: ; /* IDA 0x1eec4; jmp     short loc_1EE5E */
    goto L_1EE5E;
L_1EEC6: ; /* IDA 0x1eec6; mov     cs:word_1EE62, si */
    dl_write16(m,m->cs,61026,m->si);
L_1EECB: ; /* IDA 0x1eecb; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_1EECD: ; /* IDA 0x1eecd; cmp     dx, 8280h */
    ec_cmp(m,m->dx,33408,16);
L_1EED1: ; /* IDA 0x1eed1; jnb     short near ptr unk_1EF26 */
    if(!(m->flags&1)) goto L_1EF26;
L_1EED3: ; /* IDA 0x1eed3; cmp     bx, 8190h */
    ec_cmp(m,m->bx,33168,16);
L_1EED7: ; /* IDA 0x1eed7; jnb     short near ptr unk_1EF26 */
    if(!(m->flags&1)) goto L_1EF26;
L_1EED9: ; /* IDA 0x1eed9; cmp     dx, 8000h */
    ec_cmp(m,m->dx,32768,16);
L_1EEDD: ; /* IDA 0x1eedd; jb      short near ptr unk_1EEE5 */
    if(m->flags&1) goto L_1EEE5;
L_1EEDF: ; /* IDA 0x1eedf; cmp     bx, 8000h */
    ec_cmp(m,m->bx,32768,16);
L_1EEE3: ; /* IDA 0x1eee3; jnb     short near ptr unk_1EEF8 */
    if(!(m->flags&1)) goto L_1EEF8;
L_1EEE5: ; /* IDA 0x1eee5; cmp     dx, cs:word_1EE62 */
    ec_cmp(m,m->dx,dl_read16(m,m->cs,61026),16);
L_1EEEA: ; /* IDA 0x1eeea; jz      short near ptr unk_1EF26 */
    if(m->flags&0x40) goto L_1EF26;
L_1EEEC: ; /* IDA 0x1eeec; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1EEED: ; /* IDA 0x1eeed; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_1EEEF: ; /* IDA 0x1eeef; cmp     si, bp */
    ec_cmp(m,m->si,m->bp,16);
L_1EEF1: ; /* IDA 0x1eef1; jb      short loc_1EED9 */
    if(m->flags&1) goto L_1EED9;
L_1EEF3: ; /* IDA 0x1eef3; sub     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,1));
L_1EEF5: ; /* IDA 0x1eef5; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1EEF6: ; /* IDA 0x1eef6; jmp     short loc_1EED9 */
    goto L_1EED9;
L_1EEF8: ; /* IDA 0x1eef8; call    near ptr unk_1EFEF */
    dl_call(m,61423,61179,h);
L_1EEFB: ; /* IDA 0x1eefb; cmp     dx, 8280h */
    ec_cmp(m,m->dx,33408,16);
L_1EEFF: ; /* IDA 0x1eeff; jnb     short near ptr unk_1EF26 */
    if(!(m->flags&1)) goto L_1EF26;
L_1EF01: ; /* IDA 0x1ef01; cmp     bx, 8190h */
    ec_cmp(m,m->bx,33168,16);
L_1EF05: ; /* IDA 0x1ef05; jnb     short near ptr unk_1EF26 */
    if(!(m->flags&1)) goto L_1EF26;
L_1EF07: ; /* IDA 0x1ef07; mov     cl, [di] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
L_1EF09: ; /* IDA 0x1ef09; mov     [di], ch */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)(m->cx>>8));
L_1EF0B: ; /* IDA 0x1ef0b; cmp     dx, cs:word_1EE62 */
    ec_cmp(m,m->dx,dl_read16(m,m->cs,61026),16);
L_1EF10: ; /* IDA 0x1ef10; jz      short near ptr unk_1EF26 */
    if(m->flags&0x40) goto L_1EF26;
L_1EF12: ; /* IDA 0x1ef12; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1EF13: ; /* IDA 0x1ef13; ror     ch, 1 */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),1,8,1)));
L_1EF15: ; /* IDA 0x1ef15; jnb     short near ptr unk_1EF18 */
    if(!(m->flags&1)) goto L_1EF18;
L_1EF17: ; /* IDA 0x1ef17; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1EF18: ; /* IDA 0x1ef18; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_1EF1A: ; /* IDA 0x1ef1a; cmp     si, bp */
    ec_cmp(m,m->si,m->bp,16);
L_1EF1C: ; /* IDA 0x1ef1c; jb      short loc_1EEFB */
    if(m->flags&1) goto L_1EEFB;
L_1EF1E: ; /* IDA 0x1ef1e; sub     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,1));
L_1EF20: ; /* IDA 0x1ef20; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1EF21: ; /* IDA 0x1ef21; add     di, 50h ; 'P' */
    m->di=(uint16_t)(ec_math(m,m->di,80,16,0,0));
L_1EF24: ; /* IDA 0x1ef24; jmp     short loc_1EEFB */
    goto L_1EEFB;
L_1EF26: ; /* IDA 0x1ef26; jmp     loc_1EE5E */
    goto L_1EE5E;
L_1EF29: ; /* IDA 0x1ef29; mov     cs:word_1EE62, di */
    dl_write16(m,m->cs,61026,m->di);
L_1EF2E: ; /* IDA 0x1ef2e; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_1EF30: ; /* IDA 0x1ef30; cmp     dx, 8000h */
    ec_cmp(m,m->dx,32768,16);
L_1EF34: ; /* IDA 0x1ef34; jb      short near ptr unk_1EF89 */
    if(m->flags&1) goto L_1EF89;
L_1EF36: ; /* IDA 0x1ef36; cmp     bx, 8190h */
    ec_cmp(m,m->bx,33168,16);
L_1EF3A: ; /* IDA 0x1ef3a; jnb     short near ptr unk_1EF89 */
    if(!(m->flags&1)) goto L_1EF89;
L_1EF3C: ; /* IDA 0x1ef3c; cmp     dx, 8280h */
    ec_cmp(m,m->dx,33408,16);
L_1EF40: ; /* IDA 0x1ef40; jnb     short near ptr unk_1EF48 */
    if(!(m->flags&1)) goto L_1EF48;
L_1EF42: ; /* IDA 0x1ef42; cmp     bx, 8000h */
    ec_cmp(m,m->bx,32768,16);
L_1EF46: ; /* IDA 0x1ef46; jnb     short near ptr unk_1EF5B */
    if(!(m->flags&1)) goto L_1EF5B;
L_1EF48: ; /* IDA 0x1ef48; cmp     bx, cs:word_1EE62 */
    ec_cmp(m,m->bx,dl_read16(m,m->cs,61026),16);
L_1EF4D: ; /* IDA 0x1ef4d; jz      short near ptr unk_1EF89 */
    if(m->flags&0x40) goto L_1EF89;
L_1EF4F: ; /* IDA 0x1ef4f; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1EF50: ; /* IDA 0x1ef50; add     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,0));
L_1EF52: ; /* IDA 0x1ef52; cmp     si, ax */
    ec_cmp(m,m->si,m->ax,16);
L_1EF54: ; /* IDA 0x1ef54; jb      short loc_1EF3C */
    if(m->flags&1) goto L_1EF3C;
L_1EF56: ; /* IDA 0x1ef56; sub     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,1));
L_1EF58: ; /* IDA 0x1ef58; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1EF59: ; /* IDA 0x1ef59; jmp     short loc_1EF3C */
    goto L_1EF3C;
L_1EF5B: ; /* IDA 0x1ef5b; call    near ptr unk_1EFEF */
    dl_call(m,61423,61278,h);
L_1EF5E: ; /* IDA 0x1ef5e; cmp     dx, 8000h */
    ec_cmp(m,m->dx,32768,16);
L_1EF62: ; /* IDA 0x1ef62; jb      short near ptr unk_1EF89 */
    if(m->flags&1) goto L_1EF89;
L_1EF64: ; /* IDA 0x1ef64; cmp     bx, 8190h */
    ec_cmp(m,m->bx,33168,16);
L_1EF68: ; /* IDA 0x1ef68; jnb     short near ptr unk_1EF89 */
    if(!(m->flags&1)) goto L_1EF89;
L_1EF6A: ; /* IDA 0x1ef6a; mov     cl, [di] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
L_1EF6C: ; /* IDA 0x1ef6c; mov     [di], ch */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)(m->cx>>8));
L_1EF6E: ; /* IDA 0x1ef6e; cmp     bx, cs:word_1EE62 */
    ec_cmp(m,m->bx,dl_read16(m,m->cs,61026),16);
L_1EF73: ; /* IDA 0x1ef73; jz      short near ptr unk_1EF89 */
    if(m->flags&0x40) goto L_1EF89;
L_1EF75: ; /* IDA 0x1ef75; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1EF76: ; /* IDA 0x1ef76; add     di, 50h ; 'P' */
    m->di=(uint16_t)(ec_math(m,m->di,80,16,0,0));
L_1EF79: ; /* IDA 0x1ef79; add     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,0));
L_1EF7B: ; /* IDA 0x1ef7b; cmp     si, ax */
    ec_cmp(m,m->si,m->ax,16);
L_1EF7D: ; /* IDA 0x1ef7d; jb      short loc_1EF5E */
    if(m->flags&1) goto L_1EF5E;
L_1EF7F: ; /* IDA 0x1ef7f; sub     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,1));
L_1EF81: ; /* IDA 0x1ef81; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1EF82: ; /* IDA 0x1ef82; rol     ch, 1 */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),1,8,0)));
L_1EF84: ; /* IDA 0x1ef84; jnb     short loc_1EF5E */
    if(!(m->flags&1)) goto L_1EF5E;
L_1EF86: ; /* IDA 0x1ef86; dec     di */
    m->di=(uint16_t)(st_dec(m,m->di,16));
L_1EF87: ; /* IDA 0x1ef87; jmp     short loc_1EF5E */
    goto L_1EF5E;
L_1EF89: ; /* IDA 0x1ef89; jmp     loc_1EE5E */
    goto L_1EE5E;
L_1EF8C: ; /* IDA 0x1ef8c; mov     cs:word_1EE62, di */
    dl_write16(m,m->cs,61026,m->di);
L_1EF91: ; /* IDA 0x1ef91; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_1EF93: ; /* IDA 0x1ef93; cmp     dx, 8280h */
    ec_cmp(m,m->dx,33408,16);
L_1EF97: ; /* IDA 0x1ef97; jnb     short near ptr unk_1EFEC */
    if(!(m->flags&1)) goto L_1EFEC;
L_1EF99: ; /* IDA 0x1ef99; cmp     bx, 8190h */
    ec_cmp(m,m->bx,33168,16);
L_1EF9D: ; /* IDA 0x1ef9d; jnb     short near ptr unk_1EFEC */
    if(!(m->flags&1)) goto L_1EFEC;
L_1EF9F: ; /* IDA 0x1ef9f; cmp     dx, 8000h */
    ec_cmp(m,m->dx,32768,16);
L_1EFA3: ; /* IDA 0x1efa3; jb      short near ptr unk_1EFAB */
    if(m->flags&1) goto L_1EFAB;
L_1EFA5: ; /* IDA 0x1efa5; cmp     bx, 8000h */
    ec_cmp(m,m->bx,32768,16);
L_1EFA9: ; /* IDA 0x1efa9; jnb     short near ptr unk_1EFBE */
    if(!(m->flags&1)) goto L_1EFBE;
L_1EFAB: ; /* IDA 0x1efab; cmp     bx, cs:word_1EE62 */
    ec_cmp(m,m->bx,dl_read16(m,m->cs,61026),16);
L_1EFB0: ; /* IDA 0x1efb0; jz      short near ptr unk_1EFEC */
    if(m->flags&0x40) goto L_1EFEC;
L_1EFB2: ; /* IDA 0x1efb2; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1EFB3: ; /* IDA 0x1efb3; add     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,0));
L_1EFB5: ; /* IDA 0x1efb5; cmp     si, ax */
    ec_cmp(m,m->si,m->ax,16);
L_1EFB7: ; /* IDA 0x1efb7; jb      short loc_1EF9F */
    if(m->flags&1) goto L_1EF9F;
L_1EFB9: ; /* IDA 0x1efb9; sub     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,1));
L_1EFBB: ; /* IDA 0x1efbb; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1EFBC: ; /* IDA 0x1efbc; jmp     short loc_1EF9F */
    goto L_1EF9F;
L_1EFBE: ; /* IDA 0x1efbe; call    near ptr unk_1EFEF */
    dl_call(m,61423,61377,h);
L_1EFC1: ; /* IDA 0x1efc1; cmp     dx, 8280h */
    ec_cmp(m,m->dx,33408,16);
L_1EFC5: ; /* IDA 0x1efc5; jnb     short near ptr unk_1EFEC */
    if(!(m->flags&1)) goto L_1EFEC;
L_1EFC7: ; /* IDA 0x1efc7; cmp     bx, 8190h */
    ec_cmp(m,m->bx,33168,16);
L_1EFCB: ; /* IDA 0x1efcb; jnb     short near ptr unk_1EFEC */
    if(!(m->flags&1)) goto L_1EFEC;
L_1EFCD: ; /* IDA 0x1efcd; mov     cl, [di] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
L_1EFCF: ; /* IDA 0x1efcf; mov     [di], ch */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)(m->cx>>8));
L_1EFD1: ; /* IDA 0x1efd1; cmp     bx, cs:word_1EE62 */
    ec_cmp(m,m->bx,dl_read16(m,m->cs,61026),16);
L_1EFD6: ; /* IDA 0x1efd6; jz      short near ptr unk_1EFEC */
    if(m->flags&0x40) goto L_1EFEC;
L_1EFD8: ; /* IDA 0x1efd8; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1EFD9: ; /* IDA 0x1efd9; add     di, 50h ; 'P' */
    m->di=(uint16_t)(ec_math(m,m->di,80,16,0,0));
L_1EFDC: ; /* IDA 0x1efdc; add     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,0));
L_1EFDE: ; /* IDA 0x1efde; cmp     si, ax */
    ec_cmp(m,m->si,m->ax,16);
L_1EFE0: ; /* IDA 0x1efe0; jb      short loc_1EFC1 */
    if(m->flags&1) goto L_1EFC1;
L_1EFE2: ; /* IDA 0x1efe2; sub     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,1));
L_1EFE4: ; /* IDA 0x1efe4; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1EFE5: ; /* IDA 0x1efe5; ror     ch, 1 */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),1,8,1)));
L_1EFE7: ; /* IDA 0x1efe7; jnb     short loc_1EFC1 */
    if(!(m->flags&1)) goto L_1EFC1;
L_1EFE9: ; /* IDA 0x1efe9; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1EFEA: ; /* IDA 0x1efea; jmp     short loc_1EFC1 */
    goto L_1EFC1;
L_1EFEC: ; /* IDA 0x1efec; jmp     loc_1EE5E */
    goto L_1EE5E;
L_1EFEF: ; /* IDA 0x1efef; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_1EFF1: ; /* IDA 0x1eff1; sub     cx, 8000h */
    m->cx=(uint16_t)(ec_math(m,m->cx,32768,16,0,1));
L_1EFF5: ; /* IDA 0x1eff5; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_1EFF7: ; /* IDA 0x1eff7; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_1EFF9: ; /* IDA 0x1eff9; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_1EFFB: ; /* IDA 0x1effb; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_1EFFD: ; /* IDA 0x1effd; mov     di, cx */
    m->di=(uint16_t)(m->cx);
L_1EFFF: ; /* IDA 0x1efff; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_1F001: ; /* IDA 0x1f001; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_1F003: ; /* IDA 0x1f003; add     di, cx */
    m->di=(uint16_t)(ec_math(m,m->di,m->cx,16,0,0));
L_1F005: ; /* IDA 0x1f005; sub     dx, 8000h */
    m->dx=(uint16_t)(ec_math(m,m->dx,32768,16,0,1));
L_1F009: ; /* IDA 0x1f009; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1F00B: ; /* IDA 0x1f00b; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_1F00D: ; /* IDA 0x1f00d; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_1F00F: ; /* IDA 0x1f00f; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_1F011: ; /* IDA 0x1f011; add     di, cx */
    m->di=(uint16_t)(ec_math(m,m->di,m->cx,16,0,0));
L_1F013: ; /* IDA 0x1f013; mov     cl, 7 */
    vg_cl(m,(uint8_t)(7));
L_1F015: ; /* IDA 0x1f015; and     cl, dl */
    {uint16_t result=(((uint8_t)m->cx)&((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_1F017: ; /* IDA 0x1f017; mov     ch, 80h */
    vg_ch(m,(uint8_t)(128));
L_1F019: ; /* IDA 0x1f019; shr     ch, cl */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),(uint8_t)m->cx,8,5)));
L_1F01B: ; /* IDA 0x1f01b; add     dx, 8000h */
    m->dx=(uint16_t)(ec_math(m,m->dx,32768,16,0,0));
L_1F01F: ; /* IDA 0x1f01f; retn */
    return;
}
/* original loc_1F26E; IDA 0x1F26E; raw control flow; evidence re/107 */
static void dl_body_1F26E(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xf26e:goto L_1F26E; default:assert(0);}
L_1F26E: ; /* IDA 0x1f26e; cld */
    m->flags&=(uint16_t)~0x400u;
L_1F26F: ; /* IDA 0x1f26f; cmp     dx, si */
    ec_cmp(m,m->dx,m->si,16);
L_1F271: ; /* IDA 0x1f271; jle     short loc_1F275 */
    if(!ec_greater(m)) goto L_1F275;
L_1F273: ; /* IDA 0x1f273; xchg    dx, si */
    {uint16_t a=m->dx,b=m->si;m->dx=(uint16_t)(b);m->si=(uint16_t)(a);}
L_1F275: ; /* IDA 0x1f275; cmp     bx, di */
    ec_cmp(m,m->bx,m->di,16);
L_1F277: ; /* IDA 0x1f277; jle     short loc_1F27B */
    if(!ec_greater(m)) goto L_1F27B;
L_1F279: ; /* IDA 0x1f279; xchg    bx, di */
    {uint16_t a=m->bx,b=m->di;m->bx=(uint16_t)(b);m->di=(uint16_t)(a);}
L_1F27B: ; /* IDA 0x1f27b; cmp     dx, 280h */
    ec_cmp(m,m->dx,640,16);
L_1F27F: ; /* IDA 0x1f27f; jge     short locret_1F2B3 */
    if(!ec_less(m)) goto L_1F2B3;
L_1F281: ; /* IDA 0x1f281; cmp     dx, 0 */
    ec_cmp(m,m->dx,0,16);
L_1F284: ; /* IDA 0x1f284; jge     short loc_1F288 */
    if(!ec_less(m)) goto L_1F288;
L_1F286: ; /* IDA 0x1f286; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_1F288: ; /* IDA 0x1f288; cmp     si, 0 */
    ec_cmp(m,m->si,0,16);
L_1F28B: ; /* IDA 0x1f28b; jl      short locret_1F2B3 */
    if(ec_less(m)) goto L_1F2B3;
L_1F28D: ; /* IDA 0x1f28d; cmp     si, 280h */
    ec_cmp(m,m->si,640,16);
L_1F291: ; /* IDA 0x1f291; jl      short loc_1F296 */
    if(ec_less(m)) goto L_1F296;
L_1F293: ; /* IDA 0x1f293; mov     si, 280h */
    m->si=(uint16_t)(640);
L_1F296: ; /* IDA 0x1f296; cmp     bx, 190h */
    ec_cmp(m,m->bx,400,16);
L_1F29A: ; /* IDA 0x1f29a; jge     short locret_1F2B3 */
    if(!ec_less(m)) goto L_1F2B3;
L_1F29C: ; /* IDA 0x1f29c; cmp     bx, 0 */
    ec_cmp(m,m->bx,0,16);
L_1F29F: ; /* IDA 0x1f29f; jge     short loc_1F2A3 */
    if(!ec_less(m)) goto L_1F2A3;
L_1F2A1: ; /* IDA 0x1f2a1; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_1F2A3: ; /* IDA 0x1f2a3; cmp     di, 0 */
    ec_cmp(m,m->di,0,16);
L_1F2A6: ; /* IDA 0x1f2a6; jl      short locret_1F2B3 */
    if(ec_less(m)) goto L_1F2B3;
L_1F2A8: ; /* IDA 0x1f2a8; cmp     di, 190h */
    ec_cmp(m,m->di,400,16);
L_1F2AC: ; /* IDA 0x1f2ac; jl      short loc_1F2B4 */
    if(ec_less(m)) goto L_1F2B4;
L_1F2AE: ; /* IDA 0x1f2ae; mov     di, 190h */
    m->di=(uint16_t)(400);
L_1F2B1: ; /* IDA 0x1f2b1; jmp     short loc_1F2B4 */
    goto L_1F2B4;
L_1F2B3: ; /* IDA 0x1f2b3; retn */
    return;
L_1F2B4: ; /* IDA 0x1f2b4; push    es */
    ec_push(m,m->es);
L_1F2B5: ; /* IDA 0x1f2b5; push    bp */
    ec_push(m,m->bp);
L_1F2B6: ; /* IDA 0x1f2b6; shr     cl, 1 */
    vg_cl(m,(uint8_t)(dl_shift(m,(uint8_t)m->cx,1,8,5)));
L_1F2B8: ; /* IDA 0x1f2b8; shr     cl, 1 */
    vg_cl(m,(uint8_t)(dl_shift(m,(uint8_t)m->cx,1,8,5)));
L_1F2BA: ; /* IDA 0x1f2ba; shr     cl, 1 */
    vg_cl(m,(uint8_t)(dl_shift(m,(uint8_t)m->cx,1,8,5)));
L_1F2BC: ; /* IDA 0x1f2bc; mov     byte ptr cs:word_1F45B, cl */
    dl_write8(m,m->cs,62555,(uint8_t)m->cx);
L_1F2C1: ; /* IDA 0x1f2c1; mov     cs:byte_1F45D, ch */
    dl_write8(m,m->cs,62557,(uint8_t)(m->cx>>8));
L_1F2C6: ; /* IDA 0x1f2c6; mov     byte ptr cs:word_1F45B+1, cl */
    dl_write8(m,m->cs,62556,(uint8_t)m->cx);
L_1F2CB: ; /* IDA 0x1f2cb; mov     cs:byte_1F45E, ch */
    dl_write8(m,m->cs,62558,(uint8_t)(m->cx>>8));
L_1F2D0: ; /* IDA 0x1f2d0; mov     cs:word_1F45F, ax */
    dl_write16(m,m->cs,62559,m->ax);
L_1F2D4: ; /* IDA 0x1f2d4; mov     al, cl */
    vg_al(m,(uint8_t)((uint8_t)m->cx));
L_1F2D6: ; /* IDA 0x1f2d6; mul     ch */
    dl_mul(m,(uint8_t)(m->cx>>8),8);
L_1F2D8: ; /* IDA 0x1f2d8; mov     bp, ax */
    m->bp=(uint16_t)(m->ax);
L_1F2DA: ; /* IDA 0x1f2da; mov     ax, dx */
    m->ax=(uint16_t)(m->dx);
L_1F2DC: ; /* IDA 0x1f2dc; shl     cl, 1 */
    vg_cl(m,(uint8_t)(dl_shift(m,(uint8_t)m->cx,1,8,4)));
L_1F2DE: ; /* IDA 0x1f2de; shl     cl, 1 */
    vg_cl(m,(uint8_t)(dl_shift(m,(uint8_t)m->cx,1,8,4)));
L_1F2E0: ; /* IDA 0x1f2e0; shl     cl, 1 */
    vg_cl(m,(uint8_t)(dl_shift(m,(uint8_t)m->cx,1,8,4)));
L_1F2E2: ; /* IDA 0x1f2e2; div     cl */
    dl_div(m,(uint8_t)m->cx,8);
L_1F2E4: ; /* IDA 0x1f2e4; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1F2E6: ; /* IDA 0x1f2e6; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_1F2E8: ; /* IDA 0x1f2e8; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_1F2EA: ; /* IDA 0x1f2ea; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_1F2EC: ; /* IDA 0x1f2ec; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1F2EE: ; /* IDA 0x1f2ee; sub     byte ptr cs:word_1F45B+1, al */
    dl_write8(m,m->cs,62556,ec_math(m,dl_read8(m,m->cs,62556),(uint8_t)m->ax,8,0,1));
L_1F2F3: ; /* IDA 0x1f2f3; add     cs:word_1F45F, ax */
    dl_write16(m,m->cs,62559,ec_math(m,dl_read16(m,m->cs,62559),m->ax,16,0,0));
L_1F2F8: ; /* IDA 0x1f2f8; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_1F2FA: ; /* IDA 0x1f2fa; div     ch */
    dl_div(m,(uint8_t)(m->cx>>8),8);
L_1F2FC: ; /* IDA 0x1f2fc; sub     cs:byte_1F45E, ah */
    dl_write8(m,m->cs,62558,ec_math(m,dl_read8(m,m->cs,62558),(uint8_t)(m->ax>>8),8,0,1));
L_1F301: ; /* IDA 0x1f301; mov     al, byte ptr cs:word_1F45B */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,62555)));
L_1F305: ; /* IDA 0x1f305; mul     ah */
    dl_mul(m,(uint8_t)(m->ax>>8),8);
L_1F307: ; /* IDA 0x1f307; add     cs:word_1F45F, ax */
    dl_write16(m,m->cs,62559,ec_math(m,dl_read16(m,m->cs,62559),m->ax,16,0,0));
L_1F30C: ; /* IDA 0x1f30c; mov     ax, di */
    m->ax=(uint16_t)(m->di);
L_1F30E: ; /* IDA 0x1f30e; sub     ax, bx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->bx,16,0,1));
L_1F310: ; /* IDA 0x1f310; inc     ax */
    m->ax=(uint16_t)(ec_inc(m,m->ax,16));
L_1F311: ; /* IDA 0x1f311; mov     cs:word_1F463, ax */
    dl_write16(m,m->cs,62563,m->ax);
L_1F315: ; /* IDA 0x1f315; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_1F317: ; /* IDA 0x1f317; and     al, 0F8h */
    {uint16_t result=(((uint8_t)m->ax)&(248))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F319: ; /* IDA 0x1f319; sub     ax, dx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->dx,16,0,1));
L_1F31B: ; /* IDA 0x1f31b; jnb     short loc_1F321 */
    if(!(m->flags&1)) goto L_1F321;
L_1F31D: ; /* IDA 0x1f31d; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_1F31F: ; /* IDA 0x1f31f; jmp     short loc_1F32D */
    goto L_1F32D;
L_1F321: ; /* IDA 0x1f321; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1F323: ; /* IDA 0x1f323; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1F325: ; /* IDA 0x1f325; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1F327: ; /* IDA 0x1f327; test    dl, 7 */
    {uint16_t result=(((uint8_t)m->dx)&(7))&255;ec_logic(m,result,8);}
L_1F32A: ; /* IDA 0x1f32a; jz      short loc_1F32D */
    if(m->flags&0x40) goto L_1F32D;
L_1F32C: ; /* IDA 0x1f32c; inc     ax */
    m->ax=(uint16_t)(ec_inc(m,m->ax,16));
L_1F32D: ; /* IDA 0x1f32d; mov     cs:word_1F461, ax */
    dl_write16(m,m->cs,62561,m->ax);
L_1F331: ; /* IDA 0x1f331; mov     cx, si */
    m->cx=(uint16_t)(m->si);
L_1F333: ; /* IDA 0x1f333; mov     ch, dl */
    vg_ch(m,(uint8_t)((uint8_t)m->dx));
L_1F335: ; /* IDA 0x1f335; and     cx, 707h */
    {uint16_t result=((m->cx)&(1799))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_1F339: ; /* IDA 0x1f339; mov     ax, 0FFFFh */
    m->ax=(uint16_t)(65535);
L_1F33C: ; /* IDA 0x1f33c; shr     ah, cl */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),(uint8_t)m->cx,8,5)));
L_1F33E: ; /* IDA 0x1f33e; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_1F340: ; /* IDA 0x1f340; not     ah */
    vg_ah(m,(uint8_t)(~((uint8_t)(m->ax>>8))));
L_1F342: ; /* IDA 0x1f342; mov     cl, ch */
    vg_cl(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_1F344: ; /* IDA 0x1f344; shr     al, cl */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,(uint8_t)m->cx,8,5)));
L_1F346: ; /* IDA 0x1f346; cmp     byte ptr cs:word_1F461, 0 */
    ec_cmp(m,dl_read8(m,m->cs,62561),0,8);
L_1F34C: ; /* IDA 0x1f34c; jnz     short loc_1F352 */
    if(!(m->flags&0x40)) goto L_1F352;
L_1F34E: ; /* IDA 0x1f34e; and     al, ah */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F350: ; /* IDA 0x1f350; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1F352: ; /* IDA 0x1f352; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1F354: ; /* IDA 0x1f354; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1F356: ; /* IDA 0x1f356; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1F358: ; /* IDA 0x1f358; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1F35A: ; /* IDA 0x1f35a; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1F35C: ; /* IDA 0x1f35c; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1F35E: ; /* IDA 0x1f35e; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1F360: ; /* IDA 0x1f360; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_1F362: ; /* IDA 0x1f362; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1F364: ; /* IDA 0x1f364; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1F366: ; /* IDA 0x1f366; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_1F368: ; /* IDA 0x1f368; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_1F36A: ; /* IDA 0x1f36a; mov     si, cs:word_1F45F */
    m->si=(uint16_t)(dl_read16(m,m->cs,62559));
L_1F36F: ; /* IDA 0x1f36f; mov     dx, ax */
    m->dx=(uint16_t)(m->ax);
L_1F371: ; /* IDA 0x1f371; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_1F373: ; /* IDA 0x1f373; xchg    bl, byte ptr cs:word_1F45B+1 */
    {uint16_t a=(uint8_t)m->bx,b=dl_read8(m,m->cs,62556);gw_bl(m,(uint8_t)(b));dl_write8(m,m->cs,62556,a);}
L_1F378: ; /* IDA 0x1f378; xchg    bh, cs:byte_1F45E */
    {uint16_t a=(uint8_t)(m->bx>>8),b=dl_read8(m,m->cs,62558);gw_bh(m,(uint8_t)(b));dl_write8(m,m->cs,62558,a);}
L_1F37D: ; /* IDA 0x1f37d; mov     ch, byte ptr cs:word_1F461 */
    vg_ch(m,(uint8_t)(dl_read8(m,m->cs,62561)));
L_1F382: ; /* IDA 0x1f382; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_1F385: ; /* IDA 0x1f385; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1F387: ; /* IDA 0x1f387; call    sub_1F426 */
    dl_call(m,62502,62346,h);
L_1F38A: ; /* IDA 0x1f38a; push    si */
    ec_push(m,m->si);
L_1F38B: ; /* IDA 0x1f38b; mov     ah, 0Eh */
    vg_ah(m,(uint8_t)(14));
L_1F38D: ; /* IDA 0x1f38d; call    sub_1F43B */
    dl_call(m,62523,62352,h);
L_1F390: ; /* IDA 0x1f390; call    sub_1F3CF */
    dl_call(m,62415,62355,h);
L_1F393: ; /* IDA 0x1f393; add     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,0));
L_1F395: ; /* IDA 0x1f395; mov     ah, 0Dh */
    vg_ah(m,(uint8_t)(13));
L_1F397: ; /* IDA 0x1f397; call    sub_1F43B */
    dl_call(m,62523,62362,h);
L_1F39A: ; /* IDA 0x1f39a; call    sub_1F3CF */
    dl_call(m,62415,62365,h);
L_1F39D: ; /* IDA 0x1f39d; add     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,0));
L_1F39F: ; /* IDA 0x1f39f; mov     ah, 0Bh */
    vg_ah(m,(uint8_t)(11));
L_1F3A1: ; /* IDA 0x1f3a1; call    sub_1F43B */
    dl_call(m,62523,62372,h);
L_1F3A4: ; /* IDA 0x1f3a4; call    sub_1F3CF */
    dl_call(m,62415,62375,h);
L_1F3A7: ; /* IDA 0x1f3a7; add     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,0));
L_1F3A9: ; /* IDA 0x1f3a9; mov     ah, 7 */
    vg_ah(m,(uint8_t)(7));
L_1F3AB: ; /* IDA 0x1f3ab; call    sub_1F43B */
    dl_call(m,62523,62382,h);
L_1F3AE: ; /* IDA 0x1f3ae; call    sub_1F3CF */
    dl_call(m,62415,62385,h);
L_1F3B1: ; /* IDA 0x1f3b1; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1F3B2: ; /* IDA 0x1f3b2; add     si, cs:word_1F45B */
#if KI_DISPLAY_MUTATION == 7
    m->si=(uint16_t)(ec_math(m,m->si,(dl_read16(m,m->cs,62555)+1),16,0,0));
#else
    m->si=(uint16_t)(ec_math(m,m->si,dl_read16(m,m->cs,62555),16,0,0));
#endif
L_1F3B7: ; /* IDA 0x1f3b7; dec     bh */
    gw_bh(m,(uint8_t)(st_dec(m,(uint8_t)(m->bx>>8),8)));
L_1F3B9: ; /* IDA 0x1f3b9; jnz     short loc_1F3C2 */
    if(!(m->flags&0x40)) goto L_1F3C2;
L_1F3BB: ; /* IDA 0x1f3bb; mov     bh, cs:byte_1F45D */
    gw_bh(m,(uint8_t)(dl_read8(m,m->cs,62557)));
L_1F3C0: ; /* IDA 0x1f3c0; sub     si, bp */
    m->si=(uint16_t)(ec_math(m,m->si,m->bp,16,0,1));
L_1F3C2: ; /* IDA 0x1f3c2; add     di, 50h ; 'P' */
    m->di=(uint16_t)(ec_math(m,m->di,80,16,0,0));
L_1F3C5: ; /* IDA 0x1f3c5; dec     cs:word_1F463 */
    dl_write16(m,m->cs,62563,st_dec(m,dl_read16(m,m->cs,62563),16));
L_1F3CA: ; /* IDA 0x1f3ca; jnz     short loc_1F38A */
    if(!(m->flags&0x40)) goto L_1F38A;
L_1F3CC: ; /* IDA 0x1f3cc; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1F3CD: ; /* IDA 0x1f3cd; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1F3CE: ; /* IDA 0x1f3ce; retn */
    return;
}
/* original (no IDA function); IDA 0x1F3CF; raw control flow; evidence re/107 */
static void dl_body_1F3CF(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xf3cf:goto L_1F3CF; default:assert(0);}
L_1F3CF: ; /* IDA 0x1f3cf; push    bx */
    ec_push(m,m->bx);
L_1F3D0: ; /* IDA 0x1f3d0; push    si */
    ec_push(m,m->si);
L_1F3D1: ; /* IDA 0x1f3d1; push    di */
    ec_push(m,m->di);
L_1F3D2: ; /* IDA 0x1f3d2; push    bp */
    ec_push(m,m->bp);
L_1F3D3: ; /* IDA 0x1f3d3; mov     bp, dx */
    m->bp=(uint16_t)(m->dx);
L_1F3D5: ; /* IDA 0x1f3d5; mov     cl, ch */
    vg_cl(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_1F3D7: ; /* IDA 0x1f3d7; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_1F3DA: ; /* IDA 0x1f3da; mov     ax, bp */
    m->ax=(uint16_t)(m->bp);
L_1F3DC: ; /* IDA 0x1f3dc; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F3DD: ; /* IDA 0x1f3dd; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_1F3E0: ; /* IDA 0x1f3e0; movsb */
    vg_movsb(m);
L_1F3E1: ; /* IDA 0x1f3e1; and     cl, cl */
    {uint16_t result=(((uint8_t)m->cx)&((uint8_t)m->cx))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_1F3E3: ; /* IDA 0x1f3e3; jz      short loc_1F41F */
    if(m->flags&0x40) goto L_1F41F;
L_1F3E5: ; /* IDA 0x1f3e5; dec     bl */
    gw_bl(m,(uint8_t)(st_dec(m,(uint8_t)m->bx,8)));
L_1F3E7: ; /* IDA 0x1f3e7; jnz     short loc_1F3F3 */
    if(!(m->flags&0x40)) goto L_1F3F3;
L_1F3E9: ; /* IDA 0x1f3e9; sub     si, cs:word_1F45B */
    m->si=(uint16_t)(ec_math(m,m->si,dl_read16(m,m->cs,62555),16,0,1));
L_1F3EE: ; /* IDA 0x1f3ee; mov     bl, byte ptr cs:word_1F45B */
    gw_bl(m,(uint8_t)(dl_read8(m,m->cs,62555)));
L_1F3F3: ; /* IDA 0x1f3f3; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_1F3F5: ; /* IDA 0x1f3f5; jz      short loc_1F413 */
    if(m->flags&0x40) goto L_1F413;
L_1F3F7: ; /* IDA 0x1f3f7; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_1F3FA: ; /* IDA 0x1f3fa; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1F3FC: ; /* IDA 0x1f3fc; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F3FD: ; /* IDA 0x1f3fd; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_1F400: ; /* IDA 0x1f400; movsb */
    vg_movsb(m);
L_1F401: ; /* IDA 0x1f401; dec     bl */
    gw_bl(m,(uint8_t)(st_dec(m,(uint8_t)m->bx,8)));
L_1F403: ; /* IDA 0x1f403; jnz     short loc_1F40F */
    if(!(m->flags&0x40)) goto L_1F40F;
L_1F405: ; /* IDA 0x1f405; sub     si, cs:word_1F45B */
    m->si=(uint16_t)(ec_math(m,m->si,dl_read16(m,m->cs,62555),16,0,1));
L_1F40A: ; /* IDA 0x1f40a; mov     bl, byte ptr cs:word_1F45B */
    gw_bl(m,(uint8_t)(dl_read8(m,m->cs,62555)));
L_1F40F: ; /* IDA 0x1f40f; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_1F411: ; /* IDA 0x1f411; jnz     short loc_1F3FD */
    if(!(m->flags&0x40)) goto L_1F3FD;
L_1F413: ; /* IDA 0x1f413; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_1F416: ; /* IDA 0x1f416; mov     ax, bp */
    m->ax=(uint16_t)(m->bp);
L_1F418: ; /* IDA 0x1f418; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1F41A: ; /* IDA 0x1f41a; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F41B: ; /* IDA 0x1f41b; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_1F41E: ; /* IDA 0x1f41e; movsb */
    vg_movsb(m);
L_1F41F: ; /* IDA 0x1f41f; mov     dx, bp */
    m->dx=(uint16_t)(m->bp);
L_1F421: ; /* IDA 0x1f421; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1F422: ; /* IDA 0x1f422; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1F423: ; /* IDA 0x1f423; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1F424: ; /* IDA 0x1f424; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1F425: ; /* IDA 0x1f425; retn */
    return;
}
/* original (no IDA function); IDA 0x1F426; raw control flow; evidence re/107 */
static void dl_body_1F426(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xf426:goto L_1F426; default:assert(0);}
L_1F426: ; /* IDA 0x1f426; push    dx */
    ec_push(m,m->dx);
L_1F427: ; /* IDA 0x1f427; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1F42A: ; /* IDA 0x1f42a; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_1F42C: ; /* IDA 0x1f42c; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F42D: ; /* IDA 0x1f42d; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F42E: ; /* IDA 0x1f42e; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F430: ; /* IDA 0x1f430; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F431: ; /* IDA 0x1f431; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1F432: ; /* IDA 0x1f432; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F434: ; /* IDA 0x1f434; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F435: ; /* IDA 0x1f435; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F436: ; /* IDA 0x1f436; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F438: ; /* IDA 0x1f438; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F439: ; /* IDA 0x1f439; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1F43A: ; /* IDA 0x1f43a; retn */
    return;
}
/* original (no IDA function); IDA 0x1F43B; raw control flow; evidence re/107 */
static void dl_body_1F43B(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xf43b:goto L_1F43B; default:assert(0);}
L_1F43B: ; /* IDA 0x1f43b; push    dx */
    ec_push(m,m->dx);
L_1F43C: ; /* IDA 0x1f43c; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1F43F: ; /* IDA 0x1f43f; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_1F441: ; /* IDA 0x1f441; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F442: ; /* IDA 0x1f442; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F443: ; /* IDA 0x1f443; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1F445: ; /* IDA 0x1f445; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F446: ; /* IDA 0x1f446; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1F447: ; /* IDA 0x1f447; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1F449: ; /* IDA 0x1f449; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F44A: ; /* IDA 0x1f44a; mov     al, 10h */
    vg_al(m,(uint8_t)(16));
L_1F44C: ; /* IDA 0x1f44c; cmp     ah, 0Eh */
    ec_cmp(m,(uint8_t)(m->ax>>8),14,8);
L_1F44F: ; /* IDA 0x1f44f; jnz     short loc_1F453 */
    if(!(m->flags&0x40)) goto L_1F453;
L_1F451: ; /* IDA 0x1f451; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F453: ; /* IDA 0x1f453; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F454: ; /* IDA 0x1f454; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F455: ; /* IDA 0x1f455; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1F456: ; /* IDA 0x1f456; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_1F458: ; /* IDA 0x1f458; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F459: ; /* IDA 0x1f459; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1F45A: ; /* IDA 0x1f45a; retn */
    return;
}
/* original loc_1F465; IDA 0x1F465; raw control flow; evidence re/107 */
static void dl_body_1F465(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {
switch(entry) {case 0xf465:goto L_1F465; default:assert(0);}
L_1F465: ; /* IDA 0x1f465; push    ax */
    ec_push(m,m->ax);
L_1F466: ; /* IDA 0x1f466; push    bx */
    ec_push(m,m->bx);
L_1F467: ; /* IDA 0x1f467; push    cx */
    ec_push(m,m->cx);
L_1F468: ; /* IDA 0x1f468; push    dx */
    ec_push(m,m->dx);
L_1F469: ; /* IDA 0x1f469; push    si */
    ec_push(m,m->si);
L_1F46A: ; /* IDA 0x1f46a; push    di */
    ec_push(m,m->di);
L_1F46B: ; /* IDA 0x1f46b; cmp     dx, si */
    ec_cmp(m,m->dx,m->si,16);
L_1F46D: ; /* IDA 0x1f46d; jbe     short loc_1F471 */
    if(m->flags&(1|0x40)) goto L_1F471;
L_1F46F: ; /* IDA 0x1f46f; xchg    dx, si */
    {uint16_t a=m->dx,b=m->si;m->dx=(uint16_t)(b);m->si=(uint16_t)(a);}
L_1F471: ; /* IDA 0x1f471; cmp     bx, di */
    ec_cmp(m,m->bx,m->di,16);
L_1F473: ; /* IDA 0x1f473; jbe     short loc_1F477 */
    if(m->flags&(1|0x40)) goto L_1F477;
L_1F475: ; /* IDA 0x1f475; xchg    bx, di */
    {uint16_t a=m->bx,b=m->di;m->bx=(uint16_t)(b);m->di=(uint16_t)(a);}
L_1F477: ; /* IDA 0x1f477; mov     cx, si */
#if KI_DISPLAY_MUTATION == 9
    m->cx=(uint16_t)(m->dx);
#else
    m->cx=(uint16_t)(m->si);
#endif
L_1F479: ; /* IDA 0x1f479; sub     cx, dx */
    m->cx=(uint16_t)(ec_math(m,m->cx,m->dx,16,0,1));
L_1F47B: ; /* IDA 0x1f47b; inc     cx */
    m->cx=(uint16_t)(ec_inc(m,m->cx,16));
L_1F47C: ; /* IDA 0x1f47c; xchg    ah, al */
    {uint16_t a=(uint8_t)(m->ax>>8),b=(uint8_t)m->ax;vg_ah(m,(uint8_t)(b));vg_al(m,(uint8_t)(a));}
L_1F47E: ; /* IDA 0x1f47e; call    sub_1FB29 */
    dl_call(m,64297,62593,h);
L_1F481: ; /* IDA 0x1f481; xchg    bx, di */
    {uint16_t a=m->bx,b=m->di;m->bx=(uint16_t)(b);m->di=(uint16_t)(a);}
L_1F483: ; /* IDA 0x1f483; xchg    ah, al */
    {uint16_t a=(uint8_t)(m->ax>>8),b=(uint8_t)m->ax;vg_ah(m,(uint8_t)(b));vg_al(m,(uint8_t)(a));}
L_1F485: ; /* IDA 0x1f485; call    sub_1FB29 */
    dl_call(m,64297,62600,h);
L_1F488: ; /* IDA 0x1f488; xchg    bx, di */
    {uint16_t a=m->bx,b=m->di;m->bx=(uint16_t)(b);m->di=(uint16_t)(a);}
L_1F48A: ; /* IDA 0x1f48a; mov     cx, di */
    m->cx=(uint16_t)(m->di);
L_1F48C: ; /* IDA 0x1f48c; sub     cx, bx */
    m->cx=(uint16_t)(ec_math(m,m->cx,m->bx,16,0,1));
L_1F48E: ; /* IDA 0x1f48e; inc     cx */
    m->cx=(uint16_t)(ec_inc(m,m->cx,16));
L_1F48F: ; /* IDA 0x1f48f; xchg    ah, al */
    {uint16_t a=(uint8_t)(m->ax>>8),b=(uint8_t)m->ax;vg_ah(m,(uint8_t)(b));vg_al(m,(uint8_t)(a));}
L_1F491: ; /* IDA 0x1f491; call    sub_1FBA7 */
    dl_call(m,64423,62612,h);
L_1F494: ; /* IDA 0x1f494; xchg    dx, si */
    {uint16_t a=m->dx,b=m->si;m->dx=(uint16_t)(b);m->si=(uint16_t)(a);}
L_1F496: ; /* IDA 0x1f496; xchg    ah, al */
    {uint16_t a=(uint8_t)(m->ax>>8),b=(uint8_t)m->ax;vg_ah(m,(uint8_t)(b));vg_al(m,(uint8_t)(a));}
L_1F498: ; /* IDA 0x1f498; call    sub_1FBA7 */
    dl_call(m,64423,62619,h);
L_1F49B: ; /* IDA 0x1f49b; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1F49C: ; /* IDA 0x1f49c; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1F49D: ; /* IDA 0x1f49d; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1F49E: ; /* IDA 0x1f49e; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1F49F: ; /* IDA 0x1f49f; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1F4A0: ; /* IDA 0x1f4a0; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1F4A1: ; /* IDA 0x1f4a1; retn */
    return;
}
int ki_display_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x30f:dl_body_1030F(m,h,t);break;
case 0x337:dl_body_10337(m,h,t);break;
case 0xe993:dl_body_1E993(m,h,t);break;
case 0xe9a7:dl_body_1E9A7(m,h,t);break;
case 0xe9c1:dl_body_1E9C1(m,h,t);break;
case 0xfb29:dl_body_1FB29(m,h,t);break;
case 0xfba7:dl_body_1FBA7(m,h,t);break;
case 0xf6dc:dl_body_1F6DC(m,h,t);break;
case 0xf878:dl_body_1F878(m,h,t);break;
case 0xea1f:dl_body_1EA1F(m,h,t);break;
case 0xea26:dl_body_1EA26(m,h,t);break;
case 0xea30:dl_body_1EA30(m,h,t);break;
case 0xea3a:dl_body_1EA3A(m,h,t);break;
case 0xea44:dl_body_1EA44(m,h,t);break;
case 0xea5a:dl_body_1EA5A(m,h,t);break;
case 0xea70:dl_body_1EA70(m,h,t);break;
case 0xea96:dl_body_1EA96(m,h,t);break;
case 0xeab3:dl_body_1EAB3(m,h,t);break;
case 0xead0:dl_body_1EAD0(m,h,t);break;
case 0xedfe:dl_body_1EDFE(m,h,t);break;
case 0xefef:dl_body_1EDFE(m,h,t);break;
case 0xf26e:dl_body_1F26E(m,h,t);break;
case 0xf3cf:dl_body_1F3CF(m,h,t);break;
case 0xf426:dl_body_1F426(m,h,t);break;
case 0xf43b:dl_body_1F43B(m,h,t);break;
case 0xf465:dl_body_1F465(m,h,t);break;
default:return 0;}return 1;}
