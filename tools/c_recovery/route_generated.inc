/* Original native C route/collapse closure; re/129, spec/249. Raw loc_1491B retains its original name. */
/* Original sub_1291A; IDA 0x1291A. */
static void rt_body_1291A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1291A: ; /* IDA 0x1291A; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_1291D: ; /* IDA 0x1291D; jnb     short loc_12920 */
    if(!(m->flags&1)) goto L_12920;
L_1291F: ; /* IDA 0x1291F; retn */
    return;
L_12920: ; /* IDA 0x12920; push    ax */
    ec_push(m,m->ax);
L_12921: ; /* IDA 0x12921; push    bx */
    ec_push(m,m->bx);
L_12922: ; /* IDA 0x12922; mov     cs:byte_12919, al */
    dl_write8(m,m->cs,10521,(uint8_t)m->ax);
L_12926: ; /* IDA 0x12926; call    sub_12BA8 */
    rt_call(m,11176,10537,h);
L_12929: ; /* IDA 0x12929; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_1292C: ; /* IDA 0x1292C; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_1292E: ; /* IDA 0x1292E; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_12930: ; /* IDA 0x12930; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_12932: ; /* IDA 0x12932; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_12935: ; /* IDA 0x12935; mov     ah, [bx+3] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_12938: ; /* IDA 0x12938; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_1293A: ; /* IDA 0x1293A; sub     bx, 2240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,8768,16,0,1));
L_1293E: ; /* IDA 0x1293E; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_12940: ; /* IDA 0x12940; add     bx, 4240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16960,16,0,0));
L_12944: ; /* IDA 0x12944; cmp     ah, 0FFh */
    ec_cmp(m,(uint8_t)(m->ax>>8),255,8);
L_12947: ; /* IDA 0x12947; jz      short loc_12971 */
    if(m->flags&0x40) goto L_12971;
L_12949: ; /* IDA 0x12949; cmp     al, [si+2] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+2)),8);
L_1294C: ; /* IDA 0x1294C; jz      short loc_1296C */
    if(m->flags&0x40) goto L_1296C;
L_1294E: ; /* IDA 0x1294E; mov     al, cs:byte_12919 */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,10521)));
L_12952: ; /* IDA 0x12952; cmp     al, [si+1] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_12955: ; /* IDA 0x12955; jz      short loc_1296C */
    if(m->flags&0x40) goto L_1296C;
L_12957: ; /* IDA 0x12957; cmp     al, 18h */
    ec_cmp(m,(uint8_t)m->ax,24,8);
L_12959: ; /* IDA 0x12959; jz      short loc_1296C */
    if(m->flags&0x40) goto L_1296C;
L_1295B: ; /* IDA 0x1295B; call    sub_1ECE0 */
    rt_call(m,60640,10590,h);
L_1295E: ; /* IDA 0x1295E; and     al, 7Fh */
    {uint16_t result=(((uint8_t)m->ax)&(127))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_12960: ; /* IDA 0x12960; mov     ah, [bx+1Fh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+31))));
L_12963: ; /* IDA 0x12963; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_12965: ; /* IDA 0x12965; add     ah, 28h ; '(' */
#if KI_ROUTE_MUTATION == 13
    vg_ah(m,ec_math(m,(uint8_t)(m->ax>>8),0x27,8,0,0));
#else
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),40,8,0,0)));
#endif
L_12968: ; /* IDA 0x12968; cmp     al, ah */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
L_1296A: ; /* IDA 0x1296A; ja      short loc_12971 */
    if(!(m->flags&(1|0x40))) goto L_12971;
L_1296C: ; /* IDA 0x1296C; call    sub_12977 */
    rt_call(m,10615,10607,h);
L_1296F: ; /* IDA 0x1296F; jmp     short loc_12974 */
    goto L_12974;
L_12971: ; /* IDA 0x12971; call    sub_129C3 */
    rt_call(m,10691,10612,h);
L_12974: ; /* IDA 0x12974; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_12975: ; /* IDA 0x12975; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_12976: ; /* IDA 0x12976; retn */
    return;
}
/* Original sub_12977; IDA 0x12977. */
static void rt_body_12977(KiMachine16 *m,const KiEconomyHooks *h) {
L_12977: ; /* IDA 0x12977; push    ax */
    ec_push(m,m->ax);
L_12978: ; /* IDA 0x12978; push    cx */
    ec_push(m,m->cx);
L_12979: ; /* IDA 0x12979; push    si */
    ec_push(m,m->si);
L_1297A: ; /* IDA 0x1297A; push    di */
    ec_push(m,m->di);
L_1297B: ; /* IDA 0x1297B; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_1297D: ; /* IDA 0x1297D; sub     si, 4240h */
    m->si=(uint16_t)(ec_math(m,m->si,16960,16,0,1));
L_12981: ; /* IDA 0x12981; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_12983: ; /* IDA 0x12983; add     si, 2240h */
    m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,0));
L_12987: ; /* IDA 0x12987; mov     cx, ds */
    m->cx=(uint16_t)(m->ds);
L_12989: ; /* IDA 0x12989; lds     di, [si+1Ah] */
    {uint16_t off=(uint16_t)(m->si+26),seg=m->ds;uint16_t value=dl_read16(m,seg,off),bank=dl_read16(m,seg,(uint16_t)(off+2));m->di=(uint16_t)(value);m->ds=bank;}
L_1298C: ; /* IDA 0x1298C; dec     byte ptr [di] */
#if KI_ROUTE_MUTATION == 14
    ; /* Omitted original occupancy DEC. */
#else
    dl_write8(m,m->ds,(uint16_t)(m->di),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->di)),8));
#endif
L_1298E: ; /* IDA 0x1298E; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_12990: ; /* IDA 0x12990; mov     byte ptr [si], 8 */
    dl_write8(m,m->ds,(uint16_t)(m->si),8);
L_12993: ; /* IDA 0x12993; mov     byte ptr [si+3], 30h ; '0' */
    dl_write8(m,m->ds,(uint16_t)(m->si+3),48);
L_12997: ; /* IDA 0x12997; call    sub_14689 */
    rt_call(m,18057,10650,h);
L_1299A: ; /* IDA 0x1299A; mov     cx, 1Fh */
    m->cx=(uint16_t)(31);
L_1299D: ; /* IDA 0x1299D; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_129A1: ; /* IDA 0x129A1; cmp     al, [si+1] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_129A4: ; /* IDA 0x129A4; jz      short loc_129AE */
    if(m->flags&0x40) goto L_129AE;
L_129A6: ; /* IDA 0x129A6; cmp     al, cs:byte_12919 */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,10521),8);
L_129AB: ; /* IDA 0x129AB; jnz     short loc_129BE */
    if(!(m->flags&0x40)) goto L_129BE;
L_129AD: ; /* IDA 0x129AD; inc     cx */
    m->cx=(uint16_t)(ec_inc(m,m->cx,16));
L_129AE: ; /* IDA 0x129AE; mov     al, [si+2] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_129B1: ; /* IDA 0x129B1; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_129B3: ; /* IDA 0x129B3; push    ax */
    ec_push(m,m->ax);
L_129B4: ; /* IDA 0x129B4; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_129B6: ; /* IDA 0x129B6; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_129B8: ; /* IDA 0x129B8; call    sub_18810 */
    rt_call(m,34832,10683,h);
L_129BB: ; /* IDA 0x129BB; add     sp, 2 */
    m->sp=(uint16_t)(ec_math(m,m->sp,2,16,0,0));
