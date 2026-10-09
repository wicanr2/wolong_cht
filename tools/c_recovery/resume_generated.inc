/* Original native C scene/camera/minimap; re/116, spec/236. */
/* Original sub_11F30; IDA 0x11F30. */
static void sr_body_11F30(KiMachine16 *m,const KiEconomyHooks *h) {
L_11F30: ; /* IDA 0x11F30; mov     dx, ds:9896h */
    m->dx=(uint16_t)(dl_read16(m,m->ds,39062));
L_11F34: ; /* IDA 0x11F34; mov     bx, ds:9898h */
    m->bx=(uint16_t)(dl_read16(m,m->ds,39064));
L_11F38: ; /* IDA 0x11F38; sub     bx, 2 */
#if KI_RESUME_MUTATION == 1
    m->bx=(uint16_t)(ec_math(m,m->bx,1,16,0,1));
#else
    m->bx=(uint16_t)(ec_math(m,m->bx,2,16,0,1));
#endif
L_11F3B: ; /* IDA 0x11F3B; jb      short sub_11F5A */
    if(m->flags&1) {sr_tail(m,0x1f5a,h);return;}
L_11F3D: ; /* IDA 0x11F3D; call    sub_1D4A3 */
    sr_call(m,54435,8000,h);
L_11F40: ; /* IDA 0x11F40; cmp     byte ptr ds:98A2h, 0 */
    ec_cmp(m,dl_read8(m,m->ds,39074),0,8);
L_11F45: ; /* IDA 0x11F45; jz      short locret_11F59 */
    if(m->flags&0x40) goto L_11F59;
L_11F47: ; /* IDA 0x11F47; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_11F4A: ; /* IDA 0x11F4A; mov     ds:98A2h, ah */
    dl_write8(m,m->ds,39074,(uint8_t)(m->ax>>8));
L_11F4E: ; /* IDA 0x11F4E; call    sub_20000 */
    mx_far(m,0x1f4e,0x1f53,h);
L_11F53: ; /* IDA 0x11F53; mov     word ptr ds:988Ah, 0FFFFh */
    dl_write16(m,m->ds,39050,65535);
L_11F59: ; /* IDA 0x11F59; retn */
    return;
}
/* Original sub_11F5A; IDA 0x11F5A. */
static void sr_body_11F5A(KiMachine16 *m,const KiEconomyHooks *h) {
L_11F5A: ; /* IDA 0x11F5A; mov     dx, ds:9886h */
    m->dx=(uint16_t)(dl_read16(m,m->ds,39046));
L_11F5E: ; /* IDA 0x11F5E; mov     bx, ds:9888h */
    m->bx=(uint16_t)(dl_read16(m,m->ds,39048));
L_11F62: ; /* IDA 0x11F62; cmp     dx, ds:988Ah */
    ec_cmp(m,m->dx,dl_read16(m,m->ds,39050),16);
L_11F66: ; /* IDA 0x11F66; jnz     short loc_11F6E */
    if(!(m->flags&0x40)) goto L_11F6E;
L_11F68: ; /* IDA 0x11F68; cmp     bx, ds:988Ch */
    ec_cmp(m,m->bx,dl_read16(m,m->ds,39052),16);
L_11F6C: ; /* IDA 0x11F6C; jz      short locret_11F7E */
    if(m->flags&0x40) goto L_11F7E;
L_11F6E: ; /* IDA 0x11F6E; mov     ax, 9 */
    m->ax=(uint16_t)(9);
L_11F71: ; /* IDA 0x11F71; mov     ds:988Ah, dx */
#if KI_RESUME_MUTATION == 2
    dl_write16(m,m->ds,39052,m->dx);
#else
    dl_write16(m,m->ds,39050,m->dx);
#endif
L_11F75: ; /* IDA 0x11F75; mov     ds:988Ch, bx */
    dl_write16(m,m->ds,39052,m->bx);
L_11F79: ; /* IDA 0x11F79; call    sub_20000 */
    mx_far(m,0x1f79,0x1f7e,h);
L_11F7E: ; /* IDA 0x11F7E; retn */
    return;
}
/* Original sub_13B08; IDA 0x13B08. */
static void sr_body_13B08(KiMachine16 *m,const KiEconomyHooks *h) {
L_13B08: ; /* IDA 0x13B08; push    ax */
    ec_push(m,m->ax);
L_13B09: ; /* IDA 0x13B09; push    bx */
    ec_push(m,m->bx);
L_13B0A: ; /* IDA 0x13B0A; push    cx */
    ec_push(m,m->cx);
L_13B0B: ; /* IDA 0x13B0B; push    dx */
    ec_push(m,m->dx);
L_13B0C: ; /* IDA 0x13B0C; push    bp */
    ec_push(m,m->bp);
L_13B0D: ; /* IDA 0x13B0D; push    di */
    ec_push(m,m->di);
L_13B0E: ; /* IDA 0x13B0E; call    sub_102C2 */
    sr_call(m,706,15121,h);
L_13B11: ; /* IDA 0x13B11; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_13B13: ; /* IDA 0x13B13; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_13B15: ; /* IDA 0x13B15; mov     bp, ax */
    m->bp=(uint16_t)(m->ax);
L_13B17: ; /* IDA 0x13B17; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_13B1A: ; /* IDA 0x13B1A; call    sub_20000 */
    mx_far(m,0x3b1a,0x3b1f,h);
L_13B1F: ; /* IDA 0x13B1F; mov     al, 0 */
    vg_al(m,(uint8_t)(0));
L_13B21: ; /* IDA 0x13B21; call    sub_13D09 */
    sr_call(m,15625,15140,h);
L_13B24: ; /* IDA 0x13B24; mov     al, 6 */
    vg_al(m,(uint8_t)(6));
L_13B26: ; /* IDA 0x13B26; call    sub_10241 */
    sr_call(m,577,15145,h);
L_13B29: ; /* IDA 0x13B29; call    sub_13C99 */
    sr_call(m,15513,15148,h);
L_13B2C: ; /* IDA 0x13B2C; add     cx, 3 */
#if KI_RESUME_MUTATION == 10
    m->cx=(uint16_t)(ec_math(m,m->cx,2,16,0,0));
#else
    m->cx=(uint16_t)(ec_math(m,m->cx,3,16,0,0));
#endif
L_13B2F: ; /* IDA 0x13B2F; call    sub_13CDC */
    sr_call(m,15580,15154,h);
L_13B32: ; /* IDA 0x13B32; inc     cx */
    m->cx=(uint16_t)(ec_inc(m,m->cx,16));
L_13B33: ; /* IDA 0x13B33; and     bp, bp */
    {uint16_t result=((m->bp)&(m->bp))&65535;ec_logic(m,result,16);m->bp=(uint16_t)(result);}
L_13B35: ; /* IDA 0x13B35; jz      short loc_13B3A */
    if(m->flags&0x40) goto L_13B3A;
L_13B37: ; /* IDA 0x13B37; add     cx, 3 */
#if KI_RESUME_MUTATION == 12
    m->cx=(uint16_t)(ec_math(m,m->cx,0,16,0,0));
#else
    m->cx=(uint16_t)(ec_math(m,m->cx,3,16,0,0));
#endif
L_13B3A: ; /* IDA 0x13B3A; call    sub_13C99 */
    sr_call(m,15513,15165,h);
L_13B3D: ; /* IDA 0x13B3D; call    sub_102C2 */
    sr_call(m,706,15168,h);
L_13B40: ; /* IDA 0x13B40; call    sub_13D45 */
#if KI_RESUME_MUTATION == 11
    ; /* lost original scene restore */
#else
    sr_call(m,15685,15171,h);
#endif
L_13B43: ; /* IDA 0x13B43; call    sub_11D46 */
    sr_call(m,7494,15174,h);
L_13B46: ; /* IDA 0x13B46; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_13B49: ; /* IDA 0x13B49; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_13B4B: ; /* IDA 0x13B4B; call    sub_20000 */
    mx_far(m,0x3b4b,0x3b50,h);
L_13B50: ; /* IDA 0x13B50; call    sub_19321 */
    sr_call(m,37665,15187,h);
L_13B53: ; /* IDA 0x13B53; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_13B54: ; /* IDA 0x13B54; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_13B55: ; /* IDA 0x13B55; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_13B56: ; /* IDA 0x13B56; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_13B57: ; /* IDA 0x13B57; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_13B58: ; /* IDA 0x13B58; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_13B59: ; /* IDA 0x13B59; retn */
    return;
}
/* Original sub_15C58; IDA 0x15C58. */
static void sr_body_15C58(KiMachine16 *m,const KiEconomyHooks *h) {
L_15C58: ; /* IDA 0x15C58; mov     dx, ds:9892h */
    m->dx=(uint16_t)(dl_read16(m,m->ds,39058));
L_15C5C: ; /* IDA 0x15C5C; mov     bx, ds:9894h */
    m->bx=(uint16_t)(dl_read16(m,m->ds,39060));
L_15C60: ; /* IDA 0x15C60; cmp     dx, ds:988Eh */
    ec_cmp(m,m->dx,dl_read16(m,m->ds,39054),16);
L_15C64: ; /* IDA 0x15C64; jnz     short loc_15C6C */
    if(!(m->flags&0x40)) goto L_15C6C;
L_15C66: ; /* IDA 0x15C66; cmp     bx, ds:9890h */
    ec_cmp(m,m->bx,dl_read16(m,m->ds,39056),16);
L_15C6A: ; /* IDA 0x15C6A; jz      short locret_15CA3 */
    if(m->flags&0x40) goto L_15CA3;
L_15C6C: ; /* IDA 0x15C6C; call    sub_101B4 */
    sr_call(m,436,23663,h);
L_15C6F: ; /* IDA 0x15C6F; cmp     dx, 0FFFFh */
    ec_cmp(m,m->dx,65535,16);
L_15C72: ; /* IDA 0x15C72; jz      short loc_15C82 */
    if(m->flags&0x40) goto L_15C82;
L_15C74: ; /* IDA 0x15C74; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15C76: ; /* IDA 0x15C76; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_15C78: ; /* IDA 0x15C78; add     dx, 1B8h */
    m->dx=(uint16_t)(ec_math(m,m->dx,440,16,0,0));
L_15C7C: ; /* IDA 0x15C7C; add     bx, 28h ; '(' */
    m->bx=(uint16_t)(ec_math(m,m->bx,40,16,0,0));
L_15C7F: ; /* IDA 0x15C7F; call    sub_19541 */
    sr_call(m,38209,23682,h);
L_15C82: ; /* IDA 0x15C82; mov     dx, ds:988Eh */
    m->dx=(uint16_t)(dl_read16(m,m->ds,39054));
L_15C86: ; /* IDA 0x15C86; mov     bx, ds:9890h */
    m->bx=(uint16_t)(dl_read16(m,m->ds,39056));
L_15C8A: ; /* IDA 0x15C8A; mov     ds:9892h, dx */
#if KI_RESUME_MUTATION == 4
    ; /* lost remembered minimap X */
#else
    dl_write16(m,m->ds,39058,m->dx);
#endif
L_15C8E: ; /* IDA 0x15C8E; mov     ds:9894h, bx */
    dl_write16(m,m->ds,39060,m->bx);
L_15C92: ; /* IDA 0x15C92; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15C94: ; /* IDA 0x15C94; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_15C96: ; /* IDA 0x15C96; add     dx, 1B8h */
    m->dx=(uint16_t)(ec_math(m,m->dx,440,16,0,0));
L_15C9A: ; /* IDA 0x15C9A; add     bx, 28h ; '(' */
    m->bx=(uint16_t)(ec_math(m,m->bx,40,16,0,0));
L_15C9D: ; /* IDA 0x15C9D; call    sub_196ED */
    sr_call(m,38637,23712,h);
L_15CA0: ; /* IDA 0x15CA0; call    sub_101DB */
    sr_call(m,475,23715,h);
L_15CA3: ; /* IDA 0x15CA3; retn */
    return;
}
/* Original sub_19541; IDA 0x19541. */
static void sr_body_19541(KiMachine16 *m,const KiEconomyHooks *h) {
L_19541: ; /* IDA 0x19541; push    ds */
    ec_push(m,m->ds);
L_19542: ; /* IDA 0x19542; push    es */
    ec_push(m,m->es);
L_19543: ; /* IDA 0x19543; push    ax */
    ec_push(m,m->ax);
L_19544: ; /* IDA 0x19544; push    bx */
    ec_push(m,m->bx);
L_19545: ; /* IDA 0x19545; push    cx */
    ec_push(m,m->cx);
L_19546: ; /* IDA 0x19546; push    dx */
    ec_push(m,m->dx);
L_19547: ; /* IDA 0x19547; push    si */
    ec_push(m,m->si);
L_19548: ; /* IDA 0x19548; push    di */
    ec_push(m,m->di);
L_19549: ; /* IDA 0x19549; cld */
    m->flags&=(uint16_t)~0x400u;
L_1954A: ; /* IDA 0x1954A; mov     ds, cs:word_10D3A */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3386));
L_1954F: ; /* IDA 0x1954F; cmp     dx, 260h */
    ec_cmp(m,m->dx,608,16);
