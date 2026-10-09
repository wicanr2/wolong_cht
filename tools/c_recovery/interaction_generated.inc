/* Original world interaction; re/127, spec/247. */
static void ix_body_109D0(KiMachine16 *,const KiEconomyHooks *);
static void ix_body_11E46(KiMachine16 *,const KiEconomyHooks *);
static void ix_body_11F0E(KiMachine16 *,const KiEconomyHooks *);
static void ix_body_159B7(KiMachine16 *,const KiEconomyHooks *);
static void ix_body_159D0(KiMachine16 *,const KiEconomyHooks *);
static void ix_body_15AA2(KiMachine16 *,const KiEconomyHooks *);
static void ix_body_15E4C(KiMachine16 *,const KiEconomyHooks *);
static void ix_body_161B6(KiMachine16 *,const KiEconomyHooks *);
/* Original sub_109D0; IDA 0x109D0. */
static void ix_body_109D0(KiMachine16 *m,const KiEconomyHooks *h) {
L_109D0: ; /* IDA 0x109D0; push    ax */
    ec_push(m,m->ax);
L_109D1: ; /* IDA 0x109D1; push    bx */
    ec_push(m,m->bx);
L_109D2: ; /* IDA 0x109D2; push    cx */
    ec_push(m,m->cx);
L_109D3: ; /* IDA 0x109D3; push    dx */
    ec_push(m,m->dx);
L_109D4: ; /* IDA 0x109D4; push    si */
    ec_push(m,m->si);
L_109D5: ; /* IDA 0x109D5; mov     cl, 0 */
    vg_cl(m,(uint8_t)(0));
L_109D7: ; /* IDA 0x109D7; mov     ch, 0 */
    vg_ch(m,(uint8_t)(0));
L_109D9: ; /* IDA 0x109D9; mov     si, 1964h */
    m->si=(uint16_t)(6500);
L_109DC: ; /* IDA 0x109DC; mov     dx, 3DAh */
    m->dx=(uint16_t)(986);
L_109DF: ; /* IDA 0x109DF; in      al, dx; Video status bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
L_109E0: ; /* IDA 0x109E0; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_109E2: ; /* IDA 0x109E2; jnz     short loc_109DF */
    if(!(m->flags&0x40)) goto L_109DF;
L_109E4: ; /* IDA 0x109E4; in      al, dx; Video status bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
L_109E5: ; /* IDA 0x109E5; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_109E7: ; /* IDA 0x109E7; jz      short loc_109E4 */
    if(m->flags&0x40) goto L_109E4;
L_109E9: ; /* IDA 0x109E9; in      al, dx; Video status bits: */
#if KI_INTERACTION_MUTATION == 2
    goto L_109F3; /* Omit the second retrace wait. */
#else
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
#endif
L_109EA: ; /* IDA 0x109EA; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_109EC: ; /* IDA 0x109EC; jnz     short loc_109E9 */
    if(!(m->flags&0x40)) goto L_109E9;
L_109EE: ; /* IDA 0x109EE; in      al, dx; Video status bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
L_109EF: ; /* IDA 0x109EF; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_109F1: ; /* IDA 0x109F1; jz      short loc_109EE */
    if(m->flags&0x40) goto L_109EE;
L_109F3: ; /* IDA 0x109F3; mov     bh, cs:[si] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->si))));
L_109F6: ; /* IDA 0x109F6; mov     dl, cs:[si+1] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->si+1))));
L_109FA: ; /* IDA 0x109FA; mov     dh, cs:[si+2] */
    vg_dh(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->si+2))));