L_129BE: ; /* IDA 0x129BE; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_129BF: ; /* IDA 0x129BF; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_129C0: ; /* IDA 0x129C0; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_129C1: ; /* IDA 0x129C1; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_129C2: ; /* IDA 0x129C2; retn */
    return;
}
/* Original sub_129C3; IDA 0x129C3. */
static void rt_body_129C3(KiMachine16 *m,const KiEconomyHooks *h) {
L_129C3: ; /* IDA 0x129C3; push    ax */
    ec_push(m,m->ax);
L_129C4: ; /* IDA 0x129C4; push    cx */
    ec_push(m,m->cx);
L_129C5: ; /* IDA 0x129C5; push    dx */
    ec_push(m,m->dx);
L_129C6: ; /* IDA 0x129C6; push    si */
    ec_push(m,m->si);
L_129C7: ; /* IDA 0x129C7; push    di */
    ec_push(m,m->di);
L_129C8: ; /* IDA 0x129C8; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_129CA: ; /* IDA 0x129CA; sub     si, 4240h */
    m->si=(uint16_t)(ec_math(m,m->si,16960,16,0,1));
L_129CE: ; /* IDA 0x129CE; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_129D0: ; /* IDA 0x129D0; add     si, 2240h */
    m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,0));
L_129D4: ; /* IDA 0x129D4; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_129D7: ; /* IDA 0x129D7; jb      short loc_129E5 */
    if(m->flags&1) goto L_129E5;
L_129D9: ; /* IDA 0x129D9; mov     cx, ds */
    m->cx=(uint16_t)(m->ds);
L_129DB: ; /* IDA 0x129DB; lds     di, [si+1Ah] */
    {uint16_t off=(uint16_t)(m->si+26),seg=m->ds;uint16_t value=dl_read16(m,seg,off),bank=dl_read16(m,seg,(uint16_t)(off+2));m->di=(uint16_t)(value);m->ds=bank;}
L_129DE: ; /* IDA 0x129DE; dec     byte ptr [di] */
    dl_write8(m,m->ds,(uint16_t)(m->di),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->di)),8));
L_129E0: ; /* IDA 0x129E0; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_129E2: ; /* IDA 0x129E2; call    sub_14689 */
    rt_call(m,18057,10725,h);
L_129E5: ; /* IDA 0x129E5; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_129E7: ; /* IDA 0x129E7; mov     ah, [bx+1Ch] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+28))));
L_129EA: ; /* IDA 0x129EA; call    sub_12AD2 */
    rt_call(m,10962,10733,h);
L_129ED: ; /* IDA 0x129ED; mov     byte ptr [si], 0 */
#if KI_ROUTE_MUTATION == 15
    dl_write8(m,m->ds,m->si,8);
#else
    dl_write8(m,m->ds,(uint16_t)(m->si),0);
#endif
L_129F0: ; /* IDA 0x129F0; mov     byte ptr [bx+17h], 4 */
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),4);
L_129F4: ; /* IDA 0x129F4; mov     al, cs:byte_12919 */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,10521)));
L_129F8: ; /* IDA 0x129F8; xchg    al, [bx+1Ch] */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->ds,(uint16_t)(m->bx+28));vg_al(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->bx+28),a);}
L_129FB: ; /* IDA 0x129FB; mov     [bx+1Dh], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+29),(uint8_t)m->ax);
L_129FE: ; /* IDA 0x129FE; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_12A00: ; /* IDA 0x12A00; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_12A02: ; /* IDA 0x12A02; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12A04: ; /* IDA 0x12A04; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12A06: ; /* IDA 0x12A06; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_12A08: ; /* IDA 0x12A08; cmp     byte ptr [di], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di)),128,8);
L_12A0B: ; /* IDA 0x12A0B; jnb     short loc_12A12 */
    if(!(m->flags&1)) goto L_12A12;
L_12A0D: ; /* IDA 0x12A0D; test    byte ptr [bx], 10h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(16))&255;ec_logic(m,result,8);}
L_12A10: ; /* IDA 0x12A10; jnz     short loc_12A57 */
    if(!(m->flags&0x40)) goto L_12A57;
L_12A12: ; /* IDA 0x12A12; test    byte ptr [bx], 40h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(64))&255;ec_logic(m,result,8);}
L_12A15: ; /* IDA 0x12A15; jz      short loc_12A1E */
    if(m->flags&0x40) goto L_12A1E;
L_12A17: ; /* IDA 0x12A17; and     byte ptr [bx], 0BFh */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(191))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->bx),result);}
L_12A1A: ; /* IDA 0x12A1A; add     byte ptr [bx+1Eh], 3 */
    dl_write8(m,m->ds,(uint16_t)(m->bx+30),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->bx+30)),3,8,0,0));
L_12A1E: ; /* IDA 0x12A1E; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_12A20: ; /* IDA 0x12A20; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_12A24: ; /* IDA 0x12A24; cmp     al, [bx+1Dh] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+29)),8);
L_12A27: ; /* IDA 0x12A27; jz      short loc_12A31 */
    if(m->flags&0x40) goto L_12A31;
L_12A29: ; /* IDA 0x12A29; cmp     al, cs:byte_12919 */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,10521),8);
L_12A2E: ; /* IDA 0x12A2E; jnz     short loc_12A78 */
    if(!(m->flags&0x40)) goto L_12A78;
L_12A30: ; /* IDA 0x12A30; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_12A31: ; /* IDA 0x12A31; push    bx */
    ec_push(m,m->bx);
L_12A32: ; /* IDA 0x12A32; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_12A34: ; /* IDA 0x12A34; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_12A36: ; /* IDA 0x12A36; mov     cx, 21h ; '!' */
    m->cx=(uint16_t)(33);
L_12A39: ; /* IDA 0x12A39; add     cx, dx */
    m->cx=(uint16_t)(ec_math(m,m->cx,m->dx,16,0,0));
L_12A3B: ; /* IDA 0x12A3B; call    sub_18810 */
    rt_call(m,34832,10814,h);
L_12A3E: ; /* IDA 0x12A3E; add     sp, 2 */
    m->sp=(uint16_t)(ec_math(m,m->sp,2,16,0,0));
