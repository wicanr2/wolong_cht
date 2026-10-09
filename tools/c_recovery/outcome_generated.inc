/* Original native C battle/city outcome closure; re/130, spec/250. Nonlocal 11CB1 uses the existing transfer runner. */
/* Original sub_10CE7; IDA 0x10CE7. */
static void oc_body_10CE7(KiMachine16 *m,const KiEconomyHooks *h) {
L_10CE7: ; /* IDA 0x10CE7; push    ax */
    ec_push(m,m->ax);
L_10CE8: ; /* IDA 0x10CE8; mov     ax, 202h */
    m->ax=(uint16_t)(514);
L_10CEB: ; /* IDA 0x10CEB; call    sub_1EB11 */
    oc_call(m,60177,3310,h);
L_10CEE: ; /* IDA 0x10CEE; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10CEF: ; /* IDA 0x10CEF; retn */
    return;
}
/* Original sub_14236; IDA 0x14236. */
static void oc_body_14236(KiMachine16 *m,const KiEconomyHooks *h) {
L_14236: ; /* IDA 0x14236; push    ds */
    ec_push(m,m->ds);
L_14237: ; /* IDA 0x14237; push    cx */
    ec_push(m,m->cx);
L_14238: ; /* IDA 0x14238; push    si */
    ec_push(m,m->si);
L_14239: ; /* IDA 0x14239; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_1423E: ; /* IDA 0x1423E; mov     si, 840h */
    m->si=(uint16_t)(2112);
L_14241: ; /* IDA 0x14241; mov     cx, 0C0h */
    m->cx=(uint16_t)(192);
L_14244: ; /* IDA 0x14244; cmp     al, [si+1] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_14247: ; /* IDA 0x14247; jz      short loc_1424E */
    if(m->flags&0x40) goto L_1424E;
L_14249: ; /* IDA 0x14249; cmp     ah, [si+1] */
    ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_1424C: ; /* IDA 0x1424C; jnz     short loc_1425E */
    if(!(m->flags&0x40)) goto L_1425E;
L_1424E: ; /* IDA 0x1424E; cmp     al, [si+1Ah] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+26)),8);
L_14251: ; /* IDA 0x14251; jz      short loc_14258 */
    if(m->flags&0x40) goto L_14258;
L_14253: ; /* IDA 0x14253; cmp     ah, [si+1Ah] */
    ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->ds,(uint16_t)(m->si+26)),8);
L_14256: ; /* IDA 0x14256; jnz     short loc_1425E */
    if(!(m->flags&0x40)) goto L_1425E;
L_14258: ; /* IDA 0x14258; mov     ch, [si+1] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_1425B: ; /* IDA 0x1425B; mov     [si+1Ah], ch */
#if KI_OUTCOME_MUTATION == 19
    ; /* Omitted city owner-history write. */
#else
    dl_write8(m,m->ds,(uint16_t)(m->si+26),(uint8_t)(m->cx>>8));
#endif
L_1425E: ; /* IDA 0x1425E; add     si, 20h ; ' ' */
    m->si=(uint16_t)(ec_math(m,m->si,32,16,0,0));
L_14261: ; /* IDA 0x14261; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_14263: ; /* IDA 0x14263; jnz     short loc_14244 */
    if(!(m->flags&0x40)) goto L_14244;
L_14265: ; /* IDA 0x14265; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14266: ; /* IDA 0x14266; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_14267: ; /* IDA 0x14267; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_14268: ; /* IDA 0x14268; retn */
    return;
}
/* Original sub_1474A; IDA 0x1474A. */
static void oc_body_1474A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1474A: ; /* IDA 0x1474A; push    ax */
    ec_push(m,m->ax);
L_1474B: ; /* IDA 0x1474B; push    bx */
    ec_push(m,m->bx);
L_1474C: ; /* IDA 0x1474C; push    cx */
    ec_push(m,m->cx);
L_1474D: ; /* IDA 0x1474D; push    dx */
    ec_push(m,m->dx);
L_1474E: ; /* IDA 0x1474E; call    sub_16FD2 */
    oc_call(m,28626,18257,h);
L_14751: ; /* IDA 0x14751; cmp     byte ptr [si+6], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+6)),0,8);
L_14755: ; /* IDA 0x14755; jz      short loc_147B5 */
    if(m->flags&0x40) goto L_147B5;
L_14757: ; /* IDA 0x14757; cmp     byte ptr [si+29h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+41)),0,8);
L_1475B: ; /* IDA 0x1475B; jz      short loc_147B5 */
    if(m->flags&0x40) goto L_147B5;
L_1475D: ; /* IDA 0x1475D; and     cl, cl */
    {uint16_t result=(((uint8_t)m->cx)&((uint8_t)m->cx))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_1475F: ; /* IDA 0x1475F; jz      short loc_147AB */
    if(m->flags&0x40) goto L_147AB;
L_14761: ; /* IDA 0x14761; mov     bx, [si+0Eh] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
L_14764: ; /* IDA 0x14764; cmp     bx, 600h */
    ec_cmp(m,m->bx,1536,16);
L_14768: ; /* IDA 0x14768; jnb     short loc_14777 */
    if(!(m->flags&1)) goto L_14777;
L_1476A: ; /* IDA 0x1476A; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1476C: ; /* IDA 0x1476C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1476E: ; /* IDA 0x1476E; mov     al, [bx+841h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+2113))));
L_14772: ; /* IDA 0x14772; cmp     al, [si+1] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_14775: ; /* IDA 0x14775; jz      short loc_147AB */
#if KI_OUTCOME_MUTATION == 10
    ; /* Omitted original own-city stay branch. */
#else
    if(m->flags&0x40) goto L_147AB;
#endif
L_14777: ; /* IDA 0x14777; call    sub_1487B */
    oc_call(m,18555,18298,h);
L_1477A: ; /* IDA 0x1477A; jb      short loc_147B5 */
    if(m->flags&1) goto L_147B5;
L_1477C: ; /* IDA 0x1477C; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1477E: ; /* IDA 0x1477E; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14780: ; /* IDA 0x14780; mov     [si+14h], bx */
    dl_write16(m,m->ds,(uint16_t)(m->si+20),m->bx);
L_14783: ; /* IDA 0x14783; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14785: ; /* IDA 0x14785; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14787: ; /* IDA 0x14787; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14789: ; /* IDA 0x14789; mov     [si+20h], bl */
    dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)m->bx);
L_1478C: ; /* IDA 0x1478C; or      byte ptr [si], 2 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(2))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_1478F: ; /* IDA 0x1478F; cmp     word ptr [si+4], 12Ch */
#if KI_OUTCOME_MUTATION == 9
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),299,16);
#else
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),300,16);
#endif
L_14794: ; /* IDA 0x14794; jbe     short loc_147A5 */
    if(m->flags&(1|0x40)) goto L_147A5;
L_14796: ; /* IDA 0x14796; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14799: ; /* IDA 0x14799; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1479B: ; /* IDA 0x1479B; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1479D: ; /* IDA 0x1479D; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1479F: ; /* IDA 0x1479F; xchg    ax, bx */
    {uint16_t a=m->ax,b=m->bx;m->ax=(uint16_t)(b);m->bx=(uint16_t)(a);}
L_147A0: ; /* IDA 0x147A0; cmp     al, [bx+3] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+3)),8);
L_147A3: ; /* IDA 0x147A3; jnz     short loc_147AB */
    if(!(m->flags&0x40)) goto L_147AB;
L_147A5: ; /* IDA 0x147A5; mov     byte ptr [si+23h], 0Ah */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),10);
L_147A9: ; /* IDA 0x147A9; jmp     short loc_147AF */
    goto L_147AF;
L_147AB: ; /* IDA 0x147AB; mov     byte ptr [si+23h], 8 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),8);
L_147AF: ; /* IDA 0x147AF; clc */
    m->flags&=(uint16_t)~1u;
L_147B0: ; /* IDA 0x147B0; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_147B1: ; /* IDA 0x147B1; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_147B2: ; /* IDA 0x147B2; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_147B3: ; /* IDA 0x147B3; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_147B4: ; /* IDA 0x147B4; retn */
    return;
L_147B5: ; /* IDA 0x147B5; stc */
    m->flags|=1;
L_147B6: ; /* IDA 0x147B6; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_147B7: ; /* IDA 0x147B7; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_147B8: ; /* IDA 0x147B8; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_147B9: ; /* IDA 0x147B9; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_147BA: ; /* IDA 0x147BA; retn */
    return;
}
/* Original sub_14CF3; IDA 0x14CF3. */
static void oc_body_14CF3(KiMachine16 *m,const KiEconomyHooks *h) {
L_14CF3: ; /* IDA 0x14CF3; mov     bh, al */
    gw_bh(m,(uint8_t)((uint8_t)m->ax));
L_14CF5: ; /* IDA 0x14CF5; xchg    bh, [si+1] */
    {uint16_t a=(uint8_t)(m->bx>>8),b=dl_read8(m,m->ds,(uint16_t)(m->si+1));gw_bh(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->si+1),a);}
L_14CF8: ; /* IDA 0x14CF8; mov     [si+1Ah], bh */
    dl_write8(m,m->ds,(uint16_t)(m->si+26),(uint8_t)(m->bx>>8));
L_14CFB: ; /* IDA 0x14CFB; push    ax */
    ec_push(m,m->ax);
L_14CFC: ; /* IDA 0x14CFC; cmp     bh, 18h */
    ec_cmp(m,(uint8_t)(m->bx>>8),24,8);
L_14CFF: ; /* IDA 0x14CFF; jz      short loc_14D21 */
    if(m->flags&0x40) goto L_14D21;
L_14D01: ; /* IDA 0x14D01; call    sub_14D63 */
    oc_call(m,19811,19716,h);
L_14D04: ; /* IDA 0x14D04; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_14D06: ; /* IDA 0x14D06; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14D08: ; /* IDA 0x14D08; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14D0A: ; /* IDA 0x14D0A; dec     byte ptr [bx+23h] */
#if KI_OUTCOME_MUTATION == 11
    ; /* Omitted old city-count DEC. */
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+35),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->bx+35)),8));
#endif
L_14D0D: ; /* IDA 0x14D0D; call    sub_14DF0 */
    oc_call(m,19952,19728,h);
L_14D10: ; /* IDA 0x14D10; pushf */
    ec_push(m,m->flags);
L_14D11: ; /* IDA 0x14D11; cmp     byte ptr [bp+0FEh], 0 */
    ec_cmp(m,dl_read8(m,m->ss,(uint16_t)(m->bp+254)),0,8);
L_14D16: ; /* IDA 0x14D16; jz      short loc_14D1B */
    if(m->flags&0x40) goto L_14D1B;
L_14D18: ; /* IDA 0x14D18; call    sub_14DA4 */
#if KI_OUTCOME_MUTATION == 15
    ; /* Omitted bounded old-corps redirection call. */
#else
    oc_call(m,19876,19739,h);
#endif
L_14D1B: ; /* IDA 0x14D1B; popf */
    m->flags=ec_pop(m)|2;
L_14D1C: ; /* IDA 0x14D1C; jnb     short loc_14D21 */
    if(!(m->flags&1)) goto L_14D21;
