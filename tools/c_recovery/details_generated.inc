/* Original native C details; re/122, spec/242. */
/* Original sub_15E1E; IDA 0x15E1E. */
static void dtl_body_15E1E(KiMachine16 *m,const KiEconomyHooks *h) {
L_15E1E: ; /* IDA 0x15E1E; call    sub_101B4 */
    dtl_call(m,436,24097,h);
L_15E21: ; /* IDA 0x15E21; call    sub_15E2D */
    dtl_call(m,24109,24100,h);
L_15E24: ; /* IDA 0x15E24; mov     al, 0Fh */
    vg_al(m,(uint8_t)(15));
L_15E26: ; /* IDA 0x15E26; call    sub_15E80 */
    dtl_call(m,24192,24105,h);
L_15E29: ; /* IDA 0x15E29; call    sub_101DB */
    dtl_call(m,475,24108,h);
L_15E2C: ; /* IDA 0x15E2C; retn */
    return;
}
/* Original sub_15E2D; IDA 0x15E2D. */
static void dtl_body_15E2D(KiMachine16 *m,const KiEconomyHooks *h) {
L_15E2D: ; /* IDA 0x15E2D; or      byte ptr ds:98A6h, 2 */
#if KI_DETAILS_MUTATION == 11
    {uint16_t result=dl_read8(m,m->ds,0x98a6)|4;ec_logic(m,result,8);dl_write8(m,m->ds,0x98a6,result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,39078))|(2))&255;ec_logic(m,result,8);dl_write8(m,m->ds,39078,result);}
#endif
L_15E32: ; /* IDA 0x15E32; mov     al, 10h */
    vg_al(m,(uint8_t)(16));
L_15E34: ; /* IDA 0x15E34; mov     dx, 1Bh */
    m->dx=(uint16_t)(27);
L_15E37: ; /* IDA 0x15E37; mov     bx, 0Ah */
    m->bx=(uint16_t)(10);
L_15E3A: ; /* IDA 0x15E3A; mov     cx, 0D0Dh */
    m->cx=(uint16_t)(3341);
L_15E3D: ; /* IDA 0x15E3D; call    sub_1895D */
    dtl_call(m,35165,24128,h);
L_15E40: ; /* IDA 0x15E40; mov     al, 0 */
    vg_al(m,(uint8_t)(0));
L_15E42: ; /* IDA 0x15E42; mov     dx, 1B8h */
    m->dx=(uint16_t)(440);
L_15E45: ; /* IDA 0x15E45; mov     bx, 0C8h */
    m->bx=(uint16_t)(200);
L_15E48: ; /* IDA 0x15E48; call    sub_10337 */
    dtl_call(m,823,24139,h);
