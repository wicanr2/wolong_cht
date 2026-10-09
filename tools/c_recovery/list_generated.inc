/* Original native C city list / sort / relocation; re/118, spec/238. */
static int ll_body_15E60(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_16909(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_17400(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_1748F(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_181C0(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_1820E(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_18412(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_18458(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_18463(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_184BC(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_184DD(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_1851A(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_18546(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_18607(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_18662(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_186C9(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_18713(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_18755(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_1743B(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_1745F(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_1828F(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_1857F(KiMachine16 *,const KiEconomyHooks *);
static int ll_body_185B2(KiMachine16 *,const KiEconomyHooks *);
/* Original sub_15E60; IDA 0x15E60. */
static int ll_body_15E60(KiMachine16 *m,const KiEconomyHooks *h) {
L_15E60: ; /* IDA 0x15E60; test    cs:byte_198A6, 2 */
    {uint16_t result=((dl_read8(m,m->cs,39078))&(2))&255;ec_logic(m,result,8);}
L_15E66: ; /* IDA 0x15E66; jnz     short loc_15E69 */
    if(!(m->flags&0x40)) goto L_15E69;
L_15E68: ; /* IDA 0x15E68; retn */
    return 0;
L_15E69: ; /* IDA 0x15E69; push    ax */
    ec_push(m,m->ax);
L_15E6A: ; /* IDA 0x15E6A; push    bx */
    ec_push(m,m->bx);
L_15E6B: ; /* IDA 0x15E6B; push    dx */
    ec_push(m,m->dx);
L_15E6C: ; /* IDA 0x15E6C; mov     al, 0 */
    vg_al(m,(uint8_t)(0));
L_15E6E: ; /* IDA 0x15E6E; mov     dx, 1B8h */
    m->dx=(uint16_t)(440);
L_15E71: ; /* IDA 0x15E71; mov     bx, 0C8h */
    m->bx=(uint16_t)(200);
L_15E74: ; /* IDA 0x15E74; call    sub_10337 */
    if(ll_call(m,823,24183,h)) return 1;
L_15E77: ; /* IDA 0x15E77; mov     al, 0Fh */
    vg_al(m,(uint8_t)(15));
L_15E79: ; /* IDA 0x15E79; call    sub_15E80 */
    if(ll_call(m,24192,24188,h)) return 1;
L_15E7C: ; /* IDA 0x15E7C; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15E7D: ; /* IDA 0x15E7D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15E7E: ; /* IDA 0x15E7E; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15E7F: ; /* IDA 0x15E7F; retn */
    return 0;
}
/* Original sub_16909; IDA 0x16909. */
static int ll_body_16909(KiMachine16 *m,const KiEconomyHooks *h) {
L_16909: ; /* IDA 0x16909; push    ds */
    ec_push(m,m->ds);
L_1690A: ; /* IDA 0x1690A; call    sub_12078 */
    if(ll_call(m,8312,26893,h)) return 1;
L_1690D: ; /* IDA 0x1690D; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16912: ; /* IDA 0x16912; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_16917: ; /* IDA 0x16917; mov     ah, [si+3] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+3))));
L_1691A: ; /* IDA 0x1691A; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1691C: ; /* IDA 0x1691C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1691E: ; /* IDA 0x1691E; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16920: ; /* IDA 0x16920; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16922: ; /* IDA 0x16922; add     ax, 840h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2112,16,0,0));
L_16925: ; /* IDA 0x16925; mov     cx, 0Fh */
    m->cx=(uint16_t)(15);
L_16928: ; /* IDA 0x16928; call    sub_18853 */
    if(ll_call(m,34899,26923,h)) return 1;
L_1692B: ; /* IDA 0x1692B; call    sub_17400 */
    if(ll_call(m,29696,26926,h)) return 1;
L_1692E: ; /* IDA 0x1692E; jb      short loc_16990 */
    if(m->flags&1) goto L_16990;
L_16930: ; /* IDA 0x16930; cmp     ax, bx */
    ec_cmp(m,m->ax,m->bx,16);
L_16932: ; /* IDA 0x16932; jnz     short loc_16944 */
    if(!(m->flags&0x40)) goto L_16944;
L_16934: ; /* IDA 0x16934; push    ax */
    ec_push(m,m->ax);
L_16935: ; /* IDA 0x16935; push    bx */
    ec_push(m,m->bx);
L_16936: ; /* IDA 0x16936; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16938: ; /* IDA 0x16938; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_1693A: ; /* IDA 0x1693A; mov     cx, 3Dh ; '=' */
    m->cx=(uint16_t)(61);
L_1693D: ; /* IDA 0x1693D; call    sub_18810 */
    if(ll_call(m,34832,26944,h)) return 1;
L_16940: ; /* IDA 0x16940; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16941: ; /* IDA 0x16941; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16942: ; /* IDA 0x16942; jmp     short loc_1692B */
    goto L_1692B;
L_16944: ; /* IDA 0x16944; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16947: ; /* IDA 0x16947; call    sub_18853 */
    if(ll_call(m,34899,26954,h)) return 1;
L_1694A: ; /* IDA 0x1694A; call    sub_11D46 */
    if(ll_call(m,7494,26957,h)) return 1;
L_1694D: ; /* IDA 0x1694D; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_1694F: ; /* IDA 0x1694F; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_16951: ; /* IDA 0x16951; mov     dl, [di+16h] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+22))));
L_16954: ; /* IDA 0x16954; mov     dh, [bx+16h] */
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+22))));
L_16957: ; /* IDA 0x16957; and     dx, 0F0Fh */
    {uint16_t result=((m->dx)&(3855))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_1695B: ; /* IDA 0x1695B; cmp     dl, dh */
    ec_cmp(m,(uint8_t)m->dx,(uint8_t)(m->dx>>8),8);
L_1695D: ; /* IDA 0x1695D; jb      short loc_16969 */
#if KI_LIST_MUTATION == 9
    if(!(m->flags&1)) goto L_16969;
#else
    if(m->flags&1) goto L_16969;
#endif
L_1695F: ; /* IDA 0x1695F; mov     dx, [di+0Eh] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+14)));
L_16962: ; /* IDA 0x16962; cmp     dx, [bx+0Eh] */
    ec_cmp(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->bx+14)),16);
L_16965: ; /* IDA 0x16965; jnb     short loc_16969 */
#if KI_LIST_MUTATION == 10
    if(m->flags&1) goto L_16969;
#else
    if(!(m->flags&1)) goto L_16969;
#endif
L_16967: ; /* IDA 0x16967; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_16969: ; /* IDA 0x16969; mov     cx, 182h */
    m->cx=(uint16_t)(386);
L_1696C: ; /* IDA 0x1696C; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_1696E: ; /* IDA 0x1696E; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_16971: ; /* IDA 0x16971; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_16973: ; /* IDA 0x16973; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16975: ; /* IDA 0x16975; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16977: ; /* IDA 0x16977; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16979: ; /* IDA 0x16979; add     bx, 4240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16960,16,0,0));
L_1697D: ; /* IDA 0x1697D; push    si */
    ec_push(m,m->si);
L_1697E: ; /* IDA 0x1697E; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16980: ; /* IDA 0x16980; mov     cx, 182h */
    m->cx=(uint16_t)(386);
L_16983: ; /* IDA 0x16983; call    sub_13B08 */
    if(ll_call(m,15112,27014,h)) return 1;
L_16986: ; /* IDA 0x16986; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16987: ; /* IDA 0x16987; cmp     al, 1 */
    ec_cmp(m,(uint8_t)m->ax,1,8);
L_16989: ; /* IDA 0x16989; jz      short loc_16990 */
    if(m->flags&0x40) goto L_16990;
L_1698B: ; /* IDA 0x1698B; mov     ax, di */
    m->ax=(uint16_t)(m->di);
L_1698D: ; /* IDA 0x1698D; call    sub_133FD */
#if KI_LIST_MUTATION == 11
    ; /* lost actual relocation writer */
#else
    if(ll_call(m,13309,27024,h)) return 1;
#endif
L_16990: ; /* IDA 0x16990; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16993: ; /* IDA 0x16993; call    sub_18853 */
    if(ll_call(m,34899,27030,h)) return 1;
L_16996: ; /* IDA 0x16996; call    sub_11D46 */
    if(ll_call(m,7494,27033,h)) return 1;
L_16999: ; /* IDA 0x16999; call    sub_120D6 */
    if(ll_call(m,8406,27036,h)) return 1;
L_1699C: ; /* IDA 0x1699C; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1699D: ; /* IDA 0x1699D; retn */
    return 0;
}
/* Original sub_17400; IDA 0x17400. */
static int ll_body_17400(KiMachine16 *m,const KiEconomyHooks *h) {
L_17400: ; /* IDA 0x17400; push    ds */
    ec_push(m,m->ds);
L_17401: ; /* IDA 0x17401; push    es */
    ec_push(m,m->es);
L_17402: ; /* IDA 0x17402; push    ax */
    ec_push(m,m->ax);
L_17403: ; /* IDA 0x17403; push    cx */
    ec_push(m,m->cx);
L_17404: ; /* IDA 0x17404; push    dx */
    ec_push(m,m->dx);
L_17405: ; /* IDA 0x17405; push    si */
    ec_push(m,m->si);
L_17406: ; /* IDA 0x17406; push    di */
    ec_push(m,m->di);
L_17407: ; /* IDA 0x17407; push    bp */
    ec_push(m,m->bp);
L_17408: ; /* IDA 0x17408; mov     dx, 18h */
    m->dx=(uint16_t)(24);
L_1740B: ; /* IDA 0x1740B; mov     bx, 58h ; 'X' */
    m->bx=(uint16_t)(88);
L_1740E: ; /* IDA 0x1740E; mov     cx, 0B18h */
    m->cx=(uint16_t)(2840);
L_17411: ; /* IDA 0x17411; call    sub_181C0 */
    if(ll_call(m,33216,29716,h)) return 1;
L_17414: ; /* IDA 0x17414; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_17419: ; /* IDA 0x17419; mov     ax, 743Bh */
    m->ax=(uint16_t)(29755);
L_1741C: ; /* IDA 0x1741C; mov     bx, 745Fh */
    m->bx=(uint16_t)(29791);
L_1741F: ; /* IDA 0x1741F; mov     si, 7378h */
    m->si=(uint16_t)(29560);
L_17422: ; /* IDA 0x17422; mov     di, 73CEh */
    m->di=(uint16_t)(29646);
L_17425: ; /* IDA 0x17425; mov     cl, byte ptr cs:word_198A8+1 */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,39081)));
L_1742A: ; /* IDA 0x1742A; call    sub_1820E */
    if(ll_call(m,33294,29741,h)) return 1;
L_1742D: ; /* IDA 0x1742D; mov     byte ptr cs:word_198A8+1, cl */
    dl_write8(m,m->cs,39081,(uint8_t)m->cx);
L_17432: ; /* IDA 0x17432; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17433: ; /* IDA 0x17433; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17434: ; /* IDA 0x17434; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17435: ; /* IDA 0x17435; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17436: ; /* IDA 0x17436; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17437: ; /* IDA 0x17437; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17438: ; /* IDA 0x17438; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_17439: ; /* IDA 0x17439; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1743A: ; /* IDA 0x1743A; retn */
    return 0;
}
/* Original sub_1748F; IDA 0x1748F. */
static int ll_body_1748F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1748F: ; /* IDA 0x1748F; push    ax */
    ec_push(m,m->ax);
L_17490: ; /* IDA 0x17490; push    bx */
    ec_push(m,m->bx);
L_17491: ; /* IDA 0x17491; push    cx */
    ec_push(m,m->cx);
L_17492: ; /* IDA 0x17492; push    dx */
    ec_push(m,m->dx);
L_17493: ; /* IDA 0x17493; push    si */
    ec_push(m,m->si);
L_17494: ; /* IDA 0x17494; push    di */
    ec_push(m,m->di);
L_17495: ; /* IDA 0x17495; push    bp */
    ec_push(m,m->bp);
L_17496: ; /* IDA 0x17496; cmp     word ptr [bp+si], 0FFFFh */
#if KI_LIST_MUTATION == 2
    ec_cmp(m,0xffff,0xffff,16);
#else
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)),65535,16);
#endif
L_1749A: ; /* IDA 0x1749A; jnz     short loc_174AE */
    if(!(m->flags&0x40)) goto L_174AE;
L_1749C: ; /* IDA 0x1749C; push    ds */
    ec_push(m,m->ds);
L_1749D: ; /* IDA 0x1749D; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_1749F: ; /* IDA 0x1749F; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_174A1: ; /* IDA 0x174A1; mov     si, 73A3h */
    m->si=(uint16_t)(29603);
L_174A4: ; /* IDA 0x174A4; mov     ax, 9000h */
    m->ax=(uint16_t)(36864);
L_174A7: ; /* IDA 0x174A7; call    sub_1F6DC */
    if(ll_call(m,63196,29866,h)) return 1;
L_174AA: ; /* IDA 0x174AA; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_174AB: ; /* IDA 0x174AB; jmp     loc_17548 */
    goto L_17548;
L_174AE: ; /* IDA 0x174AE; mov     si, [bp+si] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)));
L_174B0: ; /* IDA 0x174B0; add     si, 2 */
    m->si=(uint16_t)(ec_math(m,m->si,2,16,0,0));
L_174B3: ; /* IDA 0x174B3; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_174B6: ; /* IDA 0x174B6; call    loc_10701 */
    if(ll_call(m,1793,29881,h)) return 1;
L_174B9: ; /* IDA 0x174B9; sub     si, 2 */
    m->si=(uint16_t)(ec_math(m,m->si,2,16,0,1));
L_174BC: ; /* IDA 0x174BC; push    ds */
    ec_push(m,m->ds);
L_174BD: ; /* IDA 0x174BD; push    si */
    ec_push(m,m->si);
L_174BE: ; /* IDA 0x174BE; mov     ah, [si+19h] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+25))));
L_174C1: ; /* IDA 0x174C1; cmp     ah, 0FFh */
    ec_cmp(m,(uint8_t)(m->ax>>8),255,8);
L_174C4: ; /* IDA 0x174C4; jz      short loc_174D5 */
    if(m->flags&0x40) goto L_174D5;
L_174C6: ; /* IDA 0x174C6; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_174C8: ; /* IDA 0x174C8; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_174CA: ; /* IDA 0x174CA; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_174CC: ; /* IDA 0x174CC; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_174CE: ; /* IDA 0x174CE; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_174D1: ; /* IDA 0x174D1; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_174D3: ; /* IDA 0x174D3; jmp     short loc_174DC */
    goto L_174DC;
L_174D5: ; /* IDA 0x174D5; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_174D7: ; /* IDA 0x174D7; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_174D9: ; /* IDA 0x174D9; mov     si, 73A3h */
    m->si=(uint16_t)(29603);
L_174DC: ; /* IDA 0x174DC; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_174DF: ; /* IDA 0x174DF; add     dx, 120h */
    m->dx=(uint16_t)(ec_math(m,m->dx,288,16,0,0));
L_174E3: ; /* IDA 0x174E3; call    loc_10701 */
    if(ll_call(m,1793,29926,h)) return 1;
L_174E6: ; /* IDA 0x174E6; sub     dx, 120h */
    m->dx=(uint16_t)(ec_math(m,m->dx,288,16,0,1));
L_174EA: ; /* IDA 0x174EA; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_174EB: ; /* IDA 0x174EB; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_174EC: ; /* IDA 0x174EC; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_174EE: ; /* IDA 0x174EE; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_174F0: ; /* IDA 0x174F0; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_174F2: ; /* IDA 0x174F2; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_174F4: ; /* IDA 0x174F4; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_174F6: ; /* IDA 0x174F6; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_174F8: ; /* IDA 0x174F8; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_174FA: ; /* IDA 0x174FA; add     di, bx */
    m->di=(uint16_t)(ec_math(m,m->di,m->bx,16,0,0));
L_174FC: ; /* IDA 0x174FC; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_174FE: ; /* IDA 0x174FE; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_17500: ; /* IDA 0x17500; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_17502: ; /* IDA 0x17502; add     di, dx */
    m->di=(uint16_t)(ec_math(m,m->di,m->dx,16,0,0));
L_17504: ; /* IDA 0x17504; add     di, 8 */
    m->di=(uint16_t)(ec_math(m,m->di,8,16,0,0));
L_17507: ; /* IDA 0x17507; mov     ax, [si+0Eh] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
L_1750A: ; /* IDA 0x1750A; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_1750C: ; /* IDA 0x1750C; mov     bx, 9005h */
    m->bx=(uint16_t)(36869);
L_1750F: ; /* IDA 0x1750F; call    sub_1062F */
    if(ll_call(m,1583,29970,h)) return 1;
L_17512: ; /* IDA 0x17512; add     di, 9 */
    m->di=(uint16_t)(ec_math(m,m->di,9,16,0,0));
L_17515: ; /* IDA 0x17515; mov     al, [si+10h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+16))));
L_17518: ; /* IDA 0x17518; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1751A: ; /* IDA 0x1751A; mov     bx, 9004h */
    m->bx=(uint16_t)(36868);
L_1751D: ; /* IDA 0x1751D; sub     ax, 64h ; 'd' */
    m->ax=(uint16_t)(ec_math(m,m->ax,100,16,0,1));
L_17520: ; /* IDA 0x17520; ja      short loc_17524 */
    if(!(m->flags&(1|0x40))) goto L_17524;
L_17522: ; /* IDA 0x17522; mov     bh, 9Ah */
    gw_bh(m,(uint8_t)(154));
L_17524: ; /* IDA 0x17524; cwd */
    m->dx=m->ax&0x8000?0xffff:0;
L_17525: ; /* IDA 0x17525; call    sub_1062F */
    if(ll_call(m,1583,29992,h)) return 1;
L_17528: ; /* IDA 0x17528; add     di, 7 */
    m->di=(uint16_t)(ec_math(m,m->di,7,16,0,0));
L_1752B: ; /* IDA 0x1752B; mov     al, [si+11h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+17))));
L_1752E: ; /* IDA 0x1752E; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17530: ; /* IDA 0x17530; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17532: ; /* IDA 0x17532; mov     bx, 9003h */
    m->bx=(uint16_t)(36867);
L_17535: ; /* IDA 0x17535; call    sub_1062F */
    if(ll_call(m,1583,30008,h)) return 1;
L_17538: ; /* IDA 0x17538; add     di, 6 */
    m->di=(uint16_t)(ec_math(m,m->di,6,16,0,0));
L_1753B: ; /* IDA 0x1753B; mov     al, 0Ah */
    vg_al(m,(uint8_t)(10));
L_1753D: ; /* IDA 0x1753D; mul     byte ptr [si+13h] */
    dl_mul(m,dl_read8(m,m->ds,(uint16_t)(m->si+19)),8);
L_17540: ; /* IDA 0x17540; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17542: ; /* IDA 0x17542; mov     bx, 9004h */
    m->bx=(uint16_t)(36868);
L_17545: ; /* IDA 0x17545; call    sub_1062F */
    if(ll_call(m,1583,30024,h)) return 1;
L_17548: ; /* IDA 0x17548; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17549: ; /* IDA 0x17549; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1754A: ; /* IDA 0x1754A; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1754B: ; /* IDA 0x1754B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1754C: ; /* IDA 0x1754C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1754D: ; /* IDA 0x1754D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1754E: ; /* IDA 0x1754E; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1754F: ; /* IDA 0x1754F; retn */
    return 0;
}
/* Original sub_181C0; IDA 0x181C0. */
static int ll_body_181C0(KiMachine16 *m,const KiEconomyHooks *h) {
L_181C0: ; /* IDA 0x181C0; push    ds */
    ec_push(m,m->ds);
L_181C1: ; /* IDA 0x181C1; push    ax */
    ec_push(m,m->ax);
L_181C2: ; /* IDA 0x181C2; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_181C4: ; /* IDA 0x181C4; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_181C6: ; /* IDA 0x181C6; mov     word_181B6, cx */
    dl_write16(m,m->ds,33206,m->cx);
L_181CA: ; /* IDA 0x181CA; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_181CD: ; /* IDA 0x181CD; xchg    al, cl */
    {uint16_t a=(uint8_t)m->ax,b=(uint8_t)m->cx;vg_al(m,(uint8_t)(b));vg_cl(m,(uint8_t)(a));}
L_181CF: ; /* IDA 0x181CF; shl     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,4));
L_181D1: ; /* IDA 0x181D1; mov     word_181B2, ax */
    dl_write16(m,m->ds,33202,m->ax);
L_181D4: ; /* IDA 0x181D4; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_181D6: ; /* IDA 0x181D6; mov     al, ch */
    vg_al(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_181D8: ; /* IDA 0x181D8; shl     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,4));
L_181DA: ; /* IDA 0x181DA; mov     word_181B4, ax */
    dl_write16(m,m->ds,33204,m->ax);
L_181DD: ; /* IDA 0x181DD; push    dx */
    ec_push(m,m->dx);
L_181DE: ; /* IDA 0x181DE; push    bx */
    ec_push(m,m->bx);
L_181DF: ; /* IDA 0x181DF; mov     word_181AE, dx */
    dl_write16(m,m->ds,33198,m->dx);
L_181E3: ; /* IDA 0x181E3; mov     word_181B0, bx */
    dl_write16(m,m->ds,33200,m->bx);
L_181E7: ; /* IDA 0x181E7; sub     dx, 8 */
    m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,1));
L_181EA: ; /* IDA 0x181EA; sub     bx, 8 */
    m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,1));
L_181ED: ; /* IDA 0x181ED; shr     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
L_181EF: ; /* IDA 0x181EF; shr     bx, cl */
    m->bx=(uint16_t)(dl_shift(m,m->bx,(uint8_t)m->cx,16,5));
L_181F1: ; /* IDA 0x181F1; mov     word_181AA, dx */
    dl_write16(m,m->ds,33194,m->dx);
L_181F5: ; /* IDA 0x181F5; mov     word_181AC, bx */
    dl_write16(m,m->ds,33196,m->bx);
L_181F9: ; /* IDA 0x181F9; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_181FA: ; /* IDA 0x181FA; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_181FB: ; /* IDA 0x181FB; mov     dx, 50h ; 'P' */
    m->dx=(uint16_t)(80);
L_181FE: ; /* IDA 0x181FE; mul     dx */
    dl_mul(m,m->dx,16);
L_18200: ; /* IDA 0x18200; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_18202: ; /* IDA 0x18202; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_18204: ; /* IDA 0x18204; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_18206: ; /* IDA 0x18206; add     ax, cx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->cx,16,0,0));
L_18208: ; /* IDA 0x18208; mov     word_181B8, ax */
    dl_write16(m,m->ds,33208,m->ax);
L_1820B: ; /* IDA 0x1820B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1820C: ; /* IDA 0x1820C; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1820D: ; /* IDA 0x1820D; retn */
    return 0;
}
/* Original sub_1820E; IDA 0x1820E. */
static int ll_body_1820E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1820E: ; /* IDA 0x1820E; push    bp */
    ec_push(m,m->bp);
L_1820F: ; /* IDA 0x1820F; sub     sp, 200h */
    m->sp=(uint16_t)(ec_math(m,m->sp,512,16,0,1));
L_18213: ; /* IDA 0x18213; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_18215: ; /* IDA 0x18215; mov     cs:word_181A8, ax */
    dl_write16(m,m->cs,33192,m->ax);
L_18219: ; /* IDA 0x18219; mov     cs:word_181A6, bx */
    dl_write16(m,m->cs,33190,m->bx);
L_1821E: ; /* IDA 0x1821E; mov     word ptr cs:loc_183D2+1, si */
    dl_write16(m,m->cs,33747,m->si);
L_18223: ; /* IDA 0x18223; mov     word ptr cs:loc_1831F+1, di */
    dl_write16(m,m->cs,33568,m->di);
L_18228: ; /* IDA 0x18228; mov     word ptr cs:loc_1858F+2, di */
    dl_write16(m,m->cs,34193,m->di);
L_1822D: ; /* IDA 0x1822D; mov     cs:byte_181BC, 0 */
    dl_write8(m,m->cs,33212,0);
L_18233: ; /* IDA 0x18233; xor     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_18235: ; /* IDA 0x18235; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_18237: ; /* IDA 0x18237; mov     cs:word_181BA, cx */
    dl_write16(m,m->cs,33210,m->cx);
L_1823C: ; /* IDA 0x1823C; cld */
    m->flags&=(uint16_t)~0x400u;
L_1823D: ; /* IDA 0x1823D; mov     ax, ss */
    m->ax=(uint16_t)(m->ss);
L_1823F: ; /* IDA 0x1823F; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_18241: ; /* IDA 0x18241; mov     di, bp */
    m->di=(uint16_t)(m->bp);
L_18243: ; /* IDA 0x18243; mov     cx, 100h */
    m->cx=(uint16_t)(256);
L_18246: ; /* IDA 0x18246; mov     ax, 0FFFFh */
    m->ax=(uint16_t)(65535);
L_18249: ; /* IDA 0x18249; rep stosw */
    mc_string(m,2,1);
L_1824B: ; /* IDA 0x1824B; call    sub_101B4 */
    if(ll_call(m,436,33358,h)) return 1;
L_1824E: ; /* IDA 0x1824E; call    cs:word_181A8 */
    if(ll_call(m,dl_read16(m,m->cs,33192),33363,h)) return 1;
L_18253: ; /* IDA 0x18253; mov     cs:byte_181BD, ah */
    dl_write8(m,m->cs,33213,(uint8_t)(m->ax>>8));
L_18258: ; /* IDA 0x18258; mov     al, 3Dh ; '=' */
    vg_al(m,(uint8_t)(61));
L_1825A: ; /* IDA 0x1825A; call    loc_1828F */
    if(ll_call(m,33423,33373,h)) return 1;
L_1825D: ; /* IDA 0x1825D; mov     bx, cs:word_181BA */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33210));
L_18262: ; /* IDA 0x18262; call    loc_1857F */
    if(ll_call(m,34175,33381,h)) return 1;
L_18265: ; /* IDA 0x18265; call    sub_101DB */
    if(ll_call(m,475,33384,h)) return 1;
L_18268: ; /* IDA 0x18268; call    sub_18412 */
    if(ll_call(m,33810,33387,h)) return 1;
L_1826B: ; /* IDA 0x1826B; jb      short loc_18279 */
    if(m->flags&1) goto L_18279;
L_1826D: ; /* IDA 0x1826D; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_1826F: ; /* IDA 0x1826F; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_18271: ; /* IDA 0x18271; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_18273: ; /* IDA 0x18273; add     bp, bx */
    m->bp=(uint16_t)(ec_math(m,m->bp,m->bx,16,0,0));
L_18275: ; /* IDA 0x18275; mov     bx, [bp+0] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_18278: ; /* IDA 0x18278; clc */
    m->flags&=(uint16_t)~1u;
L_18279: ; /* IDA 0x18279; pushf */
    ec_push(m,m->flags);
L_1827A: ; /* IDA 0x1827A; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_1827C: ; /* IDA 0x1827C; call    loc_1828F */
    if(ll_call(m,33423,33407,h)) return 1;
L_1827F: ; /* IDA 0x1827F; mov     cx, cs:word_181BA */
    m->cx=(uint16_t)(dl_read16(m,m->cs,33210));
L_18284: ; /* IDA 0x18284; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_18286: ; /* IDA 0x18286; popf */
    m->flags=ec_pop(m)|2;
L_18287: ; /* IDA 0x18287; lahf */
    vg_ah(m,(uint8_t)((m->flags&0xd5)|2));
L_18288: ; /* IDA 0x18288; add     sp, 200h */
    m->sp=(uint16_t)(ec_math(m,m->sp,512,16,0,0));
L_1828C: ; /* IDA 0x1828C; sahf */
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
L_1828D: ; /* IDA 0x1828D; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1828E: ; /* IDA 0x1828E; retn */
    return 0;
}
/* Original sub_18412; IDA 0x18412. */
static int ll_body_18412(KiMachine16 *m,const KiEconomyHooks *h) {
L_18412: ; /* IDA 0x18412; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_18415: ; /* IDA 0x18415; call    sub_20000 */
    mx_far(m,0x8415,0x841a,h);
L_1841A: ; /* IDA 0x1841A; mov     al, cs:byte_181BC */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,33212)));
L_1841E: ; /* IDA 0x1841E; call    cs:word_181A6 */
    if(ll_call(m,dl_read16(m,m->cs,33190),33827,h)) return 1;
L_18423: ; /* IDA 0x18423; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_18426: ; /* IDA 0x18426; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_18428: ; /* IDA 0x18428; call    sub_20000 */
    mx_far(m,0x8428,0x842d,h);
L_1842D: ; /* IDA 0x1842D; call    sub_121E7 */
    if(ll_call(m,8679,33840,h)) return 1;
L_18430: ; /* IDA 0x18430; jb      short loc_18457 */
    if(m->flags&1) goto L_18457;
L_18432: ; /* IDA 0x18432; push    cx */
    ec_push(m,m->cx);
L_18433: ; /* IDA 0x18433; push    dx */
    ec_push(m,m->dx);
L_18434: ; /* IDA 0x18434; call    sub_1E453 */
    if(ll_call(m,58451,33847,h)) return 1;
L_18437: ; /* IDA 0x18437; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_18438: ; /* IDA 0x18438; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_18439: ; /* IDA 0x18439; sub     al, 3Dh ; '=' */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,61,8,0,1)));
L_1843B: ; /* IDA 0x1843B; jb      short loc_18455 */
    if(m->flags&1) goto L_18455;
