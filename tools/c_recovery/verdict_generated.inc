/* Original native C ruler sortie / formation / sidebar; re/117, spec/237. */
/* Original sub_1461D; IDA 0x1461D. */
static void vd_body_1461D(KiMachine16 *m,const KiEconomyHooks *h) {
L_1461D: ; /* IDA 0x1461D; push    es */
    ec_push(m,m->es);
L_1461E: ; /* IDA 0x1461E; push    ax */
    ec_push(m,m->ax);
L_1461F: ; /* IDA 0x1461F; push    bx */
    ec_push(m,m->bx);
L_14620: ; /* IDA 0x14620; push    di */
    ec_push(m,m->di);
L_14621: ; /* IDA 0x14621; mov     ax, ds */
    m->ax=(uint16_t)(m->ds);
L_14623: ; /* IDA 0x14623; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_14625: ; /* IDA 0x14625; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14628: ; /* IDA 0x14628; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1462A: ; /* IDA 0x1462A; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1462C: ; /* IDA 0x1462C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1462E: ; /* IDA 0x1462E; add     ax, 4 */
    m->ax=(uint16_t)(ec_math(m,m->ax,4,16,0,0));
L_14631: ; /* IDA 0x14631; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_14633: ; /* IDA 0x14633; call    sub_14717 */
    vd_call(m,18199,17974,h);
L_14636: ; /* IDA 0x14636; call    sub_14698 */
    vd_call(m,18072,17977,h);
L_14639: ; /* IDA 0x14639; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1463A: ; /* IDA 0x1463A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1463B: ; /* IDA 0x1463B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1463C: ; /* IDA 0x1463C; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1463D: ; /* IDA 0x1463D; retn */
    return;
}
/* Original sub_14698; IDA 0x14698. */
static void vd_body_14698(KiMachine16 *m,const KiEconomyHooks *h) {
L_14698: ; /* IDA 0x14698; push    ax */
    ec_push(m,m->ax);
L_14699: ; /* IDA 0x14699; push    bx */
    ec_push(m,m->bx);
L_1469A: ; /* IDA 0x1469A; push    cx */
    ec_push(m,m->cx);
L_1469B: ; /* IDA 0x1469B; push    dx */
    ec_push(m,m->dx);
L_1469C: ; /* IDA 0x1469C; push    si */
    ec_push(m,m->si);
L_1469D: ; /* IDA 0x1469D; push    bp */
    ec_push(m,m->bp);
L_1469E: ; /* IDA 0x1469E; add     si, 28h ; '(' */
    m->si=(uint16_t)(ec_math(m,m->si,40,16,0,0));
L_146A1: ; /* IDA 0x146A1; sub     sp, 8 */
    m->sp=(uint16_t)(ec_math(m,m->sp,8,16,0,1));
L_146A4: ; /* IDA 0x146A4; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_146A6: ; /* IDA 0x146A6; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_146A8: ; /* IDA 0x146A8; mov     [bp+0], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->ax);
L_146AB: ; /* IDA 0x146AB; mov     [bp+2], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->ax);
L_146AE: ; /* IDA 0x146AE; mov     [bp+4], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp+4),m->ax);
L_146B1: ; /* IDA 0x146B1; mov     [bp+6], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp+6),m->ax);
L_146B4: ; /* IDA 0x146B4; push    si */
    ec_push(m,m->si);
L_146B5: ; /* IDA 0x146B5; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_146B8: ; /* IDA 0x146B8; mov     bl, [si+2] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_146BB: ; /* IDA 0x146BB; cmp     bl, 4 */
    ec_cmp(m,(uint8_t)m->bx,4,8);
L_146BE: ; /* IDA 0x146BE; jz      short loc_146CA */
    if(m->flags&0x40) goto L_146CA;
L_146C0: ; /* IDA 0x146C0; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_146C2: ; /* IDA 0x146C2; dec     bx */
    m->bx=(uint16_t)(st_dec(m,m->bx,16));
L_146C3: ; /* IDA 0x146C3; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_146C5: ; /* IDA 0x146C5; add     bx, bp */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->bp,16,0,0));
L_146C7: ; /* IDA 0x146C7; inc     word ptr ss:[bx] */
    dl_write16(m,m->ss,(uint16_t)(m->bx),ec_inc(m,dl_read16(m,m->ss,(uint16_t)(m->bx)),16));
L_146CA: ; /* IDA 0x146CA; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_146CD: ; /* IDA 0x146CD; loop    loc_146B8 */
    if(--m->cx) goto L_146B8;
L_146CF: ; /* IDA 0x146CF; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_146D0: ; /* IDA 0x146D0; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_146D3: ; /* IDA 0x146D3; push    bp */
    ec_push(m,m->bp);
L_146D4: ; /* IDA 0x146D4; mov     bl, [si+2] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_146D7: ; /* IDA 0x146D7; cmp     bl, 4 */
    ec_cmp(m,(uint8_t)m->bx,4,8);
L_146DA: ; /* IDA 0x146DA; jz      short loc_14702 */
    if(m->flags&0x40) goto L_14702;
L_146DC: ; /* IDA 0x146DC; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_146DE: ; /* IDA 0x146DE; dec     bx */
    m->bx=(uint16_t)(st_dec(m,m->bx,16));
L_146DF: ; /* IDA 0x146DF; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_146E1: ; /* IDA 0x146E1; add     bp, bx */
    m->bp=(uint16_t)(ec_math(m,m->bp,m->bx,16,0,0));
L_146E3: ; /* IDA 0x146E3; add     bx, di */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->di,16,0,0));
L_146E5: ; /* IDA 0x146E5; mov     ax, es:[bx] */
    m->ax=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx)));
L_146E8: ; /* IDA 0x146E8; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_146EA: ; /* IDA 0x146EA; div     word ptr [bp+0] */
    dl_div(m,dl_read16(m,m->ss,(uint16_t)(m->bp)),16);
L_146ED: ; /* IDA 0x146ED; dec     word ptr [bp+0] */
    dl_write16(m,m->ss,(uint16_t)(m->bp),st_dec(m,dl_read16(m,m->ss,(uint16_t)(m->bp)),16));
L_146F0: ; /* IDA 0x146F0; add     ax, dx */
#if KI_VERDICT_MUTATION == 5
    ; /* lost division remainder */
#else
    m->ax=(uint16_t)(ec_math(m,m->ax,m->dx,16,0,0));
