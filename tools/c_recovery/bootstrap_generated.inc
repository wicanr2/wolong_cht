/* Original native C world bootstrap; re/126, spec/246. Full 11BE0 is not registered. */
/* Original sub_1533D; IDA 0x1533D. */
static void bs_body_1533D(KiMachine16 *m,const KiEconomyHooks *h) {
#if KI_BOOTSTRAP_MUTATION == 9
const uint16_t entry_ds=m->ds;
#endif
L_1533D: ; /* IDA 0x1533D; push    ax */
    ec_push(m,m->ax);
L_1533E: ; /* IDA 0x1533E; push    bx */
    ec_push(m,m->bx);
L_1533F: ; /* IDA 0x1533F; push    cx */
    ec_push(m,m->cx);
L_15340: ; /* IDA 0x15340; push    dx */
    ec_push(m,m->dx);
L_15341: ; /* IDA 0x15341; push    si */
    ec_push(m,m->si);
L_15342: ; /* IDA 0x15342; push    di */
    ec_push(m,m->di);
L_15343: ; /* IDA 0x15343; mov     ds, word ptr ds:0D52h */
    m->ds=(uint16_t)(dl_read16(m,m->ds,3410));
L_15347: ; /* IDA 0x15347; call    sub_155A6 */
    bs_call(m,21926,21322,h);
L_1534A: ; /* IDA 0x1534A; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_1534C: ; /* IDA 0x1534C; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1534E: ; /* IDA 0x1534E; call    sub_11E17 */
    bs_call(m,7703,21329,h);
L_15351: ; /* IDA 0x15351; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_15352: ; /* IDA 0x15352; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_15353: ; /* IDA 0x15353; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_15354: ; /* IDA 0x15354; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_15355: ; /* IDA 0x15355; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_15356: ; /* IDA 0x15356; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15357: ; /* IDA 0x15357; retn */
#if KI_BOOTSTRAP_MUTATION == 9
    m->ds=entry_ds;return;
#else
    return;
#endif
}
/* Original sub_189F0; IDA 0x189F0. */
static void bs_body_189F0(KiMachine16 *m,const KiEconomyHooks *h) {
L_189F0: ; /* IDA 0x189F0; push    ds */
    ec_push(m,m->ds);
L_189F1: ; /* IDA 0x189F1; push    es */
    ec_push(m,m->es);
L_189F2: ; /* IDA 0x189F2; push    ax */
    ec_push(m,m->ax);
L_189F3: ; /* IDA 0x189F3; push    bx */
    ec_push(m,m->bx);
L_189F4: ; /* IDA 0x189F4; push    cx */
    ec_push(m,m->cx);
L_189F5: ; /* IDA 0x189F5; push    si */
    ec_push(m,m->si);
L_189F6: ; /* IDA 0x189F6; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_189FB: ; /* IDA 0x189FB; mov     cx, 0C0h */
    m->cx=(uint16_t)(192);
L_189FE: ; /* IDA 0x189FE; mov     si, 840h */
    m->si=(uint16_t)(2112);
L_18A01: ; /* IDA 0x18A01; call    sub_18A1E */
    bs_call(m,35358,35332,h);
L_18A04: ; /* IDA 0x18A04; add     si, 20h ; ' ' */
    m->si=(uint16_t)(ec_math(m,m->si,32,16,0,0));
L_18A07: ; /* IDA 0x18A07; loop    loc_18A01 */
    if(--m->cx) goto L_18A01;
L_18A09: ; /* IDA 0x18A09; mov     cx, 7Fh */
#if KI_BOOTSTRAP_MUTATION == 8
    m->cx=128;
#else
    m->cx=(uint16_t)(127);
#endif
L_18A0C: ; /* IDA 0x18A0C; mov     si, 2240h */
    m->si=(uint16_t)(8768);
L_18A0F: ; /* IDA 0x18A0F; call    sub_18AEA */
    bs_call(m,35562,35346,h);
L_18A12: ; /* IDA 0x18A12; add     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
L_18A15: ; /* IDA 0x18A15; loop    loc_18A0F */
    if(--m->cx) goto L_18A0F;
L_18A17: ; /* IDA 0x18A17; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_18A18: ; /* IDA 0x18A18; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_18A19: ; /* IDA 0x18A19; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_18A1A: ; /* IDA 0x18A1A; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_18A1B: ; /* IDA 0x18A1B; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_18A1C: ; /* IDA 0x18A1C; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_18A1D: ; /* IDA 0x18A1D; retn */
    return;
}
/* Original sub_18A1E; IDA 0x18A1E. */
static void bs_body_18A1E(KiMachine16 *m,const KiEconomyHooks *h) {
L_18A1E: ; /* IDA 0x18A1E; push    cx */
    ec_push(m,m->cx);
L_18A1F: ; /* IDA 0x18A1F; mov     ax, [si+0Ah] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+10)));
L_18A22: ; /* IDA 0x18A22; sub     ax, 2 */
    m->ax=(uint16_t)(ec_math(m,m->ax,2,16,0,1));
L_18A25: ; /* IDA 0x18A25; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18A27: ; /* IDA 0x18A27; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18A29: ; /* IDA 0x18A29; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18A2B: ; /* IDA 0x18A2B; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_18A2D: ; /* IDA 0x18A2D; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18A2F: ; /* IDA 0x18A2F; add     ax, bx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->bx,16,0,0));
L_18A31: ; /* IDA 0x18A31; add     ax, cs:word_10D44 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,3396),16,0,0));
L_18A36: ; /* IDA 0x18A36; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_18A38: ; /* IDA 0x18A38; mov     bx, [si+8] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+8)));
L_18A3B: ; /* IDA 0x18A3B; add     bx, 300h */
    m->bx=(uint16_t)(ec_math(m,m->bx,768,16,0,0));
L_18A3F: ; /* IDA 0x18A3F; mov     al, es:[bx] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->bx))));
L_18A42: ; /* IDA 0x18A42; sub     al, 0CBh */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,203,8,0,1)));
L_18A44: ; /* IDA 0x18A44; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_18A46: ; /* IDA 0x18A46; mov     cl, 3 */
#if KI_BOOTSTRAP_MUTATION == 1
    vg_cl(m,2);