L_14D1E: ; /* IDA 0x14D1E; call    sub_14FCE */
    oc_call(m,20430,19745,h);
L_14D21: ; /* IDA 0x14D21; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_14D22: ; /* IDA 0x14D22; mov     bh, al */
    gw_bh(m,(uint8_t)((uint8_t)m->ax));
L_14D24: ; /* IDA 0x14D24; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_14D26: ; /* IDA 0x14D26; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14D28: ; /* IDA 0x14D28; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14D2A: ; /* IDA 0x14D2A; inc     byte ptr [bx+23h] */
#if KI_OUTCOME_MUTATION == 12
    ; /* Omitted new city-count INC. */
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+35),ec_inc(m,dl_read8(m,m->ds,(uint16_t)(m->bx+35)),8));
#endif
L_14D2D: ; /* IDA 0x14D2D; call    sub_18A1E */
    oc_call(m,35358,19760,h);
L_14D30: ; /* IDA 0x14D30; call    sub_188CC */
    oc_call(m,35020,19763,h);
L_14D33: ; /* IDA 0x14D33; test    cs:byte_198A6, 4 */
    {uint16_t result=((dl_read8(m,m->cs,39078))&(4))&255;ec_logic(m,result,8);}
L_14D39: ; /* IDA 0x14D39; jz      short locret_14D62 */
    if(m->flags&0x40) goto L_14D62;
L_14D3B: ; /* IDA 0x14D3B; mov     cx, ds */
    m->cx=(uint16_t)(m->ds);
L_14D3D: ; /* IDA 0x14D3D; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_14D3F: ; /* IDA 0x14D3F; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_14D41: ; /* IDA 0x14D41; call    sub_15CA4 */
    oc_call(m,23716,19780,h);
L_14D44: ; /* IDA 0x14D44; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_14D46: ; /* IDA 0x14D46; mov     dx, [si+8] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+8)));
L_14D49: ; /* IDA 0x14D49; mov     bx, [si+0Ah] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+10)));
L_14D4C: ; /* IDA 0x14D4C; call    sub_19656 */
    oc_call(m,38486,19791,h);
L_14D4F: ; /* IDA 0x14D4F; call    sub_15CE0 */
    oc_call(m,23776,19794,h);
L_14D52: ; /* IDA 0x14D52; mov     dx, [si+8] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+8)));
L_14D55: ; /* IDA 0x14D55; mov     bx, [si+0Ah] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+10)));
L_14D58: ; /* IDA 0x14D58; call    sub_195C9 */
    oc_call(m,38345,19803,h);
L_14D5B: ; /* IDA 0x14D5B; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_14D5D: ; /* IDA 0x14D5D; call    sub_15C58 */
    oc_call(m,23640,19808,h);
L_14D60: ; /* IDA 0x14D60; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_14D62: ; /* IDA 0x14D62; retn */
    return;
}
/* Original sub_14D63; IDA 0x14D63. */
static void oc_body_14D63(KiMachine16 *m,const KiEconomyHooks *h) {
L_14D63: ; /* IDA 0x14D63; cmp     byte ptr [si+19h], 0FFh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+25)),255,8);
L_14D67: ; /* IDA 0x14D67; jz      short locret_14DA3 */
    if(m->flags&0x40) goto L_14DA3;
L_14D69: ; /* IDA 0x14D69; push    ax */
    ec_push(m,m->ax);
L_14D6A: ; /* IDA 0x14D6A; push    bx */
    ec_push(m,m->bx);
L_14D6B: ; /* IDA 0x14D6B; push    cx */
    ec_push(m,m->cx);
L_14D6C: ; /* IDA 0x14D6C; push    di */
    ec_push(m,m->di);
L_14D6D: ; /* IDA 0x14D6D; mov     bh, 0FFh */
    gw_bh(m,(uint8_t)(255));
L_14D6F: ; /* IDA 0x14D6F; xchg    bh, [si+19h] */
    {uint16_t a=(uint8_t)(m->bx>>8),b=dl_read8(m,m->ds,(uint16_t)(m->si+25));gw_bh(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->si+25),a);}
L_14D72: ; /* IDA 0x14D72; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_14D74: ; /* IDA 0x14D74; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14D76: ; /* IDA 0x14D76; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14D78: ; /* IDA 0x14D78; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14D7A: ; /* IDA 0x14D7A; add     bx, 4240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16960,16,0,0));
L_14D7E: ; /* IDA 0x14D7E; mov     byte ptr [bx+17h], 0 */
#if KI_OUTCOME_MUTATION == 13
    dl_write8(m,m->ds,(uint16_t)(m->bx+0x17),1);
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),0);
#endif
L_14D82: ; /* IDA 0x14D82; push    bx */
    ec_push(m,m->bx);
L_14D83: ; /* IDA 0x14D83; push    si */
    ec_push(m,m->si);
L_14D84: ; /* IDA 0x14D84; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_14D86: ; /* IDA 0x14D86; call    sub_10CE7 */
    oc_call(m,3303,19849,h);
L_14D89: ; /* IDA 0x14D89; mov     cx, 44h ; 'D' */
    m->cx=(uint16_t)(68);
L_14D8C: ; /* IDA 0x14D8C; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_14D8E: ; /* IDA 0x14D8E; call    sub_18810 */
    oc_call(m,34832,19857,h);
L_14D91: ; /* IDA 0x14D91; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14D92: ; /* IDA 0x14D92; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14D93: ; /* IDA 0x14D93; mov     cx, 1A6h */
    m->cx=(uint16_t)(422);
L_14D96: ; /* IDA 0x14D96; mov     ah, [bx+1Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
L_14D99: ; /* IDA 0x14D99; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_14D9C: ; /* IDA 0x14D9C; call    sub_18810 */
    oc_call(m,34832,19871,h);
L_14D9F: ; /* IDA 0x14D9F; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14DA0: ; /* IDA 0x14DA0; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_14DA1: ; /* IDA 0x14DA1; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14DA2: ; /* IDA 0x14DA2; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_14DA3: ; /* IDA 0x14DA3; retn */
    return;
}
/* Original sub_14DA4; IDA 0x14DA4. */
static void oc_body_14DA4(KiMachine16 *m,const KiEconomyHooks *h) {
L_14DA4: ; /* IDA 0x14DA4; push    bx */
    ec_push(m,m->bx);
L_14DA5: ; /* IDA 0x14DA5; push    si */
    ec_push(m,m->si);
L_14DA6: ; /* IDA 0x14DA6; push    bp */
    ec_push(m,m->bp);
L_14DA7: ; /* IDA 0x14DA7; mov     al, [si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14DAA: ; /* IDA 0x14DAA; mov     si, [bp+0] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_14DAD: ; /* IDA 0x14DAD; push    ax */
    ec_push(m,m->ax);
L_14DAE: ; /* IDA 0x14DAE; call    sub_1487B */
    oc_call(m,18555,19889,h);
L_14DB1: ; /* IDA 0x14DB1; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_14DB2: ; /* IDA 0x14DB2; jb      short loc_14DDC */
    if(m->flags&1) goto L_14DDC;
L_14DB4: ; /* IDA 0x14DB4; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14DB6: ; /* IDA 0x14DB6; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14DB8: ; /* IDA 0x14DB8; mov     dx, bx */
    m->dx=(uint16_t)(m->bx);
L_14DBA: ; /* IDA 0x14DBA; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14DBC: ; /* IDA 0x14DBC; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14DBE: ; /* IDA 0x14DBE; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14DC0: ; /* IDA 0x14DC0; mov     cl, [bp+0FEh] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+254))));
L_14DC4: ; /* IDA 0x14DC4; xor     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_14DC6: ; /* IDA 0x14DC6; mov     si, [bp+0] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_14DC9: ; /* IDA 0x14DC9; mov     [si+20h], bl */
    dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)m->bx);
L_14DCC: ; /* IDA 0x14DCC; mov     [si+14h], dx */
    dl_write16(m,m->ds,(uint16_t)(m->si+20),m->dx);
L_14DCF: ; /* IDA 0x14DCF; mov     byte ptr [si+0Bh], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+11),1);
L_14DD3: ; /* IDA 0x14DD3; or      byte ptr [si], 2 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(2))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_14DD6: ; /* IDA 0x14DD6; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_14DD7: ; /* IDA 0x14DD7; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_14DD8: ; /* IDA 0x14DD8; loop    loc_14DC6 */
    if(--m->cx) goto L_14DC6;
L_14DDA: ; /* IDA 0x14DDA; jmp     short loc_14DEC */
    goto L_14DEC;
L_14DDC: ; /* IDA 0x14DDC; mov     cl, [bp+0FEh] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+254))));
L_14DE0: ; /* IDA 0x14DE0; xor     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_14DE2: ; /* IDA 0x14DE2; mov     si, [bp+0] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_14DE5: ; /* IDA 0x14DE5; call    sub_1291A */
    oc_call(m,10522,19944,h);
L_14DE8: ; /* IDA 0x14DE8; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_14DE9: ; /* IDA 0x14DE9; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_14DEA: ; /* IDA 0x14DEA; loop    loc_14DE2 */
    if(--m->cx) goto L_14DE2;
L_14DEC: ; /* IDA 0x14DEC; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_14DED: ; /* IDA 0x14DED; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14DEE: ; /* IDA 0x14DEE; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14DEF: ; /* IDA 0x14DEF; retn */
    return;
}
/* Original sub_14DF0; IDA 0x14DF0. */
static void oc_body_14DF0(KiMachine16 *m,const KiEconomyHooks *h) {
L_14DF0: ; /* IDA 0x14DF0; push    bx */
    ec_push(m,m->bx);
L_14DF1: ; /* IDA 0x14DF1; push    cx */
    ec_push(m,m->cx);
L_14DF2: ; /* IDA 0x14DF2; push    di */
    ec_push(m,m->di);
L_14DF3: ; /* IDA 0x14DF3; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_14DF5: ; /* IDA 0x14DF5; sub     ax, 840h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2112,16,0,1));
L_14DF8: ; /* IDA 0x14DF8; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14DFA: ; /* IDA 0x14DFA; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14DFC: ; /* IDA 0x14DFC; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14DFE: ; /* IDA 0x14DFE; cmp     ah, [bx+3] */
    ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->ds,(uint16_t)(m->bx+3)),8);
L_14E01: ; /* IDA 0x14E01; jnz     short loc_14E4B */
    if(!(m->flags&0x40)) goto L_14E4B;
L_14E03: ; /* IDA 0x14E03; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_14E05: ; /* IDA 0x14E05; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14E07: ; /* IDA 0x14E07; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14E09: ; /* IDA 0x14E09; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_14E0B: ; /* IDA 0x14E0B; call    sub_16A3D */
    oc_call(m,27197,19982,h);
L_14E0E: ; /* IDA 0x14E0E; jb      short loc_14E50 */
    if(m->flags&1) goto L_14E50;