L_12A41: ; /* IDA 0x12A41; and     dx, dx */
    {uint16_t result=((m->dx)&(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_12A43: ; /* IDA 0x12A43; jz      short loc_12A78 */
    if(m->flags&0x40) goto L_12A78;
L_12A45: ; /* IDA 0x12A45; push    bx */
    ec_push(m,m->bx);
L_12A46: ; /* IDA 0x12A46; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_12A48: ; /* IDA 0x12A48; mov     ah, [bx+1Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
L_12A4B: ; /* IDA 0x12A4B; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_12A4E: ; /* IDA 0x12A4E; mov     cx, 19Ah */
    m->cx=(uint16_t)(410);
L_12A51: ; /* IDA 0x12A51; call    sub_18810 */
    rt_call(m,34832,10836,h);
L_12A54: ; /* IDA 0x12A54; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_12A55: ; /* IDA 0x12A55; jmp     short loc_12A78 */
    goto L_12A78;
L_12A57: ; /* IDA 0x12A57; mov     byte ptr [bx], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->bx),0);
L_12A5A: ; /* IDA 0x12A5A; mov     word ptr [bx+1Ch], 0FFFFh */
    dl_write16(m,m->ds,(uint16_t)(m->bx+28),65535);
L_12A5F: ; /* IDA 0x12A5F; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_12A63: ; /* IDA 0x12A63; cmp     al, cs:byte_12919 */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,10521),8);
L_12A68: ; /* IDA 0x12A68; jnz     short loc_12A78 */
    if(!(m->flags&0x40)) goto L_12A78;
L_12A6A: ; /* IDA 0x12A6A; push    bx */
    ec_push(m,m->bx);
L_12A6B: ; /* IDA 0x12A6B; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_12A6D: ; /* IDA 0x12A6D; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_12A6F: ; /* IDA 0x12A6F; mov     cx, 43h ; 'C' */
    m->cx=(uint16_t)(67);
L_12A72: ; /* IDA 0x12A72; call    sub_18810 */
    rt_call(m,34832,10869,h);
L_12A75: ; /* IDA 0x12A75; add     sp, 2 */
    m->sp=(uint16_t)(ec_math(m,m->sp,2,16,0,0));
L_12A78: ; /* IDA 0x12A78; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_12A79: ; /* IDA 0x12A79; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_12A7A: ; /* IDA 0x12A7A; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_12A7B: ; /* IDA 0x12A7B; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_12A7C: ; /* IDA 0x12A7C; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_12A7D: ; /* IDA 0x12A7D; retn */
    return;
}
/* Original sub_12BA8; IDA 0x12BA8. */
static void rt_body_12BA8(KiMachine16 *m,const KiEconomyHooks *h) {
L_12BA8: ; /* IDA 0x12BA8; and     byte ptr [si], 0EFh */
#if KI_ROUTE_MUTATION == 16
    {uint8_t result=dl_read8(m,m->ds,m->si)&0xdf;ec_logic(m,result,8);dl_write8(m,m->ds,m->si,result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(239))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#endif
L_12BAB: ; /* IDA 0x12BAB; test    cs:byte_198A6, 4 */
    {uint16_t result=((dl_read8(m,m->cs,39078))&(4))&255;ec_logic(m,result,8);}
L_12BB1: ; /* IDA 0x12BB1; jz      short locret_12BD8 */
    if(m->flags&0x40) goto L_12BD8;
L_12BB3: ; /* IDA 0x12BB3; push    bx */
    ec_push(m,m->bx);
L_12BB4: ; /* IDA 0x12BB4; push    dx */
    ec_push(m,m->dx);
L_12BB5: ; /* IDA 0x12BB5; mov     dx, [si+10h] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+16)));
L_12BB8: ; /* IDA 0x12BB8; mov     bx, [si+12h] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+18)));
L_12BBB: ; /* IDA 0x12BBB; call    sub_19656 */
    rt_call(m,38486,11198,h);
L_12BBE: ; /* IDA 0x12BBE; mov     dx, cs:word_1988E */
    m->dx=(uint16_t)(dl_read16(m,m->cs,39054));
L_12BC3: ; /* IDA 0x12BC3; mov     bx, cs:word_19890 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,39056));
L_12BC8: ; /* IDA 0x12BC8; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_12BCA: ; /* IDA 0x12BCA; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_12BCC: ; /* IDA 0x12BCC; add     dx, 1B8h */
    m->dx=(uint16_t)(ec_math(m,m->dx,440,16,0,0));
L_12BD0: ; /* IDA 0x12BD0; add     bx, 28h ; '(' */
    m->bx=(uint16_t)(ec_math(m,m->bx,40,16,0,0));
L_12BD3: ; /* IDA 0x12BD3; call    sub_196ED */
    rt_call(m,38637,11222,h);
L_12BD6: ; /* IDA 0x12BD6; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_12BD7: ; /* IDA 0x12BD7; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_12BD8: ; /* IDA 0x12BD8; retn */
    return;
}
/* Original sub_147BB; IDA 0x147BB. */
static void rt_body_147BB(KiMachine16 *m,const KiEconomyHooks *h) {
L_147BB: ; /* IDA 0x147BB; push    es */
    ec_push(m,m->es);
L_147BC: ; /* IDA 0x147BC; push    si */
    ec_push(m,m->si);
L_147BD: ; /* IDA 0x147BD; push    di */
    ec_push(m,m->di);
L_147BE: ; /* IDA 0x147BE; push    bp */
    ec_push(m,m->bp);
L_147BF: ; /* IDA 0x147BF; mov     ax, [si+14h] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+20)));
L_147C2: ; /* IDA 0x147C2; mov     bx, [si+0Eh] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
L_147C5: ; /* IDA 0x147C5; mov     dl, [si+1] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_147C8: ; /* IDA 0x147C8; mov     es, cs:word_19874 */
    m->es=(uint16_t)(dl_read16(m,m->cs,39028));
L_147CD: ; /* IDA 0x147CD; cmp     bx, 800h */
    ec_cmp(m,m->bx,2048,16);
L_147D1: ; /* IDA 0x147D1; jb      short loc_1482F */
    if(m->flags&1) goto L_1482F;
L_147D3: ; /* IDA 0x147D3; mov     cx, es:[bx+6] */
    m->cx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+6)));
L_147D7: ; /* IDA 0x147D7; mov     bx, es:[bx+8] */
    m->bx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+8)));
L_147DB: ; /* IDA 0x147DB; cmp     ax, bx */
    ec_cmp(m,m->ax,m->bx,16);
L_147DD: ; /* IDA 0x147DD; jz      short loc_14821 */
    if(m->flags&0x40) goto L_14821;
L_147DF: ; /* IDA 0x147DF; cmp     ax, cx */
    ec_cmp(m,m->ax,m->cx,16);
L_147E1: ; /* IDA 0x147E1; jz      short loc_14828 */
    if(m->flags&0x40) goto L_14828;
L_147E3: ; /* IDA 0x147E3; mov     bp, bx */
    m->bp=(uint16_t)(m->bx);
L_147E5: ; /* IDA 0x147E5; mov     ds, cs:word_19874 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,39028));
L_147EA: ; /* IDA 0x147EA; call    loc_1491B */
    rt_call(m,18715,18413,h);
L_147ED: ; /* IDA 0x147ED; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_147F2: ; /* IDA 0x147F2; cmp     cx, 8000h */
    ec_cmp(m,m->cx,32768,16);
L_147F6: ; /* IDA 0x147F6; jb      short loc_14806 */
    if(m->flags&1) goto L_14806;
L_147F8: ; /* IDA 0x147F8; cmp     byte ptr [si+23h], 0Ah */
#if KI_ROUTE_MUTATION == 11
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+0x23)),11,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+35)),10,8);
#endif
L_147FC: ; /* IDA 0x147FC; jb      short loc_14806 */
    if(m->flags&1) goto L_14806;
L_147FE: ; /* IDA 0x147FE; mov     al, [si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14801: ; /* IDA 0x14801; call    sub_1291A */
    rt_call(m,10522,18436,h);
L_14804: ; /* IDA 0x14804; jmp     short loc_14875 */
    goto L_14875;
L_14806: ; /* IDA 0x14806; cmp     al, 4 */
    ec_cmp(m,(uint8_t)m->ax,4,8);
L_14808: ; /* IDA 0x14808; jnz     short loc_14810 */
    if(!(m->flags&0x40)) goto L_14810;
L_1480A: ; /* IDA 0x1480A; mov     dx, es:[bx+6] */
    m->dx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+6)));
L_1480E: ; /* IDA 0x1480E; jmp     short loc_14814 */
    goto L_14814;
L_14810: ; /* IDA 0x14810; mov     dx, es:[bx+8] */
    m->dx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+8)));
L_14814: ; /* IDA 0x14814; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_14816: ; /* IDA 0x14816; or      byte ptr [si], 1 */
#if KI_ROUTE_MUTATION == 12
    ; /* Omitted original leg bit0 OR. */
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(1))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#endif
L_14819: ; /* IDA 0x14819; cmp     dx, bp */
    ec_cmp(m,m->dx,m->bp,16);