#else
    vg_cl(m,(uint8_t)(3));
#endif
L_18A48: ; /* IDA 0x18A48; div     cl */
    dl_div(m,(uint8_t)m->cx,8);
L_18A4A: ; /* IDA 0x18A4A; mul     cl */
    dl_mul(m,(uint8_t)m->cx,8);
L_18A4C: ; /* IDA 0x18A4C; add     al, 0CBh */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,203,8,0,0)));
L_18A4E: ; /* IDA 0x18A4E; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_18A51: ; /* IDA 0x18A51; cmp     ah, 18h */
    ec_cmp(m,(uint8_t)(m->ax>>8),24,8);
L_18A54: ; /* IDA 0x18A54; jnz     short loc_18A5A */
    if(!(m->flags&0x40)) goto L_18A5A;
L_18A56: ; /* IDA 0x18A56; add     al, 2 */
#if KI_BOOTSTRAP_MUTATION == 2
    vg_al(m,ec_math(m,(uint8_t)m->ax,1,8,0,0));
#else
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,2,8,0,0)));
#endif
L_18A58: ; /* IDA 0x18A58; jmp     short loc_18A63 */
    goto L_18A63;
L_18A5A: ; /* IDA 0x18A5A; cmp     ah, cs:byte_10CFF */
    ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->cs,3327),8);
L_18A5F: ; /* IDA 0x18A5F; jz      short loc_18A63 */
    if(m->flags&0x40) goto L_18A63;
L_18A61: ; /* IDA 0x18A61; add     al, 1 */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,1,8,0,0)));
L_18A63: ; /* IDA 0x18A63; mov     es:[bx], al */
    dl_write8(m,m->es,(uint16_t)(m->bx),(uint8_t)m->ax);