L_109FE: ; /* IDA 0x109FE; mov     al, ch */
    vg_al(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_10A00: ; /* IDA 0x10A00; mov     bl, cl */
    gw_bl(m,(uint8_t)((uint8_t)m->cx));
L_10A02: ; /* IDA 0x10A02; call    sub_1EBDC */
    ix_call(m,60380,2565,h);
L_10A05: ; /* IDA 0x10A05; add     si, 3 */
    m->si=(uint16_t)(ec_math(m,m->si,3,16,0,0));
L_10A08: ; /* IDA 0x10A08; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_10A0A: ; /* IDA 0x10A0A; cmp     ch, 10h */
    ec_cmp(m,(uint8_t)(m->cx>>8),16,8);
L_10A0D: ; /* IDA 0x10A0D; jb      short loc_109F3 */
    if(m->flags&1) goto L_109F3;
L_10A0F: ; /* IDA 0x10A0F; inc     cl */
    vg_cl(m,(uint8_t)(ec_inc(m,(uint8_t)m->cx,8)));
L_10A11: ; /* IDA 0x10A11; cmp     cl, 10h */
#if KI_INTERACTION_MUTATION == 1
    ec_cmp(m,(uint8_t)m->cx,15,8);
#else
    ec_cmp(m,(uint8_t)m->cx,16,8);
#endif
L_10A14: ; /* IDA 0x10A14; jbe     short loc_109D7 */
    if(m->flags&(1|0x40)) goto L_109D7;
L_10A16: ; /* IDA 0x10A16; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_10A17: ; /* IDA 0x10A17; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10A18: ; /* IDA 0x10A18; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10A19: ; /* IDA 0x10A19; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10A1A: ; /* IDA 0x10A1A; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10A1B: ; /* IDA 0x10A1B; retn */
    return;
}
/* Original sub_11E46; IDA 0x11E46. */
static void ix_body_11E46(KiMachine16 *m,const KiEconomyHooks *h) {
L_11E46: ; /* IDA 0x11E46; push    ds */
    ec_push(m,m->ds);
L_11E47: ; /* IDA 0x11E47; push    es */
    ec_push(m,m->es);
L_11E48: ; /* IDA 0x11E48; push    bp */
    ec_push(m,m->bp);
L_11E49: ; /* IDA 0x11E49; call    sub_12078 */
    ix_call(m,8312,7756,h);
L_11E4C: ; /* IDA 0x11E4C; sub     sp, 8 */
    m->sp=(uint16_t)(ec_math(m,m->sp,8,16,0,1));
L_11E4F: ; /* IDA 0x11E4F; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_11E51: ; /* IDA 0x11E51; mov     ax, cs:word_1989C */
    m->ax=(uint16_t)(dl_read16(m,m->cs,39068));
L_11E55: ; /* IDA 0x11E55; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_11E57: ; /* IDA 0x11E57; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_11E59: ; /* IDA 0x11E59; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_11E5B: ; /* IDA 0x11E5B; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_11E5D: ; /* IDA 0x11E5D; mov     dx, cx */
    m->dx=(uint16_t)(m->cx);
L_11E5F: ; /* IDA 0x11E5F; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_11E61: ; /* IDA 0x11E61; add     cx, dx */
    m->cx=(uint16_t)(ec_math(m,m->cx,m->dx,16,0,0));
L_11E63: ; /* IDA 0x11E63; mov     bx, cs:word_1989A */
    m->bx=(uint16_t)(dl_read16(m,m->cs,39066));
L_11E68: ; /* IDA 0x11E68; mov     dx, cx */
    m->dx=(uint16_t)(m->cx);
L_11E6A: ; /* IDA 0x11E6A; mov     di, bp */
    m->di=(uint16_t)(m->bp);
L_11E6C: ; /* IDA 0x11E6C; mov     word ptr [bp+0], 0FFFFh */
    dl_write16(m,m->ss,(uint16_t)(m->bp),65535);
L_11E71: ; /* IDA 0x11E71; mov     word ptr [bp+2], 0 */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),0);
L_11E76: ; /* IDA 0x11E76; add     cx, cs:word_10D44 */
    m->cx=(uint16_t)(ec_math(m,m->cx,dl_read16(m,m->cs,3396),16,0,0));
L_11E7B: ; /* IDA 0x11E7B; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_11E7D: ; /* IDA 0x11E7D; mov     cl, [bx] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx))));
L_11E7F: ; /* IDA 0x11E7F; cmp     cl, 0CBh */
#if KI_INTERACTION_MUTATION == 5
    ec_cmp(m,(uint8_t)m->cx,0xcc,8);
#else
    ec_cmp(m,(uint8_t)m->cx,203,8);