L_1481B: ; /* IDA 0x1481B; jz      short loc_1486C */
    if(m->flags&0x40) goto L_1486C;
L_1481D: ; /* IDA 0x1481D; neg     al */
    vg_al(m,(uint8_t)(ec_math(m,0,(uint8_t)m->ax,8,0,1)));
L_1481F: ; /* IDA 0x1481F; jmp     short loc_1486C */
    goto L_1486C;
L_14821: ; /* IDA 0x14821; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_14823: ; /* IDA 0x14823; or      byte ptr [si], 1 */
#if KI_ROUTE_MUTATION == 12
    ; /* Omitted original leg bit0 OR. */
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(1))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#endif
L_14826: ; /* IDA 0x14826; jmp     short loc_1486C */
    goto L_1486C;
L_14828: ; /* IDA 0x14828; mov     al, 0FCh */
    vg_al(m,(uint8_t)(252));
L_1482A: ; /* IDA 0x1482A; or      byte ptr [si], 1 */
#if KI_ROUTE_MUTATION == 12
    ; /* Omitted original leg bit0 OR. */
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(1))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#endif
L_1482D: ; /* IDA 0x1482D; jmp     short loc_1486C */
    goto L_1486C;
L_1482F: ; /* IDA 0x1482F; cmp     ax, bx */
    ec_cmp(m,m->ax,m->bx,16);
L_14831: ; /* IDA 0x14831; jz      short loc_1486F */
    if(m->flags&0x40) goto L_1486F;
L_14833: ; /* IDA 0x14833; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_14835: ; /* IDA 0x14835; mov     ds, cs:word_19874 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,39028));
L_1483A: ; /* IDA 0x1483A; call    loc_1491B */
    rt_call(m,18715,18493,h);
L_1483D: ; /* IDA 0x1483D; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_14842: ; /* IDA 0x14842; cmp     cx, 8000h */
    ec_cmp(m,m->cx,32768,16);
L_14846: ; /* IDA 0x14846; jb      short loc_14856 */
    if(m->flags&1) goto L_14856;
L_14848: ; /* IDA 0x14848; cmp     byte ptr [si+23h], 0Ah */
#if KI_ROUTE_MUTATION == 11
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+0x23)),11,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+35)),10,8);
#endif
L_1484C: ; /* IDA 0x1484C; jb      short loc_14856 */
    if(m->flags&1) goto L_14856;
L_1484E: ; /* IDA 0x1484E; mov     al, [si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14851: ; /* IDA 0x14851; call    sub_1291A */
    rt_call(m,10522,18516,h);
L_14854: ; /* IDA 0x14854; jmp     short loc_14875 */
    goto L_14875;
L_14856: ; /* IDA 0x14856; cmp     al, 4 */
    ec_cmp(m,(uint8_t)m->ax,4,8);
L_14858: ; /* IDA 0x14858; jnz     short loc_1485F */
    if(!(m->flags&0x40)) goto L_1485F;
L_1485A: ; /* IDA 0x1485A; mov     dx, es:[bx] */
    m->dx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx)));
L_1485D: ; /* IDA 0x1485D; jmp     short loc_14863 */
    goto L_14863;
L_1485F: ; /* IDA 0x1485F; mov     dx, es:[bx+2] */
    m->dx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+2)));
L_14863: ; /* IDA 0x14863; mov     [si+0Ch], dx */
    dl_write16(m,m->ds,(uint16_t)(m->si+12),m->dx);
L_14866: ; /* IDA 0x14866; mov     [si+0Eh], bx */
    dl_write16(m,m->ds,(uint16_t)(m->si+14),m->bx);
L_14869: ; /* IDA 0x14869; and     byte ptr [si], 0FEh */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(254))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_1486C: ; /* IDA 0x1486C; mov     [si+0Ah], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+10),(uint8_t)m->ax);
L_1486F: ; /* IDA 0x1486F; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_14870: ; /* IDA 0x14870; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14871: ; /* IDA 0x14871; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14872: ; /* IDA 0x14872; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_14873: ; /* IDA 0x14873; clc */
    m->flags&=(uint16_t)~1u;
L_14874: ; /* IDA 0x14874; retn */
    return;
L_14875: ; /* IDA 0x14875; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_14876: ; /* IDA 0x14876; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14877: ; /* IDA 0x14877; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14878: ; /* IDA 0x14878; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_14879: ; /* IDA 0x14879; stc */
    m->flags|=1;
L_1487A: ; /* IDA 0x1487A; retn */
    return;
}
/* Original sub_1487B; IDA 0x1487B. */
static void rt_body_1487B(KiMachine16 *m,const KiEconomyHooks *h) {
L_1487B: ; /* IDA 0x1487B; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_1487E: ; /* IDA 0x1487E; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_14880: ; /* IDA 0x14880; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14882: ; /* IDA 0x14882; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14884: ; /* IDA 0x14884; mov     al, [bx+3] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_14887: ; /* IDA 0x14887; cmp     al, 0FFh */
    ec_cmp(m,(uint8_t)m->ax,255,8);
L_14889: ; /* IDA 0x14889; jnz     short loc_1488D */
    if(!(m->flags&0x40)) goto L_1488D;
L_1488B: ; /* IDA 0x1488B; stc */
    m->flags|=1;
L_1488C: ; /* IDA 0x1488C; retn */
    return;
L_1488D: ; /* IDA 0x1488D; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1488F: ; /* IDA 0x1488F; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14891: ; /* IDA 0x14891; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14893: ; /* IDA 0x14893; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14895: ; /* IDA 0x14895; push    es */
    ec_push(m,m->es);
L_14896: ; /* IDA 0x14896; push    di */
    ec_push(m,m->di);
L_14897: ; /* IDA 0x14897; mov     dl, [si+1] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_1489A: ; /* IDA 0x1489A; mov     bx, [si+0Eh] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
L_1489D: ; /* IDA 0x1489D; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_1489F: ; /* IDA 0x1489F; mov     ds, cs:word_19874 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,39028));
L_148A4: ; /* IDA 0x148A4; cmp     bx, 800h */
    ec_cmp(m,m->bx,2048,16);
L_148A8: ; /* IDA 0x148A8; jb      short loc_148E5 */
    if(m->flags&1) goto L_148E5;
L_148AA: ; /* IDA 0x148AA; mov     es, cs:word_10D52 */
    m->es=(uint16_t)(dl_read16(m,m->cs,3410));
L_148AF: ; /* IDA 0x148AF; mov     cx, [bx+6] */
    m->cx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+6)));
L_148B2: ; /* IDA 0x148B2; mov     bx, [bx+8] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+8)));
L_148B5: ; /* IDA 0x148B5; mov     di, cx */
    m->di=(uint16_t)(m->cx);
L_148B7: ; /* IDA 0x148B7; shl     di, 1 */
    m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
L_148B9: ; /* IDA 0x148B9; shl     di, 1 */
    m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
L_148BB: ; /* IDA 0x148BB; cmp     dl, es:[di+841h] */
    ec_cmp(m,(uint8_t)m->dx,dl_read8(m,m->es,(uint16_t)(m->di+2113)),8);
L_148C0: ; /* IDA 0x148C0; jz      short loc_148C4 */
    if(m->flags&0x40) goto L_148C4;
L_148C2: ; /* IDA 0x148C2; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_148C4: ; /* IDA 0x148C4; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_148C6: ; /* IDA 0x148C6; shl     di, 1 */
    m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
L_148C8: ; /* IDA 0x148C8; shl     di, 1 */
    m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