L_19553: ; /* IDA 0x19553; jb      short loc_19558 */
    if(m->flags&1) goto L_19558;
L_19555: ; /* IDA 0x19555; mov     dx, 25Fh */
#if KI_RESUME_MUTATION == 5
    m->dx=(uint16_t)(608);
#else
    m->dx=(uint16_t)(607);
#endif
L_19558: ; /* IDA 0x19558; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1955A: ; /* IDA 0x1955A; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1955C: ; /* IDA 0x1955C; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1955E: ; /* IDA 0x1955E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19560: ; /* IDA 0x19560; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19562: ; /* IDA 0x19562; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19564: ; /* IDA 0x19564; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_19566: ; /* IDA 0x19566; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19568: ; /* IDA 0x19568; add     si, bx */
    m->si=(uint16_t)(ec_math(m,m->si,m->bx,16,0,0));
L_1956A: ; /* IDA 0x1956A; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_1956C: ; /* IDA 0x1956C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1956E: ; /* IDA 0x1956E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19570: ; /* IDA 0x19570; add     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
L_19572: ; /* IDA 0x19572; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_19574: ; /* IDA 0x19574; sub     dx, 37h ; '7' */
    m->dx=(uint16_t)(ec_math(m,m->dx,55,16,0,1));
L_19577: ; /* IDA 0x19577; sub     si, 3C0h */
    m->si=(uint16_t)(ec_math(m,m->si,960,16,0,1));