#endif
L_11E82: ; /* IDA 0x11E82; jb      short loc_11EA4 */
    if(m->flags&1) goto L_11EA4;
L_11E84: ; /* IDA 0x11E84; cmp     cl, 0D4h */
#if KI_INTERACTION_MUTATION == 6
    ec_cmp(m,(uint8_t)m->cx,0xd3,8);
#else
    ec_cmp(m,(uint8_t)m->cx,212,8);
#endif
L_11E87: ; /* IDA 0x11E87; jnb     short loc_11EA4 */
    if(!(m->flags&1)) goto L_11EA4;
L_11E89: ; /* IDA 0x11E89; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_11E8E: ; /* IDA 0x11E8E; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_11E90: ; /* IDA 0x11E90; cmp     [si+848h], bx */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+2120)),m->bx,16);
L_11E94: ; /* IDA 0x11E94; jnz     short loc_11E9C */
    if(!(m->flags&0x40)) goto L_11E9C;
L_11E96: ; /* IDA 0x11E96; cmp     [si+84Ah], ax */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+2122)),m->ax,16);
L_11E9A: ; /* IDA 0x11E9A; jz      short loc_11EA1 */
    if(m->flags&0x40) goto L_11EA1;
L_11E9C: ; /* IDA 0x11E9C; add     si, 20h ; ' ' */
    m->si=(uint16_t)(ec_math(m,m->si,32,16,0,0));
L_11E9F: ; /* IDA 0x11E9F; jmp     short loc_11E90 */
    goto L_11E90;
L_11EA1: ; /* IDA 0x11EA1; mov     [bp+0], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->si);
L_11EA4: ; /* IDA 0x11EA4; add     dx, cs:word_19872 */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->cs,39026),16,0,0));
L_11EA9: ; /* IDA 0x11EA9; mov     ds, dx */
    m->ds=(uint16_t)(m->dx);
L_11EAB: ; /* IDA 0x11EAB; mov     bl, [bx] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx))));
L_11EAD: ; /* IDA 0x11EAD; mov     [bp+2], bl */
    dl_write8(m,m->ss,(uint16_t)(m->bp+2),(uint8_t)m->bx);
L_11EB0: ; /* IDA 0x11EB0; cmp     word ptr [bp+0], 0FFFFh */
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp)),65535,16);
L_11EB5: ; /* IDA 0x11EB5; jz      short loc_11ED3 */
    if(m->flags&0x40) goto L_11ED3;
L_11EB7: ; /* IDA 0x11EB7; cmp     bl, 0 */
    ec_cmp(m,(uint8_t)m->bx,0,8);
L_11EBA: ; /* IDA 0x11EBA; jz      short loc_11EC5 */
    if(m->flags&0x40) goto L_11EC5;
L_11EBC: ; /* IDA 0x11EBC; call    sub_11F0E */
    ix_call(m,7950,7871,h);
L_11EBF: ; /* IDA 0x11EBF; jb      short loc_11F04 */
#if KI_INTERACTION_MUTATION == 7
    ; /* Ignore the city/corps popup cancellation carry. */
#else
    if(m->flags&1) goto L_11F04;
#endif
L_11EC1: ; /* IDA 0x11EC1; cmp     al, 1 */
    ec_cmp(m,(uint8_t)m->ax,1,8);
L_11EC3: ; /* IDA 0x11EC3; jz      short loc_11ED3 */
    if(m->flags&0x40) goto L_11ED3;
L_11EC5: ; /* IDA 0x11EC5; mov     si, [bp+0] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_11EC8: ; /* IDA 0x11EC8; call    sub_101B4 */
    ix_call(m,436,7883,h);
L_11ECB: ; /* IDA 0x11ECB; call    sub_17E1F */
    ix_call(m,32287,7886,h);
L_11ECE: ; /* IDA 0x11ECE; call    sub_101DB */
    ix_call(m,475,7889,h);
L_11ED1: ; /* IDA 0x11ED1; jmp     short loc_11F04 */
    goto L_11F04;
L_11ED3: ; /* IDA 0x11ED3; cmp     word ptr [bp+2], 1 */
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp+2)),1,16);
L_11ED7: ; /* IDA 0x11ED7; jb      short loc_11F04 */
    if(m->flags&1) goto L_11F04;