L_15E4B: ; /* IDA 0x15E4B; retn */
    return;
}
/* Original sub_17E1F; IDA 0x17E1F. */
static void dtl_body_17E1F(KiMachine16 *m,const KiEconomyHooks *h) {
L_17E1F: ; /* IDA 0x17E1F; push    ds */
    ec_push(m,m->ds);
L_17E20: ; /* IDA 0x17E20; push    es */
    ec_push(m,m->es);
L_17E21: ; /* IDA 0x17E21; push    ax */
    ec_push(m,m->ax);
L_17E22: ; /* IDA 0x17E22; push    bx */
    ec_push(m,m->bx);
L_17E23: ; /* IDA 0x17E23; push    cx */
    ec_push(m,m->cx);
L_17E24: ; /* IDA 0x17E24; push    dx */
    ec_push(m,m->dx);
L_17E25: ; /* IDA 0x17E25; push    di */
    ec_push(m,m->di);
L_17E26: ; /* IDA 0x17E26; push    bp */
    ec_push(m,m->bp);
L_17E27: ; /* IDA 0x17E27; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_17E2C: ; /* IDA 0x17E2C; add     si, 840h */
#if KI_DETAILS_MUTATION == 12
    m->si=ec_math(m,m->si,0x860,16,0,0);
#else
    m->si=(uint16_t)(ec_math(m,m->si,2112,16,0,0));
#endif
L_17E30: ; /* IDA 0x17E30; call    sub_17F61 */
    dtl_call(m,32609,32307,h);
L_17E33: ; /* IDA 0x17E33; call    sub_17E4A */
    dtl_call(m,32330,32310,h);
L_17E36: ; /* IDA 0x17E36; call    sub_17F1A */
    dtl_call(m,32538,32313,h);
L_17E39: ; /* IDA 0x17E39; call    sub_121E7 */
    dtl_call(m,8679,32316,h);
L_17E3C: ; /* IDA 0x17E3C; jnb     short loc_17E39 */
    if(!(m->flags&1)) goto L_17E39;
L_17E3E: ; /* IDA 0x17E3E; call    sub_17F81 */
    dtl_call(m,32641,32321,h);
L_17E41: ; /* IDA 0x17E41; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17E42: ; /* IDA 0x17E42; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17E43: ; /* IDA 0x17E43; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17E44: ; /* IDA 0x17E44; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17E45: ; /* IDA 0x17E45; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17E46: ; /* IDA 0x17E46; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17E47: ; /* IDA 0x17E47; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_17E48: ; /* IDA 0x17E48; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17E49: ; /* IDA 0x17E49; retn */
    return;
}
/* Original sub_17E4A; IDA 0x17E4A. */
static void dtl_body_17E4A(KiMachine16 *m,const KiEconomyHooks *h) {
L_17E4A: ; /* IDA 0x17E4A; push    ds */
    ec_push(m,m->ds);
L_17E4B: ; /* IDA 0x17E4B; push    si */
    ec_push(m,m->si);
L_17E4C: ; /* IDA 0x17E4C; mov     di, si */
    m->di=(uint16_t)(m->si);
L_17E4E: ; /* IDA 0x17E4E; add     si, 2 */
    m->si=(uint16_t)(ec_math(m,m->si,2,16,0,0));
L_17E51: ; /* IDA 0x17E51; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_17E54: ; /* IDA 0x17E54; mov     dx, 80h */
    m->dx=(uint16_t)(128);
L_17E57: ; /* IDA 0x17E57; mov     bx, 120h */
    m->bx=(uint16_t)(288);
L_17E5A: ; /* IDA 0x17E5A; call    sub_106FD */
    dtl_call(m,1789,32349,h);
L_17E5D: ; /* IDA 0x17E5D; push    ds */
    ec_push(m,m->ds);
L_17E5E: ; /* IDA 0x17E5E; mov     bh, [di+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+1))));
L_17E61: ; /* IDA 0x17E61; cmp     bh, 18h */
    ec_cmp(m,(uint8_t)(m->bx>>8),24,8);
L_17E64: ; /* IDA 0x17E64; jz      short loc_17E7E */
#if KI_DETAILS_MUTATION == 1
    ; /* lost neutral name branch */
#else
    if(m->flags&0x40) goto L_17E7E;
#endif
L_17E66: ; /* IDA 0x17E66; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_17E68: ; /* IDA 0x17E68; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_17E6A: ; /* IDA 0x17E6A; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_17E6C: ; /* IDA 0x17E6C; mov     ah, [bx+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_17E6F: ; /* IDA 0x17E6F; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17E71: ; /* IDA 0x17E71; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17E73: ; /* IDA 0x17E73; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17E75: ; /* IDA 0x17E75; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17E77: ; /* IDA 0x17E77; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_17E7A: ; /* IDA 0x17E7A; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17E7C: ; /* IDA 0x17E7C; jmp     short loc_17E85 */
    goto L_17E85;
L_17E7E: ; /* IDA 0x17E7E; mov     si, 7E19h */
    m->si=(uint16_t)(32281);
L_17E81: ; /* IDA 0x17E81; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_17E83: ; /* IDA 0x17E83; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_17E85: ; /* IDA 0x17E85; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_17E88: ; /* IDA 0x17E88; mov     dx, 0C0h */
    m->dx=(uint16_t)(192);
L_17E8B: ; /* IDA 0x17E8B; mov     bx, 130h */
    m->bx=(uint16_t)(304);
L_17E8E: ; /* IDA 0x17E8E; call    sub_106FD */
    dtl_call(m,1789,32401,h);
L_17E91: ; /* IDA 0x17E91; mov     si, di */
    m->si=(uint16_t)(m->di);
L_17E93: ; /* IDA 0x17E93; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17E94: ; /* IDA 0x17E94; mov     al, 0Ah */
#if KI_DETAILS_MUTATION == 2
    vg_al(m,9);
#else
    vg_al(m,(uint8_t)(10));
#endif
L_17E96: ; /* IDA 0x17E96; mul     byte ptr [si+13h] */
    dl_mul(m,dl_read8(m,m->ds,(uint16_t)(m->si+19)),8);
L_17E99: ; /* IDA 0x17E99; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17E9B: ; /* IDA 0x17E9B; mov     bx, 904h */
    m->bx=(uint16_t)(2308);
L_17E9E: ; /* IDA 0x17E9E; mov     di, 641Ah */
    m->di=(uint16_t)(25626);
L_17EA1: ; /* IDA 0x17EA1; call    sub_1062F */
    dtl_call(m,1583,32420,h);
L_17EA4: ; /* IDA 0x17EA4; mov     ax, [si+0Eh] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
L_17EA7: ; /* IDA 0x17EA7; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17EA9: ; /* IDA 0x17EA9; mov     bx, 906h */
    m->bx=(uint16_t)(2310);
L_17EAC: ; /* IDA 0x17EAC; mov     di, 6918h */
    m->di=(uint16_t)(26904);
L_17EAF: ; /* IDA 0x17EAF; call    sub_1062F */
    dtl_call(m,1583,32434,h);
L_17EB2: ; /* IDA 0x17EB2; mov     al, [si+10h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+16))));
L_17EB5: ; /* IDA 0x17EB5; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17EB7: ; /* IDA 0x17EB7; sub     ax, 64h ; 'd' */
    m->ax=(uint16_t)(ec_math(m,m->ax,100,16,0,1));
L_17EBA: ; /* IDA 0x17EBA; cwd */
    m->dx=(m->ax&0x8000)?0xffff:0;
L_17EBB: ; /* IDA 0x17EBB; mov     bx, 904h */
    m->bx=(uint16_t)(2308);
L_17EBE: ; /* IDA 0x17EBE; cmp     dh, 80h */
    ec_cmp(m,(uint8_t)(m->dx>>8),128,8);
L_17EC1: ; /* IDA 0x17EC1; jb      short loc_17EC5 */
#if KI_DETAILS_MUTATION == 3
    goto L_17EC5; /* lost negative color */
#else
    if(m->flags&1) goto L_17EC5;
#endif
L_17EC3: ; /* IDA 0x17EC3; mov     bh, 0Ah */
    gw_bh(m,(uint8_t)(10));
L_17EC5: ; /* IDA 0x17EC5; mov     di, 6E1Ah */
    m->di=(uint16_t)(28186);
L_17EC8: ; /* IDA 0x17EC8; call    sub_1062F */
    dtl_call(m,1583,32459,h);
L_17ECB: ; /* IDA 0x17ECB; mov     al, [si+11h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+17))));
L_17ECE: ; /* IDA 0x17ECE; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17ED0: ; /* IDA 0x17ED0; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17ED2: ; /* IDA 0x17ED2; mov     bx, 903h */
    m->bx=(uint16_t)(2307);
L_17ED5: ; /* IDA 0x17ED5; mov     di, 731Bh */
    m->di=(uint16_t)(29467);
L_17ED8: ; /* IDA 0x17ED8; call    sub_1062F */
    dtl_call(m,1583,32475,h);
L_17EDB: ; /* IDA 0x17EDB; mov     al, [si+16h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+22))));
L_17EDE: ; /* IDA 0x17EDE; and     al, 0Fh */
    {uint16_t result=(((uint8_t)m->ax)&(15))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17EE0: ; /* IDA 0x17EE0; mov     ah, 6 */
    vg_ah(m,(uint8_t)(6));
L_17EE2: ; /* IDA 0x17EE2; mul     ah */
    dl_mul(m,(uint8_t)(m->ax>>8),8);
L_17EE4: ; /* IDA 0x17EE4; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_17EE7: ; /* IDA 0x17EE7; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_17EE9: ; /* IDA 0x17EE9; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_17EEB: ; /* IDA 0x17EEB; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_17EED: ; /* IDA 0x17EED; mov     dx, si */
    m->dx=(uint16_t)(m->si);
L_17EEF: ; /* IDA 0x17EEF; sub     dx, 840h */
    m->dx=(uint16_t)(ec_math(m,m->dx,2112,16,0,1));
L_17EF3: ; /* IDA 0x17EF3; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_17EF5: ; /* IDA 0x17EF5; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_17EF7: ; /* IDA 0x17EF7; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_17EF9: ; /* IDA 0x17EF9; cmp     dh, [bx+3] */
    ec_cmp(m,(uint8_t)(m->dx>>8),dl_read8(m,m->ds,(uint16_t)(m->bx+3)),8);
L_17EFC: ; /* IDA 0x17EFC; jnz     short loc_17F01 */
    if(!(m->flags&0x40)) goto L_17F01;
L_17EFE: ; /* IDA 0x17EFE; mov     ax, 1Eh */
#if KI_DETAILS_MUTATION == 4
    m->ax=24;
#else
    m->ax=(uint16_t)(30);
#endif
L_17F01: ; /* IDA 0x17F01; mov     si, cs */
    m->si=(uint16_t)(m->cs);
L_17F03: ; /* IDA 0x17F03; mov     ds, si */
    m->ds=(uint16_t)(m->si);
L_17F05: ; /* IDA 0x17F05; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17F07: ; /* IDA 0x17F07; add     si, 7DF5h */
    m->si=(uint16_t)(ec_math(m,m->si,32245,16,0,0));
L_17F0B: ; /* IDA 0x17F0B; mov     dx, 0C0h */
    m->dx=(uint16_t)(192);
L_17F0E: ; /* IDA 0x17F0E; mov     bx, 120h */
    m->bx=(uint16_t)(288);
L_17F11: ; /* IDA 0x17F11; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_17F14: ; /* IDA 0x17F14; call    sub_106FD */
    dtl_call(m,1789,32535,h);
L_17F17: ; /* IDA 0x17F17; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17F18: ; /* IDA 0x17F18; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17F19: ; /* IDA 0x17F19; retn */
    return;
}
/* Original sub_17F1A; IDA 0x17F1A. */
static void dtl_body_17F1A(KiMachine16 *m,const KiEconomyHooks *h) {
L_17F1A: ; /* IDA 0x17F1A; push    ds */
    ec_push(m,m->ds);
L_17F1B: ; /* IDA 0x17F1B; push    ax */
    ec_push(m,m->ax);
L_17F1C: ; /* IDA 0x17F1C; push    bx */
    ec_push(m,m->bx);
L_17F1D: ; /* IDA 0x17F1D; push    cx */
    ec_push(m,m->cx);
L_17F1E: ; /* IDA 0x17F1E; push    dx */
    ec_push(m,m->dx);
L_17F1F: ; /* IDA 0x17F1F; push    si */
    ec_push(m,m->si);
L_17F20: ; /* IDA 0x17F20; push    di */
    ec_push(m,m->di);
L_17F21: ; /* IDA 0x17F21; mov     al, [si+16h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+22))));
L_17F24: ; /* IDA 0x17F24; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_17F26: ; /* IDA 0x17F26; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_17F28: ; /* IDA 0x17F28; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_17F2A: ; /* IDA 0x17F2A; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_17F2C: ; /* IDA 0x17F2C; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17F2E: ; /* IDA 0x17F2E; mov     dx, 1200h */
    m->dx=(uint16_t)(4608);
L_17F31: ; /* IDA 0x17F31; mul     dx */
#if KI_DETAILS_MUTATION == 5
    m->ax=(m->ax+1)%15;dl_mul(m,m->dx,16);
#else
    dl_mul(m,m->dx,16);
#endif
L_17F33: ; /* IDA 0x17F33; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_17F35: ; /* IDA 0x17F35; mov     dx, cs */
    m->dx=(uint16_t)(m->cs);
L_17F37: ; /* IDA 0x17F37; mov     ds, dx */
    m->ds=(uint16_t)(m->dx);
L_17F39: ; /* IDA 0x17F39; mov     dx, 0DFAh */
    m->dx=(uint16_t)(3578);
L_17F3C: ; /* IDA 0x17F3C; mov     bx, cs:word_1987C */
    m->bx=(uint16_t)(dl_read16(m,m->cs,39036));
L_17F41: ; /* IDA 0x17F41; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_17F43: ; /* IDA 0x17F43; mov     di, 1200h */
    m->di=(uint16_t)(4608);
L_17F46: ; /* IDA 0x17F46; call    sub_1E38C */
    dtl_call(m,58252,32585,h);
L_17F49: ; /* IDA 0x17F49; mov     ds, cs:word_1987C */
    m->ds=(uint16_t)(dl_read16(m,m->cs,39036));
L_17F4E: ; /* IDA 0x17F4E; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_17F50: ; /* IDA 0x17F50; mov     ax, 6006h */
    m->ax=(uint16_t)(24582);
L_17F53: ; /* IDA 0x17F53; mov     bx, 5A02h */
    m->bx=(uint16_t)(23042);
L_17F56: ; /* IDA 0x17F56; call    sub_1FA37 */
    dtl_call(m,64055,32601,h);
L_17F59: ; /* IDA 0x17F59; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17F5A: ; /* IDA 0x17F5A; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17F5B: ; /* IDA 0x17F5B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17F5C: ; /* IDA 0x17F5C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17F5D: ; /* IDA 0x17F5D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17F5E: ; /* IDA 0x17F5E; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17F5F: ; /* IDA 0x17F5F; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17F60: ; /* IDA 0x17F60; retn */
    return;
}
/* Original sub_17F61; IDA 0x17F61. */
static void dtl_body_17F61(KiMachine16 *m,const KiEconomyHooks *h) {
L_17F61: ; /* IDA 0x17F61; push    ds */
    ec_push(m,m->ds);
L_17F62: ; /* IDA 0x17F62; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_17F64: ; /* IDA 0x17F64; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_17F66: ; /* IDA 0x17F66; mov     al, 3Dh ; '=' */
    vg_al(m,(uint8_t)(61));
L_17F68: ; /* IDA 0x17F68; mov     dx, 0 */
    m->dx=(uint16_t)(0);
L_17F6B: ; /* IDA 0x17F6B; mov     bx, 0Fh */
    m->bx=(uint16_t)(15);
L_17F6E: ; /* IDA 0x17F6E; mov     cx, 810h */
    m->cx=(uint16_t)(2064);
L_17F71: ; /* IDA 0x17F71; call    sub_1895D */
    dtl_call(m,35165,32628,h);
L_17F74: ; /* IDA 0x17F74; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_17F76: ; /* IDA 0x17F76; mov     dx, 8 */
    m->dx=(uint16_t)(8);
L_17F79: ; /* IDA 0x17F79; mov     bx, 118h */
    m->bx=(uint16_t)(280);
L_17F7C: ; /* IDA 0x17F7C; call    sub_10337 */
    dtl_call(m,823,32639,h);
L_17F7F: ; /* IDA 0x17F7F; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17F80: ; /* IDA 0x17F80; retn */
    return;
}
/* Original sub_17F81; IDA 0x17F81. */
static void dtl_body_17F81(KiMachine16 *m,const KiEconomyHooks *h) {
L_17F81: ; /* IDA 0x17F81; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17F83: ; /* IDA 0x17F83; mov     dx, 0 */
    m->dx=(uint16_t)(0);
L_17F86: ; /* IDA 0x17F86; mov     bx, 0Fh */
    m->bx=(uint16_t)(15);
L_17F89: ; /* IDA 0x17F89; mov     cx, 810h */
    m->cx=(uint16_t)(2064);
L_17F8C: ; /* IDA 0x17F8C; call    sub_1895D */
    dtl_call(m,35165,32655,h);
L_17F8F: ; /* IDA 0x17F8F; retn */
    return;
}
/* Original sub_1807B; IDA 0x1807B. */
static void dtl_body_1807B(KiMachine16 *m,const KiEconomyHooks *h) {
L_1807B: ; /* IDA 0x1807B; push    ds */
    ec_push(m,m->ds);
L_1807C: ; /* IDA 0x1807C; push    si */
    ec_push(m,m->si);
L_1807D: ; /* IDA 0x1807D; push    di */
    ec_push(m,m->di);
L_1807E: ; /* IDA 0x1807E; mov     di, si */
    m->di=(uint16_t)(m->si);
L_18080: ; /* IDA 0x18080; mov     al, 1Fh */
    vg_al(m,(uint8_t)(31));
L_18082: ; /* IDA 0x18082; mov     dx, 1Bh */
    m->dx=(uint16_t)(27);
L_18085: ; /* IDA 0x18085; mov     bx, 0Ah */
    m->bx=(uint16_t)(10);
L_18088: ; /* IDA 0x18088; mov     cx, 0D0Dh */
    m->cx=(uint16_t)(3341);
L_1808B: ; /* IDA 0x1808B; call    sub_1895D */
    dtl_call(m,35165,32910,h);
L_1808E: ; /* IDA 0x1808E; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_18090: ; /* IDA 0x18090; mov     dx, 1B8h */
    m->dx=(uint16_t)(440);
L_18093: ; /* IDA 0x18093; mov     bx, 0C8h */
    m->bx=(uint16_t)(200);
L_18096: ; /* IDA 0x18096; call    sub_10337 */
    dtl_call(m,823,32921,h);
L_18099: ; /* IDA 0x18099; mov     ah, [si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_1809C: ; /* IDA 0x1809C; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1809E: ; /* IDA 0x1809E; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180A0: ; /* IDA 0x180A0; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180A2: ; /* IDA 0x180A2; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180A4: ; /* IDA 0x180A4; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_180A6: ; /* IDA 0x180A6; mov     al, [si+4241h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+16961))));
L_180AA: ; /* IDA 0x180AA; mov     bx, 3EB7h */
    m->bx=(uint16_t)(16055);
L_180AD: ; /* IDA 0x180AD; call    sub_107D2 */
    dtl_call(m,2002,32944,h);
L_180B0: ; /* IDA 0x180B0; mov     dx, 240h */
    m->dx=(uint16_t)(576);
L_180B3: ; /* IDA 0x180B3; mov     bx, 0D0h */
    m->bx=(uint16_t)(208);
L_180B6: ; /* IDA 0x180B6; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_180B9: ; /* IDA 0x180B9; add     si, 4242h */
    m->si=(uint16_t)(ec_math(m,m->si,16962,16,0,0));
L_180BD: ; /* IDA 0x180BD; call    sub_106FD */
    dtl_call(m,1789,32960,h);
L_180C0: ; /* IDA 0x180C0; mov     bh, [di+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+1))));
L_180C3: ; /* IDA 0x180C3; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_180C5: ; /* IDA 0x180C5; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_180C7: ; /* IDA 0x180C7; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_180C9: ; /* IDA 0x180C9; mov     ah, [bx+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_180CC: ; /* IDA 0x180CC; push    bx */
    ec_push(m,m->bx);
L_180CD: ; /* IDA 0x180CD; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_180CF: ; /* IDA 0x180CF; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180D1: ; /* IDA 0x180D1; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180D3: ; /* IDA 0x180D3; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180D5: ; /* IDA 0x180D5; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_180D8: ; /* IDA 0x180D8; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_180DA: ; /* IDA 0x180DA; mov     dx, 240h */
    m->dx=(uint16_t)(576);
L_180DD: ; /* IDA 0x180DD; mov     bx, 0E0h */
    m->bx=(uint16_t)(224);
L_180E0: ; /* IDA 0x180E0; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_180E3: ; /* IDA 0x180E3; call    sub_106FD */
    dtl_call(m,1789,32998,h);
L_180E6: ; /* IDA 0x180E6; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_180E7: ; /* IDA 0x180E7; mov     ah, [bx+3] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_180EA: ; /* IDA 0x180EA; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_180EC: ; /* IDA 0x180EC; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180EE: ; /* IDA 0x180EE; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180F0: ; /* IDA 0x180F0; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_180F2: ; /* IDA 0x180F2; add     ax, 842h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2114,16,0,0));
L_180F5: ; /* IDA 0x180F5; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_180F7: ; /* IDA 0x180F7; mov     dx, 240h */
    m->dx=(uint16_t)(576);
L_180FA: ; /* IDA 0x180FA; mov     bx, 0F0h */
    m->bx=(uint16_t)(240);
L_180FD: ; /* IDA 0x180FD; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_18100: ; /* IDA 0x18100; call    sub_106FD */
    dtl_call(m,1789,33027,h);
L_18103: ; /* IDA 0x18103; mov     si, di */
    m->si=(uint16_t)(m->di);
L_18105: ; /* IDA 0x18105; mov     ax, [si+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+4)));
L_18108: ; /* IDA 0x18108; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_1810B: ; /* IDA 0x1810B; mul     dx */
    dl_mul(m,m->dx,16);
L_1810D: ; /* IDA 0x1810D; mov     bx, 0F04h */
#if KI_DETAILS_MUTATION == 6
    m->bx=0x0904;
#else
    m->bx=(uint16_t)(3844);
#endif
L_18110: ; /* IDA 0x18110; mov     di, 5543h */
    m->di=(uint16_t)(21827);
L_18113: ; /* IDA 0x18113; call    sub_1062F */
    dtl_call(m,1583,33046,h);
L_18116: ; /* IDA 0x18116; mov     al, [si+6] */
#if KI_DETAILS_MUTATION == 7
    m->ax=dl_read16(m,m->ds,(uint16_t)(m->si+6));
#else
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+6))));
#endif
L_18119: ; /* IDA 0x18119; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_1811B: ; /* IDA 0x1811B; mov     ah, dl */
#if KI_DETAILS_MUTATION == 7
    ; /* retained word morale */