L_1957B: ; /* IDA 0x1957B; add     si, dx */
    m->si=(uint16_t)(ec_math(m,m->si,m->dx,16,0,0));
L_1957D: ; /* IDA 0x1957D; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_19580: ; /* IDA 0x19580; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_19582: ; /* IDA 0x19582; call    sub_10CC3 */
    sr_call(m,3267,38277,h);
L_19585: ; /* IDA 0x19585; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_19588: ; /* IDA 0x19588; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_1958A: ; /* IDA 0x1958A; out     dx, al; EGA: sequencer address reg */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1958B: ; /* IDA 0x1958B; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_1958D: ; /* IDA 0x1958D; call    sub_195B1 */
    sr_call(m,38321,38288,h);
L_19590: ; /* IDA 0x19590; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_19592: ; /* IDA 0x19592; call    sub_195B1 */
    sr_call(m,38321,38293,h);
L_19595: ; /* IDA 0x19595; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_19598: ; /* IDA 0x19598; call    sub_195B1 */
    sr_call(m,38321,38299,h);
L_1959B: ; /* IDA 0x1959B; mov     ax, 8 */
    m->ax=(uint16_t)(8);
L_1959E: ; /* IDA 0x1959E; call    sub_195B1 */
    sr_call(m,38321,38305,h);