L_11ED9: ; /* IDA 0x11ED9; mov     ax, cs:word_1989A */
    m->ax=(uint16_t)(dl_read16(m,m->cs,39066));
L_11EDD: ; /* IDA 0x11EDD; mov     bx, cs:word_1989C */
    m->bx=(uint16_t)(dl_read16(m,m->cs,39068));
L_11EE2: ; /* IDA 0x11EE2; mov     [bp+4], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp+4),m->ax);
L_11EE5: ; /* IDA 0x11EE5; mov     [bp+6], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+6),m->bx);
L_11EE8: ; /* IDA 0x11EE8; call    sub_171D3 */
    ix_call(m,29139,7915,h);
L_11EEB: ; /* IDA 0x11EEB; jb      short loc_11F04 */
    if(m->flags&1) goto L_11F04;
L_11EED: ; /* IDA 0x11EED; sub     bx, 2240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,8768,16,0,1));
L_11EF1: ; /* IDA 0x11EF1; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_11EF3: ; /* IDA 0x11EF3; call    sub_120D6 */
    ix_call(m,8406,7926,h);
L_11EF6: ; /* IDA 0x11EF6; call    sub_17F90 */
    ix_call(m,32656,7929,h);
L_11EF9: ; /* IDA 0x11EF9; call    sub_12078 */
    ix_call(m,8312,7932,h);
L_11EFC: ; /* IDA 0x11EFC; mov     ax, [bp+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+4)));
L_11EFF: ; /* IDA 0x11EFF; mov     bx, [bp+6] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+6)));
L_11F02: ; /* IDA 0x11F02; jmp     short loc_11EE8 */
    goto L_11EE8;
L_11F04: ; /* IDA 0x11F04; add     sp, 8 */
    m->sp=(uint16_t)(ec_math(m,m->sp,8,16,0,0));
L_11F07: ; /* IDA 0x11F07; call    sub_120D6 */
    ix_call(m,8406,7946,h);
