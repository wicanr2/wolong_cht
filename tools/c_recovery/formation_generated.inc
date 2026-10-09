/* Original native C formation; re/120, spec/240. */
/* Original sub_16C5E; IDA 0x16C5E. */
static void fc_body_16C5E(KiMachine16 *m,const KiEconomyHooks *h) {
L_16C5E: ; /* IDA 0x16C5E; push    ds */
    ec_push(m,m->ds);
L_16C5F: ; /* IDA 0x16C5F; push    es */
    ec_push(m,m->es);
L_16C60: ; /* IDA 0x16C60; push    ax */
    ec_push(m,m->ax);
L_16C61: ; /* IDA 0x16C61; push    bx */
    ec_push(m,m->bx);
L_16C62: ; /* IDA 0x16C62; push    cx */
    ec_push(m,m->cx);
L_16C63: ; /* IDA 0x16C63; push    si */
    ec_push(m,m->si);
L_16C64: ; /* IDA 0x16C64; push    di */
    ec_push(m,m->di);
L_16C65: ; /* IDA 0x16C65; call    sub_12078 */
    fc_call(m,8312,27752,h);
L_16C68: ; /* IDA 0x16C68; mov     cx, 0 */
    m->cx=(uint16_t)(0);
L_16C6B: ; /* IDA 0x16C6B; call    sub_18853 */
    fc_call(m,34899,27758,h);
L_16C6E: ; /* IDA 0x16C6E; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_16C70: ; /* IDA 0x16C70; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_16C72: ; /* IDA 0x16C72; call    sub_17663 */
    fc_call(m,30307,27765,h);
L_16C75: ; /* IDA 0x16C75; jb      short loc_16C81 */
    if(m->flags&1) goto L_16C81;
L_16C77: ; /* IDA 0x16C77; call    sub_16C92 */
    fc_call(m,27794,27770,h);
L_16C7A: ; /* IDA 0x16C7A; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_16C7C: ; /* IDA 0x16C7C; call    sub_15E80 */
#if KI_FORMATION_MUTATION == 9
    ; /* lost outer sidebar refresh */
#else
    fc_call(m,24192,27775,h);
#endif
L_16C7F: ; /* IDA 0x16C7F; jmp     short loc_16C68 */
    goto L_16C68;
L_16C81: ; /* IDA 0x16C81; call    sub_120D6 */
    fc_call(m,8406,27780,h);
L_16C84: ; /* IDA 0x16C84; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16C87: ; /* IDA 0x16C87; call    sub_18853 */
    fc_call(m,34899,27786,h);
L_16C8A: ; /* IDA 0x16C8A; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_16C8B: ; /* IDA 0x16C8B; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16C8C: ; /* IDA 0x16C8C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_16C8D: ; /* IDA 0x16C8D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16C8E: ; /* IDA 0x16C8E; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16C8F: ; /* IDA 0x16C8F; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_16C90: ; /* IDA 0x16C90; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16C91: ; /* IDA 0x16C91; retn */
    return;
}
/* Original sub_16C92; IDA 0x16C92. */
static void fc_body_16C92(KiMachine16 *m,const KiEconomyHooks *h) {
L_16C92: ; /* IDA 0x16C92; push    ds */
    ec_push(m,m->ds);
L_16C93: ; /* IDA 0x16C93; push    es */
    ec_push(m,m->es);
L_16C94: ; /* IDA 0x16C94; push    si */
    ec_push(m,m->si);
L_16C95: ; /* IDA 0x16C95; push    di */
    ec_push(m,m->di);
L_16C96: ; /* IDA 0x16C96; call    sub_101B4 */
    fc_call(m,436,27801,h);
L_16C99: ; /* IDA 0x16C99; mov     cx, 1 */
    m->cx=(uint16_t)(1);
L_16C9C: ; /* IDA 0x16C9C; call    sub_18853 */
    fc_call(m,34899,27807,h);
L_16C9F: ; /* IDA 0x16C9F; mov     ax, cs:word_10D52 */
    m->ax=(uint16_t)(dl_read16(m,m->cs,3410));
L_16CA3: ; /* IDA 0x16CA3; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_16CA5: ; /* IDA 0x16CA5; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_16CA7: ; /* IDA 0x16CA7; mov     ah, [bx+1Ch] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+28))));
L_16CAA: ; /* IDA 0x16CAA; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_16CAC: ; /* IDA 0x16CAC; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16CAE: ; /* IDA 0x16CAE; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16CB0: ; /* IDA 0x16CB0; add     ax, 4 */
    m->ax=(uint16_t)(ec_math(m,m->ax,4,16,0,0));
L_16CB3: ; /* IDA 0x16CB3; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_16CB5: ; /* IDA 0x16CB5; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16CB7: ; /* IDA 0x16CB7; sub     si, 4240h */
    m->si=(uint16_t)(ec_math(m,m->si,16960,16,0,1));
L_16CBB: ; /* IDA 0x16CBB; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_16CBD: ; /* IDA 0x16CBD; add     si, 2240h */
    m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,0));