L_195A1: ; /* IDA 0x195A1; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_195A4: ; /* IDA 0x195A4; mov     ax, 0F02h */
    m->ax=(uint16_t)(3842);
L_195A7: ; /* IDA 0x195A7; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_195A8: ; /* IDA 0x195A8; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_195A9: ; /* IDA 0x195A9; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_195AA: ; /* IDA 0x195AA; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_195AB: ; /* IDA 0x195AB; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_195AC: ; /* IDA 0x195AC; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_195AD: ; /* IDA 0x195AD; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_195AE: ; /* IDA 0x195AE; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_195AF: ; /* IDA 0x195AF; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_195B0: ; /* IDA 0x195B0; retn */
    return;
}
/* Original sub_195B1; IDA 0x195B1. */
static void sr_body_195B1(KiMachine16 *m,const KiEconomyHooks *h) {
L_195B1: ; /* IDA 0x195B1; mov     dx, 3C5h */
    m->dx=(uint16_t)(965);
L_195B4: ; /* IDA 0x195B4; out     dx, al; EGA port: sequencer data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_195B5: ; /* IDA 0x195B5; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_195B7: ; /* IDA 0x195B7; mov     cx, 0Bh */
    m->cx=(uint16_t)(11);
L_195BA: ; /* IDA 0x195BA; movsw */
    mc_string(m,0,0);
L_195BB: ; /* IDA 0x195BB; movsw */
    mc_string(m,0,0);
L_195BC: ; /* IDA 0x195BC; add     di, 4Ch ; 'L' */
    m->di=(uint16_t)(ec_math(m,m->di,76,16,0,0));
L_195BF: ; /* IDA 0x195BF; add     si, 14h */
#if KI_RESUME_MUTATION == 6
    m->si=(uint16_t)(ec_math(m,m->si,19,16,0,0));
#else
    m->si=(uint16_t)(ec_math(m,m->si,20,16,0,0));
#endif
L_195C2: ; /* IDA 0x195C2; loop    loc_195BA */
    if(--m->cx) goto L_195BA;
L_195C4: ; /* IDA 0x195C4; add     si, 0AF8h */
#if KI_RESUME_MUTATION == 7
    m->si=(uint16_t)(ec_math(m,m->si,2807,16,0,0));
#else
    m->si=(uint16_t)(ec_math(m,m->si,2808,16,0,0));
#endif
L_195C8: ; /* IDA 0x195C8; retn */
    return;
}
/* Original sub_196ED; IDA 0x196ED. */
static void sr_body_196ED(KiMachine16 *m,const KiEconomyHooks *h) {
L_196ED: ; /* IDA 0x196ED; push    ds */
    ec_push(m,m->ds);
L_196EE: ; /* IDA 0x196EE; push    es */
    ec_push(m,m->es);
L_196EF: ; /* IDA 0x196EF; push    ax */
    ec_push(m,m->ax);
L_196F0: ; /* IDA 0x196F0; push    bx */
    ec_push(m,m->bx);
L_196F1: ; /* IDA 0x196F1; push    cx */
    ec_push(m,m->cx);
L_196F2: ; /* IDA 0x196F2; push    dx */
    ec_push(m,m->dx);
L_196F3: ; /* IDA 0x196F3; push    si */
    ec_push(m,m->si);
L_196F4: ; /* IDA 0x196F4; push    di */
    ec_push(m,m->di);
L_196F5: ; /* IDA 0x196F5; push    bp */
    ec_push(m,m->bp);
L_196F6: ; /* IDA 0x196F6; cld */
    m->flags&=(uint16_t)~0x400u;
L_196F7: ; /* IDA 0x196F7; mov     ch, dl */
    vg_ch(m,(uint8_t)((uint8_t)m->dx));
L_196F9: ; /* IDA 0x196F9; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_196FB: ; /* IDA 0x196FB; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_196FD: ; /* IDA 0x196FD; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_196FF: ; /* IDA 0x196FF; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19701: ; /* IDA 0x19701; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19703: ; /* IDA 0x19703; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19705: ; /* IDA 0x19705; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_19707: ; /* IDA 0x19707; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_19709: ; /* IDA 0x19709; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1970B: ; /* IDA 0x1970B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1970D: ; /* IDA 0x1970D; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_1970F: ; /* IDA 0x1970F; and     ch, 7 */
    {uint16_t result=(((uint8_t)(m->cx>>8))&(7))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_19712: ; /* IDA 0x19712; mov     ds, cs:word_10D4C */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3404));