#else
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
#endif
L_1811D: ; /* IDA 0x1811D; mov     bx, 0F03h */
    m->bx=(uint16_t)(3843);
L_18120: ; /* IDA 0x18120; mov     di, 5549h */
    m->di=(uint16_t)(21833);
L_18123: ; /* IDA 0x18123; call    sub_1062F */
    dtl_call(m,1583,33062,h);
L_18126: ; /* IDA 0x18126; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_18127: ; /* IDA 0x18127; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_18128: ; /* IDA 0x18128; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_18129: ; /* IDA 0x18129; retn */
    return;
}
/* Original sub_1812A; IDA 0x1812A. */
static void dtl_body_1812A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1812A: ; /* IDA 0x1812A; push    bp */
    ec_push(m,m->bp);
L_1812B: ; /* IDA 0x1812B; push    ds */
    ec_push(m,m->ds);
L_1812C: ; /* IDA 0x1812C; push    si */
    ec_push(m,m->si);
L_1812D: ; /* IDA 0x1812D; add     si, 28h ; '(' */
    m->si=(uint16_t)(ec_math(m,m->si,40,16,0,0));
L_18130: ; /* IDA 0x18130; mov     di, 5A48h */
    m->di=(uint16_t)(23112);
L_18133: ; /* IDA 0x18133; mov     ax, ds */
    m->ax=(uint16_t)(m->ds);