L_16CC1: ; /* IDA 0x16CC1; mov     al, [bx+1Ch] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+28))));
L_16CC4: ; /* IDA 0x16CC4; mov     [si+1], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+1),(uint8_t)m->ax);
L_16CC7: ; /* IDA 0x16CC7; call    sub_16D56 */
    fc_call(m,27990,27850,h);
L_16CCA: ; /* IDA 0x16CCA; call    sub_14698 */
    fc_call(m,18072,27853,h);
L_16CCD: ; /* IDA 0x16CCD; jmp     short loc_16CD2 */
    goto L_16CD2;
L_16CCF: ; /* IDA 0x16CCF; call    sub_101B4 */
    fc_call(m,436,27858,h);
L_16CD2: ; /* IDA 0x16CD2; call    sub_16DFD */
    fc_call(m,28157,27861,h);
L_16CD5: ; /* IDA 0x16CD5; jmp     short loc_16CDA */
    goto L_16CDA;
L_16CD7: ; /* IDA 0x16CD7; call    sub_101B4 */
    fc_call(m,436,27866,h);
L_16CDA: ; /* IDA 0x16CDA; call    sub_16FD2 */
    fc_call(m,28626,27869,h);
L_16CDD: ; /* IDA 0x16CDD; call    sub_16D6F */
    fc_call(m,28015,27872,h);
L_16CE0: ; /* IDA 0x16CE0; call    sub_16DA8 */
    fc_call(m,28072,27875,h);
L_16CE3: ; /* IDA 0x16CE3; call    sub_101DB */
    fc_call(m,475,27878,h);
L_16CE6: ; /* IDA 0x16CE6; call    sub_121E7 */
    fc_call(m,8679,27881,h);
L_16CE9: ; /* IDA 0x16CE9; jb      short loc_16D42 */
    if(m->flags&1) goto L_16D42;
L_16CEB: ; /* IDA 0x16CEB; call    sub_1E453 */
    fc_call(m,58451,27886,h);
L_16CEE: ; /* IDA 0x16CEE; sub     al, 3Eh ; '>' */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,62,8,0,1)));
L_16CF0: ; /* IDA 0x16CF0; jb      short loc_16CD7 */
    if(m->flags&1) goto L_16CD7;
L_16CF2: ; /* IDA 0x16CF2; cmp     al, 6 */
    ec_cmp(m,(uint8_t)m->ax,6,8);
L_16CF4: ; /* IDA 0x16CF4; jz      short loc_16D14 */
    if(m->flags&0x40) goto L_16D14;
L_16CF6: ; /* IDA 0x16CF6; call    sub_14717 */
    fc_call(m,18199,27897,h);
L_16CF9: ; /* IDA 0x16CF9; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_16CFB: ; /* IDA 0x16CFB; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_16CFD: ; /* IDA 0x16CFD; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16CFF: ; /* IDA 0x16CFF; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16D01: ; /* IDA 0x16D01; mov     al, [bx+si+2Ah] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+m->si+42))));
L_16D04: ; /* IDA 0x16D04; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_16D06: ; /* IDA 0x16D06; cmp     al, 4 */
    ec_cmp(m,(uint8_t)m->ax,4,8);