L_14E10: ; /* IDA 0x14E10; sub     ax, 840h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2112,16,0,1));
L_14E13: ; /* IDA 0x14E13; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14E15: ; /* IDA 0x14E15; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14E17: ; /* IDA 0x14E17; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14E19: ; /* IDA 0x14E19; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_14E1B: ; /* IDA 0x14E1B; xchg    ah, [bx+3] */
    {uint16_t a=(uint8_t)(m->ax>>8),b=dl_read8(m,m->ds,(uint16_t)(m->bx+3));vg_ah(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->bx+3),a);}
L_14E1E: ; /* IDA 0x14E1E; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_14E20: ; /* IDA 0x14E20; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_14E22: ; /* IDA 0x14E22; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_14E24: ; /* IDA 0x14E24; mov     bl, bh */
    gw_bl(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_14E26: ; /* IDA 0x14E26; call    sub_14502 */
    oc_call(m,17666,20009,h);
L_14E29: ; /* IDA 0x14E29; mov     bx, cx */
    m->bx=(uint16_t)(m->cx);
L_14E2B: ; /* IDA 0x14E2B; cmp     bx, cs:word_10CFD */
    ec_cmp(m,m->bx,dl_read16(m,m->cs,3325),16);
L_14E30: ; /* IDA 0x14E30; jnz     short loc_14E4B */
    if(!(m->flags&0x40)) goto L_14E4B;
L_14E32: ; /* IDA 0x14E32; mov     al, [bx+3] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_14E35: ; /* IDA 0x14E35; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_14E37: ; /* IDA 0x14E37; push    ax */
    ec_push(m,m->ax);
L_14E38: ; /* IDA 0x14E38; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_14E3A: ; /* IDA 0x14E3A; call    sub_10CE7 */
    oc_call(m,3303,20029,h);
L_14E3D: ; /* IDA 0x14E3D; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_14E3F: ; /* IDA 0x14E3F; mov     cx, 1Eh */
    m->cx=(uint16_t)(30);
L_14E42: ; /* IDA 0x14E42; call    sub_18810 */
    oc_call(m,34832,20037,h);
L_14E45: ; /* IDA 0x14E45; add     sp, 2 */
    m->sp=(uint16_t)(ec_math(m,m->sp,2,16,0,0));
L_14E48: ; /* IDA 0x14E48; call    sub_15E60 */
    oc_call(m,24160,20043,h);
L_14E4B: ; /* IDA 0x14E4B; clc */
    m->flags&=(uint16_t)~1u;
L_14E4C: ; /* IDA 0x14E4C; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14E4D: ; /* IDA 0x14E4D; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_14E4E: ; /* IDA 0x14E4E; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14E4F: ; /* IDA 0x14E4F; retn */
    return;
L_14E50: ; /* IDA 0x14E50; mov     byte ptr [bx+3], 0FFh */
    dl_write8(m,m->ds,(uint16_t)(m->bx+3),255);
L_14E54: ; /* IDA 0x14E54; and     byte ptr [bx], 7Fh */
#if KI_OUTCOME_MUTATION == 14
    {uint8_t result=dl_read8(m,m->ds,m->bx)&0xbf;ec_logic(m,result,8);dl_write8(m,m->ds,m->bx,result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(127))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->bx),result);}
#endif
L_14E57: ; /* IDA 0x14E57; stc */
    m->flags|=1;
L_14E58: ; /* IDA 0x14E58; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14E59: ; /* IDA 0x14E59; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_14E5A: ; /* IDA 0x14E5A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14E5B: ; /* IDA 0x14E5B; retn */
    return;
}
/* Original sub_14FCE; IDA 0x14FCE. */
static void oc_body_14FCE(KiMachine16 *m,const KiEconomyHooks *h) {
L_14FCE: ; /* IDA 0x14FCE; push    ax */
    ec_push(m,m->ax);
L_14FCF: ; /* IDA 0x14FCF; push    bx */
    ec_push(m,m->bx);
L_14FD0: ; /* IDA 0x14FD0; push    cx */
    ec_push(m,m->cx);
L_14FD1: ; /* IDA 0x14FD1; push    di */
    ec_push(m,m->di);
L_14FD2: ; /* IDA 0x14FD2; mov     al, [si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14FD5: ; /* IDA 0x14FD5; mov     cs:byte_12919, al */
    dl_write8(m,m->cs,10521,(uint8_t)m->ax);
L_14FD9: ; /* IDA 0x14FD9; and     byte ptr [bx], 7Fh */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(127))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->bx),result);}
L_14FDC: ; /* IDA 0x14FDC; cmp     bx, cs:word_10CFD */
    ec_cmp(m,m->bx,dl_read16(m,m->cs,3325),16);
L_14FE1: ; /* IDA 0x14FE1; jnz     short loc_14FE8 */
    if(!(m->flags&0x40)) goto L_14FE8;
L_14FE3: ; /* IDA 0x14FE3; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_14FE5: ; /* IDA 0x14FE5; call    sub_11CB1 */
    oc_call(m,7345,20456,h);
L_14FE8: ; /* IDA 0x14FE8; dec     cs:byte_10D2A */
#if KI_OUTCOME_MUTATION == 17
    ; /* Omitted surviving-faction-count DEC. */
#else
    dl_write8(m,m->cs,3370,st_dec(m,dl_read8(m,m->cs,3370),8));
#endif
L_14FED: ; /* IDA 0x14FED; call    sub_15074 */
    oc_call(m,20596,20464,h);
L_14FF0: ; /* IDA 0x14FF0; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_14FF2: ; /* IDA 0x14FF2; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14FF4: ; /* IDA 0x14FF4; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14FF6: ; /* IDA 0x14FF6; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_14FF8: ; /* IDA 0x14FF8; mov     ah, cs:byte_12919 */
    vg_ah(m,(uint8_t)(dl_read8(m,m->cs,10521)));
L_14FFD: ; /* IDA 0x14FFD; call    sub_14236 */
    oc_call(m,16950,20480,h);
L_15000: ; /* IDA 0x15000; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_15002: ; /* IDA 0x15002; mov     di, 4240h */
    m->di=(uint16_t)(16960);
L_15005: ; /* IDA 0x15005; mov     ch, [bx+1] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_15008: ; /* IDA 0x15008; xor     cl, cl */
    {uint16_t result=(((uint8_t)m->cx)^((uint8_t)m->cx))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_1500A: ; /* IDA 0x1500A; cmp     byte ptr [di], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di)),128,8);
L_1500D: ; /* IDA 0x1500D; jb      short loc_15033 */
    if(m->flags&1) goto L_15033;
L_1500F: ; /* IDA 0x1500F; cmp     al, [di+1Ch] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->di+28)),8);
L_15012: ; /* IDA 0x15012; jnz     short loc_15033 */
    if(!(m->flags&0x40)) goto L_15033;
L_15014: ; /* IDA 0x15014; cmp     byte ptr [di+1Dh], 0FFh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+29)),255,8);
L_15018: ; /* IDA 0x15018; jz      short loc_1501F */
    if(m->flags&0x40) goto L_1501F;
L_1501A: ; /* IDA 0x1501A; call    sub_150D7 */
    oc_call(m,20695,20509,h);
L_1501D: ; /* IDA 0x1501D; jmp     short loc_15033 */
    goto L_15033;
L_1501F: ; /* IDA 0x1501F; cmp     ch, cl */
    ec_cmp(m,(uint8_t)(m->cx>>8),(uint8_t)m->cx,8);
L_15021: ; /* IDA 0x15021; jz      short loc_1502E */
    if(m->flags&0x40) goto L_1502E;
L_15023: ; /* IDA 0x15023; cmp     byte ptr [di+17h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+23)),0,8);
L_15027: ; /* IDA 0x15027; jz      short loc_1502E */
    if(m->flags&0x40) goto L_1502E;
L_15029: ; /* IDA 0x15029; call    sub_150B4 */
    oc_call(m,20660,20524,h);
L_1502C: ; /* IDA 0x1502C; jmp     short loc_15033 */
    goto L_15033;
L_1502E: ; /* IDA 0x1502E; mov     bx, di */
    m->bx=(uint16_t)(m->di);
L_15030: ; /* IDA 0x15030; call    sub_129C3 */
    oc_call(m,10691,20531,h);
L_15033: ; /* IDA 0x15033; add     di, 20h ; ' ' */
    m->di=(uint16_t)(ec_math(m,m->di,32,16,0,0));
L_15036: ; /* IDA 0x15036; inc     cl */
    vg_cl(m,(uint8_t)(ec_inc(m,(uint8_t)m->cx,8)));
L_15038: ; /* IDA 0x15038; cmp     cl, 7Fh */
    ec_cmp(m,(uint8_t)m->cx,127,8);
L_1503B: ; /* IDA 0x1503B; jb      short loc_1500A */
    if(m->flags&1) goto L_1500A;
L_1503D: ; /* IDA 0x1503D; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_1503F: ; /* IDA 0x1503F; push    ax */
    ec_push(m,m->ax);
L_15040: ; /* IDA 0x15040; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_15042: ; /* IDA 0x15042; call    sub_10CDE */
    oc_call(m,3294,20549,h);
L_15045: ; /* IDA 0x15045; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_15047: ; /* IDA 0x15047; mov     cx, 24h ; '$' */
    m->cx=(uint16_t)(36);
L_1504A: ; /* IDA 0x1504A; call    sub_18810 */
    oc_call(m,34832,20557,h);
L_1504D: ; /* IDA 0x1504D; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1504E: ; /* IDA 0x1504E; mov     di, 0 */
    m->di=(uint16_t)(0);
L_15051: ; /* IDA 0x15051; mov     ch, al */
    vg_ch(m,(uint8_t)((uint8_t)m->ax));
L_15053: ; /* IDA 0x15053; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_15055: ; /* IDA 0x15055; cmp     byte ptr [di], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di)),128,8);
L_15058: ; /* IDA 0x15058; jb      short loc_15066 */
    if(m->flags&1) goto L_15066;
L_1505A: ; /* IDA 0x1505A; cmp     ch, [di+19h] */
    ec_cmp(m,(uint8_t)(m->cx>>8),dl_read8(m,m->ds,(uint16_t)(m->di+25)),8);
L_1505D: ; /* IDA 0x1505D; jnz     short loc_15066 */
    if(!(m->flags&0x40)) goto L_15066;
L_1505F: ; /* IDA 0x1505F; mov     byte ptr [di+19h], 0FFh */
    dl_write8(m,m->ds,(uint16_t)(m->di+25),255);
L_15063: ; /* IDA 0x15063; call    sub_13669 */
    oc_call(m,13929,20582,h);
L_15066: ; /* IDA 0x15066; add     di, 40h ; '@' */
    m->di=(uint16_t)(ec_math(m,m->di,64,16,0,0));
L_15069: ; /* IDA 0x15069; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_1506B: ; /* IDA 0x1506B; cmp     al, 16h */
    ec_cmp(m,(uint8_t)m->ax,22,8);
L_1506D: ; /* IDA 0x1506D; jb      short loc_15055 */
    if(m->flags&1) goto L_15055;
L_1506F: ; /* IDA 0x1506F; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_15070: ; /* IDA 0x15070; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_15071: ; /* IDA 0x15071; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15072: ; /* IDA 0x15072; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15073: ; /* IDA 0x15073; retn */
    return;
}
/* Original sub_15074; IDA 0x15074. */
static void oc_body_15074(KiMachine16 *m,const KiEconomyHooks *h) {
L_15074: ; /* IDA 0x15074; cmp     byte ptr [bx+2Ah], 0FFh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+42)),255,8);
L_15078: ; /* IDA 0x15078; jz      short locret_150B3 */
    if(m->flags&0x40) goto L_150B3;
L_1507A: ; /* IDA 0x1507A; push    ax */
    ec_push(m,m->ax);
L_1507B: ; /* IDA 0x1507B; push    bx */
    ec_push(m,m->bx);
L_1507C: ; /* IDA 0x1507C; push    cx */
    ec_push(m,m->cx);
L_1507D: ; /* IDA 0x1507D; push    di */
    ec_push(m,m->di);
L_1507E: ; /* IDA 0x1507E; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_15080: ; /* IDA 0x15080; mov     bh, 0FFh */
    gw_bh(m,(uint8_t)(255));
L_15082: ; /* IDA 0x15082; xchg    bh, [di+2Ah] */
    {uint16_t a=(uint8_t)(m->bx>>8),b=dl_read8(m,m->ds,(uint16_t)(m->di+42));gw_bh(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->di+42),a);}