L_1843D: ; /* IDA 0x1843D; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1843F: ; /* IDA 0x1843F; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18441: ; /* IDA 0x18441; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_18443: ; /* IDA 0x18443; cmp     bx, 0Ah */
    ec_cmp(m,m->bx,10,16);
L_18446: ; /* IDA 0x18446; jb      short loc_18450 */
    if(m->flags&1) goto L_18450;
L_18448: ; /* IDA 0x18448; sub     bx, 0Ah */
    m->bx=(uint16_t)(ec_math(m,m->bx,10,16,0,1));
L_1844B: ; /* IDA 0x1844B; call    loc_1857F */
    if(ll_call(m,34175,33870,h)) return 1;
L_1844E: ; /* IDA 0x1844E; jmp     short loc_18455 */
    goto L_18455;
L_18450: ; /* IDA 0x18450; call    cs:funcs_18450[bx] */
    if(ll_call(m,dl_read16(m,m->cs,(uint16_t)(m->bx+33881)),33877,h)) return 1;
L_18455: ; /* IDA 0x18455; jmp     short sub_18412 */
    m->ip=0x8412;if(h&&h->enter) h->enter(m,0x8412,h->user);goto L_18412;
L_18457: ; /* IDA 0x18457; stc */
    m->flags|=1;