L_16D08: ; /* IDA 0x16D08; jb      short loc_16D0C */
    if(m->flags&1) goto L_16D0C;
L_16D0A: ; /* IDA 0x16D0A; mov     al, 1 */
#if KI_FORMATION_MUTATION == 6
    vg_al(m,2);
#else
    vg_al(m,(uint8_t)(1));
#endif
L_16D0C: ; /* IDA 0x16D0C; mov     [bx+si+2Ah], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+m->si+42),(uint8_t)m->ax);
L_16D0F: ; /* IDA 0x16D0F; call    sub_14698 */
    fc_call(m,18072,27922,h);
L_16D12: ; /* IDA 0x16D12; jmp     short loc_16CD7 */
    goto L_16CD7;
L_16D14: ; /* IDA 0x16D14; cmp     byte ptr [si+29h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+41)),0,8);
L_16D18: ; /* IDA 0x16D18; jnz     short loc_16D24 */
#if KI_FORMATION_MUTATION == 7
    goto L_16D24; /* lost empty leader refusal */
#else
    if(!(m->flags&0x40)) goto L_16D24;
#endif
L_16D1A: ; /* IDA 0x16D1A; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_16D1C: ; /* IDA 0x16D1C; mov     cx, 3Eh ; '>' */
    m->cx=(uint16_t)(62);
L_16D1F: ; /* IDA 0x16D1F; call    sub_18810 */
    fc_call(m,34832,27938,h);
L_16D22: ; /* IDA 0x16D22; jmp     short loc_16CCF */
    goto L_16CCF;
L_16D24: ; /* IDA 0x16D24; call    sub_16F26 */
#if KI_FORMATION_MUTATION == 8
    ; /* lost committed corps initialization */
#else
    fc_call(m,28454,27943,h);
#endif
L_16D27: ; /* IDA 0x16D27; sub     di, 4 */
    m->di=(uint16_t)(ec_math(m,m->di,4,16,0,1));
L_16D2A: ; /* IDA 0x16D2A; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_16D2C: ; /* IDA 0x16D2C; sub     bx, 2240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,8768,16,0,1));
L_16D30: ; /* IDA 0x16D30; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16D32: ; /* IDA 0x16D32; mov     cx, 19Bh */
    m->cx=(uint16_t)(411);
L_16D35: ; /* IDA 0x16D35; mov     ah, [bx+425Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16990))));
L_16D39: ; /* IDA 0x16D39; mov     al, [bx+4241h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16961))));
L_16D3D: ; /* IDA 0x16D3D; call    sub_18810 */
    fc_call(m,34832,27968,h);
L_16D40: ; /* IDA 0x16D40; jmp     short loc_16D48 */
    goto L_16D48;
L_16D42: ; /* IDA 0x16D42; call    sub_14717 */
#if KI_FORMATION_MUTATION == 5
    ; /* lost cancellation pool return */
#else
    fc_call(m,18199,27973,h);
#endif
L_16D45: ; /* IDA 0x16D45; sub     di, 4 */
#if KI_FORMATION_MUTATION == 10
    fc_call(m,0x6fd2,0x6d48,h);m->di=ec_math(m,m->di,4,16,0,1);
#else
    m->di=(uint16_t)(ec_math(m,m->di,4,16,0,1));
#endif
L_16D48: ; /* IDA 0x16D48; call    sub_16E80 */
    fc_call(m,28288,27979,h);
L_16D4B: ; /* IDA 0x16D4B; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16D4E: ; /* IDA 0x16D4E; call    sub_18853 */
    fc_call(m,34899,27985,h);