L_18135: ; /* IDA 0x18135; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_18137: ; /* IDA 0x18137; mov     ds, cs:word_10D50 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3408));
L_1813C: ; /* IDA 0x1813C; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_1813F: ; /* IDA 0x1813F; push    si */
    ec_push(m,m->si);
L_18140: ; /* IDA 0x18140; mov     ah, es:[si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->si+2))));
L_18144: ; /* IDA 0x18144; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_18146: ; /* IDA 0x18146; dec     ah */
    vg_ah(m,(uint8_t)(st_dec(m,(uint8_t)(m->ax>>8),8)));
L_18148: ; /* IDA 0x18148; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1814A: ; /* IDA 0x1814A; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1814C: ; /* IDA 0x1814C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1814E: ; /* IDA 0x1814E; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_18150: ; /* IDA 0x18150; add     si, 12C0h */
#if KI_DETAILS_MUTATION == 8
    m->si=ec_math(m,m->si,0x15c0,16,0,0);
#else
    m->si=(uint16_t)(ec_math(m,m->si,4800,16,0,0));
#endif
L_18154: ; /* IDA 0x18154; mov     bx, di */
    m->bx=(uint16_t)(m->di);
L_18156: ; /* IDA 0x18156; sub     bx, 7 */
#if KI_DETAILS_MUTATION == 9
    m->bx=ec_math(m,m->bx,4,16,0,1);