L_15085: ; /* IDA 0x15085; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_15087: ; /* IDA 0x15087; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_15089: ; /* IDA 0x15089; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1508B: ; /* IDA 0x1508B; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1508D: ; /* IDA 0x1508D; add     bx, 4240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16960,16,0,0));
L_15091: ; /* IDA 0x15091; mov     byte ptr [bx+17h], 0 */
#if KI_OUTCOME_MUTATION == 18
    dl_write8(m,m->ds,(uint16_t)(m->bx+0x17),1);
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),0);
#endif
L_15095: ; /* IDA 0x15095; push    bx */
    ec_push(m,m->bx);
L_15096: ; /* IDA 0x15096; push    di */
    ec_push(m,m->di);
L_15097: ; /* IDA 0x15097; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_15099: ; /* IDA 0x15099; mov     cx, 45h ; 'E' */
    m->cx=(uint16_t)(69);
L_1509C: ; /* IDA 0x1509C; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_1509E: ; /* IDA 0x1509E; call    sub_18810 */
    oc_call(m,34832,20641,h);
L_150A1: ; /* IDA 0x150A1; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_150A2: ; /* IDA 0x150A2; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_150A3: ; /* IDA 0x150A3; mov     cx, 1A7h */
    m->cx=(uint16_t)(423);
L_150A6: ; /* IDA 0x150A6; mov     ah, [bx+1Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
L_150A9: ; /* IDA 0x150A9; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_150AC: ; /* IDA 0x150AC; call    sub_18810 */
    oc_call(m,34832,20655,h);
L_150AF: ; /* IDA 0x150AF; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_150B0: ; /* IDA 0x150B0; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_150B1: ; /* IDA 0x150B1; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_150B2: ; /* IDA 0x150B2; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_150B3: ; /* IDA 0x150B3; retn */
    return;
}
/* Original sub_150B4; IDA 0x150B4. */
static void oc_body_150B4(KiMachine16 *m,const KiEconomyHooks *h) {
L_150B4: ; /* IDA 0x150B4; push    si */
    ec_push(m,m->si);
L_150B5: ; /* IDA 0x150B5; cmp     byte ptr [di+17h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+23)),0,8);
L_150B9: ; /* IDA 0x150B9; jz      short loc_150D1 */
    if(m->flags&0x40) goto L_150D1;
L_150BB: ; /* IDA 0x150BB; mov     si, di */
    m->si=(uint16_t)(m->di);
L_150BD: ; /* IDA 0x150BD; sub     si, 4240h */
    m->si=(uint16_t)(ec_math(m,m->si,16960,16,0,1));
L_150C1: ; /* IDA 0x150C1; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_150C3: ; /* IDA 0x150C3; add     si, 2240h */
    m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,0));
L_150C7: ; /* IDA 0x150C7; call    sub_12BA8 */
    oc_call(m,11176,20682,h);
L_150CA: ; /* IDA 0x150CA; call    sub_17028 */
    oc_call(m,28712,20685,h);
L_150CD: ; /* IDA 0x150CD; mov     byte ptr [di+17h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->di+23),0);
L_150D1: ; /* IDA 0x150D1; mov     byte ptr [di+1Ch], 0FFh */
    dl_write8(m,m->ds,(uint16_t)(m->di+28),255);
L_150D5: ; /* IDA 0x150D5; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_150D6: ; /* IDA 0x150D6; retn */
    return;
}
/* Original sub_15130; IDA 0x15130. */
static void oc_body_15130(KiMachine16 *m,const KiEconomyHooks *h) {
L_15130: ; /* IDA 0x15130; push    bx */
    ec_push(m,m->bx);
L_15131: ; /* IDA 0x15131; push    cx */
    ec_push(m,m->cx);
L_15132: ; /* IDA 0x15132; push    dx */
    ec_push(m,m->dx);
L_15133: ; /* IDA 0x15133; push    bp */
    ec_push(m,m->bp);
L_15134: ; /* IDA 0x15134; sub     sp, 10h */
    m->sp=(uint16_t)(ec_math(m,m->sp,16,16,0,1));
L_15137: ; /* IDA 0x15137; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_15139: ; /* IDA 0x15139; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_1513B: ; /* IDA 0x1513B; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1513D: ; /* IDA 0x1513D; jnz     short loc_15141 */
    if(!(m->flags&0x40)) goto L_15141;
L_1513F: ; /* IDA 0x1513F; mov     al, 3 */
#if KI_OUTCOME_MUTATION == 1
    vg_al(m,2);
#else
    vg_al(m,(uint8_t)(3));
#endif
L_15141: ; /* IDA 0x15141; call    sub_15285 */
    oc_call(m,21125,20804,h);
L_15144: ; /* IDA 0x15144; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_15146: ; /* IDA 0x15146; call    sub_152D7 */
    oc_call(m,21207,20809,h);
L_15149: ; /* IDA 0x15149; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_1514B: ; /* IDA 0x1514B; xchg    si, di */
    {uint16_t a=m->si,b=m->di;m->si=(uint16_t)(b);m->di=(uint16_t)(a);}
L_1514D: ; /* IDA 0x1514D; call    sub_15285 */
    oc_call(m,21125,20816,h);
L_15150: ; /* IDA 0x15150; call    sub_152D7 */
    oc_call(m,21207,20819,h);
L_15153: ; /* IDA 0x15153; xchg    si, di */
    {uint16_t a=m->si,b=m->di;m->si=(uint16_t)(b);m->di=(uint16_t)(a);}
L_15155: ; /* IDA 0x15155; add     cx, 8 */
#if KI_OUTCOME_MUTATION == 5
    m->cx=ec_math(m,m->cx,7,16,0,0);
#else
    m->cx=(uint16_t)(ec_math(m,m->cx,8,16,0,0));
#endif
L_15158: ; /* IDA 0x15158; add     dx, 8 */
#if KI_OUTCOME_MUTATION == 5
    m->dx=ec_math(m,m->dx,7,16,0,0);
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,0));
#endif
L_1515B: ; /* IDA 0x1515B; cmp     cx, dx */
    ec_cmp(m,m->cx,m->dx,16);
L_1515D: ; /* IDA 0x1515D; pushf */
    ec_push(m,m->flags);
L_1515E: ; /* IDA 0x1515E; jnb     short loc_15162 */
    if(!(m->flags&1)) goto L_15162;
L_15160: ; /* IDA 0x15160; xchg    dx, cx */
    {uint16_t a=m->dx,b=m->cx;m->dx=(uint16_t)(b);m->cx=(uint16_t)(a);}
L_15162: ; /* IDA 0x15162; mov     ax, cx */
    m->ax=(uint16_t)(m->cx);
L_15164: ; /* IDA 0x15164; mov     bx, dx */
    m->bx=(uint16_t)(m->dx);
L_15166: ; /* IDA 0x15166; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_15168: ; /* IDA 0x15168; mov     cx, 3 */
    m->cx=(uint16_t)(3);
L_1516B: ; /* IDA 0x1516B; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1516D: ; /* IDA 0x1516D; rcl     dx, 1 */
    m->dx=(uint16_t)(tx_rcl(m,m->dx,16));
L_1516F: ; /* IDA 0x1516F; loop    loc_1516B */
    if(--m->cx) goto L_1516B;
L_15171: ; /* IDA 0x15171; div     bx */
    dl_div(m,m->bx,16);
L_15173: ; /* IDA 0x15173; cmp     ax, 64h ; 'd' */
#if KI_OUTCOME_MUTATION == 6
    ec_cmp(m,m->ax,99,16);
#else
    ec_cmp(m,m->ax,100,16);
#endif
L_15176: ; /* IDA 0x15176; jb      short loc_1517B */
    if(m->flags&1) goto L_1517B;
L_15178: ; /* IDA 0x15178; mov     ax, 64h ; 'd' */
#if KI_OUTCOME_MUTATION == 6
    m->ax=99;
#else
    m->ax=(uint16_t)(100);
#endif
L_1517B: ; /* IDA 0x1517B; popf */
    m->flags=ec_pop(m)|2;
L_1517C: ; /* IDA 0x1517C; jnb     short loc_15189 */
    if(!(m->flags&1)) goto L_15189;
L_1517E: ; /* IDA 0x1517E; xchg    si, di */
    {uint16_t a=m->si,b=m->di;m->si=(uint16_t)(b);m->di=(uint16_t)(a);}
L_15180: ; /* IDA 0x15180; call    sub_151B3 */
    oc_call(m,20915,20867,h);
L_15183: ; /* IDA 0x15183; xchg    si, di */
    {uint16_t a=m->si,b=m->di;m->si=(uint16_t)(b);m->di=(uint16_t)(a);}
L_15185: ; /* IDA 0x15185; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_15187: ; /* IDA 0x15187; jmp     short loc_1518E */
    goto L_1518E;
L_15189: ; /* IDA 0x15189; call    sub_151B3 */
    oc_call(m,20915,20876,h);
L_1518C: ; /* IDA 0x1518C; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1518E: ; /* IDA 0x1518E; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_15190: ; /* IDA 0x15190; mov     cl, al */
    vg_cl(m,(uint8_t)((uint8_t)m->ax));
L_15192: ; /* IDA 0x15192; call    sub_1474A */
    oc_call(m,18250,20885,h);
L_15195: ; /* IDA 0x15195; jnb     short loc_1519A */
    if(!(m->flags&1)) goto L_1519A;