#endif
L_146F2: ; /* IDA 0x146F2; cmp     ax, 64h ; 'd' */
    ec_cmp(m,m->ax,100,16);
L_146F5: ; /* IDA 0x146F5; jbe     short loc_146FC */
    if(m->flags&(1|0x40)) goto L_146FC;
L_146F7: ; /* IDA 0x146F7; mov     ax, 64h ; 'd' */
#if KI_VERDICT_MUTATION == 6
    m->ax=(uint16_t)(99);
#else
    m->ax=(uint16_t)(100);
#endif
L_146FA: ; /* IDA 0x146FA; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_146FC: ; /* IDA 0x146FC; sub     es:[bx], ax */
    dl_write16(m,m->es,(uint16_t)(m->bx),ec_math(m,dl_read16(m,m->es,(uint16_t)(m->bx)),m->ax,16,0,1));
L_146FF: ; /* IDA 0x146FF; mov     [si+1], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+1),(uint8_t)m->ax);
L_14702: ; /* IDA 0x14702; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_14703: ; /* IDA 0x14703; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_14706: ; /* IDA 0x14706; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_14708: ; /* IDA 0x14708; jnz     short loc_146D3 */
    if(!(m->flags&0x40)) goto L_146D3;
L_1470A: ; /* IDA 0x1470A; add     sp, 8 */
    m->sp=(uint16_t)(ec_math(m,m->sp,8,16,0,0));
L_1470D: ; /* IDA 0x1470D; cmp     ch, 6 */
    ec_cmp(m,(uint8_t)(m->cx>>8),6,8);
L_14710: ; /* IDA 0x14710; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_14711: ; /* IDA 0x14711; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14712: ; /* IDA 0x14712; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_14713: ; /* IDA 0x14713; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_14714: ; /* IDA 0x14714; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14715: ; /* IDA 0x14715; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_14716: ; /* IDA 0x14716; retn */
    return;
}
/* Original sub_14717; IDA 0x14717. */
static void vd_body_14717(KiMachine16 *m,const KiEconomyHooks *h) {
L_14717: ; /* IDA 0x14717; push    ax */
    ec_push(m,m->ax);
L_14718: ; /* IDA 0x14718; push    bx */
    ec_push(m,m->bx);
L_14719: ; /* IDA 0x14719; push    cx */
    ec_push(m,m->cx);
L_1471A: ; /* IDA 0x1471A; push    dx */
    ec_push(m,m->dx);
L_1471B: ; /* IDA 0x1471B; push    si */
    ec_push(m,m->si);
L_1471C: ; /* IDA 0x1471C; push    di */
    ec_push(m,m->di);
L_1471D: ; /* IDA 0x1471D; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_14720: ; /* IDA 0x14720; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_14722: ; /* IDA 0x14722; add     si, 28h ; '(' */
    m->si=(uint16_t)(ec_math(m,m->si,40,16,0,0));
L_14725: ; /* IDA 0x14725; mov     bl, [si+2] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_14728: ; /* IDA 0x14728; cmp     bl, 4 */
    ec_cmp(m,(uint8_t)m->bx,4,8);
L_1472B: ; /* IDA 0x1472B; jz      short loc_1473E */
    if(m->flags&0x40) goto L_1473E;
L_1472D: ; /* IDA 0x1472D; dec     bx */
    m->bx=(uint16_t)(st_dec(m,m->bx,16));
L_1472E: ; /* IDA 0x1472E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_14730: ; /* IDA 0x14730; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_14732: ; /* IDA 0x14732; xchg    al, [si+1] */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->ds,(uint16_t)(m->si+1));vg_al(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->si+1),a);}
L_14735: ; /* IDA 0x14735; mov     dx, es:[bx+di] */
    m->dx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+m->di)));
L_14738: ; /* IDA 0x14738; call    sub_155EC */
    vd_call(m,21996,18235,h);
L_1473B: ; /* IDA 0x1473B; mov     es:[bx+di], ax */
    dl_write16(m,m->es,(uint16_t)(m->bx+m->di),m->ax);
L_1473E: ; /* IDA 0x1473E; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_14741: ; /* IDA 0x14741; loop    loc_14725 */
    if(--m->cx) goto L_14725;