L_19717: ; /* IDA 0x19717; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_19719: ; /* IDA 0x19719; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_1971C: ; /* IDA 0x1971C; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1971E: ; /* IDA 0x1971E; call    sub_10CAC */
    sr_call(m,3244,38689,h);
L_19721: ; /* IDA 0x19721; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_19723: ; /* IDA 0x19723; call    sub_19752 */
    sr_call(m,38738,38694,h);
L_19726: ; /* IDA 0x19726; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_19729: ; /* IDA 0x19729; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1972B: ; /* IDA 0x1972B; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1972C: ; /* IDA 0x1972C; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1972D: ; /* IDA 0x1972D; mov     al, 10h */
    vg_al(m,(uint8_t)(16));
L_1972F: ; /* IDA 0x1972F; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_19730: ; /* IDA 0x19730; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_19731: ; /* IDA 0x19731; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_19733: ; /* IDA 0x19733; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_19734: ; /* IDA 0x19734; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_19736: ; /* IDA 0x19736; call    sub_19752 */
    sr_call(m,38738,38713,h);
L_19739: ; /* IDA 0x19739; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_1973B: ; /* IDA 0x1973B; call    sub_19752 */
    sr_call(m,38738,38718,h);
L_1973E: ; /* IDA 0x1973E; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_19740: ; /* IDA 0x19740; call    sub_19752 */
    sr_call(m,38738,38723,h);