L_15197: ; /* IDA 0x15197; or      ah, 1 */
    {uint16_t result=(((uint8_t)(m->ax>>8))|(1))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1519A: ; /* IDA 0x1519A; xchg    si, di */
    {uint16_t a=m->si,b=m->di;m->si=(uint16_t)(b);m->di=(uint16_t)(a);}
L_1519C: ; /* IDA 0x1519C; mov     cl, al */
    vg_cl(m,(uint8_t)((uint8_t)m->ax));
L_1519E: ; /* IDA 0x1519E; xor     cl, 1 */
    {uint16_t result=(((uint8_t)m->cx)^(1))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_151A1: ; /* IDA 0x151A1; call    sub_1474A */
    oc_call(m,18250,20900,h);
L_151A4: ; /* IDA 0x151A4; jnb     short loc_151A9 */
    if(!(m->flags&1)) goto L_151A9;
L_151A6: ; /* IDA 0x151A6; or      ah, 2 */
    {uint16_t result=(((uint8_t)(m->ax>>8))|(2))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_151A9: ; /* IDA 0x151A9; xchg    si, di */
    {uint16_t a=m->si,b=m->di;m->si=(uint16_t)(b);m->di=(uint16_t)(a);}
L_151AB: ; /* IDA 0x151AB; add     sp, 10h */
    m->sp=(uint16_t)(ec_math(m,m->sp,16,16,0,0));
L_151AE: ; /* IDA 0x151AE; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_151AF: ; /* IDA 0x151AF; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_151B0: ; /* IDA 0x151B0; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_151B1: ; /* IDA 0x151B1; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_151B2: ; /* IDA 0x151B2; retn */
    return;
}
/* Original sub_151B3; IDA 0x151B3. */
static void oc_body_151B3(KiMachine16 *m,const KiEconomyHooks *h) {
L_151B3: ; /* IDA 0x151B3; push    ax */
    ec_push(m,m->ax);
L_151B4: ; /* IDA 0x151B4; push    bx */
    ec_push(m,m->bx);
L_151B5: ; /* IDA 0x151B5; push    cx */
    ec_push(m,m->cx);
L_151B6: ; /* IDA 0x151B6; push    dx */
    ec_push(m,m->dx);
L_151B7: ; /* IDA 0x151B7; push    si */
    ec_push(m,m->si);
L_151B8: ; /* IDA 0x151B8; push    di */
    ec_push(m,m->di);
L_151B9: ; /* IDA 0x151B9; mov     dl, cs:byte_10D34 */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,3380)));
L_151BE: ; /* IDA 0x151BE; cmp     dl, 0C0h */
    ec_cmp(m,(uint8_t)m->dx,192,8);
L_151C1: ; /* IDA 0x151C1; jnb     short loc_151EB */
    if(!(m->flags&1)) goto L_151EB;
L_151C3: ; /* IDA 0x151C3; mov     ah, 3Fh ; '?' */
    vg_ah(m,(uint8_t)(63));
L_151C5: ; /* IDA 0x151C5; sub     ah, al */
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),(uint8_t)m->ax,8,0,1)));
L_151C7: ; /* IDA 0x151C7; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_151C9: ; /* IDA 0x151C9; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_151CB: ; /* IDA 0x151CB; mov     bx, cs:word_10D32 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,3378));
L_151D0: ; /* IDA 0x151D0; sub     [bx+13h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+19),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->bx+19)),(uint8_t)(m->ax>>8),8,0,1));
L_151D3: ; /* IDA 0x151D3; jnb     short loc_151D9 */
    if(!(m->flags&1)) goto L_151D9;
L_151D5: ; /* IDA 0x151D5; mov     byte ptr [bx+13h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->bx+19),0);
L_151D9: ; /* IDA 0x151D9; sub     [bx+10h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+16),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->bx+16)),(uint8_t)(m->ax>>8),8,0,1));
L_151DC: ; /* IDA 0x151DC; jnb     short loc_151E2 */
    if(!(m->flags&1)) goto L_151E2;
L_151DE: ; /* IDA 0x151DE; mov     byte ptr [bx+10h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->bx+16),0);
L_151E2: ; /* IDA 0x151E2; sub     [bx+11h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+17),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->bx+17)),(uint8_t)(m->ax>>8),8,0,1));
L_151E5: ; /* IDA 0x151E5; jnb     short loc_151EB */
    if(!(m->flags&1)) goto L_151EB;
L_151E7: ; /* IDA 0x151E7; mov     byte ptr [bx+11h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->bx+17),0);
L_151EB: ; /* IDA 0x151EB; add     si, 28h ; '(' */
    m->si=(uint16_t)(ec_math(m,m->si,40,16,0,0));
L_151EE: ; /* IDA 0x151EE; add     di, 28h ; '(' */
    m->di=(uint16_t)(ec_math(m,m->di,40,16,0,0));
L_151F1: ; /* IDA 0x151F1; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_151F3: ; /* IDA 0x151F3; mov     dx, bx */
    m->dx=(uint16_t)(m->bx);
L_151F5: ; /* IDA 0x151F5; mov     ch, al */
    vg_ch(m,(uint8_t)((uint8_t)m->ax));
L_151F7: ; /* IDA 0x151F7; mov     cl, 6 */
    vg_cl(m,(uint8_t)(6));
L_151F9: ; /* IDA 0x151F9; call    sub_1ECE0 */
    oc_call(m,60640,20988,h);
L_151FC: ; /* IDA 0x151FC; and     al, 7 */
    {uint16_t result=(((uint8_t)m->ax)&(7))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_151FE: ; /* IDA 0x151FE; add     al, 2 */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,2,8,0,0)));
L_15200: ; /* IDA 0x15200; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_15203: ; /* IDA 0x15203; sub     ah, al */
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),(uint8_t)m->ax,8,0,1)));
L_15205: ; /* IDA 0x15205; ja      short loc_15210 */
    if(!(m->flags&(1|0x40))) goto L_15210;
L_15207: ; /* IDA 0x15207; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_15209: ; /* IDA 0x15209; cmp     cl, 6 */
    ec_cmp(m,(uint8_t)m->cx,6,8);
L_1520C: ; /* IDA 0x1520C; jnz     short loc_15210 */
    if(!(m->flags&0x40)) goto L_15210;
L_1520E: ; /* IDA 0x1520E; mov     ah, 1 */
#if KI_OUTCOME_MUTATION == 8
    vg_ah(m,0);
#else
    vg_ah(m,(uint8_t)(1));
#endif
L_15210: ; /* IDA 0x15210; add     bl, ah */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,(uint8_t)(m->ax>>8),8,0,0)));
L_15212: ; /* IDA 0x15212; adc     bh, 0 */
    gw_bh(m,(uint8_t)(ec_math(m,(uint8_t)(m->bx>>8),0,8,m->flags&1,0)));
L_15215: ; /* IDA 0x15215; mov     [si+1], ah */
    dl_write8(m,m->ds,(uint16_t)(m->si+1),(uint8_t)(m->ax>>8));
L_15218: ; /* IDA 0x15218; call    sub_1ECE0 */
    oc_call(m,60640,21019,h);
L_1521B: ; /* IDA 0x1521B; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1521D: ; /* IDA 0x1521D; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_1521F: ; /* IDA 0x1521F; div     ch */
    dl_div(m,(uint8_t)(m->cx>>8),8);
L_15221: ; /* IDA 0x15221; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_15223: ; /* IDA 0x15223; add     al, 8 */
#if KI_OUTCOME_MUTATION == 7
    vg_al(m,ec_math(m,(uint8_t)m->ax,7,8,0,0));
#else
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,8,8,0,0)));
#endif
L_15225: ; /* IDA 0x15225; mov     ah, [di+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+1))));
L_15228: ; /* IDA 0x15228; sub     ah, al */
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),(uint8_t)m->ax,8,0,1)));
L_1522A: ; /* IDA 0x1522A; ja      short loc_15235 */
    if(!(m->flags&(1|0x40))) goto L_15235;
L_1522C: ; /* IDA 0x1522C; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1522E: ; /* IDA 0x1522E; cmp     cl, 6 */
    ec_cmp(m,(uint8_t)m->cx,6,8);
L_15231: ; /* IDA 0x15231; jnz     short loc_15235 */
    if(!(m->flags&0x40)) goto L_15235;
L_15233: ; /* IDA 0x15233; mov     ah, 1 */
#if KI_OUTCOME_MUTATION == 8
    vg_ah(m,0);
#else
    vg_ah(m,(uint8_t)(1));
#endif
L_15235: ; /* IDA 0x15235; add     dl, ah */
    vg_dl(m,(uint8_t)(ec_math(m,(uint8_t)m->dx,(uint8_t)(m->ax>>8),8,0,0)));
L_15237: ; /* IDA 0x15237; adc     dh, 0 */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),0,8,m->flags&1,0)));
L_1523A: ; /* IDA 0x1523A; mov     [di+1], ah */
    dl_write8(m,m->ds,(uint16_t)(m->di+1),(uint8_t)(m->ax>>8));
L_1523D: ; /* IDA 0x1523D; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_15240: ; /* IDA 0x15240; add     di, 4 */
    m->di=(uint16_t)(ec_math(m,m->di,4,16,0,0));
L_15243: ; /* IDA 0x15243; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_15245: ; /* IDA 0x15245; jnz     short loc_151F9 */
    if(!(m->flags&0x40)) goto L_151F9;
L_15247: ; /* IDA 0x15247; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_15248: ; /* IDA 0x15248; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_15249: ; /* IDA 0x15249; xchg    bx, [si+4] */
    {uint16_t a=m->bx,b=dl_read16(m,m->ds,(uint16_t)(m->si+4));m->bx=(uint16_t)(b);dl_write16(m,m->ds,(uint16_t)(m->si+4),a);}
L_1524C: ; /* IDA 0x1524C; xchg    dx, [di+4] */
    {uint16_t a=m->dx,b=dl_read16(m,m->ds,(uint16_t)(m->di+4));m->dx=(uint16_t)(b);dl_write16(m,m->ds,(uint16_t)(m->di+4),a);}
L_1524F: ; /* IDA 0x1524F; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_15251: ; /* IDA 0x15251; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_15253: ; /* IDA 0x15253; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_15255: ; /* IDA 0x15255; jz      short loc_15265 */
    if(m->flags&0x40) goto L_15265;
L_15257: ; /* IDA 0x15257; cmp     byte ptr [si+6], 64h ; 'd' */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+6)),100,8);
L_1525B: ; /* IDA 0x1525B; jb      short loc_15265 */
    if(m->flags&1) goto L_15265;
L_1525D: ; /* IDA 0x1525D; mov     al, [si+6] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+6))));
L_15260: ; /* IDA 0x15260; mul     word ptr [si+4] */
    dl_mul(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),16);
L_15263: ; /* IDA 0x15263; div     bx */
    dl_div(m,m->bx,16);
L_15265: ; /* IDA 0x15265; mov     [si+6], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+6),(uint8_t)m->ax);
L_15268: ; /* IDA 0x15268; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_1526A: ; /* IDA 0x1526A; and     cx, cx */
    {uint16_t result=((m->cx)&(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_1526C: ; /* IDA 0x1526C; jz      short loc_1527D */
    if(m->flags&0x40) goto L_1527D;
L_1526E: ; /* IDA 0x1526E; cmp     byte ptr [di+6], 64h ; 'd' */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+6)),100,8);
L_15272: ; /* IDA 0x15272; jb      short loc_1527D */
    if(m->flags&1) goto L_1527D;
L_15274: ; /* IDA 0x15274; mov     al, 64h ; 'd' */
    vg_al(m,(uint8_t)(100));