L_14743: ; /* IDA 0x14743; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14744: ; /* IDA 0x14744; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14745: ; /* IDA 0x14745; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_14746: ; /* IDA 0x14746; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_14747: ; /* IDA 0x14747; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14748: ; /* IDA 0x14748; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_14749: ; /* IDA 0x14749; retn */
    return;
}
/* Original sub_15E80; IDA 0x15E80. */
static void vd_body_15E80(KiMachine16 *m,const KiEconomyHooks *h) {
L_15E80: ; /* IDA 0x15E80; test    cs:byte_198A6, 2 */
#if KI_VERDICT_MUTATION == 10
    {uint16_t result=((dl_read8(m,m->cs,39078))&(0))&255;ec_logic(m,result,8);}
#else
    {uint16_t result=((dl_read8(m,m->cs,39078))&(2))&255;ec_logic(m,result,8);}
#endif
L_15E86: ; /* IDA 0x15E86; jnz     short loc_15E89 */
    if(!(m->flags&0x40)) goto L_15E89;
L_15E88: ; /* IDA 0x15E88; retn */
    return;
L_15E89: ; /* IDA 0x15E89; push    ds */
    ec_push(m,m->ds);
L_15E8A: ; /* IDA 0x15E8A; push    si */
    ec_push(m,m->si);
L_15E8B: ; /* IDA 0x15E8B; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_15E90: ; /* IDA 0x15E90; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_15E95: ; /* IDA 0x15E95; add     si, 0 */
    m->si=(uint16_t)(ec_math(m,m->si,0,16,0,0));
L_15E98: ; /* IDA 0x15E98; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_15E9A: ; /* IDA 0x15E9A; jnb     short loc_15E9F */
    if(!(m->flags&1)) goto L_15E9F;
L_15E9C: ; /* IDA 0x15E9C; call    sub_15EB7 */
    vd_call(m,24247,24223,h);
L_15E9F: ; /* IDA 0x15E9F; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_15EA1: ; /* IDA 0x15EA1; jnb     short loc_15EA6 */
    if(!(m->flags&1)) goto L_15EA6;
L_15EA3: ; /* IDA 0x15EA3; call    sub_15F27 */
    vd_call(m,24359,24230,h);
L_15EA6: ; /* IDA 0x15EA6; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_15EA8: ; /* IDA 0x15EA8; jnb     short loc_15EAD */
    if(!(m->flags&1)) goto L_15EAD;
L_15EAA: ; /* IDA 0x15EAA; call    sub_15F5D */
    vd_call(m,24413,24237,h);
L_15EAD: ; /* IDA 0x15EAD; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_15EAF: ; /* IDA 0x15EAF; jnb     short loc_15EB4 */
    if(!(m->flags&1)) goto L_15EB4;
L_15EB1: ; /* IDA 0x15EB1; call    sub_15F7F */
    vd_call(m,24447,24244,h);
L_15EB4: ; /* IDA 0x15EB4; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_15EB5: ; /* IDA 0x15EB5; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_15EB6: ; /* IDA 0x15EB6; retn */
    return;
}
/* Original sub_15EB7; IDA 0x15EB7. */
static void vd_body_15EB7(KiMachine16 *m,const KiEconomyHooks *h) {
L_15EB7: ; /* IDA 0x15EB7; push    ax */
    ec_push(m,m->ax);
L_15EB8: ; /* IDA 0x15EB8; push    bx */
    ec_push(m,m->bx);
L_15EB9: ; /* IDA 0x15EB9; push    dx */
    ec_push(m,m->dx);
L_15EBA: ; /* IDA 0x15EBA; push    ds */
    ec_push(m,m->ds);
L_15EBB: ; /* IDA 0x15EBB; push    si */
    ec_push(m,m->si);
L_15EBC: ; /* IDA 0x15EBC; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_15EBF: ; /* IDA 0x15EBF; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_15EC1: ; /* IDA 0x15EC1; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15EC3: ; /* IDA 0x15EC3; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15EC5: ; /* IDA 0x15EC5; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15EC7: ; /* IDA 0x15EC7; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_15EC9: ; /* IDA 0x15EC9; mov     al, [si+4241h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+16961))));
L_15ECD: ; /* IDA 0x15ECD; mov     bx, 3EB7h */
    m->bx=(uint16_t)(16055);
L_15ED0: ; /* IDA 0x15ED0; call    sub_107D2 */
    vd_call(m,2002,24275,h);
L_15ED3: ; /* IDA 0x15ED3; add     si, 4242h */
    m->si=(uint16_t)(ec_math(m,m->si,16962,16,0,0));
L_15ED7: ; /* IDA 0x15ED7; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_15EDA: ; /* IDA 0x15EDA; mov     dx, 240h */
    m->dx=(uint16_t)(576);
L_15EDD: ; /* IDA 0x15EDD; mov     bx, 0D0h */
    m->bx=(uint16_t)(208);
L_15EE0: ; /* IDA 0x15EE0; call    sub_106FD */
    vd_call(m,1789,24291,h);
L_15EE3: ; /* IDA 0x15EE3; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_15EE4: ; /* IDA 0x15EE4; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_15EE5: ; /* IDA 0x15EE5; push    ds */
    ec_push(m,m->ds);
L_15EE6: ; /* IDA 0x15EE6; push    si */
    ec_push(m,m->si);
L_15EE7: ; /* IDA 0x15EE7; mov     ah, [si+3] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+3))));
L_15EEA: ; /* IDA 0x15EEA; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_15EEC: ; /* IDA 0x15EEC; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15EEE: ; /* IDA 0x15EEE; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15EF0: ; /* IDA 0x15EF0; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15EF2: ; /* IDA 0x15EF2; add     ax, 842h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2114,16,0,0));
L_15EF5: ; /* IDA 0x15EF5; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_15EF7: ; /* IDA 0x15EF7; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_15EFA: ; /* IDA 0x15EFA; mov     dx, 240h */
    m->dx=(uint16_t)(576);
L_15EFD: ; /* IDA 0x15EFD; mov     bx, 0E0h */
    m->bx=(uint16_t)(224);
L_15F00: ; /* IDA 0x15F00; call    sub_106FD */
    vd_call(m,1789,24323,h);
L_15F03: ; /* IDA 0x15F03; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_15F04: ; /* IDA 0x15F04; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_15F05: ; /* IDA 0x15F05; push    si */
    ec_push(m,m->si);
L_15F06: ; /* IDA 0x15F06; mov     ah, [si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_15F09: ; /* IDA 0x15F09; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_15F0B: ; /* IDA 0x15F0B; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15F0D: ; /* IDA 0x15F0D; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15F0F: ; /* IDA 0x15F0F; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_15F11: ; /* IDA 0x15F11; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_15F14: ; /* IDA 0x15F14; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_15F16: ; /* IDA 0x15F16; mov     ax, 0F03h */
    m->ax=(uint16_t)(3843);
L_15F19: ; /* IDA 0x15F19; mov     dx, 240h */
    m->dx=(uint16_t)(576);
L_15F1C: ; /* IDA 0x15F1C; mov     bx, 0F0h */
    m->bx=(uint16_t)(240);
L_15F1F: ; /* IDA 0x15F1F; call    sub_106FD */
    vd_call(m,1789,24354,h);
L_15F22: ; /* IDA 0x15F22; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_15F23: ; /* IDA 0x15F23; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15F24: ; /* IDA 0x15F24; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15F25: ; /* IDA 0x15F25; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15F26: ; /* IDA 0x15F26; retn */
    return;
}
/* Original sub_15F27; IDA 0x15F27. */
static void vd_body_15F27(KiMachine16 *m,const KiEconomyHooks *h) {
L_15F27: ; /* IDA 0x15F27; push    ax */
    ec_push(m,m->ax);
L_15F28: ; /* IDA 0x15F28; push    bx */
    ec_push(m,m->bx);
L_15F29: ; /* IDA 0x15F29; push    cx */
    ec_push(m,m->cx);
L_15F2A: ; /* IDA 0x15F2A; push    dx */
    ec_push(m,m->dx);
L_15F2B: ; /* IDA 0x15F2B; mov     al, cs:byte_10D00 */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3328)));
L_15F2F: ; /* IDA 0x15F2F; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_15F31: ; /* IDA 0x15F31; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_15F33: ; /* IDA 0x15F33; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_15F35: ; /* IDA 0x15F35; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_15F37: ; /* IDA 0x15F37; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_15F39: ; /* IDA 0x15F39; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_15F3B: ; /* IDA 0x15F3B; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_15F3D: ; /* IDA 0x15F3D; add     cx, ax */
    m->cx=(uint16_t)(ec_math(m,m->cx,m->ax,16,0,0));