#else
    m->bx=(uint16_t)(ec_math(m,m->bx,7,16,0,1));
#endif
L_18159: ; /* IDA 0x18159; mov     ax, 1003h */
    m->ax=(uint16_t)(4099);
L_1815C: ; /* IDA 0x1815C; call    sub_1F9B0 */
    dtl_call(m,63920,33119,h);
L_1815F: ; /* IDA 0x1815F; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_18160: ; /* IDA 0x18160; mov     al, es:[si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->si+1))));
L_18164: ; /* IDA 0x18164; mov     ah, 0Ah */
    vg_ah(m,(uint8_t)(10));
L_18166: ; /* IDA 0x18166; mul     ah */
    dl_mul(m,(uint8_t)(m->ax>>8),8);
L_18168: ; /* IDA 0x18168; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_1816A: ; /* IDA 0x1816A; mov     bx, 904h */
    m->bx=(uint16_t)(2308);
L_1816D: ; /* IDA 0x1816D; call    sub_1062F */
    dtl_call(m,1583,33136,h);
L_18170: ; /* IDA 0x18170; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_18173: ; /* IDA 0x18173; add     di, 500h */
    m->di=(uint16_t)(ec_math(m,m->di,1280,16,0,0));
L_18177: ; /* IDA 0x18177; loop    loc_1813F */
    if(--m->cx) goto L_1813F;