L_148CA: ; /* IDA 0x148CA; cmp     dl, es:[di+841h] */
    ec_cmp(m,(uint8_t)m->dx,dl_read8(m,m->es,(uint16_t)(m->di+2113)),8);
L_148CF: ; /* IDA 0x148CF; jz      short loc_148D7 */
    if(m->flags&0x40) goto L_148D7;
L_148D1: ; /* IDA 0x148D1; cmp     bx, cx */
    ec_cmp(m,m->bx,m->cx,16);
L_148D3: ; /* IDA 0x148D3; jz      short loc_14903 */
    if(m->flags&0x40) goto L_14903;
L_148D5: ; /* IDA 0x148D5; mov     bx, cx */
    m->bx=(uint16_t)(m->cx);
L_148D7: ; /* IDA 0x148D7; call    loc_1491B */
    rt_call(m,18715,18650,h);
L_148DA: ; /* IDA 0x148DA; jb      short loc_1490C */
    if(m->flags&1) goto L_1490C;
L_148DC: ; /* IDA 0x148DC; cmp     ax, 0FFFCh */
    ec_cmp(m,m->ax,65532,16);
L_148DF: ; /* IDA 0x148DF; jnz     short loc_148F1 */
    if(!(m->flags&0x40)) goto L_148F1;
L_148E1: ; /* IDA 0x148E1; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_148E2: ; /* IDA 0x148E2; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_148E3: ; /* IDA 0x148E3; jmp     short loc_148F1 */
    goto L_148F1;
L_148E5: ; /* IDA 0x148E5; call    loc_1491B */
    rt_call(m,18715,18664,h);
L_148E8: ; /* IDA 0x148E8; jb      short loc_1490C */
    if(m->flags&1) goto L_1490C;
L_148EA: ; /* IDA 0x148EA; cmp     ax, 4 */
    ec_cmp(m,m->ax,4,16);
L_148ED: ; /* IDA 0x148ED; jnz     short loc_148F1 */
    if(!(m->flags&0x40)) goto L_148F1;
L_148EF: ; /* IDA 0x148EF; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_148F0: ; /* IDA 0x148F0; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_148F1: ; /* IDA 0x148F1; mov     bx, [bx+6] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+6)));
L_148F4: ; /* IDA 0x148F4; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_148F9: ; /* IDA 0x148F9; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_148FB: ; /* IDA 0x148FB; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_148FD: ; /* IDA 0x148FD; cmp     dl, [bx+841h] */
    ec_cmp(m,(uint8_t)m->dx,dl_read8(m,m->ds,(uint16_t)(m->bx+2113)),8);
L_14901: ; /* IDA 0x14901; jz      short loc_14917 */
    if(m->flags&0x40) goto L_14917;
L_14903: ; /* IDA 0x14903; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_14908: ; /* IDA 0x14908; stc */
    m->flags|=1;
L_14909: ; /* IDA 0x14909; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1490A: ; /* IDA 0x1490A; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1490B: ; /* IDA 0x1490B; retn */
    return;
L_1490C: ; /* IDA 0x1490C; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_14911: ; /* IDA 0x14911; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14913: ; /* IDA 0x14913; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14915: ; /* IDA 0x14915; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_14917: ; /* IDA 0x14917; clc */
    m->flags&=(uint16_t)~1u;
L_14918: ; /* IDA 0x14918; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14919: ; /* IDA 0x14919; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1491A: ; /* IDA 0x1491A; retn */
    return;
}
/* Original loc_1491B raw entry, complete through original RET; IDA 0x1491B. */
static void rt_body_1491B(KiMachine16 *m,const KiEconomyHooks *h) {
L_1491B: ; /* IDA 0x1491B; cmp     ax, bx */
    ec_cmp(m,m->ax,m->bx,16);
L_1491D: ; /* IDA 0x1491D; jz      short loc_14923 */
    if(m->flags&0x40) goto L_14923;
L_1491F: ; /* IDA 0x1491F; cmp     ax, cx */
    ec_cmp(m,m->ax,m->cx,16);
L_14921: ; /* IDA 0x14921; jnz     short loc_14925 */
    if(!(m->flags&0x40)) goto L_14925;
L_14923: ; /* IDA 0x14923; stc */
    m->flags|=1;
L_14924: ; /* IDA 0x14924; retn */
    return;
L_14925: ; /* IDA 0x14925; push    es */
    ec_push(m,m->es);
L_14926: ; /* IDA 0x14926; push    dx */
    ec_push(m,m->dx);
L_14927: ; /* IDA 0x14927; push    si */
    ec_push(m,m->si);
L_14928: ; /* IDA 0x14928; push    di */
    ec_push(m,m->di);
L_14929: ; /* IDA 0x14929; push    bp */
    ec_push(m,m->bp);
L_1492A: ; /* IDA 0x1492A; cld */
    m->flags&=(uint16_t)~0x400u;
L_1492B: ; /* IDA 0x1492B; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1492D: ; /* IDA 0x1492D; mov     ax, ds */
    m->ax=(uint16_t)(m->ds);
L_1492F: ; /* IDA 0x1492F; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_14931: ; /* IDA 0x14931; mov     cs:word_149B8, bx */
    dl_write16(m,m->cs,18872,m->bx);
L_14936: ; /* IDA 0x14936; mov     cs:word_149BE, cx */
    dl_write16(m,m->cs,18878,m->cx);
L_1493B: ; /* IDA 0x1493B; mov     cs:byte_149D2, dl */
    dl_write8(m,m->cs,18898,(uint8_t)m->dx);
L_14940: ; /* IDA 0x14940; mov     di, 8000h */
    m->di=(uint16_t)(32768);
L_14943: ; /* IDA 0x14943; mov     cx, 400h */
    m->cx=(uint16_t)(1024);
L_14946: ; /* IDA 0x14946; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_14948: ; /* IDA 0x14948; rep stosw */
    mc_string(m,2,1);
L_1494A: ; /* IDA 0x1494A; mov     es, cs:word_10D52 */
    m->es=(uint16_t)(dl_read16(m,m->cs,3410));
L_1494F: ; /* IDA 0x1494F; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_14951: ; /* IDA 0x14951; mov     bx, 8800h */
    m->bx=(uint16_t)(34816);
L_14954: ; /* IDA 0x14954; mov     di, 87F8h */
    m->di=(uint16_t)(34808);
L_14957: ; /* IDA 0x14957; mov     cx, 1 */
    m->cx=(uint16_t)(1);
L_1495A: ; /* IDA 0x1495A; jmp     short near ptr byte_149B6 */
    goto L_149B6;
L_1495C: ; /* IDA 0x1495C; mov     ax, 0FFFFh */
    m->ax=(uint16_t)(65535);
L_1495F: ; /* IDA 0x1495F; mov     cl, ch */
    vg_cl(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_14961: ; /* IDA 0x14961; xchg    bx, di */
    {uint16_t a=m->bx,b=m->di;m->bx=(uint16_t)(b);m->di=(uint16_t)(a);}
L_14963: ; /* IDA 0x14963; mov     dx, bx */
    m->dx=(uint16_t)(m->bx);
L_14965: ; /* IDA 0x14965; add     bl, 2 */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,2,8,0,0)));
L_14968: ; /* IDA 0x14968; cmp     ax, [bx] */
    ec_cmp(m,m->ax,dl_read16(m,m->ds,(uint16_t)(m->bx)),16);
L_1496A: ; /* IDA 0x1496A; jbe     short near ptr unk_1496E */
    if(m->flags&(1|0x40)) goto L_1496E;