L_15F3F: ; /* IDA 0x15F3F; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_15F41: ; /* IDA 0x15F41; add     ax, cx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->cx,16,0,0));
L_15F43: ; /* IDA 0x15F43; add     ax, 9Fh */
    m->ax=(uint16_t)(ec_math(m,m->ax,159,16,0,0));
L_15F46: ; /* IDA 0x15F46; mov     ch, 0A0h */
    vg_ch(m,(uint8_t)(160));
L_15F48: ; /* IDA 0x15F48; div     ch */
    dl_div(m,(uint8_t)(m->cx>>8),8);
L_15F4A: ; /* IDA 0x15F4A; mov     cl, al */
    vg_cl(m,(uint8_t)((uint8_t)m->ax));
L_15F4C: ; /* IDA 0x15F4C; mov     dx, 1C8h */
    m->dx=(uint16_t)(456);
L_15F4F: ; /* IDA 0x15F4F; mov     bx, 124h */
    m->bx=(uint16_t)(292);
L_15F52: ; /* IDA 0x15F52; mov     ax, 0A00h */
    m->ax=(uint16_t)(2560);
L_15F55: ; /* IDA 0x15F55; call    sub_10AAA */
    vd_call(m,2730,24408,h);
L_15F58: ; /* IDA 0x15F58; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15F59: ; /* IDA 0x15F59; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_15F5A: ; /* IDA 0x15F5A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15F5B: ; /* IDA 0x15F5B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15F5C: ; /* IDA 0x15F5C; retn */
    return;
}
/* Original sub_15F5D; IDA 0x15F5D. */
static void vd_body_15F5D(KiMachine16 *m,const KiEconomyHooks *h) {
L_15F5D: ; /* IDA 0x15F5D; push    ax */
    ec_push(m,m->ax);
L_15F5E: ; /* IDA 0x15F5E; push    bx */
    ec_push(m,m->bx);
L_15F5F: ; /* IDA 0x15F5F; push    dx */
    ec_push(m,m->dx);
L_15F60: ; /* IDA 0x15F60; push    di */
    ec_push(m,m->di);
L_15F61: ; /* IDA 0x15F61; mov     al, [si+22h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+34))));
L_15F64: ; /* IDA 0x15F64; cbw */
    m->ax=(uint16_t)(int16_t)(int8_t)m->ax;
L_15F65: ; /* IDA 0x15F65; mov     dx, ax */
    m->dx=(uint16_t)(m->ax);
L_15F67: ; /* IDA 0x15F67; mov     ax, [si+20h] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+32)));
L_15F6A: ; /* IDA 0x15F6A; mov     di, 61C6h */
    m->di=(uint16_t)(25030);
L_15F6D: ; /* IDA 0x15F6D; mov     bx, 0F07h */
    m->bx=(uint16_t)(3847);
L_15F70: ; /* IDA 0x15F70; cmp     dh, 80h */
    ec_cmp(m,(uint8_t)(m->dx>>8),128,8);
L_15F73: ; /* IDA 0x15F73; jb      short loc_15F77 */
    if(m->flags&1) goto L_15F77;
L_15F75: ; /* IDA 0x15F75; mov     bh, 0Ah */
    gw_bh(m,(uint8_t)(10));
L_15F77: ; /* IDA 0x15F77; call    sub_1062F */
    vd_call(m,1583,24442,h);
L_15F7A: ; /* IDA 0x15F7A; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_15F7B: ; /* IDA 0x15F7B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15F7C: ; /* IDA 0x15F7C; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15F7D: ; /* IDA 0x15F7D; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15F7E: ; /* IDA 0x15F7E; retn */
    return;
}
/* Original sub_15F7F; IDA 0x15F7F. */
static void vd_body_15F7F(KiMachine16 *m,const KiEconomyHooks *h) {
L_15F7F: ; /* IDA 0x15F7F; push    ax */
    ec_push(m,m->ax);
L_15F80: ; /* IDA 0x15F80; push    bx */
    ec_push(m,m->bx);
L_15F81: ; /* IDA 0x15F81; push    cx */
    ec_push(m,m->cx);
L_15F82: ; /* IDA 0x15F82; push    dx */
    ec_push(m,m->dx);
L_15F83: ; /* IDA 0x15F83; push    si */
    ec_push(m,m->si);
L_15F84: ; /* IDA 0x15F84; push    di */
    ec_push(m,m->di);
L_15F85: ; /* IDA 0x15F85; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_15F88: ; /* IDA 0x15F88; mov     cx, 3 */
    m->cx=(uint16_t)(3);
L_15F8B: ; /* IDA 0x15F8B; mov     di, 66C7h */
    m->di=(uint16_t)(26311);
L_15F8E: ; /* IDA 0x15F8E; mov     ax, [si] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_15F90: ; /* IDA 0x15F90; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_15F93: ; /* IDA 0x15F93; mul     dx */
    dl_mul(m,m->dx,16);
L_15F95: ; /* IDA 0x15F95; mov     bx, 0F06h */
    m->bx=(uint16_t)(3846);
L_15F98: ; /* IDA 0x15F98; call    sub_1062F */
    vd_call(m,1583,24475,h);
L_15F9B: ; /* IDA 0x15F9B; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_15F9C: ; /* IDA 0x15F9C; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_15F9D: ; /* IDA 0x15F9D; add     di, 500h */
    m->di=(uint16_t)(ec_math(m,m->di,1280,16,0,0));
L_15FA1: ; /* IDA 0x15FA1; loop    loc_15F8E */
    if(--m->cx) goto L_15F8E;
L_15FA3: ; /* IDA 0x15FA3; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_15FA4: ; /* IDA 0x15FA4; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_15FA5: ; /* IDA 0x15FA5; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15FA6: ; /* IDA 0x15FA6; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_15FA7: ; /* IDA 0x15FA7; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15FA8: ; /* IDA 0x15FA8; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15FA9: ; /* IDA 0x15FA9; retn */
    return;
}
/* Original sub_1699E; IDA 0x1699E. */
static void vd_body_1699E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1699E: ; /* IDA 0x1699E; push    ds */
    ec_push(m,m->ds);
L_1699F: ; /* IDA 0x1699F; call    sub_12078 */
    vd_call(m,8312,27042,h);
L_169A2: ; /* IDA 0x169A2; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_169A5: ; /* IDA 0x169A5; call    sub_20000 */
    mx_far(m,0x69a5,0x69aa,h);