L_15276: ; /* IDA 0x15276; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_15278: ; /* IDA 0x15278; mul     word ptr [di+4] */
    dl_mul(m,dl_read16(m,m->ds,(uint16_t)(m->di+4)),16);
L_1527B: ; /* IDA 0x1527B; div     cx */
    dl_div(m,m->cx,16);
L_1527D: ; /* IDA 0x1527D; mov     [di+6], al */
    dl_write8(m,m->ds,(uint16_t)(m->di+6),(uint8_t)m->ax);
L_15280: ; /* IDA 0x15280; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15281: ; /* IDA 0x15281; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_15282: ; /* IDA 0x15282; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15283: ; /* IDA 0x15283; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15284: ; /* IDA 0x15284; retn */
    return;
}
/* Original sub_15285; IDA 0x15285. */
static void oc_body_15285(KiMachine16 *m,const KiEconomyHooks *h) {
L_15285: ; /* IDA 0x15285; push    ax */
    ec_push(m,m->ax);
L_15286: ; /* IDA 0x15286; push    bx */
    ec_push(m,m->bx);
L_15287: ; /* IDA 0x15287; push    cx */
    ec_push(m,m->cx);
L_15288: ; /* IDA 0x15288; push    di */
    ec_push(m,m->di);
L_15289: ; /* IDA 0x15289; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1528B: ; /* IDA 0x1528B; push    ax */
    ec_push(m,m->ax);
L_1528C: ; /* IDA 0x1528C; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1528E: ; /* IDA 0x1528E; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_15290: ; /* IDA 0x15290; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_15292: ; /* IDA 0x15292; add     si, 28h ; '(' */
    m->si=(uint16_t)(ec_math(m,m->si,40,16,0,0));
L_15295: ; /* IDA 0x15295; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_15297: ; /* IDA 0x15297; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_1529A: ; /* IDA 0x1529A; mov     bl, [si+2] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_1529D: ; /* IDA 0x1529D; dec     bl */
    gw_bl(m,(uint8_t)(st_dec(m,(uint8_t)m->bx,8)));
L_1529F: ; /* IDA 0x1529F; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_152A1: ; /* IDA 0x152A1; mov     al, [si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_152A4: ; /* IDA 0x152A4; mul     byte ptr cs:[bx+di+5120h] */
    dl_mul(m,dl_read8(m,m->cs,(uint16_t)(m->bx+m->di+20768)),8);
L_152A9: ; /* IDA 0x152A9; add     dx, ax */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->ax,16,0,0));
L_152AB: ; /* IDA 0x152AB; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_152AE: ; /* IDA 0x152AE; loop    loc_1529A */
    if(--m->cx) goto L_1529A;
L_152B0: ; /* IDA 0x152B0; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_152B1: ; /* IDA 0x152B1; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_152B3: ; /* IDA 0x152B3; jnz     short loc_152C0 */
    if(!(m->flags&0x40)) goto L_152C0;
L_152B5: ; /* IDA 0x152B5; mov     bx, cs:word_10D32 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,3378));
L_152BA: ; /* IDA 0x152BA; add     dl, [bx+13h] */
    vg_dl(m,(uint8_t)(ec_math(m,(uint8_t)m->dx,dl_read8(m,m->ds,(uint16_t)(m->bx+19)),8,0,0)));
L_152BD: ; /* IDA 0x152BD; adc     dh, 0 */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),0,8,m->flags&1,0)));
L_152C0: ; /* IDA 0x152C0; sub     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,1));
L_152C3: ; /* IDA 0x152C3; mov     al, [si+6] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+6))));
L_152C6: ; /* IDA 0x152C6; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_152C8: ; /* IDA 0x152C8; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_152CA: ; /* IDA 0x152CA; shr     al, 1 */
#if KI_OUTCOME_MUTATION == 2
    ; /* Omitted third morale SHR, giving division by four. */
#else
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
#endif
L_152CC: ; /* IDA 0x152CC; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_152CE: ; /* IDA 0x152CE; mul     dx */
    dl_mul(m,m->dx,16);
L_152D0: ; /* IDA 0x152D0; mov     dx, ax */
    m->dx=(uint16_t)(m->ax);
L_152D2: ; /* IDA 0x152D2; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_152D3: ; /* IDA 0x152D3; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_152D4: ; /* IDA 0x152D4; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_152D5: ; /* IDA 0x152D5; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_152D6: ; /* IDA 0x152D6; retn */
    return;
}
/* Original sub_152D7; IDA 0x152D7. */
static void oc_body_152D7(KiMachine16 *m,const KiEconomyHooks *h) {
L_152D7: ; /* IDA 0x152D7; push    ax */
    ec_push(m,m->ax);
L_152D8: ; /* IDA 0x152D8; push    bx */
    ec_push(m,m->bx);
L_152D9: ; /* IDA 0x152D9; push    cx */
    ec_push(m,m->cx);
L_152DA: ; /* IDA 0x152DA; push    dx */
    ec_push(m,m->dx);
L_152DB: ; /* IDA 0x152DB; mov     bh, [si+2] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_152DE: ; /* IDA 0x152DE; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_152E0: ; /* IDA 0x152E0; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_152E2: ; /* IDA 0x152E2; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_152E4: ; /* IDA 0x152E4; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_152E6: ; /* IDA 0x152E6; push    ax */
    ec_push(m,m->ax);
L_152E7: ; /* IDA 0x152E7; mov     dx, [bx+4251h] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+16977)));
L_152EB: ; /* IDA 0x152EB; mov     ch, dl */
    vg_ch(m,(uint8_t)((uint8_t)m->dx));
L_152ED: ; /* IDA 0x152ED; shl     ch, 1 */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),1,8,4)));
L_152EF: ; /* IDA 0x152EF; cmp     dl, dh */
    ec_cmp(m,(uint8_t)m->dx,(uint8_t)(m->dx>>8),8);
L_152F1: ; /* IDA 0x152F1; jb      short loc_152FC */
    if(m->flags&1) goto L_152FC;
L_152F3: ; /* IDA 0x152F3; call    sub_1ECE0 */
    oc_call(m,60640,21238,h);
L_152F6: ; /* IDA 0x152F6; and     al, 3 */
#if KI_OUTCOME_MUTATION == 4
    {uint8_t result=(uint8_t)m->ax&1;ec_logic(m,result,8);vg_al(m,result);}
#else
    {uint16_t result=(((uint8_t)m->ax)&(3))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
#endif
L_152F8: ; /* IDA 0x152F8; jnz     short loc_15308 */
    if(!(m->flags&0x40)) goto L_15308;
L_152FA: ; /* IDA 0x152FA; mov     dh, dl */
    vg_dh(m,(uint8_t)((uint8_t)m->dx));
L_152FC: ; /* IDA 0x152FC; mov     ch, dh */
    vg_ch(m,(uint8_t)((uint8_t)(m->dx>>8)));
L_152FE: ; /* IDA 0x152FE; shr     dh, 1 */
    vg_dh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->dx>>8),1,8,5)));
L_15300: ; /* IDA 0x15300; shr     dh, 1 */
    vg_dh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->dx>>8),1,8,5)));
L_15302: ; /* IDA 0x15302; sub     ch, dh */
    vg_ch(m,(uint8_t)(ec_math(m,(uint8_t)(m->cx>>8),(uint8_t)(m->dx>>8),8,0,1)));
L_15304: ; /* IDA 0x15304; add     ch, [bx+4252h] */
    vg_ch(m,(uint8_t)(ec_math(m,(uint8_t)(m->cx>>8),dl_read8(m,m->ds,(uint16_t)(m->bx+16978)),8,0,0)));
L_15308: ; /* IDA 0x15308; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15309: ; /* IDA 0x15309; add     bx, 424Eh */
    m->bx=(uint16_t)(ec_math(m,m->bx,16974,16,0,0));
L_1530D: ; /* IDA 0x1530D; xlat */
    vg_al(m,dl_read8(m,m->ds,(uint16_t)(m->bx+(uint8_t)m->ax)));
L_1530E: ; /* IDA 0x1530E; mov     cl, 4 */
#if KI_OUTCOME_MUTATION == 3
    vg_cl(m,3);
#else
    vg_cl(m,(uint8_t)(4));
#endif
L_15310: ; /* IDA 0x15310; shr     al, cl */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,(uint8_t)m->cx,8,5)));
L_15312: ; /* IDA 0x15312; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_15314: ; /* IDA 0x15314; mov     bl, 10h */
    gw_bl(m,(uint8_t)(16));