L_1496C: ; /* IDA 0x1496C; mov     ax, [bx] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx)));
L_1496E: ; /* IDA 0x1496E; add     bx, 8 */
    m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,0));
L_14971: ; /* IDA 0x14971; and     bh, 0FBh */
    {uint16_t result=(((uint8_t)(m->bx>>8))&(251))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_14974: ; /* IDA 0x14974; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_14976: ; /* IDA 0x14976; jnz     short loc_14968 */
    if(!(m->flags&0x40)) goto L_14968;
L_14978: ; /* IDA 0x14978; mov     bx, di */
    m->bx=(uint16_t)(m->di);
L_1497A: ; /* IDA 0x1497A; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_1497C: ; /* IDA 0x1497C; xchg    cl, ch */
    {uint16_t a=(uint8_t)m->cx,b=(uint8_t)(m->cx>>8);vg_cl(m,(uint8_t)(b));vg_ch(m,(uint8_t)(a));}
L_1497E: ; /* IDA 0x1497E; mov     si, [di] */
    m->si=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di)));
L_14980: ; /* IDA 0x14980; mov     dx, [di+2] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+2)));
L_14983: ; /* IDA 0x14983; cmp     ax, dx */
    ec_cmp(m,m->ax,m->dx,16);
L_14985: ; /* IDA 0x14985; jz      short near ptr byte_149B6 */
#if KI_ROUTE_MUTATION == 9
    if(m->flags&(1|0x40)) goto L_149B6;
#else
    if(m->flags&0x40) goto L_149B6;
#endif
L_14987: ; /* IDA 0x14987; mov     [bx], si */
    dl_write16(m,m->ds,(uint16_t)(m->bx),m->si);
L_14989: ; /* IDA 0x14989; mov     [bx+2], dx */
    dl_write16(m,m->ds,(uint16_t)(m->bx+2),m->dx);
L_1498C: ; /* IDA 0x1498C; mov     dx, [di+4] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+4)));
L_1498F: ; /* IDA 0x1498F; mov     [bx+4], dx */
    dl_write16(m,m->ds,(uint16_t)(m->bx+4),m->dx);
L_14992: ; /* IDA 0x14992; mov     dx, [di+6] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+6)));
L_14995: ; /* IDA 0x14995; mov     [bx+6], dx */
    dl_write16(m,m->ds,(uint16_t)(m->bx+6),m->dx);
L_14998: ; /* IDA 0x14998; add     bx, 8 */
    m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,0));
L_1499B: ; /* IDA 0x1499B; and     bh, 0FBh */
    {uint16_t result=(((uint8_t)(m->bx>>8))&(251))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1499E: ; /* IDA 0x1499E; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_149A0: ; /* IDA 0x149A0; add     di, 8 */
    m->di=(uint16_t)(ec_math(m,m->di,8,16,0,0));
L_149A3: ; /* IDA 0x149A3; and     di, 0FBFFh */
    {uint16_t result=((m->di)&(64511))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_149A7: ; /* IDA 0x149A7; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_149A9: ; /* IDA 0x149A9; jnz     short loc_1497E */
    if(!(m->flags&0x40)) goto L_1497E;
L_149AB: ; /* IDA 0x149AB; cmp     ch, 0 */
    ec_cmp(m,(uint8_t)(m->cx>>8),0,8);
L_149AE: ; /* IDA 0x149AE; jnz     short loc_1495C */
    if(!(m->flags&0x40)) goto L_1495C;
L_149B0: ; /* IDA 0x149B0; mov     bx, 800h */
    m->bx=(uint16_t)(2048);
L_149B3: ; /* IDA 0x149B3; stc */
    m->flags|=1;
L_149B4: ; /* IDA 0x149B4; jmp     short near ptr byte_149D3+36h */
    goto L_14A09;
L_149B6: ; /* IDA 0x149B6; cmp     si, 1234h */
#if KI_ROUTE_MUTATION == 1
    ec_cmp(m,m->si,dl_read16(m,m->cs,0x49be),16);
#else
    ec_cmp(m,m->si,dl_read16(m,m->cs,0x49b8),16);
#endif
L_149BA: ; /* IDA 0x149BA; jz      short near ptr byte_149D3+2Dh */
    if(m->flags&0x40) goto L_14A00;
L_149BC: ; /* IDA 0x149BC; cmp     si, 1234h */
#if KI_ROUTE_MUTATION == 2
    ec_cmp(m,m->si,dl_read16(m,m->cs,0x49b8),16);
#else
    ec_cmp(m,m->si,dl_read16(m,m->cs,0x49be),16);
#endif
L_149C0: ; /* IDA 0x149C0; jz      short near ptr byte_149D3+2Dh */
    if(m->flags&0x40) goto L_14A00;
L_149C2: ; /* IDA 0x149C2; push    ax */
    ec_push(m,m->ax);
L_149C3: ; /* IDA 0x149C3; cmp     si, 600h */
    ec_cmp(m,m->si,1536,16);
L_149C7: ; /* IDA 0x149C7; jnb     short near ptr byte_149D3+10h */
    if(!(m->flags&1)) goto L_149E3;
L_149C9: ; /* IDA 0x149C9; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_149CB: ; /* IDA 0x149CB; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_149CD: ; /* IDA 0x149CD; cmp     byte ptr es:[si+841h], 12h */
#if KI_ROUTE_MUTATION == 3
    ec_cmp(m,dl_read8(m,m->es,(uint16_t)(m->si+0x841)),0x12,8);
#else
    ec_cmp(m,dl_read8(m,m->es,(uint16_t)(m->si+0x841)),dl_read8(m,m->cs,0x49d2),8);
#endif
L_149D3: ; /* IDA 0x149D3; jz      short near ptr unk_149DC */
    if(m->flags&0x40) goto L_149DC;
L_149D5: ; /* IDA 0x149D5; add     dx, 0A6h */
#if KI_ROUTE_MUTATION == 5
    m->dx=ec_math(m,m->dx,0xa5,16,0,0);
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,166,16,0,0));
#endif
L_149D9: ; /* IDA 0x149D9; or      dh, 80h */
#if KI_ROUTE_MUTATION == 6
    {uint8_t result=(uint8_t)(m->dx>>8)|0x40;ec_logic(m,result,8);vg_dh(m,result);}
#else
    {uint16_t result=(((uint8_t)(m->dx>>8))|(128))&255;ec_logic(m,result,8);vg_dh(m,(uint8_t)(result));}
#endif
L_149DC: ; /* IDA 0x149DC; add     dx, 4 */
#if KI_ROUTE_MUTATION == 4
    m->dx=ec_math(m,m->dx,5,16,0,0);
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,4,16,0,0));
#endif
L_149DF: ; /* IDA 0x149DF; shr     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,5));
L_149E1: ; /* IDA 0x149E1; shr     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,5));
L_149E3: ; /* IDA 0x149E3; mov     ax, [si] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_149E5: ; /* IDA 0x149E5; call    sub_14A0F */
    rt_call(m,18959,18920,h);
L_149E8: ; /* IDA 0x149E8; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_149E9: ; /* IDA 0x149E9; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_149EA: ; /* IDA 0x149EA; mov     ax, [si] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_149EC: ; /* IDA 0x149EC; call    sub_14A0F */
    rt_call(m,18959,18927,h);
L_149EF: ; /* IDA 0x149EF; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_149F0: ; /* IDA 0x149F0; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_149F1: ; /* IDA 0x149F1; mov     ax, [si] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_149F3: ; /* IDA 0x149F3; call    sub_14A0F */
    rt_call(m,18959,18934,h);