L_18179: ; /* IDA 0x18179; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1817A: ; /* IDA 0x1817A; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1817B: ; /* IDA 0x1817B; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1817C: ; /* IDA 0x1817C; retn */
    return;
}
/* Original sub_1817D; IDA 0x1817D. */
static void dtl_body_1817D(KiMachine16 *m,const KiEconomyHooks *h) {
L_1817D: ; /* IDA 0x1817D; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1817F: ; /* IDA 0x1817F; mov     dx, 1Bh */
    m->dx=(uint16_t)(27);
L_18182: ; /* IDA 0x18182; mov     bx, 0Ah */
    m->bx=(uint16_t)(10);
L_18185: ; /* IDA 0x18185; mov     cx, 0D0Dh */
    m->cx=(uint16_t)(3341);
L_18188: ; /* IDA 0x18188; call    sub_1895D */
    dtl_call(m,35165,33163,h);
L_1818B: ; /* IDA 0x1818B; test    cs:byte_198A6, 2 */
    {uint16_t result=((dl_read8(m,m->cs,39078))&(2))&255;ec_logic(m,result,8);}
L_18191: ; /* IDA 0x18191; jz      short locret_181A2 */
#if KI_DETAILS_MUTATION == 10
    goto L_181A2; /* lost own-faction restoration */
#else
    if(m->flags&0x40) goto L_181A2;
#endif
L_18193: ; /* IDA 0x18193; push    ds */
    ec_push(m,m->ds);
L_18194: ; /* IDA 0x18194; push    es */
    ec_push(m,m->es);
L_18195: ; /* IDA 0x18195; push    si */
    ec_push(m,m->si);
L_18196: ; /* IDA 0x18196; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_18198: ; /* IDA 0x18198; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1819A: ; /* IDA 0x1819A; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1819C: ; /* IDA 0x1819C; call    sub_15E1E */
    dtl_call(m,24094,33183,h);
L_1819F: ; /* IDA 0x1819F; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_181A0: ; /* IDA 0x181A0; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_181A1: ; /* IDA 0x181A1; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_181A2: ; /* IDA 0x181A2; retn */
    return;
}
int ki_details_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x5e1e:dtl_body_15E1E(m,h);break;
case 0x5e2d:dtl_body_15E2D(m,h);break;
case 0x7e1f:dtl_body_17E1F(m,h);break;
case 0x7e4a:dtl_body_17E4A(m,h);break;
case 0x7f1a:dtl_body_17F1A(m,h);break;
case 0x7f61:dtl_body_17F61(m,h);break;
case 0x7f81:dtl_body_17F81(m,h);break;
case 0x807b:dtl_body_1807B(m,h);break;
case 0x812a:dtl_body_1812A(m,h);break;
case 0x817d:dtl_body_1817D(m,h);break;
default:return 0;}return 1;}
static int dtl_known(uint16_t t) {switch(t) {case 0x5e1e: case 0x5e2d: case 0x7e1f: case 0x7e4a: case 0x7f1a: case 0x7f61: case 0x7f81: case 0x807b: case 0x812a: case 0x817d: return 1;default:return 0;}}