L_19743: ; /* IDA 0x19743; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_19745: ; /* IDA 0x19745; call    sub_19752 */
    sr_call(m,38738,38728,h);
L_19748: ; /* IDA 0x19748; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_19749: ; /* IDA 0x19749; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1974A: ; /* IDA 0x1974A; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1974B: ; /* IDA 0x1974B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1974C: ; /* IDA 0x1974C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1974D: ; /* IDA 0x1974D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1974E: ; /* IDA 0x1974E; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1974F: ; /* IDA 0x1974F; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_19750: ; /* IDA 0x19750; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_19751: ; /* IDA 0x19751; retn */
    return;
}
/* Original sub_19752; IDA 0x19752. */
static void sr_body_19752(KiMachine16 *m,const KiEconomyHooks *h) {
L_19752: ; /* IDA 0x19752; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_19755: ; /* IDA 0x19755; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_19756: ; /* IDA 0x19756; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_19758: ; /* IDA 0x19758; mov     bp, 0Bh */
#if KI_RESUME_MUTATION == 9
    m->bp=(uint16_t)(10);
#else
    m->bp=(uint16_t)(11);
#endif
L_1975B: ; /* IDA 0x1975B; lodsw */
    mc_string(m,1,0);
L_1975C: ; /* IDA 0x1975C; mov     dh, [si] */
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
L_1975E: ; /* IDA 0x1975E; xor     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_19760: ; /* IDA 0x19760; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_19761: ; /* IDA 0x19761; xchg    al, ah */
    {uint16_t a=(uint8_t)m->ax,b=(uint8_t)(m->ax>>8);vg_al(m,(uint8_t)(b));vg_ah(m,(uint8_t)(a));}
L_19763: ; /* IDA 0x19763; and     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))&((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_19765: ; /* IDA 0x19765; jz      short loc_19771 */
    if(m->flags&0x40) goto L_19771;
L_19767: ; /* IDA 0x19767; mov     cl, ch */
    vg_cl(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_19769: ; /* IDA 0x19769; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1976B: ; /* IDA 0x1976B; rcr     dx, 1 */
#if KI_RESUME_MUTATION == 8
    ; /* lost rotation through carry */
#else
    m->dx=hr_rcr(m,m->dx);
#endif
L_1976D: ; /* IDA 0x1976D; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_1976F: ; /* IDA 0x1976F; jnz     short loc_19769 */
    if(!(m->flags&0x40)) goto L_19769;
L_19771: ; /* IDA 0x19771; mov     cl, es:[di] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_19774: ; /* IDA 0x19774; mov     es:[di], ah */
    dl_write8(m,m->es,(uint16_t)(m->di),(uint8_t)(m->ax>>8));
L_19777: ; /* IDA 0x19777; mov     cl, es:[di+1] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di+1))));
L_1977B: ; /* IDA 0x1977B; mov     es:[di+1], al */
    dl_write8(m,m->es,(uint16_t)(m->di+1),(uint8_t)m->ax);
L_1977F: ; /* IDA 0x1977F; mov     cl, es:[di+2] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di+2))));
L_19783: ; /* IDA 0x19783; mov     es:[di+2], dh */
    dl_write8(m,m->es,(uint16_t)(m->di+2),(uint8_t)(m->dx>>8));
L_19787: ; /* IDA 0x19787; mov     cl, es:[di+3] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di+3))));
L_1978B: ; /* IDA 0x1978B; mov     es:[di+3], dl */
    dl_write8(m,m->es,(uint16_t)(m->di+3),(uint8_t)m->dx);
L_1978F: ; /* IDA 0x1978F; add     di, 50h ; 'P' */
    m->di=(uint16_t)(ec_math(m,m->di,80,16,0,0));
L_19792: ; /* IDA 0x19792; dec     bp */
    m->bp=(uint16_t)(st_dec(m,m->bp,16));
L_19793: ; /* IDA 0x19793; jnz     short loc_1975B */
    if(!(m->flags&0x40)) goto L_1975B;
L_19795: ; /* IDA 0x19795; retn */
    return;
}
/* Original sub_1D4A3; IDA 0x1D4A3. */
static void sr_body_1D4A3(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D4A3: ; /* IDA 0x1D4A3; push    ds */
    ec_push(m,m->ds);
L_1D4A4: ; /* IDA 0x1D4A4; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4A6: ; /* IDA 0x1D4A6; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4A8: ; /* IDA 0x1D4A8; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4AA: ; /* IDA 0x1D4AA; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_1D4AC: ; /* IDA 0x1D4AC; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4AE: ; /* IDA 0x1D4AE; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4B0: ; /* IDA 0x1D4B0; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_1D4B2: ; /* IDA 0x1D4B2; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4B4: ; /* IDA 0x1D4B4; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4B6: ; /* IDA 0x1D4B6; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4B8: ; /* IDA 0x1D4B8; mov     ds, cs:word_1D84E */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55374));
L_1D4BD: ; /* IDA 0x1D4BD; test    byte ptr [bx], 10h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(16))&255;ec_logic(m,result,8);}
L_1D4C0: ; /* IDA 0x1D4C0; jnz     short loc_1D4C5 */
    if(!(m->flags&0x40)) goto L_1D4C5;
L_1D4C2: ; /* IDA 0x1D4C2; or      byte ptr [bx], 0E0h */
#if KI_RESUME_MUTATION == 3
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))|(96))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->bx),result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))|(224))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->bx),result);}
#endif
L_1D4C5: ; /* IDA 0x1D4C5; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1D4C6: ; /* IDA 0x1D4C6; retn */
    return;
}
int ki_resume_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x1f30:sr_body_11F30(m,h);break;
case 0x1f5a:sr_body_11F5A(m,h);break;
case 0x3b08:sr_body_13B08(m,h);break;
case 0x5c58:sr_body_15C58(m,h);break;
case 0x9541:sr_body_19541(m,h);break;
case 0x95b1:sr_body_195B1(m,h);break;
case 0x96ed:sr_body_196ED(m,h);break;
case 0x9752:sr_body_19752(m,h);break;
case 0xd4a3:sr_body_1D4A3(m,h);break;
default:return 0;}return 1;}