L_11F0A: ; /* IDA 0x11F0A; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_11F0B: ; /* IDA 0x11F0B; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_11F0C: ; /* IDA 0x11F0C; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_11F0D: ; /* IDA 0x11F0D; retn */
    return;
}
/* Original sub_11F0E; IDA 0x11F0E. */
static void ix_body_11F0E(KiMachine16 *m,const KiEconomyHooks *h) {
L_11F0E: ; /* IDA 0x11F0E; mov     ax, cs:word_19898 */
    m->ax=(uint16_t)(dl_read16(m,m->cs,39064));
L_11F12: ; /* IDA 0x11F12; cmp     al, 16h */
#if KI_INTERACTION_MUTATION == 3
    ec_cmp(m,(uint8_t)m->ax,21,8);
#else
    ec_cmp(m,(uint8_t)m->ax,22,8);
#endif
L_11F14: ; /* IDA 0x11F14; jbe     short loc_11F18 */
    if(m->flags&(1|0x40)) goto L_11F18;
L_11F16: ; /* IDA 0x11F16; mov     al, 16h */
#if KI_INTERACTION_MUTATION == 3
    vg_al(m,21);
#else
    vg_al(m,(uint8_t)(22));
#endif
L_11F18: ; /* IDA 0x11F18; mov     dx, cs:word_19896 */
    m->dx=(uint16_t)(dl_read16(m,m->cs,39062));
L_11F1D: ; /* IDA 0x11F1D; cmp     dl, 23h ; '#' */
#if KI_INTERACTION_MUTATION == 4
    ec_cmp(m,(uint8_t)m->dx,34,8);
#else
    ec_cmp(m,(uint8_t)m->dx,35,8);
#endif
L_11F20: ; /* IDA 0x11F20; jbe     short loc_11F24 */
    if(m->flags&(1|0x40)) goto L_11F24;
L_11F22: ; /* IDA 0x11F22; mov     dl, 23h ; '#' */
#if KI_INTERACTION_MUTATION == 4
    vg_dl(m,34);
#else
    vg_dl(m,(uint8_t)(35));
#endif
L_11F24: ; /* IDA 0x11F24; mov     dh, al */
    vg_dh(m,(uint8_t)((uint8_t)m->ax));
L_11F26: ; /* IDA 0x11F26; mov     ax, 102h */
    m->ax=(uint16_t)(258);
L_11F29: ; /* IDA 0x11F29; mov     cx, 50h ; 'P' */
    m->cx=(uint16_t)(80);
L_11F2C: ; /* IDA 0x11F2C; call    sub_193E9 */
    ix_call(m,37865,7983,h);
L_11F2F: ; /* IDA 0x11F2F; retn */
    return;
}
/* Original sub_159B7; IDA 0x159B7. */
static void ix_body_159B7(KiMachine16 *m,const KiEconomyHooks *h) {
L_159B7: ; /* IDA 0x159B7; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_159B9: ; /* IDA 0x159B9; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_159BB: ; /* IDA 0x159BB; shl     bx, 1 */
#if KI_INTERACTION_MUTATION == 8
    m->bx=0; /* Wrong but valid live table[0] index. */
#else
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
#endif
L_159BD: ; /* IDA 0x159BD; call    sub_101B4 */
    ix_call(m,436,22976,h);
L_159C0: ; /* IDA 0x159C0; call    funcs_159C0[bx] */
#if KI_INTERACTION_MUTATION == 9
    ix_call(m,0x59d0,0x59c4,h); /* Hard-coded nullsub_1. */
#else
    ix_call(m,dl_read16(m,m->ds,(uint16_t)(m->bx+23046)),22980,h);
#endif
L_159C4: ; /* IDA 0x159C4; call    sub_11C8D */
    ix_call(m,7309,22983,h);
L_159C7: ; /* IDA 0x159C7; call    sub_101DB */
    ix_call(m,475,22986,h);
L_159CA: ; /* IDA 0x159CA; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_159CC: ; /* IDA 0x159CC; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_159CE: ; /* IDA 0x159CE; mov     es, ax */
    m->es=(uint16_t)(m->ax);
m->ip=0x59d0; /* Original fall-through to the separate nullsub_1 RET. */
if(h&&h->enter)h->enter(m,0x59d0,h->user);
ix_body_159D0(m,h); /* Body tail only; the caller consumes one guest return address. */
}
/* Original nullsub_1; IDA 0x159D0. */
static void ix_body_159D0(KiMachine16 *m,const KiEconomyHooks *h) {
L_159D0: ; /* IDA 0x159D0; retn */
    return;
}
/* Original sub_15AA2; IDA 0x15AA2. */
static void ix_body_15AA2(KiMachine16 *m,const KiEconomyHooks *h) {
L_15AA2: ; /* IDA 0x15AA2; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_15AA4: ; /* IDA 0x15AA4; mov     dx, 1Bh */
    m->dx=(uint16_t)(27);
L_15AA7: ; /* IDA 0x15AA7; mov     bx, 0 */
    m->bx=(uint16_t)(0);
L_15AAA: ; /* IDA 0x15AAA; mov     cx, 0A0Dh */
    m->cx=(uint16_t)(2573);
L_15AAD: ; /* IDA 0x15AAD; call    sub_1895D */
    ix_call(m,35165,23216,h);
L_15AB0: ; /* IDA 0x15AB0; and     byte_198A6, 0FBh */
#if KI_INTERACTION_MUTATION == 10
    {uint16_t result=dl_read8(m,m->ds,0x98a6)&0xfd;ec_logic(m,result,8);dl_write8(m,m->ds,0x98a6,result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,39078))&(251))&255;ec_logic(m,result,8);dl_write8(m,m->ds,39078,result);}
#endif
L_15AB5: ; /* IDA 0x15AB5; retn */
    return;
}
/* Original sub_15E4C; IDA 0x15E4C. */
static void ix_body_15E4C(KiMachine16 *m,const KiEconomyHooks *h) {
L_15E4C: ; /* IDA 0x15E4C; and     byte ptr ds:98A6h, 0FDh */
#if KI_INTERACTION_MUTATION == 11
    {uint16_t result=dl_read8(m,m->ds,0x98a6)&0xfe;ec_logic(m,result,8);dl_write8(m,m->ds,0x98a6,result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,39078))&(253))&255;ec_logic(m,result,8);dl_write8(m,m->ds,39078,result);}