L_169AA: ; /* IDA 0x169AA; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_169AF: ; /* IDA 0x169AF; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_169B4: ; /* IDA 0x169B4; xor     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_169B6: ; /* IDA 0x169B6; mov     al, 0Fh */
    vg_al(m,(uint8_t)(15));
L_169B8: ; /* IDA 0x169B8; sub     al, [si+28h] */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+40)),8,0,1)));
L_169BB: ; /* IDA 0x169BB; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_169BD: ; /* IDA 0x169BD; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_169BF: ; /* IDA 0x169BF; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_169C1: ; /* IDA 0x169C1; mov     bx, [si+21h] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+33)));
L_169C4: ; /* IDA 0x169C4; cmp     bh, 80h */
    ec_cmp(m,(uint8_t)(m->bx>>8),128,8);
L_169C7: ; /* IDA 0x169C7; jnb     short loc_169D2 */
#if KI_VERDICT_MUTATION == 2
    if(m->flags&1) goto L_169D2;
#else
    if(!(m->flags&1)) goto L_169D2;
#endif
L_169C9: ; /* IDA 0x169C9; and     bh, 7Fh */
    {uint16_t result=(((uint8_t)(m->bx>>8))&(127))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_169CC: ; /* IDA 0x169CC; cmp     bx, ax */
    ec_cmp(m,m->bx,m->ax,16);
L_169CE: ; /* IDA 0x169CE; jnb     short loc_169D2 */
    if(!(m->flags&1)) goto L_169D2;
L_169D0: ; /* IDA 0x169D0; mov     dl, 1 */
    vg_dl(m,(uint8_t)(1));
L_169D2: ; /* IDA 0x169D2; mov     ax, [si+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+4)));
L_169D5: ; /* IDA 0x169D5; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_169D7: ; /* IDA 0x169D7; add     ax, [si+6] */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->ds,(uint16_t)(m->si+6)),16,0,0));
L_169DA: ; /* IDA 0x169DA; adc     bx, 0 */
    m->bx=(uint16_t)(ec_math(m,m->bx,0,16,m->flags&1,0));
L_169DD: ; /* IDA 0x169DD; add     ax, [si+8] */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->ds,(uint16_t)(m->si+8)),16,0,0));
L_169E0: ; /* IDA 0x169E0; adc     bx, 0 */
    m->bx=(uint16_t)(ec_math(m,m->bx,0,16,m->flags&1,0));
L_169E3: ; /* IDA 0x169E3; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_169E5: ; /* IDA 0x169E5; jnz     short loc_169EE */
    if(!(m->flags&0x40)) goto L_169EE;
L_169E7: ; /* IDA 0x169E7; cmp     ax, 258h */
#if KI_VERDICT_MUTATION == 1
    ec_cmp(m,m->ax,601,16);
#else
    ec_cmp(m,m->ax,600,16);
#endif
L_169EA: ; /* IDA 0x169EA; jnb     short loc_169EE */
    if(!(m->flags&1)) goto L_169EE;
L_169EC: ; /* IDA 0x169EC; mov     dl, 1 */
    vg_dl(m,(uint8_t)(1));
L_169EE: ; /* IDA 0x169EE; mov     al, dl */
    vg_al(m,(uint8_t)((uint8_t)m->dx));
L_169F0: ; /* IDA 0x169F0; call    sub_187FF */
    vd_call(m,34815,27123,h);
L_169F3: ; /* IDA 0x169F3; add     bx, 4240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16960,16,0,0));
L_169F7: ; /* IDA 0x169F7; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_169F9: ; /* IDA 0x169F9; mov     cx, 18Ch */
    m->cx=(uint16_t)(396);
L_169FC: ; /* IDA 0x169FC; call    sub_13B08 */
    vd_call(m,15112,27135,h);
L_169FF: ; /* IDA 0x169FF; and     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)&((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_16A01: ; /* IDA 0x16A01; jnz     short loc_16A10 */
    if(!(m->flags&0x40)) goto L_16A10;
L_16A03: ; /* IDA 0x16A03; call    sub_16E8F */
    vd_call(m,28303,27142,h);
L_16A06: ; /* IDA 0x16A06; jb      short loc_16A10 */
    if(m->flags&1) goto L_16A10;
L_16A08: ; /* IDA 0x16A08; and     byte ptr [di], 0FBh */
#if KI_VERDICT_MUTATION == 9
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->di)))&(255))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->di),result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->di)))&(251))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->di),result);}
#endif
L_16A0B: ; /* IDA 0x16A0B; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_16A0D: ; /* IDA 0x16A0D; call    sub_15E80 */
    vd_call(m,24192,27152,h);
L_16A10: ; /* IDA 0x16A10; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16A13: ; /* IDA 0x16A13; call    sub_18853 */
    vd_call(m,34899,27158,h);
L_16A16: ; /* IDA 0x16A16; call    sub_11D46 */
    vd_call(m,7494,27161,h);
L_16A19: ; /* IDA 0x16A19; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_16A1C: ; /* IDA 0x16A1C; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_16A1E: ; /* IDA 0x16A1E; call    sub_20000 */
    mx_far(m,0x6a1e,0x6a23,h);
L_16A23: ; /* IDA 0x16A23; call    sub_120D6 */
    vd_call(m,8406,27174,h);
L_16A26: ; /* IDA 0x16A26; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16A27: ; /* IDA 0x16A27; retn */
    return;
}
/* Original sub_16E8F; IDA 0x16E8F. */
static void vd_body_16E8F(KiMachine16 *m,const KiEconomyHooks *h) {
L_16E8F: ; /* IDA 0x16E8F; push    ax */
    ec_push(m,m->ax);
L_16E90: ; /* IDA 0x16E90; push    bx */
    ec_push(m,m->bx);
L_16E91: ; /* IDA 0x16E91; push    si */
    ec_push(m,m->si);
L_16E92: ; /* IDA 0x16E92; mov     di, si */
    m->di=(uint16_t)(m->si);
L_16E94: ; /* IDA 0x16E94; sub     si, 4240h */
    m->si=(uint16_t)(ec_math(m,m->si,16960,16,0,1));
L_16E98: ; /* IDA 0x16E98; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_16E9A: ; /* IDA 0x16E9A; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_16E9C: ; /* IDA 0x16E9C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16E9E: ; /* IDA 0x16E9E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16EA0: ; /* IDA 0x16EA0; mov     [si+2242h], bh */
    dl_write8(m,m->ds,(uint16_t)(m->si+8770),(uint8_t)(m->bx>>8));
L_16EA4: ; /* IDA 0x16EA4; mov     bh, [di+1Ch] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+28))));
L_16EA7: ; /* IDA 0x16EA7; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_16EA9: ; /* IDA 0x16EA9; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16EAB: ; /* IDA 0x16EAB; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16EAD: ; /* IDA 0x16EAD; add     si, 2240h */
    m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,0));