L_16D51: ; /* IDA 0x16D51; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_16D52: ; /* IDA 0x16D52; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16D53: ; /* IDA 0x16D53; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_16D54: ; /* IDA 0x16D54; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16D55: ; /* IDA 0x16D55; retn */
    return;
}
/* Original sub_16D56; IDA 0x16D56. */
static void fc_body_16D56(KiMachine16 *m,const KiEconomyHooks *h) {
L_16D56: ; /* IDA 0x16D56; mov     byte ptr [si+2Ah], 1 */
#if KI_FORMATION_MUTATION == 1
    dl_write8(m,m->ds,(uint16_t)(m->si+42),2);
#else
    dl_write8(m,m->ds,(uint16_t)(m->si+42),1);
#endif
L_16D5A: ; /* IDA 0x16D5A; mov     byte ptr [si+2Eh], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+46),1);
L_16D5E: ; /* IDA 0x16D5E; mov     byte ptr [si+32h], 3 */
    dl_write8(m,m->ds,(uint16_t)(m->si+50),3);
L_16D62: ; /* IDA 0x16D62; mov     byte ptr [si+36h], 3 */
    dl_write8(m,m->ds,(uint16_t)(m->si+54),3);
L_16D66: ; /* IDA 0x16D66; mov     byte ptr [si+3Ah], 2 */
    dl_write8(m,m->ds,(uint16_t)(m->si+58),2);
L_16D6A: ; /* IDA 0x16D6A; mov     byte ptr [si+3Eh], 2 */
    dl_write8(m,m->ds,(uint16_t)(m->si+62),2);
L_16D6E: ; /* IDA 0x16D6E; retn */
    return;
}
/* Original sub_16D6F; IDA 0x16D6F. */
static void fc_body_16D6F(KiMachine16 *m,const KiEconomyHooks *h) {
L_16D6F: ; /* IDA 0x16D6F; push    ds */
    ec_push(m,m->ds);
L_16D70: ; /* IDA 0x16D70; push    es */
    ec_push(m,m->es);
L_16D71: ; /* IDA 0x16D71; push    si */
    ec_push(m,m->si);
L_16D72: ; /* IDA 0x16D72; push    di */
    ec_push(m,m->di);
L_16D73: ; /* IDA 0x16D73; push    si */
    ec_push(m,m->si);
L_16D74: ; /* IDA 0x16D74; mov     si, di */
    m->si=(uint16_t)(m->di);
L_16D76: ; /* IDA 0x16D76; mov     cx, 3 */
    m->cx=(uint16_t)(3);
L_16D79: ; /* IDA 0x16D79; mov     di, 43A7h */
    m->di=(uint16_t)(17319);
L_16D7C: ; /* IDA 0x16D7C; mov     ax, [si] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_16D7E: ; /* IDA 0x16D7E; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_16D81: ; /* IDA 0x16D81; mul     dx */
#if KI_FORMATION_MUTATION == 2
    dl_mul(m,9,16);
#else
    dl_mul(m,m->dx,16);
#endif
L_16D83: ; /* IDA 0x16D83; mov     bx, 906h */
    m->bx=(uint16_t)(2310);
L_16D86: ; /* IDA 0x16D86; call    sub_1062F */
    fc_call(m,1583,28041,h);
L_16D89: ; /* IDA 0x16D89; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_16D8A: ; /* IDA 0x16D8A; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_16D8B: ; /* IDA 0x16D8B; add     di, 500h */
    m->di=(uint16_t)(ec_math(m,m->di,1280,16,0,0));
L_16D8F: ; /* IDA 0x16D8F; loop    loc_16D7C */
    if(--m->cx) goto L_16D7C;
L_16D91: ; /* IDA 0x16D91; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16D92: ; /* IDA 0x16D92; mov     ax, [si+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+4)));
L_16D95: ; /* IDA 0x16D95; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_16D98: ; /* IDA 0x16D98; mul     dx */
    dl_mul(m,m->dx,16);
L_16D9A: ; /* IDA 0x16D9A; mov     bx, 0F04h */
    m->bx=(uint16_t)(3844);
L_16D9D: ; /* IDA 0x16D9D; mov     di, 2FA7h */
    m->di=(uint16_t)(12199);
L_16DA0: ; /* IDA 0x16DA0; call    sub_1062F */
    fc_call(m,1583,28067,h);
L_16DA3: ; /* IDA 0x16DA3; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_16DA4: ; /* IDA 0x16DA4; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16DA5: ; /* IDA 0x16DA5; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_16DA6: ; /* IDA 0x16DA6; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16DA7: ; /* IDA 0x16DA7; retn */
    return;
}
/* Original sub_16DA8; IDA 0x16DA8. */
static void fc_body_16DA8(KiMachine16 *m,const KiEconomyHooks *h) {
L_16DA8: ; /* IDA 0x16DA8; push    ds */
    ec_push(m,m->ds);
L_16DA9: ; /* IDA 0x16DA9; push    es */
    ec_push(m,m->es);
L_16DAA: ; /* IDA 0x16DAA; push    si */
    ec_push(m,m->si);
L_16DAB: ; /* IDA 0x16DAB; push    di */
    ec_push(m,m->di);
L_16DAC: ; /* IDA 0x16DAC; add     si, 28h ; '(' */
    m->si=(uint16_t)(ec_math(m,m->si,40,16,0,0));
L_16DAF: ; /* IDA 0x16DAF; mov     di, 3C1Dh */
    m->di=(uint16_t)(15389);
L_16DB2: ; /* IDA 0x16DB2; mov     ax, ds */
    m->ax=(uint16_t)(m->ds);
L_16DB4: ; /* IDA 0x16DB4; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_16DB6: ; /* IDA 0x16DB6; mov     ds, cs:word_10D50 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3408));
L_16DBB: ; /* IDA 0x16DBB; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_16DBE: ; /* IDA 0x16DBE; push    si */
    ec_push(m,m->si);
L_16DBF: ; /* IDA 0x16DBF; mov     ah, es:[si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->si+2))));
L_16DC3: ; /* IDA 0x16DC3; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_16DC5: ; /* IDA 0x16DC5; dec     ah */
    vg_ah(m,(uint8_t)(st_dec(m,(uint8_t)(m->ax>>8),8)));
L_16DC7: ; /* IDA 0x16DC7; shr     ax, 1 */
#if KI_FORMATION_MUTATION == 3
    m->ax=dl_shift(m,m->ax,2,16,5);
#else
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
#endif
L_16DC9: ; /* IDA 0x16DC9; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_16DCB: ; /* IDA 0x16DCB; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16DCD: ; /* IDA 0x16DCD; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_16DCF: ; /* IDA 0x16DCF; add     si, 15C0h */
    m->si=(uint16_t)(ec_math(m,m->si,5568,16,0,0));
L_16DD3: ; /* IDA 0x16DD3; mov     bx, di */
    m->bx=(uint16_t)(m->di);
L_16DD5: ; /* IDA 0x16DD5; sub     bx, 4 */
    m->bx=(uint16_t)(ec_math(m,m->bx,4,16,0,1));
L_16DD8: ; /* IDA 0x16DD8; mov     ax, 1003h */
    m->ax=(uint16_t)(4099);
L_16DDB: ; /* IDA 0x16DDB; call    sub_1F9B0 */
    fc_call(m,63920,28126,h);
L_16DDE: ; /* IDA 0x16DDE; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16DDF: ; /* IDA 0x16DDF; mov     al, es:[si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->si+1))));
L_16DE3: ; /* IDA 0x16DE3; mov     ah, 0Ah */
    vg_ah(m,(uint8_t)(10));
L_16DE5: ; /* IDA 0x16DE5; mul     ah */
    dl_mul(m,(uint8_t)(m->ax>>8),8);
L_16DE7: ; /* IDA 0x16DE7; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_16DE9: ; /* IDA 0x16DE9; mov     bx, 904h */
    m->bx=(uint16_t)(2308);
L_16DEC: ; /* IDA 0x16DEC; call    sub_1062F */
    fc_call(m,1583,28143,h);
L_16DEF: ; /* IDA 0x16DEF; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_16DF2: ; /* IDA 0x16DF2; add     di, 500h */
    m->di=(uint16_t)(ec_math(m,m->di,1280,16,0,0));
L_16DF6: ; /* IDA 0x16DF6; loop    loc_16DBE */
    if(--m->cx) goto L_16DBE;
L_16DF8: ; /* IDA 0x16DF8; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_16DF9: ; /* IDA 0x16DF9; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16DFA: ; /* IDA 0x16DFA; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_16DFB: ; /* IDA 0x16DFB; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16DFC: ; /* IDA 0x16DFC; retn */
    return;
}
/* Original sub_16DFD; IDA 0x16DFD. */
static void fc_body_16DFD(KiMachine16 *m,const KiEconomyHooks *h) {
L_16DFD: ; /* IDA 0x16DFD; push    ax */
    ec_push(m,m->ax);
L_16DFE: ; /* IDA 0x16DFE; push    bx */
    ec_push(m,m->bx);
L_16DFF: ; /* IDA 0x16DFF; push    cx */
    ec_push(m,m->cx);
L_16E00: ; /* IDA 0x16E00; push    dx */
    ec_push(m,m->dx);
L_16E01: ; /* IDA 0x16E01; push    si */
    ec_push(m,m->si);
L_16E02: ; /* IDA 0x16E02; push    di */
    ec_push(m,m->di);
L_16E03: ; /* IDA 0x16E03; mov     al, 3Dh ; '=' */
    vg_al(m,(uint8_t)(61));
L_16E05: ; /* IDA 0x16E05; mov     dx, 9 */
    m->dx=(uint16_t)(9);
L_16E08: ; /* IDA 0x16E08; mov     bx, 5 */
    m->bx=(uint16_t)(5);
L_16E0B: ; /* IDA 0x16E0B; mov     cx, 0C0Fh */
    m->cx=(uint16_t)(3087);
L_16E0E: ; /* IDA 0x16E0E; call    sub_1895D */
    fc_call(m,35165,28177,h);
L_16E11: ; /* IDA 0x16E11; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_16E13: ; /* IDA 0x16E13; mov     dx, 98h */
    m->dx=(uint16_t)(152);
L_16E16: ; /* IDA 0x16E16; mov     bx, 78h ; 'x' */
    m->bx=(uint16_t)(120);
L_16E19: ; /* IDA 0x16E19; call    sub_10337 */
    fc_call(m,823,28188,h);
L_16E1C: ; /* IDA 0x16E1C; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_16E1E: ; /* IDA 0x16E1E; sub     ax, 2240h */
    m->ax=(uint16_t)(ec_math(m,m->ax,8768,16,0,1));
L_16E21: ; /* IDA 0x16E21; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16E23: ; /* IDA 0x16E23; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_16E25: ; /* IDA 0x16E25; mov     al, [si+4241h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+16961))));
L_16E29: ; /* IDA 0x16E29; mov     bx, 2593h */
    m->bx=(uint16_t)(9619);
L_16E2C: ; /* IDA 0x16E2C; call    sub_107D2 */
    fc_call(m,2002,28207,h);
L_16E2F: ; /* IDA 0x16E2F; mov     dx, 128h */
    m->dx=(uint16_t)(296);
L_16E32: ; /* IDA 0x16E32; mov     bx, 80h */
    m->bx=(uint16_t)(128);
L_16E35: ; /* IDA 0x16E35; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_16E38: ; /* IDA 0x16E38; add     si, 4242h */
    m->si=(uint16_t)(ec_math(m,m->si,16962,16,0,0));
L_16E3C: ; /* IDA 0x16E3C; call    sub_106FD */
    fc_call(m,1789,28223,h);
L_16E3F: ; /* IDA 0x16E3F; mov     al, [di+19h] */
#if KI_FORMATION_MUTATION == 4
    vg_al(m,dl_read8(m,m->ds,(uint16_t)(m->di+24)));
#else
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+25))));
#endif
L_16E42: ; /* IDA 0x16E42; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_16E44: ; /* IDA 0x16E44; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_16E46: ; /* IDA 0x16E46; mov     bx, 0F03h */
    m->bx=(uint16_t)(3843);