#endif
L_15E51: ; /* IDA 0x15E51; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_15E53: ; /* IDA 0x15E53; mov     dx, 1Bh */
    m->dx=(uint16_t)(27);
L_15E56: ; /* IDA 0x15E56; mov     bx, 0Ah */
    m->bx=(uint16_t)(10);
L_15E59: ; /* IDA 0x15E59; mov     cx, 0D0Dh */
    m->cx=(uint16_t)(3341);
L_15E5C: ; /* IDA 0x15E5C; call    sub_1895D */
    ix_call(m,35165,24159,h);
L_15E5F: ; /* IDA 0x15E5F; retn */
    return;
}
/* Original sub_161B6; IDA 0x161B6. */
static void ix_body_161B6(KiMachine16 *m,const KiEconomyHooks *h) {
L_161B6: ; /* IDA 0x161B6; and     byte ptr ds:98A6h, 0FEh */
#if KI_INTERACTION_MUTATION == 12
    {uint16_t result=dl_read8(m,m->ds,0x98a6)&0xfb;ec_logic(m,result,8);dl_write8(m,m->ds,0x98a6,result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,39078))&(254))&255;ec_logic(m,result,8);dl_write8(m,m->ds,39078,result);}
#endif
L_161BB: ; /* IDA 0x161BB; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_161BD: ; /* IDA 0x161BD; mov     dx, 0 */
    m->dx=(uint16_t)(0);
L_161C0: ; /* IDA 0x161C0; mov     bx, 0 */
    m->bx=(uint16_t)(0);
L_161C3: ; /* IDA 0x161C3; mov     cx, 21Bh */
    m->cx=(uint16_t)(539);
L_161C6: ; /* IDA 0x161C6; call    sub_1895D */
    ix_call(m,35165,25033,h);
L_161C9: ; /* IDA 0x161C9; retn */
    return;
}
int ki_interaction_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x9d0:ix_body_109D0(m,h);break;
case 0x1e46:ix_body_11E46(m,h);break;
case 0x1f0e:ix_body_11F0E(m,h);break;
case 0x59b7:ix_body_159B7(m,h);break;
case 0x59d0:ix_body_159D0(m,h);break;
case 0x5aa2:ix_body_15AA2(m,h);break;
case 0x5e4c:ix_body_15E4C(m,h);break;
case 0x61b6:ix_body_161B6(m,h);break;
default:return 0;}return 1;}
static int ix_known(uint16_t t) {switch(t) {case 0x9d0: case 0x1e46: case 0x1f0e: case 0x59b7: case 0x59d0: case 0x5aa2: case 0x5e4c: case 0x61b6: return 1;default:return 0;}}
void sub_109D0(KiMachine16 *m,const KiEconomyHooks *h) {ki_interaction_body(m,0x9d0,h);m->ip=ec_pop(m);}
void sub_11E46(KiMachine16 *m,const KiEconomyHooks *h) {ki_interaction_body(m,0x1e46,h);m->ip=ec_pop(m);}
void sub_11F0E(KiMachine16 *m,const KiEconomyHooks *h) {ki_interaction_body(m,0x1f0e,h);m->ip=ec_pop(m);}
void sub_159B7(KiMachine16 *m,const KiEconomyHooks *h) {ki_interaction_body(m,0x59b7,h);m->ip=ec_pop(m);}
void nullsub_1(KiMachine16 *m,const KiEconomyHooks *h) {ki_interaction_body(m,0x59d0,h);m->ip=ec_pop(m);}
void sub_15AA2(KiMachine16 *m,const KiEconomyHooks *h) {ki_interaction_body(m,0x5aa2,h);m->ip=ec_pop(m);}
void sub_15E4C(KiMachine16 *m,const KiEconomyHooks *h) {ki_interaction_body(m,0x5e4c,h);m->ip=ec_pop(m);}
void sub_161B6(KiMachine16 *m,const KiEconomyHooks *h) {ki_interaction_body(m,0x61b6,h);m->ip=ec_pop(m);}