m->ip=0x8458;if(h&&h->enter) h->enter(m,0x8458,h->user);return ll_body_18458(m,h);
}
/* Original nullsub_2; IDA 0x18458. */
static int ll_body_18458(KiMachine16 *m,const KiEconomyHooks *h) {
L_18458: ; /* IDA 0x18458; retn */
    return 0;
}
/* Original sub_18463; IDA 0x18463. */
static int ll_body_18463(KiMachine16 *m,const KiEconomyHooks *h) {
L_18463: ; /* IDA 0x18463; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_18465: ; /* IDA 0x18465; sub     dx, cs:word_181B0 */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->cs,33200),16,0,1));
L_1846A: ; /* IDA 0x1846A; sub     dx, 10h */
    m->dx=(uint16_t)(ec_math(m,m->dx,16,16,0,1));
L_1846D: ; /* IDA 0x1846D; mov     cl, 4 */
    vg_cl(m,(uint8_t)(4));
L_1846F: ; /* IDA 0x1846F; shr     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
L_18471: ; /* IDA 0x18471; add     dl, cs:byte_181BC */
    vg_dl(m,(uint8_t)(ec_math(m,(uint8_t)m->dx,dl_read8(m,m->cs,33212),8,0,0)));
L_18476: ; /* IDA 0x18476; cmp     dl, cs:byte_181BD */
    ec_cmp(m,(uint8_t)m->dx,dl_read8(m,m->cs,33213),8);
L_1847B: ; /* IDA 0x1847B; jb      short loc_1848C */
    if(m->flags&1) goto L_1848C;
L_1847D: ; /* IDA 0x1847D; retn */
    return 0;
L_1847E: ; /* IDA 0x1847E; call    sub_184BC */
    if(ll_call(m,33980,33921,h)) return 1;
L_18481: ; /* IDA 0x18481; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_18484: ; /* IDA 0x18484; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_18486: ; /* IDA 0x18486; call    sub_20000 */
    mx_far(m,0x8486,0x848b,h);
L_1848B: ; /* IDA 0x1848B; retn */
    return 0;
L_1848C: ; /* IDA 0x1848C; push    ax */
    ec_push(m,m->ax);