L_16E49: ; /* IDA 0x16E49; mov     di, 34A8h */
    m->di=(uint16_t)(13480);
L_16E4C: ; /* IDA 0x16E4C; call    sub_1062F */
    fc_call(m,1583,28239,h);
L_16E4F: ; /* IDA 0x16E4F; mov     al, 3Eh ; '>' */
    vg_al(m,(uint8_t)(62));
L_16E51: ; /* IDA 0x16E51; mov     bx, 0C0h */
    m->bx=(uint16_t)(192);
L_16E54: ; /* IDA 0x16E54; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_16E57: ; /* IDA 0x16E57; push    ax */
    ec_push(m,m->ax);
L_16E58: ; /* IDA 0x16E58; push    bx */
    ec_push(m,m->bx);
L_16E59: ; /* IDA 0x16E59; push    cx */
    ec_push(m,m->cx);
L_16E5A: ; /* IDA 0x16E5A; mov     dx, 0C8h */
    m->dx=(uint16_t)(200);
L_16E5D: ; /* IDA 0x16E5D; mov     cx, 203h */
    m->cx=(uint16_t)(515);
L_16E60: ; /* IDA 0x16E60; call    sub_1E3D7 */
    fc_call(m,58327,28259,h);
L_16E63: ; /* IDA 0x16E63; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_16E64: ; /* IDA 0x16E64; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16E65: ; /* IDA 0x16E65; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16E66: ; /* IDA 0x16E66; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_16E68: ; /* IDA 0x16E68; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_16E6B: ; /* IDA 0x16E6B; loop    loc_16E57 */
    if(--m->cx) goto L_16E57;