L_15316: ; /* IDA 0x15316; sub     bl, al */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,(uint8_t)m->ax,8,0,1)));
L_15318: ; /* IDA 0x15318; mov     al, ch */
    vg_al(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_1531A: ; /* IDA 0x1531A; mov     ah, bh */
    vg_ah(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_1531C: ; /* IDA 0x1531C; shl     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,4));
L_1531E: ; /* IDA 0x1531E; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_15320: ; /* IDA 0x15320; div     bx */
    dl_div(m,m->bx,16);
L_15322: ; /* IDA 0x15322; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_15324: ; /* IDA 0x15324; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15325: ; /* IDA 0x15325; mul     dx */
    dl_mul(m,m->dx,16);
L_15327: ; /* IDA 0x15327; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_15329: ; /* IDA 0x15329; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_1532B: ; /* IDA 0x1532B; mov     dl, dh */
    vg_dl(m,(uint8_t)((uint8_t)(m->dx>>8)));
L_1532D: ; /* IDA 0x1532D; xor     dh, dh */
    {uint16_t result=(((uint8_t)(m->dx>>8))^((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);vg_dh(m,(uint8_t)(result));}
L_1532F: ; /* IDA 0x1532F; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15331: ; /* IDA 0x15331; rcr     ax, 1 */
    m->ax=(uint16_t)(hr_rcr(m,m->ax));
L_15333: ; /* IDA 0x15333; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15335: ; /* IDA 0x15335; rcr     ax, 1 */
    m->ax=(uint16_t)(hr_rcr(m,m->ax));
L_15337: ; /* IDA 0x15337; mov     dx, ax */
    m->dx=(uint16_t)(m->ax);
L_15339: ; /* IDA 0x15339; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1533A: ; /* IDA 0x1533A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1533B: ; /* IDA 0x1533B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1533C: ; /* IDA 0x1533C; retn */
    return;
}
/* Original sub_15CA4; IDA 0x15CA4. */
static void oc_body_15CA4(KiMachine16 *m,const KiEconomyHooks *h) {
L_15CA4: ; /* IDA 0x15CA4; mov     dx, ds:9892h */
    m->dx=(uint16_t)(dl_read16(m,m->ds,39058));
L_15CA8: ; /* IDA 0x15CA8; mov     bx, ds:9894h */
    m->bx=(uint16_t)(dl_read16(m,m->ds,39060));
L_15CAC: ; /* IDA 0x15CAC; cmp     dx, 0FFFFh */
    ec_cmp(m,m->dx,65535,16);
L_15CAF: ; /* IDA 0x15CAF; jz      short locret_15CC5 */
    if(m->flags&0x40) goto L_15CC5;
L_15CB1: ; /* IDA 0x15CB1; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15CB3: ; /* IDA 0x15CB3; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_15CB5: ; /* IDA 0x15CB5; add     dx, 1B8h */
    m->dx=(uint16_t)(ec_math(m,m->dx,440,16,0,0));
L_15CB9: ; /* IDA 0x15CB9; add     bx, 28h ; '(' */
    m->bx=(uint16_t)(ec_math(m,m->bx,40,16,0,0));
L_15CBC: ; /* IDA 0x15CBC; call    sub_19541 */
    oc_call(m,38209,23743,h);
L_15CBF: ; /* IDA 0x15CBF; mov     word ptr ds:9892h, 0FFFFh */
    dl_write16(m,m->ds,39058,65535);
L_15CC5: ; /* IDA 0x15CC5; retn */
    return;
}
/* Original sub_15CE0; IDA 0x15CE0. */
static void oc_body_15CE0(KiMachine16 *m,const KiEconomyHooks *h) {
L_15CE0: ; /* IDA 0x15CE0; push    ax */
    ec_push(m,m->ax);
L_15CE1: ; /* IDA 0x15CE1; push    bx */
    ec_push(m,m->bx);
L_15CE2: ; /* IDA 0x15CE2; push    dx */
    ec_push(m,m->dx);
L_15CE3: ; /* IDA 0x15CE3; mov     bl, [si+1] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_15CE6: ; /* IDA 0x15CE6; mov     ah, 0Fh */
    vg_ah(m,(uint8_t)(15));
L_15CE8: ; /* IDA 0x15CE8; cmp     bl, 18h */
    ec_cmp(m,(uint8_t)m->bx,24,8);
L_15CEB: ; /* IDA 0x15CEB; jz      short loc_15D01 */
    if(m->flags&0x40) goto L_15D01;
L_15CED: ; /* IDA 0x15CED; mov     ah, 0ACh */
    vg_ah(m,(uint8_t)(172));
L_15CEF: ; /* IDA 0x15CEF; cmp     bl, cs:byte_10CFF */
    ec_cmp(m,(uint8_t)m->bx,dl_read8(m,m->cs,3327),8);
L_15CF4: ; /* IDA 0x15CF4; jz      short loc_15D01 */
    if(m->flags&0x40) goto L_15D01;
L_15CF6: ; /* IDA 0x15CF6; mov     ah, 0F3h */
    vg_ah(m,(uint8_t)(243));
L_15CF8: ; /* IDA 0x15CF8; cmp     bl, cs:byte_198A7 */
    ec_cmp(m,(uint8_t)m->bx,dl_read8(m,m->cs,39079),8);
L_15CFD: ; /* IDA 0x15CFD; jz      short loc_15D01 */
    if(m->flags&0x40) goto L_15D01;
L_15CFF: ; /* IDA 0x15CFF; mov     ah, 83h */
    vg_ah(m,(uint8_t)(131));
L_15D01: ; /* IDA 0x15D01; mov     dx, [si+8] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+8)));
L_15D04: ; /* IDA 0x15D04; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15D06: ; /* IDA 0x15D06; add     dx, 1B6h */
    m->dx=(uint16_t)(ec_math(m,m->dx,438,16,0,0));
L_15D0A: ; /* IDA 0x15D0A; mov     bx, [si+0Ah] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+10)));
L_15D0D: ; /* IDA 0x15D0D; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_15D0F: ; /* IDA 0x15D0F; add     bx, 26h ; '&' */
    m->bx=(uint16_t)(ec_math(m,m->bx,38,16,0,0));
L_15D12: ; /* IDA 0x15D12; call    sub_15D19 */
    oc_call(m,23833,23829,h);
L_15D15: ; /* IDA 0x15D15; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15D16: ; /* IDA 0x15D16; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15D17: ; /* IDA 0x15D17; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15D18: ; /* IDA 0x15D18; retn */
    return;
}
/* Original sub_17028; IDA 0x17028. */
static void oc_body_17028(KiMachine16 *m,const KiEconomyHooks *h) {
L_17028: ; /* IDA 0x17028; push    ds */
    ec_push(m,m->ds);
L_17029: ; /* IDA 0x17029; push    ax */
    ec_push(m,m->ax);
L_1702A: ; /* IDA 0x1702A; push    si */
    ec_push(m,m->si);
L_1702B: ; /* IDA 0x1702B; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1702D: ; /* IDA 0x1702D; xchg    al, [si] */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->ds,(uint16_t)(m->si));vg_al(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->si),a);}
L_1702F: ; /* IDA 0x1702F; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_17031: ; /* IDA 0x17031; jnz     short loc_17038 */
    if(!(m->flags&0x40)) goto L_17038;
L_17033: ; /* IDA 0x17033; lds     si, [si+1Ah] */
    {uint16_t off=(uint16_t)(m->si+26),seg=m->ds;uint16_t value=dl_read16(m,seg,off),bank=dl_read16(m,seg,(uint16_t)(off+2));m->si=(uint16_t)(value);m->ds=bank;}
L_17036: ; /* IDA 0x17036; dec     byte ptr [si] */
#if KI_OUTCOME_MUTATION == 20
    ; /* Omitted army occupancy DEC. */
#else
    dl_write8(m,m->ds,(uint16_t)(m->si),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->si)),8));
#endif
L_17038: ; /* IDA 0x17038; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17039: ; /* IDA 0x17039; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1703A: ; /* IDA 0x1703A; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1703B: ; /* IDA 0x1703B; retn */
    return;
}
/* Original sub_188CC; IDA 0x188CC. */
static void oc_body_188CC(KiMachine16 *m,const KiEconomyHooks *h) {
L_188CC: ; /* IDA 0x188CC; push    ax */
    ec_push(m,m->ax);
L_188CD: ; /* IDA 0x188CD; push    bx */
    ec_push(m,m->bx);
L_188CE: ; /* IDA 0x188CE; push    cx */
    ec_push(m,m->cx);
L_188CF: ; /* IDA 0x188CF; push    dx */
    ec_push(m,m->dx);
L_188D0: ; /* IDA 0x188D0; push    di */
    ec_push(m,m->di);
L_188D1: ; /* IDA 0x188D1; mov     dl, 1 */
    vg_dl(m,(uint8_t)(1));
L_188D3: ; /* IDA 0x188D3; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_188D5: ; /* IDA 0x188D5; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_188D7: ; /* IDA 0x188D7; sub     ax, 840h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2112,16,0,1));
L_188DA: ; /* IDA 0x188DA; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_188DC: ; /* IDA 0x188DC; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_188DE: ; /* IDA 0x188DE; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_188E0: ; /* IDA 0x188E0; mov     al, [si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_188E3: ; /* IDA 0x188E3; add     di, 1Ch */
    m->di=(uint16_t)(ec_math(m,m->di,28,16,0,0));
L_188E6: ; /* IDA 0x188E6; mov     cx, 4 */
    m->cx=(uint16_t)(4);
L_188E9: ; /* IDA 0x188E9; push    cx */
    ec_push(m,m->cx);
L_188EA: ; /* IDA 0x188EA; push    di */
    ec_push(m,m->di);
L_188EB: ; /* IDA 0x188EB; mov     bh, [di] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
L_188ED: ; /* IDA 0x188ED; cmp     bh, 0FFh */
    ec_cmp(m,(uint8_t)(m->bx>>8),255,8);
L_188F0: ; /* IDA 0x188F0; jz      short loc_188FF */
    if(m->flags&0x40) goto L_188FF;
L_188F2: ; /* IDA 0x188F2; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_188F4: ; /* IDA 0x188F4; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_188F6: ; /* IDA 0x188F6; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_188F8: ; /* IDA 0x188F8; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_188FA: ; /* IDA 0x188FA; call    sub_1890A */
    oc_call(m,35082,35069,h);
L_188FD: ; /* IDA 0x188FD; shl     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,4)));
L_188FF: ; /* IDA 0x188FF; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_18900: ; /* IDA 0x18900; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_18901: ; /* IDA 0x18901; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_18902: ; /* IDA 0x18902; loop    loc_188E9 */
    if(--m->cx) goto L_188E9;
L_18904: ; /* IDA 0x18904; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_18905: ; /* IDA 0x18905; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_18906: ; /* IDA 0x18906; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_18907: ; /* IDA 0x18907; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_18908: ; /* IDA 0x18908; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_18909: ; /* IDA 0x18909; retn */
    return;
}
/* Original sub_1890A; IDA 0x1890A. */
static void oc_body_1890A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1890A: ; /* IDA 0x1890A; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_1890C: ; /* IDA 0x1890C; add     di, 85Ch */
    m->di=(uint16_t)(ec_math(m,m->di,2140,16,0,0));
L_18910: ; /* IDA 0x18910; mov     dh, 1 */
    vg_dh(m,(uint8_t)(1));
L_18912: ; /* IDA 0x18912; mov     cx, 4 */
    m->cx=(uint16_t)(4);
L_18915: ; /* IDA 0x18915; cmp     ah, [di] */
    ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->ds,(uint16_t)(m->di)),8);
L_18917: ; /* IDA 0x18917; jz      short loc_1891F */
    if(m->flags&0x40) goto L_1891F;
L_18919: ; /* IDA 0x18919; shl     dh, 1 */
    vg_dh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->dx>>8),1,8,4)));
L_1891B: ; /* IDA 0x1891B; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1891C: ; /* IDA 0x1891C; loop    loc_18915 */
    if(--m->cx) goto L_18915;
L_1891E: ; /* IDA 0x1891E; retn */
    return;
L_1891F: ; /* IDA 0x1891F; cmp     al, [bx+841h] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+2113)),8);
L_18923: ; /* IDA 0x18923; jz      short loc_1893D */
    if(m->flags&0x40) goto L_1893D;
L_18925: ; /* IDA 0x18925; test    [bx+840h], dh */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx+2112)))&((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);}
L_18929: ; /* IDA 0x18929; jnz     short loc_18933 */
    if(!(m->flags&0x40)) goto L_18933;
L_1892B: ; /* IDA 0x1892B; inc     byte ptr [bx+85Bh] */
#if KI_OUTCOME_MUTATION == 16
    ; /* Omitted neighboring city-count INC. */
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+2139),ec_inc(m,dl_read8(m,m->ds,(uint16_t)(m->bx+2139)),8));
#endif
L_1892F: ; /* IDA 0x1892F; or      [bx+840h], dh */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx+2112)))|((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->bx+2112),result);}
L_18933: ; /* IDA 0x18933; test    [si], dl */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&((uint8_t)m->dx))&255;ec_logic(m,result,8);}
L_18935: ; /* IDA 0x18935; jnz     short locret_1893C */
    if(!(m->flags&0x40)) goto L_1893C;
L_18937: ; /* IDA 0x18937; inc     byte ptr [si+1Bh] */
    dl_write8(m,m->ds,(uint16_t)(m->si+27),ec_inc(m,dl_read8(m,m->ds,(uint16_t)(m->si+27)),8));
L_1893A: ; /* IDA 0x1893A; or      [si], dl */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|((uint8_t)m->dx))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_1893C: ; /* IDA 0x1893C; retn */
    return;