L_18A66: ; /* IDA 0x18A66; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_18A68: ; /* IDA 0x18A68; mov     al, [si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_18A6B: ; /* IDA 0x18A6B; cmp     al, cs:byte_10CFF */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
L_18A70: ; /* IDA 0x18A70; jnz     short loc_18A74 */
    if(!(m->flags&0x40)) goto L_18A74;
L_18A72: ; /* IDA 0x18A72; mov     ah, 0Ah */
#if KI_BOOTSTRAP_MUTATION == 3
    vg_ah(m,0);
#else
    vg_ah(m,(uint8_t)(10));
#endif
L_18A74: ; /* IDA 0x18A74; mov     al, [si+16h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+22))));
L_18A77: ; /* IDA 0x18A77; and     al, 0Fh */
    {uint16_t result=(((uint8_t)m->ax)&(15))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_18A79: ; /* IDA 0x18A79; jz      short loc_18A99 */
    if(m->flags&0x40) goto L_18A99;
L_18A7B: ; /* IDA 0x18A7B; cmp     al, 3 */
    ec_cmp(m,(uint8_t)m->ax,3,8);
L_18A7D: ; /* IDA 0x18A7D; jz      short loc_18AB5 */
    if(m->flags&0x40) goto L_18AB5;
L_18A7F: ; /* IDA 0x18A7F; sub     bx, 181h */
    m->bx=(uint16_t)(ec_math(m,m->bx,385,16,0,1));
L_18A83: ; /* IDA 0x18A83; call    sub_18AD1 */
    bs_call(m,35537,35462,h);
L_18A86: ; /* IDA 0x18A86; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_18A87: ; /* IDA 0x18A87; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_18A88: ; /* IDA 0x18A88; call    sub_18AD1 */
    bs_call(m,35537,35467,h);
L_18A8B: ; /* IDA 0x18A8B; add     bx, 300h */
    m->bx=(uint16_t)(ec_math(m,m->bx,768,16,0,0));
L_18A8F: ; /* IDA 0x18A8F; call    sub_18AD1 */
    bs_call(m,35537,35474,h);
L_18A92: ; /* IDA 0x18A92; dec     bx */
    m->bx=(uint16_t)(st_dec(m,m->bx,16));
L_18A93: ; /* IDA 0x18A93; dec     bx */
    m->bx=(uint16_t)(st_dec(m,m->bx,16));
L_18A94: ; /* IDA 0x18A94; call    sub_18AD1 */
    bs_call(m,35537,35479,h);
L_18A97: ; /* IDA 0x18A97; jmp     short loc_18ACF */
    goto L_18ACF;
L_18A99: ; /* IDA 0x18A99; sub     bx, 302h */
    m->bx=(uint16_t)(ec_math(m,m->bx,770,16,0,1));
L_18A9D: ; /* IDA 0x18A9D; call    sub_18AD1 */
    bs_call(m,35537,35488,h);
L_18AA0: ; /* IDA 0x18AA0; add     bx, 4 */
    m->bx=(uint16_t)(ec_math(m,m->bx,4,16,0,0));
L_18AA3: ; /* IDA 0x18AA3; call    sub_18AD1 */
    bs_call(m,35537,35494,h);
L_18AA6: ; /* IDA 0x18AA6; add     bx, 600h */
    m->bx=(uint16_t)(ec_math(m,m->bx,1536,16,0,0));
L_18AAA: ; /* IDA 0x18AAA; call    sub_18AD1 */
    bs_call(m,35537,35501,h);
L_18AAD: ; /* IDA 0x18AAD; sub     bx, 4 */
    m->bx=(uint16_t)(ec_math(m,m->bx,4,16,0,1));
L_18AB0: ; /* IDA 0x18AB0; call    sub_18AD1 */
    bs_call(m,35537,35507,h);
L_18AB3: ; /* IDA 0x18AB3; jmp     short loc_18ACF */
    goto L_18ACF;
L_18AB5: ; /* IDA 0x18AB5; sub     bx, 180h */
    m->bx=(uint16_t)(ec_math(m,m->bx,384,16,0,1));
L_18AB9: ; /* IDA 0x18AB9; call    sub_18AD1 */
    bs_call(m,35537,35516,h);
L_18ABC: ; /* IDA 0x18ABC; add     bx, 17Fh */
    m->bx=(uint16_t)(ec_math(m,m->bx,383,16,0,0));
L_18AC0: ; /* IDA 0x18AC0; call    sub_18AD1 */
    bs_call(m,35537,35523,h);
L_18AC3: ; /* IDA 0x18AC3; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_18AC4: ; /* IDA 0x18AC4; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_18AC5: ; /* IDA 0x18AC5; call    sub_18AD1 */
    bs_call(m,35537,35528,h);
L_18AC8: ; /* IDA 0x18AC8; add     bx, 17Fh */
    m->bx=(uint16_t)(ec_math(m,m->bx,383,16,0,0));
L_18ACC: ; /* IDA 0x18ACC; call    sub_18AD1 */
    bs_call(m,35537,35535,h);
L_18ACF: ; /* IDA 0x18ACF; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_18AD0: ; /* IDA 0x18AD0; retn */
    return;
}
/* Original sub_18AD1; IDA 0x18AD1. */
static void bs_body_18AD1(KiMachine16 *m,const KiEconomyHooks *h) {
L_18AD1: ; /* IDA 0x18AD1; mov     al, es:[bx] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->bx))));
L_18AD4: ; /* IDA 0x18AD4; sub     al, 0DEh */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,222,8,0,1)));
L_18AD6: ; /* IDA 0x18AD6; jb      short locret_18AE9 */
    if(m->flags&1) goto L_18AE9;