L_16E6D: ; /* IDA 0x16E6D; mov     dx, 118h */
    m->dx=(uint16_t)(280);
L_16E70: ; /* IDA 0x16E70; mov     bx, 110h */
    m->bx=(uint16_t)(272);
L_16E73: ; /* IDA 0x16E73; mov     cx, 20Bh */
    m->cx=(uint16_t)(523);
L_16E76: ; /* IDA 0x16E76; call    sub_1E3D7 */
    fc_call(m,58327,28281,h);
L_16E79: ; /* IDA 0x16E79; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_16E7A: ; /* IDA 0x16E7A; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16E7B: ; /* IDA 0x16E7B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_16E7C: ; /* IDA 0x16E7C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_16E7D: ; /* IDA 0x16E7D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16E7E: ; /* IDA 0x16E7E; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16E7F: ; /* IDA 0x16E7F; retn */
    return;
}
/* Original sub_16E80; IDA 0x16E80. */
static void fc_body_16E80(KiMachine16 *m,const KiEconomyHooks *h) {
L_16E80: ; /* IDA 0x16E80; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_16E82: ; /* IDA 0x16E82; mov     dx, 9 */
    m->dx=(uint16_t)(9);
L_16E85: ; /* IDA 0x16E85; mov     bx, 5 */
    m->bx=(uint16_t)(5);
L_16E88: ; /* IDA 0x16E88; mov     cx, 0C0Fh */
    m->cx=(uint16_t)(3087);
L_16E8B: ; /* IDA 0x16E8B; call    sub_1895D */
    fc_call(m,35165,28302,h);
L_16E8E: ; /* IDA 0x16E8E; retn */
    return;
}
int ki_formation_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x6c5e:fc_body_16C5E(m,h);break;
case 0x6c92:fc_body_16C92(m,h);break;
case 0x6d56:fc_body_16D56(m,h);break;
case 0x6d6f:fc_body_16D6F(m,h);break;
case 0x6da8:fc_body_16DA8(m,h);break;
case 0x6dfd:fc_body_16DFD(m,h);break;
case 0x6e80:fc_body_16E80(m,h);break;
default:return 0;}return 1;}
static int fc_known(uint16_t t) {switch(t) {case 0x6c5e: case 0x6c92: case 0x6d56: case 0x6d6f: case 0x6da8: case 0x6dfd: case 0x6e80: return 1;default:return 0;}}