L_1848D: ; /* IDA 0x1848D; push    bx */
    ec_push(m,m->bx);
L_1848E: ; /* IDA 0x1848E; push    dx */
    ec_push(m,m->dx);
L_1848F: ; /* IDA 0x1848F; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_18492: ; /* IDA 0x18492; call    sub_20000 */
    mx_far(m,0x8492,0x8497,h);
L_18497: ; /* IDA 0x18497; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_18498: ; /* IDA 0x18498; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_18499: ; /* IDA 0x18499; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1849A: ; /* IDA 0x1849A; sub     di, 8 */
    m->di=(uint16_t)(ec_math(m,m->di,8,16,0,1));
L_1849D: ; /* IDA 0x1849D; and     di, 0FFF0h */
    {uint16_t result=((m->di)&(65520))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_184A0: ; /* IDA 0x184A0; add     di, 8 */
    m->di=(uint16_t)(ec_math(m,m->di,8,16,0,0));
L_184A3: ; /* IDA 0x184A3; call    sub_184BC */
    if(ll_call(m,33980,33958,h)) return 1;
L_184A6: ; /* IDA 0x184A6; mov     si, dx */
    m->si=(uint16_t)(m->dx);
L_184A8: ; /* IDA 0x184A8; call    sub_121E7 */
    if(ll_call(m,8679,33963,h)) return 1;
L_184AB: ; /* IDA 0x184AB; jb      short loc_1847E */
    if(m->flags&1) goto L_1847E;
L_184AD: ; /* IDA 0x184AD; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_184B0: ; /* IDA 0x184B0; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_184B2: ; /* IDA 0x184B2; call    sub_20000 */
    mx_far(m,0x84b2,0x84b7,h);
L_184B7: ; /* IDA 0x184B7; pop     ax */
#if KI_LIST_MUTATION == 6
    ; /* lost removal of selection-loop return address */
#else
    m->ax=(uint16_t)(ec_pop(m));
#endif
L_184B8: ; /* IDA 0x184B8; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_184BA: ; /* IDA 0x184BA; clc */
    m->flags&=(uint16_t)~1u;
L_184BB: ; /* IDA 0x184BB; retn */
    return 0;
}
/* Original sub_184BC; IDA 0x184BC. */
static int ll_body_184BC(KiMachine16 *m,const KiEconomyHooks *h) {
L_184BC: ; /* IDA 0x184BC; push    bx */
    ec_push(m,m->bx);
L_184BD: ; /* IDA 0x184BD; push    dx */
    ec_push(m,m->dx);
L_184BE: ; /* IDA 0x184BE; push    si */
    ec_push(m,m->si);
L_184BF: ; /* IDA 0x184BF; push    di */
    ec_push(m,m->di);
L_184C0: ; /* IDA 0x184C0; mov     bx, di */
    m->bx=(uint16_t)(m->di);
L_184C2: ; /* IDA 0x184C2; mov     dx, cs:word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33198));
L_184C7: ; /* IDA 0x184C7; add     dx, 10h */
    m->dx=(uint16_t)(ec_math(m,m->dx,16,16,0,0));
L_184CA: ; /* IDA 0x184CA; mov     si, cs:word_181B2 */
    m->si=(uint16_t)(dl_read16(m,m->cs,33202));
L_184CF: ; /* IDA 0x184CF; sub     si, 10h */
    m->si=(uint16_t)(ec_math(m,m->si,16,16,0,1));
L_184D2: ; /* IDA 0x184D2; mov     di, 10h */
#if KI_LIST_MUTATION == 12
    m->di=(uint16_t)(15);
#else
    m->di=(uint16_t)(16);
#endif
L_184D5: ; /* IDA 0x184D5; call    sub_10B46 */
    if(ll_call(m,2886,34008,h)) return 1;
L_184D8: ; /* IDA 0x184D8; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_184D9: ; /* IDA 0x184D9; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_184DA: ; /* IDA 0x184DA; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_184DB: ; /* IDA 0x184DB; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_184DC: ; /* IDA 0x184DC; retn */
    return 0;
}
/* Original sub_184DD; IDA 0x184DD. */
static int ll_body_184DD(KiMachine16 *m,const KiEconomyHooks *h) {
L_184DD: ; /* IDA 0x184DD; push    si */
    ec_push(m,m->si);
L_184DE: ; /* IDA 0x184DE; mov     si, 0FFFFh */
    m->si=(uint16_t)(65535);
L_184E1: ; /* IDA 0x184E1; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_184E4: ; /* IDA 0x184E4; call    sub_20000 */
    mx_far(m,0x84e4,0x84e9,h);
L_184E9: ; /* IDA 0x184E9; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_184EC: ; /* IDA 0x184EC; call    sub_20000 */
    mx_far(m,0x84ec,0x84f1,h);
L_184F1: ; /* IDA 0x184F1; mov     ax, dx */
    m->ax=(uint16_t)(m->dx);
L_184F3: ; /* IDA 0x184F3; call    sub_18713 */
    if(ll_call(m,34579,34038,h)) return 1;
L_184F6: ; /* IDA 0x184F6; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_184F8: ; /* IDA 0x184F8; cmp     ax, si */
    ec_cmp(m,m->ax,m->si,16);
L_184FA: ; /* IDA 0x184FA; jz      short loc_18503 */
    if(m->flags&0x40) goto L_18503;
L_184FC: ; /* IDA 0x184FC; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_184FE: ; /* IDA 0x184FE; call    cs:word_181A6 */
    if(ll_call(m,dl_read16(m,m->cs,33190),34051,h)) return 1;
L_18503: ; /* IDA 0x18503; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_18505: ; /* IDA 0x18505; jnz     short loc_184E9 */
    if(!(m->flags&0x40)) goto L_184E9;
L_18507: ; /* IDA 0x18507; mov     dl, al */
    vg_dl(m,(uint8_t)((uint8_t)m->ax));
L_18509: ; /* IDA 0x18509; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_1850C: ; /* IDA 0x1850C; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_1850E: ; /* IDA 0x1850E; call    sub_20000 */
    mx_far(m,0x850e,0x8513,h);
L_18513: ; /* IDA 0x18513; mov     cs:byte_181BC, dl */
    dl_write8(m,m->cs,33212,(uint8_t)m->dx);
L_18518: ; /* IDA 0x18518; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_18519: ; /* IDA 0x18519; retn */
    return 0;
}
/* Original sub_1851A; IDA 0x1851A. */
static int ll_body_1851A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1851A: ; /* IDA 0x1851A; mov     dl, cs:byte_181BC */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,33212)));
L_1851F: ; /* IDA 0x1851F; push    dx */
    ec_push(m,m->dx);
L_18520: ; /* IDA 0x18520; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_18523: ; /* IDA 0x18523; call    sub_20000 */
    mx_far(m,0x8523,0x8528,h);
L_18528: ; /* IDA 0x18528; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_18529: ; /* IDA 0x18529; and     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)&((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_1852B: ; /* IDA 0x1852B; jz      short loc_1852F */
    if(m->flags&0x40) goto L_1852F;
L_1852D: ; /* IDA 0x1852D; dec     dl */
    vg_dl(m,(uint8_t)(st_dec(m,(uint8_t)m->dx,8)));
L_1852F: ; /* IDA 0x1852F; mov     al, dl */
    vg_al(m,(uint8_t)((uint8_t)m->dx));
L_18531: ; /* IDA 0x18531; call    cs:word_181A6 */
    if(ll_call(m,dl_read16(m,m->cs,33190),34102,h)) return 1;
L_18536: ; /* IDA 0x18536; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_18539: ; /* IDA 0x18539; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_1853B: ; /* IDA 0x1853B; call    sub_20000 */
    mx_far(m,0x853b,0x8540,h);
L_18540: ; /* IDA 0x18540; mov     cs:byte_181BC, dl */
    dl_write8(m,m->cs,33212,(uint8_t)m->dx);
L_18545: ; /* IDA 0x18545; retn */
    return 0;
}
/* Original sub_18546; IDA 0x18546. */
static int ll_body_18546(KiMachine16 *m,const KiEconomyHooks *h) {
L_18546: ; /* IDA 0x18546; mov     dl, cs:byte_181BC */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,33212)));
L_1854B: ; /* IDA 0x1854B; push    dx */
    ec_push(m,m->dx);
L_1854C: ; /* IDA 0x1854C; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_1854F: ; /* IDA 0x1854F; call    sub_20000 */
    mx_far(m,0x854f,0x8554,h);
L_18554: ; /* IDA 0x18554; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_18555: ; /* IDA 0x18555; mov     al, byte ptr cs:word_181B6+1 */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,33207)));
L_18559: ; /* IDA 0x18559; dec     al */
    vg_al(m,(uint8_t)(st_dec(m,(uint8_t)m->ax,8)));
L_1855B: ; /* IDA 0x1855B; mov     dh, cs:byte_181BE */
    vg_dh(m,(uint8_t)(dl_read8(m,m->cs,33214)));
L_18560: ; /* IDA 0x18560; sub     dh, al */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),(uint8_t)m->ax,8,0,1)));
L_18562: ; /* IDA 0x18562; cmp     dl, dh */
    ec_cmp(m,(uint8_t)m->dx,(uint8_t)(m->dx>>8),8);
L_18564: ; /* IDA 0x18564; jnb     short loc_18568 */
    if(!(m->flags&1)) goto L_18568;
L_18566: ; /* IDA 0x18566; inc     dl */
#if KI_LIST_MUTATION == 7
    ; /* lost scroll increment */
#else
    vg_dl(m,(uint8_t)(ec_inc(m,(uint8_t)m->dx,8)));
#endif
L_18568: ; /* IDA 0x18568; mov     al, dl */
    vg_al(m,(uint8_t)((uint8_t)m->dx));
L_1856A: ; /* IDA 0x1856A; call    cs:word_181A6 */
    if(ll_call(m,dl_read16(m,m->cs,33190),34159,h)) return 1;
L_1856F: ; /* IDA 0x1856F; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_18572: ; /* IDA 0x18572; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_18574: ; /* IDA 0x18574; call    sub_20000 */
    mx_far(m,0x8574,0x8579,h);
L_18579: ; /* IDA 0x18579; mov     cs:byte_181BC, dl */
    dl_write8(m,m->cs,33212,(uint8_t)m->dx);