L_16EB1: ; /* IDA 0x16EB1; call    sub_16EC9 */
    vd_call(m,28361,28340,h);
L_16EB4: ; /* IDA 0x16EB4; jb      short loc_16EC5 */
    if(m->flags&1) goto L_16EC5;
L_16EB6: ; /* IDA 0x16EB6; call    sub_16F26 */
    vd_call(m,28454,28345,h);
L_16EB9: ; /* IDA 0x16EB9; call    sub_1461D */
    vd_call(m,17949,28348,h);
L_16EBC: ; /* IDA 0x16EBC; call    sub_16FD2 */
    vd_call(m,28626,28351,h);
L_16EBF: ; /* IDA 0x16EBF; mov     di, si */
    m->di=(uint16_t)(m->si);
L_16EC1: ; /* IDA 0x16EC1; or      byte ptr [si], 4 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(4))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_16EC4: ; /* IDA 0x16EC4; clc */
    m->flags&=(uint16_t)~1u;
L_16EC5: ; /* IDA 0x16EC5; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16EC6: ; /* IDA 0x16EC6; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16EC7: ; /* IDA 0x16EC7; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16EC8: ; /* IDA 0x16EC8; retn */
    return;
}
/* Original sub_16EC9; IDA 0x16EC9. */
static void vd_body_16EC9(KiMachine16 *m,const KiEconomyHooks *h) {
L_16EC9: ; /* IDA 0x16EC9; push    ax */
    ec_push(m,m->ax);
L_16ECA: ; /* IDA 0x16ECA; push    cx */
    ec_push(m,m->cx);
L_16ECB: ; /* IDA 0x16ECB; push    dx */
    ec_push(m,m->dx);
L_16ECC: ; /* IDA 0x16ECC; push    si */
    ec_push(m,m->si);
L_16ECD: ; /* IDA 0x16ECD; push    di */
    ec_push(m,m->di);
L_16ECE: ; /* IDA 0x16ECE; push    bp */
    ec_push(m,m->bp);
L_16ECF: ; /* IDA 0x16ECF; sub     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,1));
L_16ED2: ; /* IDA 0x16ED2; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_16ED4: ; /* IDA 0x16ED4; mov     cx, 3 */
    m->cx=(uint16_t)(3);
L_16ED7: ; /* IDA 0x16ED7; mov     ax, [bx+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+4)));
L_16EDA: ; /* IDA 0x16EDA; mov     [bp+0], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->ax);
L_16EDD: ; /* IDA 0x16EDD; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_16EDE: ; /* IDA 0x16EDE; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_16EDF: ; /* IDA 0x16EDF; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_16EE0: ; /* IDA 0x16EE0; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_16EE1: ; /* IDA 0x16EE1; loop    loc_16ED7 */
    if(--m->cx) goto L_16ED7;
L_16EE3: ; /* IDA 0x16EE3; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_16EE5: ; /* IDA 0x16EE5; mov     dx, 6C4Ch */
    m->dx=(uint16_t)(27724);
L_16EE8: ; /* IDA 0x16EE8; add     si, 2Ah ; '*' */
    m->si=(uint16_t)(ec_math(m,m->si,42,16,0,0));
L_16EEB: ; /* IDA 0x16EEB; mov     cl, 6 */
    vg_cl(m,(uint8_t)(6));
L_16EED: ; /* IDA 0x16EED; mov     bx, dx */
    m->bx=(uint16_t)(m->dx);
L_16EEF: ; /* IDA 0x16EEF; mov     ch, 3 */
    vg_ch(m,(uint8_t)(3));
L_16EF1: ; /* IDA 0x16EF1; mov     al, cs:[bx] */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->bx))));
L_16EF4: ; /* IDA 0x16EF4; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_16EF6: ; /* IDA 0x16EF6; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_16EF8: ; /* IDA 0x16EF8; dec     di */
    m->di=(uint16_t)(st_dec(m,m->di,16));
L_16EF9: ; /* IDA 0x16EF9; shl     di, 1 */
    m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
L_16EFB: ; /* IDA 0x16EFB; cmp     word ptr [bp+di], 32h ; '2' */
#if KI_VERDICT_MUTATION == 3
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp+m->di)),51,16);
#else
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp+m->di)),50,16);
#endif
L_16EFE: ; /* IDA 0x16EFE; jnb     short loc_16F07 */
    if(!(m->flags&1)) goto L_16F07;
L_16F00: ; /* IDA 0x16F00; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_16F01: ; /* IDA 0x16F01; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_16F03: ; /* IDA 0x16F03; jnz     short loc_16EF1 */
    if(!(m->flags&0x40)) goto L_16EF1;
L_16F05: ; /* IDA 0x16F05; jmp     short loc_16F19 */
    goto L_16F19;
L_16F07: ; /* IDA 0x16F07; sub     word ptr [bp+di], 32h ; '2' */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),ec_math(m,dl_read16(m,m->ss,(uint16_t)(m->bp+m->di)),50,16,0,1));
L_16F0A: ; /* IDA 0x16F0A; mov     [si], al */
    dl_write8(m,m->ds,(uint16_t)(m->si),(uint8_t)m->ax);
L_16F0C: ; /* IDA 0x16F0C; add     si, 4 */
    m->si=(uint16_t)(ec_math(m,m->si,4,16,0,0));
L_16F0F: ; /* IDA 0x16F0F; add     dx, 3 */
    m->dx=(uint16_t)(ec_math(m,m->dx,3,16,0,0));
L_16F12: ; /* IDA 0x16F12; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_16F14: ; /* IDA 0x16F14; jnz     short loc_16EED */
    if(!(m->flags&0x40)) goto L_16EED;
L_16F16: ; /* IDA 0x16F16; clc */
    m->flags&=(uint16_t)~1u;
L_16F17: ; /* IDA 0x16F17; jmp     short loc_16F1A */
    goto L_16F1A;
L_16F19: ; /* IDA 0x16F19; stc */
    m->flags|=1;
L_16F1A: ; /* IDA 0x16F1A; lahf */
    vg_ah(m,(uint8_t)((m->flags&0xd5)|2));