L_1893D: ; /* IDA 0x1893D; test    [bx+840h], dh */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx+2112)))&((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);}
L_18941: ; /* IDA 0x18941; jz      short loc_1894F */
    if(m->flags&0x40) goto L_1894F;
L_18943: ; /* IDA 0x18943; dec     byte ptr [bx+85Bh] */
    dl_write8(m,m->ds,(uint16_t)(m->bx+2139),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->bx+2139)),8));
L_18947: ; /* IDA 0x18947; not     dh */
    vg_dh(m,(uint8_t)(~((uint8_t)(m->dx>>8))));
L_18949: ; /* IDA 0x18949; and     [bx+840h], dh */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx+2112)))&((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->bx+2112),result);}
L_1894D: ; /* IDA 0x1894D; not     dh */
    vg_dh(m,(uint8_t)(~((uint8_t)(m->dx>>8))));
L_1894F: ; /* IDA 0x1894F; test    [si], dl */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&((uint8_t)m->dx))&255;ec_logic(m,result,8);}
L_18951: ; /* IDA 0x18951; jz      short locret_1895C */
    if(m->flags&0x40) goto L_1895C;
L_18953: ; /* IDA 0x18953; dec     byte ptr [si+1Bh] */
    dl_write8(m,m->ds,(uint16_t)(m->si+27),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->si+27)),8));
L_18956: ; /* IDA 0x18956; not     dl */
    vg_dl(m,(uint8_t)(~((uint8_t)m->dx)));
L_18958: ; /* IDA 0x18958; and     [si], dl */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&((uint8_t)m->dx))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_1895A: ; /* IDA 0x1895A; not     dl */
    vg_dl(m,(uint8_t)(~((uint8_t)m->dx)));
L_1895C: ; /* IDA 0x1895C; retn */
    return;
}
/* Original sub_195C9; IDA 0x195C9. */
static void oc_body_195C9(KiMachine16 *m,const KiEconomyHooks *h) {
L_195C9: ; /* IDA 0x195C9; push    ds */
    ec_push(m,m->ds);
L_195CA: ; /* IDA 0x195CA; push    es */
    ec_push(m,m->es);
L_195CB: ; /* IDA 0x195CB; push    ax */
    ec_push(m,m->ax);
L_195CC: ; /* IDA 0x195CC; push    bx */
    ec_push(m,m->bx);
L_195CD: ; /* IDA 0x195CD; push    cx */
    ec_push(m,m->cx);
L_195CE: ; /* IDA 0x195CE; push    dx */
    ec_push(m,m->dx);
L_195CF: ; /* IDA 0x195CF; push    si */
    ec_push(m,m->si);
L_195D0: ; /* IDA 0x195D0; push    di */
    ec_push(m,m->di);
L_195D1: ; /* IDA 0x195D1; cld */
    m->flags&=(uint16_t)~0x400u;
L_195D2: ; /* IDA 0x195D2; mov     es, cs:word_10D3A */
    m->es=(uint16_t)(dl_read16(m,m->cs,3386));
L_195D7: ; /* IDA 0x195D7; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_195D9: ; /* IDA 0x195D9; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_195DB: ; /* IDA 0x195DB; cmp     dx, 0B2h */
    ec_cmp(m,m->dx,178,16);
L_195DF: ; /* IDA 0x195DF; jbe     short loc_195E4 */
    if(m->flags&(1|0x40)) goto L_195E4;
L_195E1: ; /* IDA 0x195E1; mov     dx, 0B2h */
    m->dx=(uint16_t)(178);
L_195E4: ; /* IDA 0x195E4; add     dx, 1B6h */
    m->dx=(uint16_t)(ec_math(m,m->dx,438,16,0,0));
L_195E8: ; /* IDA 0x195E8; add     bx, 26h ; '&' */
    m->bx=(uint16_t)(ec_math(m,m->bx,38,16,0,0));
L_195EB: ; /* IDA 0x195EB; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_195ED: ; /* IDA 0x195ED; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_195EF: ; /* IDA 0x195EF; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_195F1: ; /* IDA 0x195F1; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_195F3: ; /* IDA 0x195F3; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_195F5: ; /* IDA 0x195F5; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_195F7: ; /* IDA 0x195F7; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_195F9: ; /* IDA 0x195F9; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_195FB: ; /* IDA 0x195FB; add     di, bx */
    m->di=(uint16_t)(ec_math(m,m->di,m->bx,16,0,0));
L_195FD: ; /* IDA 0x195FD; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_195FF: ; /* IDA 0x195FF; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19601: ; /* IDA 0x19601; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19603: ; /* IDA 0x19603; add     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
L_19605: ; /* IDA 0x19605; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_19607: ; /* IDA 0x19607; sub     dx, 37h ; '7' */
    m->dx=(uint16_t)(ec_math(m,m->dx,55,16,0,1));
L_1960A: ; /* IDA 0x1960A; sub     di, 3C0h */
    m->di=(uint16_t)(ec_math(m,m->di,960,16,0,1));
L_1960E: ; /* IDA 0x1960E; add     di, dx */
    m->di=(uint16_t)(ec_math(m,m->di,m->dx,16,0,0));
L_19610: ; /* IDA 0x19610; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_19613: ; /* IDA 0x19613; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_19615: ; /* IDA 0x19615; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_19618: ; /* IDA 0x19618; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_1961A: ; /* IDA 0x1961A; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1961B: ; /* IDA 0x1961B; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1961C: ; /* IDA 0x1961C; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1961E: ; /* IDA 0x1961E; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1961F: ; /* IDA 0x1961F; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_19621: ; /* IDA 0x19621; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_19622: ; /* IDA 0x19622; mov     al, 0 */
    vg_al(m,(uint8_t)(0));
L_19624: ; /* IDA 0x19624; call    sub_1963F */
    oc_call(m,38463,38439,h);
L_19627: ; /* IDA 0x19627; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_19629: ; /* IDA 0x19629; call    sub_1963F */
    oc_call(m,38463,38444,h);
L_1962C: ; /* IDA 0x1962C; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_1962E: ; /* IDA 0x1962E; call    sub_1963F */
    oc_call(m,38463,38449,h);
L_19631: ; /* IDA 0x19631; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_19633: ; /* IDA 0x19633; call    sub_1963F */
    oc_call(m,38463,38454,h);
L_19636: ; /* IDA 0x19636; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_19637: ; /* IDA 0x19637; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_19638: ; /* IDA 0x19638; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_19639: ; /* IDA 0x19639; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1963A: ; /* IDA 0x1963A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1963B: ; /* IDA 0x1963B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1963C: ; /* IDA 0x1963C; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1963D: ; /* IDA 0x1963D; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1963E: ; /* IDA 0x1963E; retn */
    return;
}
/* Original sub_1963F; IDA 0x1963F. */
static void oc_body_1963F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1963F: ; /* IDA 0x1963F; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_19642: ; /* IDA 0x19642; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_19643: ; /* IDA 0x19643; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_19645: ; /* IDA 0x19645; mov     cx, 4 */
    m->cx=(uint16_t)(4);
L_19648: ; /* IDA 0x19648; movsw */
    mc_string(m,0,0);
L_19649: ; /* IDA 0x19649; add     si, 4Eh ; 'N' */
    m->si=(uint16_t)(ec_math(m,m->si,78,16,0,0));
L_1964C: ; /* IDA 0x1964C; add     di, 16h */
    m->di=(uint16_t)(ec_math(m,m->di,22,16,0,0));
L_1964F: ; /* IDA 0x1964F; loop    loc_19648 */
    if(--m->cx) goto L_19648;
L_19651: ; /* IDA 0x19651; add     di, 0BA0h */
    m->di=(uint16_t)(ec_math(m,m->di,2976,16,0,0));
L_19655: ; /* IDA 0x19655; retn */
    return;
}
int ki_outcome_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0xce7:oc_body_10CE7(m,h);break;
case 0x4236:oc_body_14236(m,h);break;
case 0x474a:oc_body_1474A(m,h);break;
case 0x4cf3:oc_body_14CF3(m,h);break;
case 0x4d63:oc_body_14D63(m,h);break;
case 0x4da4:oc_body_14DA4(m,h);break;
case 0x4df0:oc_body_14DF0(m,h);break;
case 0x4fce:oc_body_14FCE(m,h);break;
case 0x5074:oc_body_15074(m,h);break;
case 0x50b4:oc_body_150B4(m,h);break;
case 0x5130:oc_body_15130(m,h);break;
case 0x51b3:oc_body_151B3(m,h);break;
case 0x5285:oc_body_15285(m,h);break;
case 0x52d7:oc_body_152D7(m,h);break;
case 0x5ca4:oc_body_15CA4(m,h);break;
case 0x5ce0:oc_body_15CE0(m,h);break;
case 0x7028:oc_body_17028(m,h);break;
case 0x88cc:oc_body_188CC(m,h);break;
case 0x890a:oc_body_1890A(m,h);break;
case 0x95c9:oc_body_195C9(m,h);break;
case 0x963f:oc_body_1963F(m,h);break;
default:return 0;}return 1;}
static int oc_known(uint16_t t) {switch(t) {case 0xce7: case 0x4236: case 0x474a: case 0x4cf3: case 0x4d63: case 0x4da4: case 0x4df0: case 0x4fce: case 0x5074: case 0x50b4: case 0x5130: case 0x51b3: case 0x5285: case 0x52d7: case 0x5ca4: case 0x5ce0: case 0x7028: case 0x88cc: case 0x890a: case 0x95c9: case 0x963f: return 1;default:return 0;}}
void sub_10CE7(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0xce7,h);m->ip=ec_pop(m);}
void sub_14236(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x4236,h);m->ip=ec_pop(m);}
void sub_1474A(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x474a,h);m->ip=ec_pop(m);}
void sub_14CF3(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x4cf3,h);m->ip=ec_pop(m);}
void sub_14D63(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x4d63,h);m->ip=ec_pop(m);}
void sub_14DA4(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x4da4,h);m->ip=ec_pop(m);}
void sub_14DF0(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x4df0,h);m->ip=ec_pop(m);}
void sub_14FCE(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x4fce,h);m->ip=ec_pop(m);}
void sub_15074(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x5074,h);m->ip=ec_pop(m);}
void sub_150B4(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x50b4,h);m->ip=ec_pop(m);}
void sub_15130(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x5130,h);m->ip=ec_pop(m);}
void sub_151B3(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x51b3,h);m->ip=ec_pop(m);}
void sub_15285(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x5285,h);m->ip=ec_pop(m);}
void sub_152D7(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x52d7,h);m->ip=ec_pop(m);}
void sub_15CA4(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x5ca4,h);m->ip=ec_pop(m);}
void sub_15CE0(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x5ce0,h);m->ip=ec_pop(m);}
void sub_17028(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x7028,h);m->ip=ec_pop(m);}
void sub_188CC(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x88cc,h);m->ip=ec_pop(m);}
void sub_1890A(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x890a,h);m->ip=ec_pop(m);}
void sub_195C9(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x95c9,h);m->ip=ec_pop(m);}
void sub_1963F(KiMachine16 *m,const KiEconomyHooks *h) {ki_outcome_body(m,0x963f,h);m->ip=ec_pop(m);}