L_1857E: ; /* IDA 0x1857E; retn */
    return 0;
}
/* Original sub_18607; IDA 0x18607. */
static int ll_body_18607(KiMachine16 *m,const KiEconomyHooks *h) {
L_18607: ; /* IDA 0x18607; push    ax */
    ec_push(m,m->ax);
L_18608: ; /* IDA 0x18608; push    bx */
    ec_push(m,m->bx);
L_18609: ; /* IDA 0x18609; push    dx */
    ec_push(m,m->dx);
L_1860A: ; /* IDA 0x1860A; push    si */
    ec_push(m,m->si);
L_1860B: ; /* IDA 0x1860B; push    di */
    ec_push(m,m->di);
L_1860C: ; /* IDA 0x1860C; mov     di, cs:word_181B0 */
    m->di=(uint16_t)(dl_read16(m,m->cs,33200));
L_18611: ; /* IDA 0x18611; add     di, 20h ; ' ' */
    m->di=(uint16_t)(ec_math(m,m->di,32,16,0,0));
L_18614: ; /* IDA 0x18614; add     di, cs:word_181A4 */
    m->di=(uint16_t)(ec_math(m,m->di,dl_read16(m,m->cs,33188),16,0,0));
L_18619: ; /* IDA 0x18619; jb      short loc_18639 */
    if(m->flags&1) goto L_18639;
L_1861B: ; /* IDA 0x1861B; push    ax */
    ec_push(m,m->ax);
L_1861C: ; /* IDA 0x1861C; mov     bl, cs:byte_181BF */
    gw_bl(m,(uint8_t)(dl_read8(m,m->cs,33215)));
L_18621: ; /* IDA 0x18621; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_18623: ; /* IDA 0x18623; dec     bx */
    m->bx=(uint16_t)(st_dec(m,m->bx,16));
L_18624: ; /* IDA 0x18624; add     bx, di */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->di,16,0,0));
L_18626: ; /* IDA 0x18626; mov     dx, cs:word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33198));
L_1862B: ; /* IDA 0x1862B; mov     si, cs:word_181AE */
    m->si=(uint16_t)(dl_read16(m,m->cs,33198));
L_18630: ; /* IDA 0x18630; add     si, 0Fh */
    m->si=(uint16_t)(ec_math(m,m->si,15,16,0,0));
L_18633: ; /* IDA 0x18633; mov     ah, 0 */
    vg_ah(m,(uint8_t)(0));
L_18635: ; /* IDA 0x18635; call    sub_1F1A3 */
    if(ll_call(m,61859,34360,h)) return 1;
L_18638: ; /* IDA 0x18638; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_18639: ; /* IDA 0x18639; mov     dx, cs:word_181B4 */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33204));
L_1863E: ; /* IDA 0x1863E; sub     dx, 30h ; '0' */
    m->dx=(uint16_t)(ec_math(m,m->dx,48,16,0,1));
L_18641: ; /* IDA 0x18641; mul     dl */
    dl_mul(m,(uint8_t)m->dx,8);
L_18643: ; /* IDA 0x18643; div     cs:byte_181BE */
    dl_div(m,dl_read8(m,m->cs,33214),8);
L_18648: ; /* IDA 0x18648; and     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))&((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1864A: ; /* IDA 0x1864A; jz      short loc_1864E */
    if(m->flags&0x40) goto L_1864E;
L_1864C: ; /* IDA 0x1864C; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_1864E: ; /* IDA 0x1864E; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_18650: ; /* IDA 0x18650; mov     cs:word_181A4, ax */
    dl_write16(m,m->cs,33188,m->ax);
L_18654: ; /* IDA 0x18654; mov     cl, cs:byte_181BF */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,33215)));
L_18659: ; /* IDA 0x18659; call    sub_18662 */
    if(ll_call(m,34402,34396,h)) return 1;
L_1865C: ; /* IDA 0x1865C; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1865D: ; /* IDA 0x1865D; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1865E: ; /* IDA 0x1865E; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1865F: ; /* IDA 0x1865F; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_18660: ; /* IDA 0x18660; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_18661: ; /* IDA 0x18661; retn */
    return 0;
}
/* Original sub_18662; IDA 0x18662. */
static int ll_body_18662(KiMachine16 *m,const KiEconomyHooks *h) {
L_18662: ; /* IDA 0x18662; push    ds */
    ec_push(m,m->ds);
L_18663: ; /* IDA 0x18663; push    es */
    ec_push(m,m->es);
L_18664: ; /* IDA 0x18664; push    ax */
    ec_push(m,m->ax);
L_18665: ; /* IDA 0x18665; push    bx */
    ec_push(m,m->bx);
L_18666: ; /* IDA 0x18666; push    cx */
    ec_push(m,m->cx);
L_18667: ; /* IDA 0x18667; push    dx */
    ec_push(m,m->dx);
L_18668: ; /* IDA 0x18668; push    si */
    ec_push(m,m->si);
L_18669: ; /* IDA 0x18669; push    di */
    ec_push(m,m->di);
L_1866A: ; /* IDA 0x1866A; cld */
    m->flags&=(uint16_t)~0x400u;
L_1866B: ; /* IDA 0x1866B; add     ax, cs:word_181B0 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,33200),16,0,0));
L_18670: ; /* IDA 0x18670; add     ax, 20h ; ' ' */
    m->ax=(uint16_t)(ec_math(m,m->ax,32,16,0,0));
L_18673: ; /* IDA 0x18673; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18675: ; /* IDA 0x18675; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18677: ; /* IDA 0x18677; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18679: ; /* IDA 0x18679; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1867B: ; /* IDA 0x1867B; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_1867D: ; /* IDA 0x1867D; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1867F: ; /* IDA 0x1867F; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18681: ; /* IDA 0x18681; add     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
L_18683: ; /* IDA 0x18683; mov     ax, cs:word_181AA */
    m->ax=(uint16_t)(dl_read16(m,m->cs,33194));
L_18687: ; /* IDA 0x18687; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18689: ; /* IDA 0x18689; inc     ax */
    m->ax=(uint16_t)(ec_inc(m,m->ax,16));
L_1868A: ; /* IDA 0x1868A; add     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
L_1868C: ; /* IDA 0x1868C; mov     si, 86FBh */
    m->si=(uint16_t)(34555);
L_1868F: ; /* IDA 0x1868F; sub     cl, 2 */
    vg_cl(m,(uint8_t)(ec_math(m,(uint8_t)m->cx,2,8,0,1)));
L_18692: ; /* IDA 0x18692; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_18694: ; /* IDA 0x18694; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_18696: ; /* IDA 0x18696; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_18699: ; /* IDA 0x18699; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1869B: ; /* IDA 0x1869B; call    sub_10CAC */
    if(ll_call(m,3244,34462,h)) return 1;
L_1869E: ; /* IDA 0x1869E; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_186A1: ; /* IDA 0x186A1; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_186A3: ; /* IDA 0x186A3; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_186A4: ; /* IDA 0x186A4; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_186A5: ; /* IDA 0x186A5; mov     al, 10h */
    vg_al(m,(uint8_t)(16));
L_186A7: ; /* IDA 0x186A7; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_186A8: ; /* IDA 0x186A8; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_186A9: ; /* IDA 0x186A9; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_186AB: ; /* IDA 0x186AB; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_186AC: ; /* IDA 0x186AC; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_186AE: ; /* IDA 0x186AE; call    sub_186C9 */
    if(ll_call(m,34505,34481,h)) return 1;
L_186B1: ; /* IDA 0x186B1; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_186B3: ; /* IDA 0x186B3; call    sub_186C9 */
    if(ll_call(m,34505,34486,h)) return 1;
L_186B6: ; /* IDA 0x186B6; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_186B8: ; /* IDA 0x186B8; call    sub_186C9 */
    if(ll_call(m,34505,34491,h)) return 1;
L_186BB: ; /* IDA 0x186BB; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_186BD: ; /* IDA 0x186BD; call    sub_186C9 */
    if(ll_call(m,34505,34496,h)) return 1;
L_186C0: ; /* IDA 0x186C0; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_186C1: ; /* IDA 0x186C1; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_186C2: ; /* IDA 0x186C2; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_186C3: ; /* IDA 0x186C3; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_186C4: ; /* IDA 0x186C4; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_186C5: ; /* IDA 0x186C5; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_186C6: ; /* IDA 0x186C6; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_186C7: ; /* IDA 0x186C7; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_186C8: ; /* IDA 0x186C8; retn */
    return 0;
}
/* Original sub_186C9; IDA 0x186C9. */
static int ll_body_186C9(KiMachine16 *m,const KiEconomyHooks *h) {
L_186C9: ; /* IDA 0x186C9; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_186CC: ; /* IDA 0x186CC; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_186CD: ; /* IDA 0x186CD; mov     ch, cl */
    vg_ch(m,(uint8_t)((uint8_t)m->cx));
L_186CF: ; /* IDA 0x186CF; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_186D1: ; /* IDA 0x186D1; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_186D4: ; /* IDA 0x186D4; movsb */
    vg_movsb(m);
L_186D5: ; /* IDA 0x186D5; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_186D8: ; /* IDA 0x186D8; movsb */
    vg_movsb(m);
L_186D9: ; /* IDA 0x186D9; add     di, 4Eh ; 'N' */
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
L_186DC: ; /* IDA 0x186DC; lodsw */
    mc_string(m,1,0);
L_186DD: ; /* IDA 0x186DD; mov     dl, es:[di] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_186E0: ; /* IDA 0x186E0; mov     es:[di], al */
    dl_write8(m,m->es,(uint16_t)(m->di),(uint8_t)m->ax);
L_186E3: ; /* IDA 0x186E3; mov     dl, es:[di+1] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di+1))));
L_186E7: ; /* IDA 0x186E7; mov     es:[di+1], ah */
    dl_write8(m,m->es,(uint16_t)(m->di+1),(uint8_t)(m->ax>>8));
L_186EB: ; /* IDA 0x186EB; add     di, 50h ; 'P' */
    m->di=(uint16_t)(ec_math(m,m->di,80,16,0,0));
L_186EE: ; /* IDA 0x186EE; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_186F0: ; /* IDA 0x186F0; jnz     short loc_186DD */
    if(!(m->flags&0x40)) goto L_186DD;
L_186F2: ; /* IDA 0x186F2; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_186F5: ; /* IDA 0x186F5; movsb */
    vg_movsb(m);
L_186F6: ; /* IDA 0x186F6; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_186F9: ; /* IDA 0x186F9; movsb */
    vg_movsb(m);
L_186FA: ; /* IDA 0x186FA; retn */
    return 0;
}
/* Original sub_18713; IDA 0x18713. */
static int ll_body_18713(KiMachine16 *m,const KiEconomyHooks *h) {
L_18713: ; /* IDA 0x18713; mov     dx, cs:word_181B0 */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33200));
L_18718: ; /* IDA 0x18718; add     dx, 20h ; ' ' */
    m->dx=(uint16_t)(ec_math(m,m->dx,32,16,0,0));
L_1871B: ; /* IDA 0x1871B; sub     ax, dx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->dx,16,0,1));
L_1871D: ; /* IDA 0x1871D; jnb     short loc_18721 */
    if(!(m->flags&1)) goto L_18721;
L_1871F: ; /* IDA 0x1871F; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_18721: ; /* IDA 0x18721; mov     dx, cs:word_181B4 */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33204));
L_18726: ; /* IDA 0x18726; sub     dx, 30h ; '0' */
    m->dx=(uint16_t)(ec_math(m,m->dx,48,16,0,1));
L_18729: ; /* IDA 0x18729; cmp     ax, dx */
    ec_cmp(m,m->ax,m->dx,16);
L_1872B: ; /* IDA 0x1872B; jb      short loc_1872F */
    if(m->flags&1) goto L_1872F;
L_1872D: ; /* IDA 0x1872D; mov     ax, dx */
    m->ax=(uint16_t)(m->dx);