L_16F1B: ; /* IDA 0x16F1B; add     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,0));
L_16F1E: ; /* IDA 0x16F1E; sahf */
#if KI_VERDICT_MUTATION == 4
    ; /* lost carry restoration */
#else
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
#endif
L_16F1F: ; /* IDA 0x16F1F; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_16F20: ; /* IDA 0x16F20; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_16F21: ; /* IDA 0x16F21; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16F22: ; /* IDA 0x16F22; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_16F23: ; /* IDA 0x16F23; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_16F24: ; /* IDA 0x16F24; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16F25: ; /* IDA 0x16F25; retn */
    return;
}
/* Original sub_16F26; IDA 0x16F26. */
static void vd_body_16F26(KiMachine16 *m,const KiEconomyHooks *h) {
L_16F26: ; /* IDA 0x16F26; push    ds */
    ec_push(m,m->ds);
L_16F27: ; /* IDA 0x16F27; push    es */
    ec_push(m,m->es);
L_16F28: ; /* IDA 0x16F28; push    si */
    ec_push(m,m->si);
L_16F29: ; /* IDA 0x16F29; push    di */
    ec_push(m,m->di);
L_16F2A: ; /* IDA 0x16F2A; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_16F2C: ; /* IDA 0x16F2C; sub     bx, 2240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,8768,16,0,1));
L_16F30: ; /* IDA 0x16F30; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16F32: ; /* IDA 0x16F32; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16F34: ; /* IDA 0x16F34; mov     [si+2], bh */
    dl_write8(m,m->ds,(uint16_t)(m->si+2),(uint8_t)(m->bx>>8));
L_16F37: ; /* IDA 0x16F37; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16F39: ; /* IDA 0x16F39; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16F3B: ; /* IDA 0x16F3B; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16F3D: ; /* IDA 0x16F3D; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_16F3F: ; /* IDA 0x16F3F; mov     byte ptr [di+4257h], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->di+16983),1);
L_16F44: ; /* IDA 0x16F44; mov     bh, [di+425Ch] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+16988))));
L_16F48: ; /* IDA 0x16F48; mov     [si+1], bh */
    dl_write8(m,m->ds,(uint16_t)(m->si+1),(uint8_t)(m->bx>>8));
L_16F4B: ; /* IDA 0x16F4B; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_16F4D: ; /* IDA 0x16F4D; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16F4F: ; /* IDA 0x16F4F; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16F51: ; /* IDA 0x16F51; mov     ah, [bx+3] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_16F54: ; /* IDA 0x16F54; mov     [si+20h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)(m->ax>>8));
L_16F57: ; /* IDA 0x16F57; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_16F5A: ; /* IDA 0x16F5A; jnb     short loc_16F5F */
    if(!(m->flags&1)) goto L_16F5F;
L_16F5C: ; /* IDA 0x16F5C; inc     byte ptr [bx+14h] */
    dl_write8(m,m->ds,(uint16_t)(m->bx+20),ec_inc(m,dl_read8(m,m->ds,(uint16_t)(m->bx+20)),8));
L_16F5F: ; /* IDA 0x16F5F; mov     byte ptr [si], 0C0h */
    dl_write8(m,m->ds,(uint16_t)(m->si),192);
L_16F62: ; /* IDA 0x16F62; mov     byte ptr [si+8], 4 */
    dl_write8(m,m->ds,(uint16_t)(m->si+8),4);
L_16F66: ; /* IDA 0x16F66; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_16F68: ; /* IDA 0x16F68; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16F6A: ; /* IDA 0x16F6A; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16F6C: ; /* IDA 0x16F6C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_16F6E: ; /* IDA 0x16F6E; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_16F70: ; /* IDA 0x16F70; add     di, 840h */
    m->di=(uint16_t)(ec_math(m,m->di,2112,16,0,0));
L_16F74: ; /* IDA 0x16F74; call    sub_16F86 */
    vd_call(m,28550,28535,h);
L_16F77: ; /* IDA 0x16F77; mov     al, [bx+1Dh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+29))));
L_16F7A: ; /* IDA 0x16F7A; mov     [si+6], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+6),(uint8_t)m->ax);
L_16F7D: ; /* IDA 0x16F7D; mov     byte ptr [si+23h], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),1);
L_16F81: ; /* IDA 0x16F81; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_16F82: ; /* IDA 0x16F82; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16F83: ; /* IDA 0x16F83; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_16F84: ; /* IDA 0x16F84; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16F85: ; /* IDA 0x16F85; retn */
    return;
}
/* Original sub_16F86; IDA 0x16F86. */
static void vd_body_16F86(KiMachine16 *m,const KiEconomyHooks *h) {
L_16F86: ; /* IDA 0x16F86; push    ax */
    ec_push(m,m->ax);
L_16F87: ; /* IDA 0x16F87; push    bx */
    ec_push(m,m->bx);
L_16F88: ; /* IDA 0x16F88; push    cx */
    ec_push(m,m->cx);
L_16F89: ; /* IDA 0x16F89; mov     ax, [di+8] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+8)));
L_16F8C: ; /* IDA 0x16F8C; mov     bx, [di+0Ah] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+10)));
L_16F8F: ; /* IDA 0x16F8F; mov     cx, di */
    m->cx=(uint16_t)(m->di);
L_16F91: ; /* IDA 0x16F91; sub     cx, 840h */
    m->cx=(uint16_t)(ec_math(m,m->cx,2112,16,0,1));
L_16F95: ; /* IDA 0x16F95; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_16F97: ; /* IDA 0x16F97; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_16F99: ; /* IDA 0x16F99; mov     [si+0Eh], cx */
    dl_write16(m,m->ds,(uint16_t)(m->si+14),m->cx);
L_16F9C: ; /* IDA 0x16F9C; mov     [si+14h], cx */
    dl_write16(m,m->ds,(uint16_t)(m->si+20),m->cx);
L_16F9F: ; /* IDA 0x16F9F; mov     [si+10h], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+16),m->ax);
L_16FA2: ; /* IDA 0x16FA2; mov     [si+16h], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+22),m->ax);
L_16FA5: ; /* IDA 0x16FA5; mov     [si+12h], bx */
    dl_write16(m,m->ds,(uint16_t)(m->si+18),m->bx);
L_16FA8: ; /* IDA 0x16FA8; mov     [si+18h], bx */
    dl_write16(m,m->ds,(uint16_t)(m->si+24),m->bx);
L_16FAB: ; /* IDA 0x16FAB; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_16FAD: ; /* IDA 0x16FAD; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16FAF: ; /* IDA 0x16FAF; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16FB1: ; /* IDA 0x16FB1; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16FB3: ; /* IDA 0x16FB3; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_16FB5: ; /* IDA 0x16FB5; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16FB7: ; /* IDA 0x16FB7; add     bx, cx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->cx,16,0,0));
L_16FB9: ; /* IDA 0x16FB9; add     bx, cs:word_19872 */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->cs,39026),16,0,0));
L_16FBE: ; /* IDA 0x16FBE; mov     [si+1Ah], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+26),m->ax);
L_16FC1: ; /* IDA 0x16FC1; mov     [si+1Ch], bx */
    dl_write16(m,m->ds,(uint16_t)(m->si+28),m->bx);