L_18AD8: ; /* IDA 0x18AD8; cmp     al, 14h */
#if KI_BOOTSTRAP_MUTATION == 4
    ec_cmp(m,(uint8_t)m->ax,21,8);
#else
    ec_cmp(m,(uint8_t)m->ax,20,8);
#endif
L_18ADA: ; /* IDA 0x18ADA; jnb     short locret_18AE9 */
    if(!(m->flags&1)) goto L_18AE9;
L_18ADC: ; /* IDA 0x18ADC; cmp     al, 0Ah */
    ec_cmp(m,(uint8_t)m->ax,10,8);
L_18ADE: ; /* IDA 0x18ADE; jb      short loc_18AE2 */
    if(m->flags&1) goto L_18AE2;
L_18AE0: ; /* IDA 0x18AE0; sub     al, 0Ah */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,10,8,0,1)));
L_18AE2: ; /* IDA 0x18AE2; add     al, 0DEh */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,222,8,0,0)));
L_18AE4: ; /* IDA 0x18AE4; add     al, ah */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8,0,0)));
L_18AE6: ; /* IDA 0x18AE6; mov     es:[bx], al */
    dl_write8(m,m->es,(uint16_t)(m->bx),(uint8_t)m->ax);
L_18AE9: ; /* IDA 0x18AE9; retn */
    return;
}
/* Original sub_18AEA; IDA 0x18AEA. */
static void bs_body_18AEA(KiMachine16 *m,const KiEconomyHooks *h) {
L_18AEA: ; /* IDA 0x18AEA; cmp     byte ptr [si], 80h */
#if KI_BOOTSTRAP_MUTATION == 5
    ec_cmp(m,dl_read8(m,m->ds,m->si),0x81,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
#endif
L_18AED: ; /* IDA 0x18AED; jb      short locret_18B11 */
    if(m->flags&1) goto L_18B11;
L_18AEF: ; /* IDA 0x18AEF; mov     ax, [si+12h] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+18)));
L_18AF2: ; /* IDA 0x18AF2; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18AF4: ; /* IDA 0x18AF4; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18AF6: ; /* IDA 0x18AF6; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18AF8: ; /* IDA 0x18AF8; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_18AFA: ; /* IDA 0x18AFA; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_18AFC: ; /* IDA 0x18AFC; add     ax, bx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->bx,16,0,0));
L_18AFE: ; /* IDA 0x18AFE; add     ax, cs:word_19872 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,39026),16,0,0));
L_18B03: ; /* IDA 0x18B03; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_18B05: ; /* IDA 0x18B05; mov     bx, [si+10h] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+16)));
L_18B08: ; /* IDA 0x18B08; mov     [si+1Ch], ax */
#if KI_BOOTSTRAP_MUTATION == 6
    dl_write16(m,m->ds,(uint16_t)(m->si+28),0xb000);
#else
    dl_write16(m,m->ds,(uint16_t)(m->si+28),m->ax);
#endif
L_18B0B: ; /* IDA 0x18B0B; mov     [si+1Ah], bx */
    dl_write16(m,m->ds,(uint16_t)(m->si+26),m->bx);
L_18B0E: ; /* IDA 0x18B0E; inc     byte ptr es:[bx] */
#if KI_BOOTSTRAP_MUTATION == 7
    dl_write8(m,m->es,m->bx,1);
#else
    dl_write8(m,m->es,(uint16_t)(m->bx),ec_inc(m,dl_read8(m,m->es,(uint16_t)(m->bx)),8));
#endif
L_18B11: ; /* IDA 0x18B11; retn */
    return;
}
/* Original sub_11BE0 bounded prefix, not complete function; IDA 0x11BE0. */
static void bs_body_11BE0(KiMachine16 *m,const KiEconomyHooks *h) {
L_11BE0: ; /* IDA 0x11BE0; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_11BE2: ; /* IDA 0x11BE2; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_11BE4: ; /* IDA 0x11BE4; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_11BE6: ; /* IDA 0x11BE6; call    sub_189F0 */
    bs_call(m,35312,7145,h);
L_11BE9: ; /* IDA 0x11BE9; call    sub_1533D */
    bs_call(m,21309,7148,h);
L_11BEC: ; /* IDA 0x11BEC; mov     cs:word_19903, ss */
    dl_write16(m,m->cs,39171,m->ss);
L_11BF1: ; /* IDA 0x11BF1; mov     cs:word_19901, sp */
#if KI_BOOTSTRAP_MUTATION == 10
    dl_write16(m,m->cs,0x9901,(uint16_t)(m->sp+2));
#else
    dl_write16(m,m->cs,39169,m->sp);
#endif
m->ip=0x1bf6; /* Explicit control point; no virtual RET or POP. */
}
int ki_bootstrap_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x533d:bs_body_1533D(m,h);break;
case 0x89f0:bs_body_189F0(m,h);break;
case 0x8a1e:bs_body_18A1E(m,h);break;
case 0x8ad1:bs_body_18AD1(m,h);break;
case 0x8aea:bs_body_18AEA(m,h);break;
default:return 0;}return 1;}
static int bs_known(uint16_t t) {switch(t) {case 0x533d: case 0x89f0: case 0x8a1e: case 0x8ad1: case 0x8aea: return 1;default:return 0;}}
void sub_1533D(KiMachine16 *m,const KiEconomyHooks *h) {ki_bootstrap_body(m,0x533d,h);m->ip=ec_pop(m);}
void sub_189F0(KiMachine16 *m,const KiEconomyHooks *h) {ki_bootstrap_body(m,0x89f0,h);m->ip=ec_pop(m);}
void sub_18A1E(KiMachine16 *m,const KiEconomyHooks *h) {ki_bootstrap_body(m,0x8a1e,h);m->ip=ec_pop(m);}
void sub_18AD1(KiMachine16 *m,const KiEconomyHooks *h) {ki_bootstrap_body(m,0x8ad1,h);m->ip=ec_pop(m);}
void sub_18AEA(KiMachine16 *m,const KiEconomyHooks *h) {ki_bootstrap_body(m,0x8aea,h);m->ip=ec_pop(m);}
void raw_entry_11BE0_prefix(KiMachine16 *m,const KiEconomyHooks *h) {bs_body_11BE0(m,h);}