L_1872F: ; /* IDA 0x1872F; mul     cs:byte_181BE */
    dl_mul(m,dl_read8(m,m->cs,33214),8);
L_18734: ; /* IDA 0x18734; div     dl */
#if KI_LIST_MUTATION == 8
    dl_div(m,(uint8_t)m->dx+1,8);
#else
    dl_div(m,(uint8_t)m->dx,8);
#endif
L_18736: ; /* IDA 0x18736; mov     dl, byte ptr cs:word_181B6+1 */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,33207)));
L_1873B: ; /* IDA 0x1873B; dec     dl */
    vg_dl(m,(uint8_t)(st_dec(m,(uint8_t)m->dx,8)));
L_1873D: ; /* IDA 0x1873D; mov     dh, dl */
    vg_dh(m,(uint8_t)((uint8_t)m->dx));
L_1873F: ; /* IDA 0x1873F; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_18741: ; /* IDA 0x18741; sub     al, dl */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,(uint8_t)m->dx,8,0,1)));
L_18743: ; /* IDA 0x18743; jnb     short loc_18747 */
    if(!(m->flags&1)) goto L_18747;
L_18745: ; /* IDA 0x18745; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_18747: ; /* IDA 0x18747; mov     ah, cs:byte_181BE */
    vg_ah(m,(uint8_t)(dl_read8(m,m->cs,33214)));
L_1874C: ; /* IDA 0x1874C; sub     ah, dh */
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),(uint8_t)(m->dx>>8),8,0,1)));
L_1874E: ; /* IDA 0x1874E; cmp     al, ah */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
L_18750: ; /* IDA 0x18750; jbe     short locret_18754 */
    if(m->flags&(1|0x40)) goto L_18754;
L_18752: ; /* IDA 0x18752; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_18754: ; /* IDA 0x18754; retn */
    return 0;
}
/* Original sub_18755; IDA 0x18755. */
static int ll_body_18755(KiMachine16 *m,const KiEconomyHooks *h) {
L_18755: ; /* IDA 0x18755; push    ds */
    ec_push(m,m->ds);
L_18756: ; /* IDA 0x18756; push    ax */
    ec_push(m,m->ax);
L_18757: ; /* IDA 0x18757; push    bx */
    ec_push(m,m->bx);
L_18758: ; /* IDA 0x18758; push    cx */
    ec_push(m,m->cx);
L_18759: ; /* IDA 0x18759; push    si */
    ec_push(m,m->si);
L_1875A: ; /* IDA 0x1875A; push    di */
    ec_push(m,m->di);
L_1875B: ; /* IDA 0x1875B; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_1875D: ; /* IDA 0x1875D; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1875F: ; /* IDA 0x1875F; mov     dx, 0D6Dh */
    m->dx=(uint16_t)(3437);
L_18762: ; /* IDA 0x18762; mov     bx, word_1987C */
    m->bx=(uint16_t)(dl_read16(m,m->ds,39036));
L_18766: ; /* IDA 0x18766; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_18768: ; /* IDA 0x18768; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_1876A: ; /* IDA 0x1876A; mov     ax, cx */
    m->ax=(uint16_t)(m->cx);
L_1876C: ; /* IDA 0x1876C; mov     di, 2800h */
    m->di=(uint16_t)(10240);
L_1876F: ; /* IDA 0x1876F; call    sub_1E38C */
    if(ll_call(m,58252,34674,h)) return 1;
L_18772: ; /* IDA 0x18772; mov     ds, cs:word_1987C */
    m->ds=(uint16_t)(dl_read16(m,m->cs,39036));
L_18777: ; /* IDA 0x18777; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_18779: ; /* IDA 0x18779; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_1877B: ; /* IDA 0x1877B; mov     ax, 2028h */
    m->ax=(uint16_t)(8232);
L_1877E: ; /* IDA 0x1877E; call    sub_1FA37 */
    if(ll_call(m,64055,34689,h)) return 1;
L_18781: ; /* IDA 0x18781; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_18783: ; /* IDA 0x18783; mov     bx, dx */
    m->bx=(uint16_t)(m->dx);
L_18785: ; /* IDA 0x18785; mov     cx, 450h */
    m->cx=(uint16_t)(1104);
L_18788: ; /* IDA 0x18788; mov     al, 6 */
    vg_al(m,(uint8_t)(6));
L_1878A: ; /* IDA 0x1878A; call    sub_1E3D7 */
    if(ll_call(m,58327,34701,h)) return 1;
L_1878D: ; /* IDA 0x1878D; mov     cx, 5 */
    m->cx=(uint16_t)(5);
L_18790: ; /* IDA 0x18790; mov     dx, 150h */
    m->dx=(uint16_t)(336);
L_18793: ; /* IDA 0x18793; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_18795: ; /* IDA 0x18795; push    cx */
    ec_push(m,m->cx);
L_18796: ; /* IDA 0x18796; push    dx */
    ec_push(m,m->dx);
L_18797: ; /* IDA 0x18797; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_18799: ; /* IDA 0x18799; mov     cx, 404h */
    m->cx=(uint16_t)(1028);
L_1879C: ; /* IDA 0x1879C; call    sub_1E3D7 */
    if(ll_call(m,58327,34719,h)) return 1;
L_1879F: ; /* IDA 0x1879F; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_187A1: ; /* IDA 0x187A1; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_187A2: ; /* IDA 0x187A2; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_187A3: ; /* IDA 0x187A3; add     dx, 20h ; ' ' */
    m->dx=(uint16_t)(ec_math(m,m->dx,32,16,0,0));
L_187A6: ; /* IDA 0x187A6; loop    loc_18795 */
    if(--m->cx) goto L_18795;
L_187A8: ; /* IDA 0x187A8; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_187A9: ; /* IDA 0x187A9; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_187AA: ; /* IDA 0x187AA; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_187AB: ; /* IDA 0x187AB; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_187AC: ; /* IDA 0x187AC; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_187AD: ; /* IDA 0x187AD; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_187AE: ; /* IDA 0x187AE; retn */
    return 0;
}
/* Original unnamed code; IDA 0x1743B. */
static int ll_body_1743B(KiMachine16 *m,const KiEconomyHooks *h) {
L_1743B: ; /* IDA 0x1743B; push    si */
    ec_push(m,m->si);
L_1743C: ; /* IDA 0x1743C; push    di */
    ec_push(m,m->di);
L_1743D: ; /* IDA 0x1743D; mov     si, 840h */
    m->si=(uint16_t)(2112);
L_17440: ; /* IDA 0x17440; xor     di, di */
    {uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_17442: ; /* IDA 0x17442; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_17446: ; /* IDA 0x17446; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17448: ; /* IDA 0x17448; cmp     al, [si+1] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_1744B: ; /* IDA 0x1744B; jnz     short loc_17453 */
    if(!(m->flags&0x40)) goto L_17453;
L_1744D: ; /* IDA 0x1744D; mov     [bp+di], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->si);
L_1744F: ; /* IDA 0x1744F; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_17451: ; /* IDA 0x17451; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_17452: ; /* IDA 0x17452; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_17453: ; /* IDA 0x17453; add     si, 20h ; ' ' */
    m->si=(uint16_t)(ec_math(m,m->si,32,16,0,0));
L_17456: ; /* IDA 0x17456; cmp     si, 2040h */
#if KI_LIST_MUTATION == 1
    ec_cmp(m,m->si,8224,16);
#else
    ec_cmp(m,m->si,8256,16);
#endif
L_1745A: ; /* IDA 0x1745A; jnz     short loc_17448 */
    if(!(m->flags&0x40)) goto L_17448;
L_1745C: ; /* IDA 0x1745C; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1745D: ; /* IDA 0x1745D; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1745E: ; /* IDA 0x1745E; retn */
    return 0;
}
/* Original unnamed code; IDA 0x1745F. */
static int ll_body_1745F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1745F: ; /* IDA 0x1745F; push    ax */
    ec_push(m,m->ax);
L_17460: ; /* IDA 0x17460; push    bx */
    ec_push(m,m->bx);
L_17461: ; /* IDA 0x17461; push    dx */
    ec_push(m,m->dx);
L_17462: ; /* IDA 0x17462; push    si */
    ec_push(m,m->si);
L_17463: ; /* IDA 0x17463; call    sub_18607 */
    if(ll_call(m,34311,29798,h)) return 1;
L_17466: ; /* IDA 0x17466; mov     dx, cs:word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33198));
L_1746B: ; /* IDA 0x1746B; add     dx, 18h */
    m->dx=(uint16_t)(ec_math(m,m->dx,24,16,0,0));
L_1746E: ; /* IDA 0x1746E; mov     bx, cs:word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33200));
L_17473: ; /* IDA 0x17473; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17476: ; /* IDA 0x17476; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17478: ; /* IDA 0x17478; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1747A: ; /* IDA 0x1747A; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1747C: ; /* IDA 0x1747C; mov     ax, 0Ah */
    m->ax=(uint16_t)(10);
L_1747F: ; /* IDA 0x1747F; call    sub_1748F */
    if(ll_call(m,29839,29826,h)) return 1;
L_17482: ; /* IDA 0x17482; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17483: ; /* IDA 0x17483; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17484: ; /* IDA 0x17484; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17487: ; /* IDA 0x17487; dec     ax */
    m->ax=(uint16_t)(st_dec(m,m->ax,16));
L_17488: ; /* IDA 0x17488; jnz     short loc_1747F */
    if(!(m->flags&0x40)) goto L_1747F;
L_1748A: ; /* IDA 0x1748A; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1748B: ; /* IDA 0x1748B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1748C: ; /* IDA 0x1748C; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1748D: ; /* IDA 0x1748D; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1748E: ; /* IDA 0x1748E; retn */
    return 0;
}
/* Original loc_1828F; IDA 0x1828F. */
static int ll_body_1828F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1828F: ; /* IDA 0x1828F; push    bx */
    ec_push(m,m->bx);
L_18290: ; /* IDA 0x18290; push    cx */
    ec_push(m,m->cx);
L_18291: ; /* IDA 0x18291; push    dx */
    ec_push(m,m->dx);
L_18292: ; /* IDA 0x18292; mov     dx, cs:word_181AA */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33194));
L_18297: ; /* IDA 0x18297; mov     bx, cs:word_181AC */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33196));
L_1829C: ; /* IDA 0x1829C; sub     bx, 2 */
    m->bx=(uint16_t)(ec_math(m,m->bx,2,16,0,1));
L_1829F: ; /* IDA 0x1829F; mov     cx, cs:word_181B6 */
    m->cx=(uint16_t)(dl_read16(m,m->cs,33206));
L_182A4: ; /* IDA 0x182A4; add     cx, 101h */
    m->cx=(uint16_t)(ec_math(m,m->cx,257,16,0,0));
L_182A8: ; /* IDA 0x182A8; call    sub_1895D */
    if(ll_call(m,35165,33451,h)) return 1;
L_182AB: ; /* IDA 0x182AB; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_182AD: ; /* IDA 0x182AD; jnz     short loc_182B3 */
    if(!(m->flags&0x40)) goto L_182B3;