L_16FC4: ; /* IDA 0x16FC4; mov     cx, ds */
    m->cx=(uint16_t)(m->ds);
L_16FC6: ; /* IDA 0x16FC6; mov     ds, bx */
    m->ds=(uint16_t)(m->bx);
L_16FC8: ; /* IDA 0x16FC8; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_16FCA: ; /* IDA 0x16FCA; inc     byte ptr [bx] */
#if KI_VERDICT_MUTATION == 7
    ; /* lost occupancy increment */
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx),ec_inc(m,dl_read8(m,m->ds,(uint16_t)(m->bx)),8));
#endif
L_16FCC: ; /* IDA 0x16FCC; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_16FCE: ; /* IDA 0x16FCE; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_16FCF: ; /* IDA 0x16FCF; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16FD0: ; /* IDA 0x16FD0; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16FD1: ; /* IDA 0x16FD1; retn */
    return;
}
/* Original sub_16FD2; IDA 0x16FD2. */
static void vd_body_16FD2(KiMachine16 *m,const KiEconomyHooks *h) {
L_16FD2: ; /* IDA 0x16FD2; push    ax */
    ec_push(m,m->ax);
L_16FD3: ; /* IDA 0x16FD3; push    bx */
    ec_push(m,m->bx);
L_16FD4: ; /* IDA 0x16FD4; push    cx */
    ec_push(m,m->cx);
L_16FD5: ; /* IDA 0x16FD5; push    dx */
    ec_push(m,m->dx);
L_16FD6: ; /* IDA 0x16FD6; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_16FD8: ; /* IDA 0x16FD8; add     bx, 28h ; '(' */
    m->bx=(uint16_t)(ec_math(m,m->bx,40,16,0,0));
L_16FDB: ; /* IDA 0x16FDB; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_16FDD: ; /* IDA 0x16FDD; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_16FE0: ; /* IDA 0x16FE0; cmp     byte ptr [bx+2], 1 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+2)),1,8);
L_16FE4: ; /* IDA 0x16FE4; jz      short loc_16FE9 */
    if(m->flags&0x40) goto L_16FE9;
L_16FE6: ; /* IDA 0x16FE6; or      ch, 1 */
    {uint16_t result=(((uint8_t)(m->cx>>8))|(1))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_16FE9: ; /* IDA 0x16FE9; add     al, [bx+1] */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+1)),8,0,0)));
L_16FEC: ; /* IDA 0x16FEC; adc     ah, 0 */
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),0,8,m->flags&1,0)));
L_16FEF: ; /* IDA 0x16FEF; add     bx, 4 */
    m->bx=(uint16_t)(ec_math(m,m->bx,4,16,0,0));
L_16FF2: ; /* IDA 0x16FF2; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_16FF4: ; /* IDA 0x16FF4; jnz     short loc_16FE0 */
    if(!(m->flags&0x40)) goto L_16FE0;
L_16FF6: ; /* IDA 0x16FF6; mov     [si+4], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+4),m->ax);
L_16FF9: ; /* IDA 0x16FF9; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_16FFB: ; /* IDA 0x16FFB; cmp     bx, 12Ch */
#if KI_VERDICT_MUTATION == 8
    ec_cmp(m,m->ax,300,16);
#else
    ec_cmp(m,m->bx,300,16);
#endif
L_16FFF: ; /* IDA 0x16FFF; jbe     short loc_17003 */
    if(m->flags&(1|0x40)) goto L_17003;
L_17001: ; /* IDA 0x17001; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_17003: ; /* IDA 0x17003; mov     [si+1Eh], ch */
    dl_write8(m,m->ds,(uint16_t)(m->si+30),(uint8_t)(m->cx>>8));
L_17006: ; /* IDA 0x17006; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_17009: ; /* IDA 0x17009; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_1700B: ; /* IDA 0x1700B; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1700D: ; /* IDA 0x1700D; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1700F: ; /* IDA 0x1700F; mov     al, [bx+3Eh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+62))));
L_17012: ; /* IDA 0x17012; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_17014: ; /* IDA 0x17014; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_17016: ; /* IDA 0x17016; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_17018: ; /* IDA 0x17018; add     al, ah */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8,0,0)));
L_1701A: ; /* IDA 0x1701A; mov     [si+9], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+9),(uint8_t)m->ax);
L_1701D: ; /* IDA 0x1701D; mov     byte ptr [si+0Bh], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+11),1);
L_17021: ; /* IDA 0x17021; mov     ax, dx */
    m->ax=(uint16_t)(m->dx);
L_17023: ; /* IDA 0x17023; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17024: ; /* IDA 0x17024; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17025: ; /* IDA 0x17025; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17026: ; /* IDA 0x17026; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17027: ; /* IDA 0x17027; retn */
    return;
}
int ki_verdict_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x461d:vd_body_1461D(m,h);break;
case 0x4698:vd_body_14698(m,h);break;
case 0x4717:vd_body_14717(m,h);break;
case 0x5e80:vd_body_15E80(m,h);break;
case 0x5eb7:vd_body_15EB7(m,h);break;
case 0x5f27:vd_body_15F27(m,h);break;
case 0x5f5d:vd_body_15F5D(m,h);break;
case 0x5f7f:vd_body_15F7F(m,h);break;
case 0x699e:vd_body_1699E(m,h);break;
case 0x6e8f:vd_body_16E8F(m,h);break;
case 0x6ec9:vd_body_16EC9(m,h);break;
case 0x6f26:vd_body_16F26(m,h);break;
case 0x6f86:vd_body_16F86(m,h);break;
case 0x6fd2:vd_body_16FD2(m,h);break;
default:return 0;}return 1;}