L_149F6: ; /* IDA 0x149F6; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_149F7: ; /* IDA 0x149F7; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_149F8: ; /* IDA 0x149F8; mov     ax, [si] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_149FA: ; /* IDA 0x149FA; call    sub_14A0F */
    rt_call(m,18959,18941,h);
L_149FD: ; /* IDA 0x149FD; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_149FE: ; /* IDA 0x149FE; jmp     short loc_149A0 */
    goto L_149A0;
L_14A00: ; /* IDA 0x14A00; mov     ax, [di+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+4)));
L_14A03: ; /* IDA 0x14A03; mov     bx, [di+6] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+6)));
L_14A06: ; /* IDA 0x14A06; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_14A08: ; /* IDA 0x14A08; clc */
    m->flags&=(uint16_t)~1u;
L_14A09: ; /* IDA 0x14A09; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_14A0A: ; /* IDA 0x14A0A; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14A0B: ; /* IDA 0x14A0B; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14A0C: ; /* IDA 0x14A0C; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_14A0D: ; /* IDA 0x14A0D; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_14A0E: ; /* IDA 0x14A0E; retn */
    return;
}
/* Original sub_14A0F; IDA 0x14A0F. */
static void rt_body_14A0F(KiMachine16 *m,const KiEconomyHooks *h) {
L_14A0F: ; /* IDA 0x14A0F; cmp     byte ptr [si-8000h], 0 */
#if KI_ROUTE_MUTATION == 8
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+0x8000)),1,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+32768)),0,8);
#endif
L_14A14: ; /* IDA 0x14A14; jnz     short loc_14A79 */
    if(!(m->flags&0x40)) goto L_14A79;
L_14A16: ; /* IDA 0x14A16; push    ax */
    ec_push(m,m->ax);
L_14A17: ; /* IDA 0x14A17; push    dx */
    ec_push(m,m->dx);
L_14A18: ; /* IDA 0x14A18; push    di */
    ec_push(m,m->di);
L_14A19: ; /* IDA 0x14A19; push    bx */
    ec_push(m,m->bx);
L_14A1A: ; /* IDA 0x14A1A; mov     byte ptr [si-8000h], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+32768),1);
L_14A1F: ; /* IDA 0x14A1F; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_14A21: ; /* IDA 0x14A21; and     bx, 3FFFh */
    {uint16_t result=((m->bx)&(16383))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_14A25: ; /* IDA 0x14A25; and     ah, 0C0h */
    {uint16_t result=(((uint8_t)(m->ax>>8))&(192))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_14A28: ; /* IDA 0x14A28; cmp     ah, 40h ; '@' */
    ec_cmp(m,(uint8_t)(m->ax>>8),64,8);
L_14A2B: ; /* IDA 0x14A2B; jb      short loc_14A75 */
    if(m->flags&1) goto L_14A75;
L_14A2D: ; /* IDA 0x14A2D; ja      short loc_14A37 */
    if(!(m->flags&(1|0x40))) goto L_14A37;
L_14A2F: ; /* IDA 0x14A2F; mov     di, [bx+8] */
    m->di=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+8)));
L_14A32: ; /* IDA 0x14A32; mov     bp, 0FFFCh */
    m->bp=(uint16_t)(65532);
L_14A35: ; /* IDA 0x14A35; jmp     short loc_14A3D */
    goto L_14A3D;
L_14A37: ; /* IDA 0x14A37; mov     di, [bx+6] */
    m->di=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+6)));
L_14A3A: ; /* IDA 0x14A3A; mov     bp, 4 */
#if KI_ROUTE_MUTATION == 10
    m->bp=0xfffc;
#else
    m->bp=(uint16_t)(4);
#endif
L_14A3D: ; /* IDA 0x14A3D; add     dl, [bx+4] */
#if KI_ROUTE_MUTATION == 7
    vg_dl(m,ec_math(m,(uint8_t)m->dx,dl_read8(m,m->ds,(uint16_t)(m->bx+5)),8,0,0));
#else
    vg_dl(m,(uint8_t)(ec_math(m,(uint8_t)m->dx,dl_read8(m,m->ds,(uint16_t)(m->bx+4)),8,0,0)));
#endif
L_14A40: ; /* IDA 0x14A40; adc     dh, 0 */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),0,8,m->flags&1,0)));
L_14A43: ; /* IDA 0x14A43; mov     ax, [di] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di)));
L_14A45: ; /* IDA 0x14A45; and     ah, 3Fh */
    {uint16_t result=(((uint8_t)(m->ax>>8))&(63))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_14A48: ; /* IDA 0x14A48; cmp     ax, bx */
    ec_cmp(m,m->ax,m->bx,16);
L_14A4A: ; /* IDA 0x14A4A; jz      short loc_14A50 */
    if(m->flags&0x40) goto L_14A50;
L_14A4C: ; /* IDA 0x14A4C; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_14A4D: ; /* IDA 0x14A4D; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_14A4E: ; /* IDA 0x14A4E; jmp     short loc_14A43 */
    goto L_14A43;
L_14A50: ; /* IDA 0x14A50; cmp     byte ptr [di-8000h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+32768)),0,8);
L_14A55: ; /* IDA 0x14A55; jnz     short loc_14A75 */
    if(!(m->flags&0x40)) goto L_14A75;
L_14A57: ; /* IDA 0x14A57; mov     byte ptr [di-8000h], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->di+32768),1);
L_14A5C: ; /* IDA 0x14A5C; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14A5D: ; /* IDA 0x14A5D; and     di, 0FFF8h */
    {uint16_t result=((m->di)&(65528))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_14A60: ; /* IDA 0x14A60; mov     [bx], di */
    dl_write16(m,m->ds,(uint16_t)(m->bx),m->di);
L_14A62: ; /* IDA 0x14A62; mov     [bx+2], dx */
    dl_write16(m,m->ds,(uint16_t)(m->bx+2),m->dx);
L_14A65: ; /* IDA 0x14A65; mov     [bx+4], bp */
    dl_write16(m,m->ds,(uint16_t)(m->bx+4),m->bp);
L_14A68: ; /* IDA 0x14A68; mov     [bx+6], ax */
    dl_write16(m,m->ds,(uint16_t)(m->bx+6),m->ax);
L_14A6B: ; /* IDA 0x14A6B; add     bx, 8 */
    m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,0));
L_14A6E: ; /* IDA 0x14A6E; and     bh, 0FBh */
    {uint16_t result=(((uint8_t)(m->bx>>8))&(251))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_14A71: ; /* IDA 0x14A71; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_14A73: ; /* IDA 0x14A73; jmp     short loc_14A76 */
    goto L_14A76;
L_14A75: ; /* IDA 0x14A75; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14A76: ; /* IDA 0x14A76; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14A77: ; /* IDA 0x14A77; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_14A78: ; /* IDA 0x14A78; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_14A79: ; /* IDA 0x14A79; clc */
    m->flags&=(uint16_t)~1u;
L_14A7A: ; /* IDA 0x14A7A; retn */
    return;
}
/* Original sub_19656; IDA 0x19656. */
static void rt_body_19656(KiMachine16 *m,const KiEconomyHooks *h) {
L_19656: ; /* IDA 0x19656; push    ds */
    ec_push(m,m->ds);
L_19657: ; /* IDA 0x19657; push    es */
    ec_push(m,m->es);
L_19658: ; /* IDA 0x19658; push    ax */
    ec_push(m,m->ax);
L_19659: ; /* IDA 0x19659; push    bx */
    ec_push(m,m->bx);
L_1965A: ; /* IDA 0x1965A; push    cx */
    ec_push(m,m->cx);
L_1965B: ; /* IDA 0x1965B; push    dx */
    ec_push(m,m->dx);
L_1965C: ; /* IDA 0x1965C; push    si */
    ec_push(m,m->si);
L_1965D: ; /* IDA 0x1965D; push    di */
    ec_push(m,m->di);
L_1965E: ; /* IDA 0x1965E; cld */
    m->flags&=(uint16_t)~0x400u;
L_1965F: ; /* IDA 0x1965F; mov     ds, cs:word_10D3A */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3386));
L_19664: ; /* IDA 0x19664; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_19666: ; /* IDA 0x19666; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_19668: ; /* IDA 0x19668; cmp     dx, 0B2h */
    ec_cmp(m,m->dx,178,16);