L_182AF: ; /* IDA 0x182AF; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_182B0: ; /* IDA 0x182B0; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_182B1: ; /* IDA 0x182B1; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_182B2: ; /* IDA 0x182B2; retn */
    return 0;
L_182B3: ; /* IDA 0x182B3; push    ds */
    ec_push(m,m->ds);
L_182B4: ; /* IDA 0x182B4; push    ax */
    ec_push(m,m->ax);
L_182B5: ; /* IDA 0x182B5; push    si */
    ec_push(m,m->si);
L_182B6: ; /* IDA 0x182B6; push    di */
    ec_push(m,m->di);
L_182B7: ; /* IDA 0x182B7; mov     dx, cs */
    m->dx=(uint16_t)(m->cs);
L_182B9: ; /* IDA 0x182B9; mov     ds, dx */
    m->ds=(uint16_t)(m->dx);
L_182BB: ; /* IDA 0x182BB; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_182BD: ; /* IDA 0x182BD; mov     dx, word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33198));
L_182C1: ; /* IDA 0x182C1; add     dx, 10h */
    m->dx=(uint16_t)(ec_math(m,m->dx,16,16,0,0));
L_182C4: ; /* IDA 0x182C4; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_182C8: ; /* IDA 0x182C8; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_182CB: ; /* IDA 0x182CB; mov     cx, word_181B6 */
    m->cx=(uint16_t)(dl_read16(m,m->ds,33206));
L_182CF: ; /* IDA 0x182CF; sub     cx, 101h */
    m->cx=(uint16_t)(ec_math(m,m->cx,257,16,0,1));
L_182D3: ; /* IDA 0x182D3; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_182D5: ; /* IDA 0x182D5; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_182D7: ; /* IDA 0x182D7; call    sub_1E3D7 */
    if(ll_call(m,58327,33498,h)) return 1;
L_182DA: ; /* IDA 0x182DA; mov     dx, word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33198));
L_182DE: ; /* IDA 0x182DE; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_182E2: ; /* IDA 0x182E2; add     bx, 20h ; ' ' */
    m->bx=(uint16_t)(ec_math(m,m->bx,32,16,0,0));
L_182E5: ; /* IDA 0x182E5; mov     cx, word_181B6 */
    m->cx=(uint16_t)(dl_read16(m,m->ds,33206));
L_182E9: ; /* IDA 0x182E9; mov     cl, 1 */
    vg_cl(m,(uint8_t)(1));
L_182EB: ; /* IDA 0x182EB; sub     ch, 3 */
    vg_ch(m,(uint8_t)(ec_math(m,(uint8_t)(m->cx>>8),3,8,0,1)));
L_182EE: ; /* IDA 0x182EE; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_182F0: ; /* IDA 0x182F0; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_182F2: ; /* IDA 0x182F2; call    sub_1E3D7 */
    if(ll_call(m,58327,33525,h)) return 1;
L_182F5: ; /* IDA 0x182F5; mov     dx, word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33198));
L_182F9: ; /* IDA 0x182F9; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_182FD: ; /* IDA 0x182FD; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_18300: ; /* IDA 0x18300; mov     cx, 202h */
    m->cx=(uint16_t)(514);
L_18303: ; /* IDA 0x18303; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_18305: ; /* IDA 0x18305; call    sub_1E3D7 */
    if(ll_call(m,58327,33544,h)) return 1;
L_18308: ; /* IDA 0x18308; mov     dx, word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33198));
L_1830C: ; /* IDA 0x1830C; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_18310: ; /* IDA 0x18310; add     bx, word_181B4 */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->ds,33204),16,0,0));
L_18314: ; /* IDA 0x18314; sub     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,1));
L_18317: ; /* IDA 0x18317; mov     cx, 202h */
    m->cx=(uint16_t)(514);
L_1831A: ; /* IDA 0x1831A; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_1831C: ; /* IDA 0x1831C; call    sub_1E3D7 */
    if(ll_call(m,58327,33567,h)) return 1;
L_1831F: ; /* IDA 0x1831F; mov     si, 1234h */
    m->si=dl_read16(m,m->cs,0x8320);
L_18322: ; /* IDA 0x18322; mov     cx, [si] */
    m->cx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_18324: ; /* IDA 0x18324; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_18325: ; /* IDA 0x18325; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_18326: ; /* IDA 0x18326; push    cx */
    ec_push(m,m->cx);
L_18327: ; /* IDA 0x18327; mov     dx, [si] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_18329: ; /* IDA 0x18329; add     dx, word_181AE */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,33198),16,0,0));
L_1832D: ; /* IDA 0x1832D; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_18331: ; /* IDA 0x18331; mov     cx, [si+2] */
    m->cx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+2)));
L_18334: ; /* IDA 0x18334; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_18336: ; /* IDA 0x18336; call    sub_1E3D7 */
    if(ll_call(m,58327,33593,h)) return 1;
L_18339: ; /* IDA 0x18339; add     si, 8 */
    m->si=(uint16_t)(ec_math(m,m->si,8,16,0,0));
L_1833C: ; /* IDA 0x1833C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1833D: ; /* IDA 0x1833D; loop    loc_18326 */
    if(--m->cx) goto L_18326;
L_1833F: ; /* IDA 0x1833F; mov     ax, di */
    m->ax=(uint16_t)(m->di);
L_18341: ; /* IDA 0x18341; mov     al, byte ptr word_181B6+1 */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,33207)));
L_18344: ; /* IDA 0x18344; dec     al */
    vg_al(m,(uint8_t)(st_dec(m,(uint8_t)m->ax,8)));
L_18346: ; /* IDA 0x18346; cmp     ah, al */
    ec_cmp(m,(uint8_t)(m->ax>>8),(uint8_t)m->ax,8);
L_18348: ; /* IDA 0x18348; jnb     short near ptr unk_1834C */
    if(!(m->flags&1)) goto L_1834C;
L_1834A: ; /* IDA 0x1834A; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_1834C: ; /* IDA 0x1834C; mov     byte_181BE, ah */
    dl_write8(m,m->ds,33214,(uint8_t)(m->ax>>8));
L_18350: ; /* IDA 0x18350; mov     bl, ah */
    gw_bl(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_18352: ; /* IDA 0x18352; mov     dx, word_181B4 */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33204));
L_18356: ; /* IDA 0x18356; sub     dx, 30h ; '0' */
    m->dx=(uint16_t)(ec_math(m,m->dx,48,16,0,1));
L_18359: ; /* IDA 0x18359; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1835B: ; /* IDA 0x1835B; mul     dx */
    dl_mul(m,m->dx,16);
L_1835D: ; /* IDA 0x1835D; div     bl */
    dl_div(m,(uint8_t)m->bx,8);
L_1835F: ; /* IDA 0x1835F; mov     byte_181BF, al */
    dl_write8(m,m->ds,33215,(uint8_t)m->ax);
L_18362: ; /* IDA 0x18362; mov     word_181A4, 0FFFFh */
    dl_write16(m,m->ds,33188,65535);
L_18368: ; /* IDA 0x18368; mov     dx, word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33198));
L_1836C: ; /* IDA 0x1836C; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_18370: ; /* IDA 0x18370; mov     si, dx */
    m->si=(uint16_t)(m->dx);
L_18372: ; /* IDA 0x18372; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_18374: ; /* IDA 0x18374; dec     si */
    m->si=(uint16_t)(st_dec(m,m->si,16));
L_18375: ; /* IDA 0x18375; add     di, 0Fh */
    m->di=(uint16_t)(ec_math(m,m->di,15,16,0,0));
L_18378: ; /* IDA 0x18378; add     si, word_181B2 */
    m->si=(uint16_t)(ec_math(m,m->si,dl_read16(m,m->ds,33202),16,0,0));
L_1837C: ; /* IDA 0x1837C; mov     ah, 0 */
    vg_ah(m,(uint8_t)(0));
L_1837E: ; /* IDA 0x1837E; call    sub_1F1A3 */
    if(ll_call(m,61859,33665,h)) return 1;
L_18381: ; /* IDA 0x18381; mov     dx, word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33198));
L_18385: ; /* IDA 0x18385; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_18389: ; /* IDA 0x18389; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_1838C: ; /* IDA 0x1838C; mov     si, dx */
    m->si=(uint16_t)(m->dx);
L_1838E: ; /* IDA 0x1838E; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_18390: ; /* IDA 0x18390; add     si, 0Fh */
    m->si=(uint16_t)(ec_math(m,m->si,15,16,0,0));
L_18393: ; /* IDA 0x18393; add     di, word_181B4 */
    m->di=(uint16_t)(ec_math(m,m->di,dl_read16(m,m->ds,33204),16,0,0));
L_18397: ; /* IDA 0x18397; sub     di, 11h */
    m->di=(uint16_t)(ec_math(m,m->di,17,16,0,1));
L_1839A: ; /* IDA 0x1839A; mov     ah, 0 */
    vg_ah(m,(uint8_t)(0));
L_1839C: ; /* IDA 0x1839C; call    sub_1F1A3 */
    if(ll_call(m,61859,33695,h)) return 1;
L_1839F: ; /* IDA 0x1839F; mov     dx, word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33198));
L_183A3: ; /* IDA 0x183A3; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_183A7: ; /* IDA 0x183A7; add     dx, 10h */
    m->dx=(uint16_t)(ec_math(m,m->dx,16,16,0,0));
L_183AA: ; /* IDA 0x183AA; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_183AD: ; /* IDA 0x183AD; mov     si, dx */
    m->si=(uint16_t)(m->dx);
L_183AF: ; /* IDA 0x183AF; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_183B1: ; /* IDA 0x183B1; add     si, word_181B2 */
    m->si=(uint16_t)(ec_math(m,m->si,dl_read16(m,m->ds,33202),16,0,0));
L_183B5: ; /* IDA 0x183B5; add     di, word_181B4 */
    m->di=(uint16_t)(ec_math(m,m->di,dl_read16(m,m->ds,33204),16,0,0));
L_183B9: ; /* IDA 0x183B9; sub     si, 11h */
    m->si=(uint16_t)(ec_math(m,m->si,17,16,0,1));
L_183BC: ; /* IDA 0x183BC; sub     di, 11h */
    m->di=(uint16_t)(ec_math(m,m->di,17,16,0,1));
L_183BF: ; /* IDA 0x183BF; mov     ah, 9 */
    vg_ah(m,(uint8_t)(9));
L_183C1: ; /* IDA 0x183C1; call    sub_1F1A3 */
    if(ll_call(m,61859,33732,h)) return 1;
L_183C4: ; /* IDA 0x183C4; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_183C7: ; /* IDA 0x183C7; mov     dx, word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->ds,33198));
L_183CB: ; /* IDA 0x183CB; add     dx, 18h */
    m->dx=(uint16_t)(ec_math(m,m->dx,24,16,0,0));
L_183CE: ; /* IDA 0x183CE; mov     bx, word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,33200));
L_183D2: ; /* IDA 0x183D2; mov     si, 1234h */
#if KI_LIST_MUTATION == 3
    m->si=0x7378+2;
#else
    m->si=dl_read16(m,m->cs,0x83d3);