L_1966C: ; /* IDA 0x1966C; jbe     short loc_19671 */
    if(m->flags&(1|0x40)) goto L_19671;
L_1966E: ; /* IDA 0x1966E; mov     dx, 0B2h */
    m->dx=(uint16_t)(178);
L_19671: ; /* IDA 0x19671; add     dx, 1B6h */
    m->dx=(uint16_t)(ec_math(m,m->dx,438,16,0,0));
L_19675: ; /* IDA 0x19675; add     bx, 26h ; '&' */
    m->bx=(uint16_t)(ec_math(m,m->bx,38,16,0,0));
L_19678: ; /* IDA 0x19678; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1967A: ; /* IDA 0x1967A; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1967C: ; /* IDA 0x1967C; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1967E: ; /* IDA 0x1967E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19680: ; /* IDA 0x19680; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19682: ; /* IDA 0x19682; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19684: ; /* IDA 0x19684; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_19686: ; /* IDA 0x19686; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19688: ; /* IDA 0x19688; add     si, bx */
    m->si=(uint16_t)(ec_math(m,m->si,m->bx,16,0,0));
L_1968A: ; /* IDA 0x1968A; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_1968C: ; /* IDA 0x1968C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1968E: ; /* IDA 0x1968E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19690: ; /* IDA 0x19690; add     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
L_19692: ; /* IDA 0x19692; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_19694: ; /* IDA 0x19694; sub     dx, 37h ; '7' */
    m->dx=(uint16_t)(ec_math(m,m->dx,55,16,0,1));
L_19697: ; /* IDA 0x19697; sub     si, 3C0h */
    m->si=(uint16_t)(ec_math(m,m->si,960,16,0,1));
L_1969B: ; /* IDA 0x1969B; add     si, dx */
    m->si=(uint16_t)(ec_math(m,m->si,m->dx,16,0,0));
L_1969D: ; /* IDA 0x1969D; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_196A0: ; /* IDA 0x196A0; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_196A2: ; /* IDA 0x196A2; call    sub_10CC3 */
    rt_call(m,3267,38565,h);
L_196A5: ; /* IDA 0x196A5; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_196A8: ; /* IDA 0x196A8; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_196AA: ; /* IDA 0x196AA; out     dx, al; EGA: sequencer address reg */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_196AB: ; /* IDA 0x196AB; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_196AD: ; /* IDA 0x196AD; call    sub_196CF */
    rt_call(m,38607,38576,h);
L_196B0: ; /* IDA 0x196B0; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_196B2: ; /* IDA 0x196B2; call    sub_196CF */
    rt_call(m,38607,38581,h);
L_196B5: ; /* IDA 0x196B5; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_196B7: ; /* IDA 0x196B7; call    sub_196CF */
    rt_call(m,38607,38586,h);
L_196BA: ; /* IDA 0x196BA; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_196BC: ; /* IDA 0x196BC; call    sub_196CF */
    rt_call(m,38607,38591,h);
L_196BF: ; /* IDA 0x196BF; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_196C2: ; /* IDA 0x196C2; mov     ax, 0F02h */
    m->ax=(uint16_t)(3842);
L_196C5: ; /* IDA 0x196C5; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_196C6: ; /* IDA 0x196C6; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_196C7: ; /* IDA 0x196C7; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_196C8: ; /* IDA 0x196C8; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_196C9: ; /* IDA 0x196C9; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_196CA: ; /* IDA 0x196CA; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_196CB: ; /* IDA 0x196CB; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_196CC: ; /* IDA 0x196CC; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_196CD: ; /* IDA 0x196CD; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_196CE: ; /* IDA 0x196CE; retn */
    return;
}
/* Original sub_196CF; IDA 0x196CF. */
static void rt_body_196CF(KiMachine16 *m,const KiEconomyHooks *h) {
L_196CF: ; /* IDA 0x196CF; mov     dx, 3C5h */
    m->dx=(uint16_t)(965);
L_196D2: ; /* IDA 0x196D2; out     dx, al; EGA port: sequencer data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_196D3: ; /* IDA 0x196D3; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_196D5: ; /* IDA 0x196D5; mov     cx, 4 */
    m->cx=(uint16_t)(4);
L_196D8: ; /* IDA 0x196D8; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_196DB: ; /* IDA 0x196DB; movsb */
    vg_movsb(m);
L_196DC: ; /* IDA 0x196DC; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_196DF: ; /* IDA 0x196DF; movsb */
    vg_movsb(m);
L_196E0: ; /* IDA 0x196E0; add     di, 4Eh ; 'N' */
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
L_196E3: ; /* IDA 0x196E3; add     si, 16h */
    m->si=(uint16_t)(ec_math(m,m->si,22,16,0,0));
L_196E6: ; /* IDA 0x196E6; loop    loc_196D8 */
    if(--m->cx) goto L_196D8;
L_196E8: ; /* IDA 0x196E8; add     si, 0BA0h */
    m->si=(uint16_t)(ec_math(m,m->si,2976,16,0,0));
L_196EC: ; /* IDA 0x196EC; retn */
    return;
}
int ki_route_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x291a:rt_body_1291A(m,h);break;
case 0x2977:rt_body_12977(m,h);break;
case 0x29c3:rt_body_129C3(m,h);break;
case 0x2ba8:rt_body_12BA8(m,h);break;
case 0x47bb:rt_body_147BB(m,h);break;
case 0x487b:rt_body_1487B(m,h);break;
case 0x491b:rt_body_1491B(m,h);break;
case 0x4a0f:rt_body_14A0F(m,h);break;
case 0x9656:rt_body_19656(m,h);break;
case 0x96cf:rt_body_196CF(m,h);break;
default:return 0;}return 1;}
static int rt_known(uint16_t t) {switch(t) {case 0x291a: case 0x2977: case 0x29c3: case 0x2ba8: case 0x47bb: case 0x487b: case 0x491b: case 0x4a0f: case 0x9656: case 0x96cf: return 1;default:return 0;}}
void sub_1291A(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x291a,h);m->ip=ec_pop(m);}
void sub_12977(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x2977,h);m->ip=ec_pop(m);}
void sub_129C3(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x29c3,h);m->ip=ec_pop(m);}
void sub_12BA8(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x2ba8,h);m->ip=ec_pop(m);}
void sub_147BB(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x47bb,h);m->ip=ec_pop(m);}
void sub_1487B(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x487b,h);m->ip=ec_pop(m);}
void raw_entry_1491B(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x491b,h);m->ip=ec_pop(m);}
void sub_14A0F(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x4a0f,h);m->ip=ec_pop(m);}
void sub_19656(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x9656,h);m->ip=ec_pop(m);}
void sub_196CF(KiMachine16 *m,const KiEconomyHooks *h) {ki_route_body(m,0x96cf,h);m->ip=ec_pop(m);}