#endif
L_183D5: ; /* IDA 0x183D5; call    sub_106F5 */
    if(ll_call(m,1781,33752,h)) return 1;
L_183D8: ; /* IDA 0x183D8; mov     ds, word_10D4A */
    m->ds=(uint16_t)(dl_read16(m,m->ds,3402));
L_183DC: ; /* IDA 0x183DC; mov     si, 80h */
    m->si=(uint16_t)(128);
L_183DF: ; /* IDA 0x183DF; mov     bx, cs:word_181B8 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33208));
L_183E4: ; /* IDA 0x183E4; add     bx, 500h */
    m->bx=(uint16_t)(ec_math(m,m->bx,1280,16,0,0));
L_183E8: ; /* IDA 0x183E8; mov     ax, 1001h */
    m->ax=(uint16_t)(4097);
L_183EB: ; /* IDA 0x183EB; call    sub_1FA37 */
    if(ll_call(m,64055,33774,h)) return 1;
L_183EE: ; /* IDA 0x183EE; mov     ax, cs:word_181B4 */
    m->ax=(uint16_t)(dl_read16(m,m->cs,33204));
L_183F2: ; /* IDA 0x183F2; sub     ax, 10h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16,16,0,1));
L_183F5: ; /* IDA 0x183F5; mov     dx, 50h ; 'P' */
    m->dx=(uint16_t)(80);
L_183F8: ; /* IDA 0x183F8; mul     dx */
    dl_mul(m,m->dx,16);
L_183FA: ; /* IDA 0x183FA; mov     si, 100h */
    m->si=(uint16_t)(256);
L_183FD: ; /* IDA 0x183FD; mov     bx, cs:word_181B8 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33208));
L_18402: ; /* IDA 0x18402; add     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
L_18404: ; /* IDA 0x18404; mov     ax, 1001h */
    m->ax=(uint16_t)(4097);
L_18407: ; /* IDA 0x18407; call    sub_1FA37 */
    if(ll_call(m,64055,33802,h)) return 1;
L_1840A: ; /* IDA 0x1840A; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1840B: ; /* IDA 0x1840B; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1840C: ; /* IDA 0x1840C; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1840D: ; /* IDA 0x1840D; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1840E: ; /* IDA 0x1840E; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1840F: ; /* IDA 0x1840F; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_18410: ; /* IDA 0x18410; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_18411: ; /* IDA 0x18411; retn */
    return 0;
}
/* Original loc_1857F; IDA 0x1857F. */
static int ll_body_1857F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1857F: ; /* IDA 0x1857F; call    sub_101B4 */
    if(ll_call(m,436,34178,h)) return 1;
L_18582: ; /* IDA 0x18582; mov     cs:word_181BA, bx */
    dl_write16(m,m->cs,33210,m->bx);
L_18587: ; /* IDA 0x18587; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_18589: ; /* IDA 0x18589; jz      short loc_185A9 */
    if(m->flags&0x40) goto L_185A9;
L_1858B: ; /* IDA 0x1858B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1858D: ; /* IDA 0x1858D; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1858F: ; /* IDA 0x1858F; add     bx, 1234h */
    m->bx=ec_math(m,m->bx,dl_read16(m,m->cs,0x8591),16,0,0);
L_18593: ; /* IDA 0x18593; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_18594: ; /* IDA 0x18594; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_18595: ; /* IDA 0x18595; mov     si, cs:[bx+4] */
    m->si=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->bx+4)));
L_18599: ; /* IDA 0x18599; mov     ax, cs:[bx+6] */
    m->ax=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->bx+6)));
L_1859D: ; /* IDA 0x1859D; mov     ch, cs:byte_181BD */
    vg_ch(m,(uint8_t)(dl_read8(m,m->cs,33213)));
L_185A2: ; /* IDA 0x185A2; call    loc_185B2 */
    if(ll_call(m,34226,34213,h)) return 1;
L_185A5: ; /* IDA 0x185A5; call    sub_101DB */
    if(ll_call(m,475,34216,h)) return 1;
L_185A8: ; /* IDA 0x185A8; retn */
    return 0;
L_185A9: ; /* IDA 0x185A9; call    cs:word_181A8 */
    if(ll_call(m,dl_read16(m,m->cs,33192),34222,h)) return 1;
L_185AE: ; /* IDA 0x185AE; call    sub_101DB */
    if(ll_call(m,475,34225,h)) return 1;
L_185B1: ; /* IDA 0x185B1; retn */
    return 0;
}
/* Original loc_185B2; IDA 0x185B2. */
static int ll_body_185B2(KiMachine16 *m,const KiEconomyHooks *h) {
L_185B2: ; /* IDA 0x185B2; and     cs:byte_185E8, 0FEh */
    {uint16_t result=((dl_read8(m,m->cs,34280))&(254))&255;ec_logic(m,result,8);dl_write8(m,m->cs,34280,result);}
L_185B8: ; /* IDA 0x185B8; or      cs:byte_185E8, ah */
    {uint16_t result=((dl_read8(m,m->cs,34280))|((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);dl_write8(m,m->cs,34280,result);}
L_185BD: ; /* IDA 0x185BD; and     cs:byte_185EF, 0FEh */
    {uint16_t result=((dl_read8(m,m->cs,34287))&(254))&255;ec_logic(m,result,8);dl_write8(m,m->cs,34287,result);}
L_185C3: ; /* IDA 0x185C3; or      cs:byte_185EF, ah */
    {uint16_t result=((dl_read8(m,m->cs,34287))|((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);dl_write8(m,m->cs,34287,result);}
L_185C8: ; /* IDA 0x185C8; and     cs:byte_185F3, 0FEh */
    {uint16_t result=((dl_read8(m,m->cs,34291))&(254))&255;ec_logic(m,result,8);dl_write8(m,m->cs,34291,result);}
L_185CE: ; /* IDA 0x185CE; or      cs:byte_185F3, ah */
    {uint16_t result=((dl_read8(m,m->cs,34291))|((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);dl_write8(m,m->cs,34291,result);}
L_185D3: ; /* IDA 0x185D3; mov     cs:byte_185F1, al */
    dl_write8(m,m->cs,34289,(uint8_t)m->ax);
L_185D7: ; /* IDA 0x185D7; push    di */
    ec_push(m,m->di);
L_185D8: ; /* IDA 0x185D8; xor     di, di */
    {uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_185DA: ; /* IDA 0x185DA; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_185DC: ; /* IDA 0x185DC; jz      short loc_18605 */
    if(m->flags&0x40) goto L_18605;
L_185DE: ; /* IDA 0x185DE; push    bp */
    ec_push(m,m->bp);
L_185DF: ; /* IDA 0x185DF; mov     cl, ch */
    vg_cl(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_185E1: ; /* IDA 0x185E1; add     bp, di */
    m->bp=(uint16_t)(ec_math(m,m->bp,m->di,16,0,0));
L_185E3: ; /* IDA 0x185E3; mov     bx, [bp+0] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_185E6: ; /* IDA 0x185E6; mov     dx, bx */
    m->dx=(uint16_t)(m->bx);
L_185E8: ; /* IDA 0x185E8; mov     ax, [bx+si] */
#if KI_LIST_MUTATION == 4
    m->ax=dl_read16(m,m->ds,(uint16_t)(m->bx+m->si));
#else
    ll_sort_load(m,0x85e8);
#endif
L_185EA: ; /* IDA 0x185EA; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_185EB: ; /* IDA 0x185EB; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_185EC: ; /* IDA 0x185EC; mov     bx, [bp+0] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_185EF: ; /* IDA 0x185EF; cmp     ax, [bx+si] */
    ll_sort_compare(m);
L_185F1: ; /* IDA 0x185F1; jbe     short near ptr byte_185F4+4 */
#if KI_LIST_MUTATION == 5
    if(m->flags&1) goto L_185F8;
#else
    if(ll_sort_skip(m)) goto L_185F8;
#endif
L_185F3: ; /* IDA 0x185F3; mov     ax, [bx+si] */
    ll_sort_load(m,0x85f3);
L_185F5: ; /* IDA 0x185F5; xchg    dx, [bp+0] */
    {uint16_t a=m->dx,b=dl_read16(m,m->ss,(uint16_t)(m->bp));m->dx=(uint16_t)(b);dl_write16(m,m->ss,(uint16_t)(m->bp),a);}
L_185F8: ; /* IDA 0x185F8; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_185F9: ; /* IDA 0x185F9; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_185FA: ; /* IDA 0x185FA; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_185FC: ; /* IDA 0x185FC; jnz     short loc_185EC */
    if(!(m->flags&0x40)) goto L_185EC;
L_185FE: ; /* IDA 0x185FE; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_185FF: ; /* IDA 0x185FF; mov     [bp+di], dx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->dx);
L_18601: ; /* IDA 0x18601; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_18602: ; /* IDA 0x18602; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_18603: ; /* IDA 0x18603; jmp     short loc_185DA */
    goto L_185DA;
L_18605: ; /* IDA 0x18605; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_18606: ; /* IDA 0x18606; retn */
    return 0;
}
int ki_list_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h,int *escaped) {switch(t) {
case 0x5e60:*escaped=ll_body_15E60(m,h);break;
case 0x6909:*escaped=ll_body_16909(m,h);break;
case 0x7400:*escaped=ll_body_17400(m,h);break;
case 0x748f:*escaped=ll_body_1748F(m,h);break;
case 0x81c0:*escaped=ll_body_181C0(m,h);break;
case 0x820e:*escaped=ll_body_1820E(m,h);break;
case 0x8412:*escaped=ll_body_18412(m,h);break;
case 0x8458:*escaped=ll_body_18458(m,h);break;
case 0x8463:*escaped=ll_body_18463(m,h);break;
case 0x84bc:*escaped=ll_body_184BC(m,h);break;
case 0x84dd:*escaped=ll_body_184DD(m,h);break;
case 0x851a:*escaped=ll_body_1851A(m,h);break;
case 0x8546:*escaped=ll_body_18546(m,h);break;
case 0x8607:*escaped=ll_body_18607(m,h);break;
case 0x8662:*escaped=ll_body_18662(m,h);break;
case 0x86c9:*escaped=ll_body_186C9(m,h);break;
case 0x8713:*escaped=ll_body_18713(m,h);break;
case 0x8755:*escaped=ll_body_18755(m,h);break;
case 0x743b:*escaped=ll_body_1743B(m,h);break;
case 0x745f:*escaped=ll_body_1745F(m,h);break;
case 0x828f:*escaped=ll_body_1828F(m,h);break;
case 0x857f:*escaped=ll_body_1857F(m,h);break;
case 0x85b2:*escaped=ll_body_185B2(m,h);break;
default:return 0;}return 1;}
static int ll_known(uint16_t t) {switch(t) {case 0x5e60: case 0x6909: case 0x7400: case 0x748f: case 0x81c0: case 0x820e: case 0x8412: case 0x8458: case 0x8463: case 0x84bc: case 0x84dd: case 0x851a: case 0x8546: case 0x8607: case 0x8662: case 0x86c9: case 0x8713: case 0x8755: case 0x743b: case 0x745f: case 0x828f: case 0x857f: case 0x85b2: return 1;default:return 0;}}
