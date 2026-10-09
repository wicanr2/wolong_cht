/* Native C original overlay producers; re/115, spec/235. Original locators retained. */
/* Original sub_11CC9; IDA 0x11CC9. */
static void ov_body_11CC9(KiMachine16 *m,const KiEconomyHooks *h) {
L_11CC9: ; /* IDA 0x11CC9; call    sub_12AF4 */
    ov_call(m,10996,7372,h);
L_11CCC: ; /* IDA 0x11CCC; call    sub_12533 */
    ov_call(m,9523,7375,h);
L_11CCF: ; /* IDA 0x11CCF; retn */
    return;
}
/* Original sub_12533; IDA 0x12533. */
static void ov_body_12533(KiMachine16 *m,const KiEconomyHooks *h) {
L_12533: ; /* IDA 0x12533; mov     es, word_1987A */
    m->es=(uint16_t)(dl_read16(m,m->ds,39034));
L_12537: ; /* IDA 0x12537; mov     ds, word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->ds,3410));
L_1253B: ; /* IDA 0x1253B; mov     si, 2040h */
    m->si=(uint16_t)(8256);
L_1253E: ; /* IDA 0x1253E; mov     bp, 20h ; ' ' */
    m->bp=(uint16_t)(32);
L_12541: ; /* IDA 0x12541; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_12544: ; /* IDA 0x12544; jb      short loc_12598 */
    if(m->flags&1) goto L_12598;
L_12546: ; /* IDA 0x12546; mov     bl, [si+0Eh] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+14))));
L_12549: ; /* IDA 0x12549; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1254B: ; /* IDA 0x1254B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1254D: ; /* IDA 0x1254D; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1254F: ; /* IDA 0x1254F; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_12551: ; /* IDA 0x12551; mov     al, [si+0Fh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+15))));
L_12554: ; /* IDA 0x12554; add     bl, al */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,(uint8_t)m->ax,8,0,0)));
L_12556: ; /* IDA 0x12556; test    byte ptr [si], 1 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(1))&255;ec_logic(m,result,8);}
L_12559: ; /* IDA 0x12559; jz      short loc_12565 */
    if(m->flags&0x40) goto L_12565;
L_1255B: ; /* IDA 0x1255B; and     byte ptr [si], 0FEh */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(254))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_1255E: ; /* IDA 0x1255E; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_12560: ; /* IDA 0x12560; and     al, 7 */
#if KI_OVERLAY_MUTATION == 5
    {uint16_t result=(((uint8_t)m->ax)&(3))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
#else
    {uint16_t result=(((uint8_t)m->ax)&(7))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
#endif
L_12562: ; /* IDA 0x12562; mov     [si+0Fh], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+15),(uint8_t)m->ax);
L_12565: ; /* IDA 0x12565; mov     bl, cs:[bx-67A6h] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->cs,(uint16_t)(m->bx+39002))));
L_1256A: ; /* IDA 0x1256A; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1256C: ; /* IDA 0x1256C; shl     bx, 1 */
#if KI_OVERLAY_MUTATION == 6
    ; /* lost metadata stride shift */
#else
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
#endif
L_1256E: ; /* IDA 0x1256E; xor     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_12570: ; /* IDA 0x12570; mov     cl, es:[bx] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->bx))));
L_12573: ; /* IDA 0x12573; mov     al, es:[bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->bx+1))));
L_12577: ; /* IDA 0x12577; mov     di, es:[bx+2] */
    m->di=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+2)));
L_1257B: ; /* IDA 0x1257B; add     di, 100h */
#if KI_OVERLAY_MUTATION == 7
    m->di=(uint16_t)(ec_math(m,m->di,128,16,0,0));
#else
    m->di=(uint16_t)(ec_math(m,m->di,256,16,0,0));
#endif
L_1257F: ; /* IDA 0x1257F; mov     dx, cx */
    m->dx=(uint16_t)(m->cx);
L_12581: ; /* IDA 0x12581; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_12583: ; /* IDA 0x12583; sub     dx, [si+2] */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+2)),16,0,1));
L_12586: ; /* IDA 0x12586; neg     dx */
    m->dx=(uint16_t)(ec_math(m,0,m->dx,16,0,1));
L_12588: ; /* IDA 0x12588; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_1258A: ; /* IDA 0x1258A; mov     bh, ch */
    gw_bh(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_1258C: ; /* IDA 0x1258C; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1258E: ; /* IDA 0x1258E; sub     bx, [si+4] */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->si+4)),16,0,1));
L_12591: ; /* IDA 0x12591; neg     bx */
    m->bx=(uint16_t)(ec_math(m,0,m->bx,16,0,1));
L_12593: ; /* IDA 0x12593; mov     ah, 5 */
    vg_ah(m,(uint8_t)(5));
L_12595: ; /* IDA 0x12595; call    loc_1D51F */
    ov_call(m,54559,9624,h);
L_12598: ; /* IDA 0x12598; add     si, 10h */
    m->si=(uint16_t)(ec_math(m,m->si,16,16,0,0));
L_1259B: ; /* IDA 0x1259B; dec     bp */
    m->bp=(uint16_t)(st_dec(m,m->bp,16));
L_1259C: ; /* IDA 0x1259C; jnz     short loc_12541 */
    if(!(m->flags&0x40)) goto L_12541;
L_1259E: ; /* IDA 0x1259E; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_125A0: ; /* IDA 0x125A0; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_125A2: ; /* IDA 0x125A2; retn */
    return;
}
/* Original sub_12AF4; IDA 0x12AF4. */
static void ov_body_12AF4(KiMachine16 *m,const KiEconomyHooks *h) {
L_12AF4: ; /* IDA 0x12AF4; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_12AF9: ; /* IDA 0x12AF9; mov     es, cs:word_1987A */
    m->es=(uint16_t)(dl_read16(m,m->cs,39034));
L_12AFE: ; /* IDA 0x12AFE; mov     si, 2240h */
    m->si=(uint16_t)(8768);
L_12B01: ; /* IDA 0x12B01; mov     bp, 7Fh */
    m->bp=(uint16_t)(127);
L_12B04: ; /* IDA 0x12B04; test    byte ptr [si], 20h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(32))&255;ec_logic(m,result,8);}
L_12B07: ; /* IDA 0x12B07; jz      short loc_12B0C */
    if(m->flags&0x40) goto L_12B0C;
L_12B09: ; /* IDA 0x12B09; call    sub_12B3C */
    ov_call(m,11068,11020,h);
L_12B0C: ; /* IDA 0x12B0C; add     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
L_12B0F: ; /* IDA 0x12B0F; dec     bp */
    m->bp=(uint16_t)(st_dec(m,m->bp,16));
L_12B10: ; /* IDA 0x12B10; jnz     short loc_12B04 */
    if(!(m->flags&0x40)) goto L_12B04;
L_12B12: ; /* IDA 0x12B12; mov     si, 2240h */
    m->si=(uint16_t)(8768);
L_12B15: ; /* IDA 0x12B15; mov     cx, 7Fh */
    m->cx=(uint16_t)(127);
L_12B18: ; /* IDA 0x12B18; cmp     byte ptr [si], 0C0h */
#if KI_OVERLAY_MUTATION == 1
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),191,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),192,8);
#endif
L_12B1B: ; /* IDA 0x12B1B; jb      short loc_12B20 */
    if(m->flags&1) goto L_12B20;
L_12B1D: ; /* IDA 0x12B1D; call    sub_12B2A */
    ov_call(m,11050,11040,h);
L_12B20: ; /* IDA 0x12B20; add     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
L_12B23: ; /* IDA 0x12B23; loop    loc_12B18 */
    if(--m->cx) goto L_12B18;
L_12B25: ; /* IDA 0x12B25; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_12B27: ; /* IDA 0x12B27; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_12B29: ; /* IDA 0x12B29; retn */
    return;
}
/* Original sub_12B2A; IDA 0x12B2A. */
static void ov_body_12B2A(KiMachine16 *m,const KiEconomyHooks *h) {
L_12B2A: ; /* IDA 0x12B2A; mov     dx, [si+10h] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+16)));
L_12B2D: ; /* IDA 0x12B2D; mov     bl, [si+12h] */
#if KI_OVERLAY_MUTATION == 2
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+19))));
#else
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+18))));
#endif
L_12B30: ; /* IDA 0x12B30; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_12B32: ; /* IDA 0x12B32; mov     al, [si+9] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+9))));
L_12B35: ; /* IDA 0x12B35; add     al, [si+8] */
#if KI_OVERLAY_MUTATION == 3
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+7)),8,0,0)));
#else
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+8)),8,0,0)));
#endif
L_12B38: ; /* IDA 0x12B38; call    sub_1D4C7 */
    ov_call(m,54471,11067,h);
L_12B3B: ; /* IDA 0x12B3B; retn */
    return;
}
/* Original sub_12B3C; IDA 0x12B3C. */
static void ov_body_12B3C(KiMachine16 *m,const KiEconomyHooks *h) {
L_12B3C: ; /* IDA 0x12B3C; or      byte ptr [si], 10h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(16))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_12B3F: ; /* IDA 0x12B3F; mov     bl, [si+3] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+3))));
L_12B42: ; /* IDA 0x12B42; and     bl, 3 */
    {uint16_t result=(((uint8_t)m->bx)&(3))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_12B45: ; /* IDA 0x12B45; mov     al, [si+21h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+33))));
L_12B48: ; /* IDA 0x12B48; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_12B4A: ; /* IDA 0x12B4A; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_12B4C: ; /* IDA 0x12B4C; add     bl, al */
#if KI_OVERLAY_MUTATION == 4
    ; /* lost formation pattern contribution */
#else
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,(uint8_t)m->ax,8,0,0)));
#endif
L_12B4E: ; /* IDA 0x12B4E; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_12B50: ; /* IDA 0x12B50; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_12B52: ; /* IDA 0x12B52; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_12B54: ; /* IDA 0x12B54; mov     cl, es:[bx] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->bx))));
L_12B57: ; /* IDA 0x12B57; mov     al, es:[bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->bx+1))));
L_12B5B: ; /* IDA 0x12B5B; mov     di, es:[bx+2] */
    m->di=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+2)));
L_12B5F: ; /* IDA 0x12B5F; add     di, 100h */
    m->di=(uint16_t)(ec_math(m,m->di,256,16,0,0));
L_12B63: ; /* IDA 0x12B63; mov     dl, cl */
    vg_dl(m,(uint8_t)((uint8_t)m->cx));
L_12B65: ; /* IDA 0x12B65; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_12B67: ; /* IDA 0x12B67; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_12B69: ; /* IDA 0x12B69; shr     bl, 1 */
    gw_bl(m,(uint8_t)(dl_shift(m,(uint8_t)m->bx,1,8,5)));
L_12B6B: ; /* IDA 0x12B6B; xor     dh, dh */
    {uint16_t result=(((uint8_t)(m->dx>>8))^((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);vg_dh(m,(uint8_t)(result));}
L_12B6D: ; /* IDA 0x12B6D; mov     bh, dh */
    gw_bh(m,(uint8_t)((uint8_t)(m->dx>>8)));
L_12B6F: ; /* IDA 0x12B6F; sub     dx, [si+10h] */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+16)),16,0,1));
L_12B72: ; /* IDA 0x12B72; sub     bx, [si+12h] */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->si+18)),16,0,1));
L_12B75: ; /* IDA 0x12B75; neg     dx */
    m->dx=(uint16_t)(ec_math(m,0,m->dx,16,0,1));
L_12B77: ; /* IDA 0x12B77; neg     bx */
    m->bx=(uint16_t)(ec_math(m,0,m->bx,16,0,1));
L_12B79: ; /* IDA 0x12B79; xor     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_12B7B: ; /* IDA 0x12B7B; mov     ah, 3 */
    vg_ah(m,(uint8_t)(3));
L_12B7D: ; /* IDA 0x12B7D; call    loc_1D51F */
    ov_call(m,54559,11136,h);
L_12B80: ; /* IDA 0x12B80; test    cs:byte_198A6, 4 */
    {uint16_t result=((dl_read8(m,m->cs,39078))&(4))&255;ec_logic(m,result,8);}
L_12B86: ; /* IDA 0x12B86; jz      short locret_12BA7 */
    if(m->flags&0x40) goto L_12BA7;
L_12B88: ; /* IDA 0x12B88; mov     dx, [si+10h] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+16)));
L_12B8B: ; /* IDA 0x12B8B; mov     bx, [si+12h] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+18)));
L_12B8E: ; /* IDA 0x12B8E; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_12B90: ; /* IDA 0x12B90; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_12B92: ; /* IDA 0x12B92; add     dx, 1B6h */
    m->dx=(uint16_t)(ec_math(m,m->dx,438,16,0,0));
L_12B96: ; /* IDA 0x12B96; add     bx, 26h ; '&' */
    m->bx=(uint16_t)(ec_math(m,m->bx,38,16,0,0));
L_12B99: ; /* IDA 0x12B99; mov     ah, 3Fh ; '?' */
    vg_ah(m,(uint8_t)(63));
L_12B9B: ; /* IDA 0x12B9B; mov     al, [si+3] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+3))));
L_12B9E: ; /* IDA 0x12B9E; test    al, 1 */
    {uint16_t result=(((uint8_t)m->ax)&(1))&255;ec_logic(m,result,8);}
L_12BA0: ; /* IDA 0x12BA0; jz      short loc_12BA4 */
    if(m->flags&0x40) goto L_12BA4;
L_12BA2: ; /* IDA 0x12BA2; mov     ah, 0FAh */
    vg_ah(m,(uint8_t)(250));
L_12BA4: ; /* IDA 0x12BA4; call    sub_15D19 */
    ov_call(m,23833,11175,h);
L_12BA7: ; /* IDA 0x12BA7; retn */
    return;
}
/* Original sub_15D19; IDA 0x15D19. */
static void ov_body_15D19(KiMachine16 *m,const KiEconomyHooks *h) {
L_15D19: ; /* IDA 0x15D19; push    ds */
    ec_push(m,m->ds);
L_15D1A: ; /* IDA 0x15D1A; push    cx */
    ec_push(m,m->cx);
L_15D1B: ; /* IDA 0x15D1B; mov     cx, 0A0C8h */
    m->cx=(uint16_t)(41160);
L_15D1E: ; /* IDA 0x15D1E; mov     ds, cx */
    m->ds=(uint16_t)(m->cx);
L_15D20: ; /* IDA 0x15D20; mov     cl, dl */
    vg_cl(m,(uint8_t)((uint8_t)m->dx));
L_15D22: ; /* IDA 0x15D22; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15D24: ; /* IDA 0x15D24; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15D26: ; /* IDA 0x15D26; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_15D28: ; /* IDA 0x15D28; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_15D2A: ; /* IDA 0x15D2A; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_15D2C: ; /* IDA 0x15D2C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_15D2E: ; /* IDA 0x15D2E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_15D30: ; /* IDA 0x15D30; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_15D32: ; /* IDA 0x15D32; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_15D34: ; /* IDA 0x15D34; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_15D36: ; /* IDA 0x15D36; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_15D38: ; /* IDA 0x15D38; and     cl, 7 */
    {uint16_t result=(((uint8_t)m->cx)&(7))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_15D3B: ; /* IDA 0x15D3B; call    sub_10CAC */
    ov_call(m,3244,23870,h);
L_15D3E: ; /* IDA 0x15D3E; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_15D41: ; /* IDA 0x15D41; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_15D43: ; /* IDA 0x15D43; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_15D44: ; /* IDA 0x15D44; push    ax */
    ec_push(m,m->ax);
L_15D45: ; /* IDA 0x15D45; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_15D47: ; /* IDA 0x15D47; mov     ah, 60h ; '`' */
#if KI_OVERLAY_MUTATION == 12
    vg_ah(m,(uint8_t)(64));
#else
    vg_ah(m,(uint8_t)(96));
#endif
L_15D49: ; /* IDA 0x15D49; shr     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,5));
L_15D4B: ; /* IDA 0x15D4B; mov     ch, [bx+50h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+80))));
L_15D4E: ; /* IDA 0x15D4E; mov     [bx+50h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+80),(uint8_t)(m->ax>>8));
L_15D51: ; /* IDA 0x15D51; mov     ch, [bx+51h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+81))));
L_15D54: ; /* IDA 0x15D54; mov     [bx+51h], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+81),(uint8_t)m->ax);
L_15D57: ; /* IDA 0x15D57; mov     ch, [bx+0A0h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+160))));
L_15D5B: ; /* IDA 0x15D5B; mov     [bx+0A0h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+160),(uint8_t)(m->ax>>8));
L_15D5F: ; /* IDA 0x15D5F; mov     ch, [bx+0A1h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+161))));
L_15D63: ; /* IDA 0x15D63; mov     [bx+0A1h], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+161),(uint8_t)m->ax);
L_15D67: ; /* IDA 0x15D67; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_15D68: ; /* IDA 0x15D68; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_15D6B: ; /* IDA 0x15D6B; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_15D6D: ; /* IDA 0x15D6D; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_15D6F: ; /* IDA 0x15D6F; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_15D71: ; /* IDA 0x15D71; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_15D73: ; /* IDA 0x15D73; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_15D75: ; /* IDA 0x15D75; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_15D76: ; /* IDA 0x15D76; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_15D78: ; /* IDA 0x15D78; mov     ah, 0F0h */
    vg_ah(m,(uint8_t)(240));
L_15D7A: ; /* IDA 0x15D7A; shr     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,5));
L_15D7C: ; /* IDA 0x15D7C; mov     ch, [bx] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx))));
L_15D7E: ; /* IDA 0x15D7E; mov     [bx], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx),(uint8_t)(m->ax>>8));
L_15D80: ; /* IDA 0x15D80; mov     ch, [bx+1] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_15D83: ; /* IDA 0x15D83; mov     [bx+1], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+1),(uint8_t)m->ax);
L_15D86: ; /* IDA 0x15D86; mov     ch, [bx+0F0h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+240))));
L_15D8A: ; /* IDA 0x15D8A; mov     [bx+0F0h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+240),(uint8_t)(m->ax>>8));
L_15D8E: ; /* IDA 0x15D8E; mov     ch, [bx+0F1h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+241))));
L_15D92: ; /* IDA 0x15D92; mov     [bx+0F1h], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+241),(uint8_t)m->ax);
L_15D96: ; /* IDA 0x15D96; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_15D98: ; /* IDA 0x15D98; mov     ah, 90h */
    vg_ah(m,(uint8_t)(144));
L_15D9A: ; /* IDA 0x15D9A; shr     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,5));
L_15D9C: ; /* IDA 0x15D9C; mov     ch, [bx+50h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+80))));
L_15D9F: ; /* IDA 0x15D9F; mov     [bx+50h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+80),(uint8_t)(m->ax>>8));
L_15DA2: ; /* IDA 0x15DA2; mov     ch, [bx+51h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+81))));
L_15DA5: ; /* IDA 0x15DA5; mov     [bx+51h], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+81),(uint8_t)m->ax);
L_15DA8: ; /* IDA 0x15DA8; mov     ch, [bx+0A0h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+160))));
L_15DAC: ; /* IDA 0x15DAC; mov     [bx+0A0h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+160),(uint8_t)(m->ax>>8));
L_15DB0: ; /* IDA 0x15DB0; mov     ch, [bx+0A1h] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+161))));
L_15DB4: ; /* IDA 0x15DB4; mov     [bx+0A1h], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+161),(uint8_t)m->ax);
L_15DB8: ; /* IDA 0x15DB8; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_15DB9: ; /* IDA 0x15DB9; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_15DBA: ; /* IDA 0x15DBA; retn */
    return;
}
/* Original loc_1D51F; IDA 0x1D51F. */
static void ov_body_1D51F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D51F: ; /* IDA 0x1D51F; push    si */
    ec_push(m,m->si);
L_1D520: ; /* IDA 0x1D520; push    bp */
    ec_push(m,m->bp);
L_1D521: ; /* IDA 0x1D521; mov     cs:byte_1D5AB, ah */
    dl_write8(m,m->cs,54699,(uint8_t)(m->ax>>8));
L_1D526: ; /* IDA 0x1D526; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1D528: ; /* IDA 0x1D528; mov     bp, cx */
    m->bp=(uint16_t)(m->cx);
L_1D52A: ; /* IDA 0x1D52A; mov     si, cs:word_1D854 */
    m->si=(uint16_t)(dl_read16(m,m->cs,55380));
L_1D52F: ; /* IDA 0x1D52F; sub     dx, si */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->si,16,0,1));
L_1D531: ; /* IDA 0x1D531; jge     short loc_1D53F */
    if(!ec_less(m)) goto L_1D53F;
L_1D533: ; /* IDA 0x1D533; neg     dx */
    m->dx=(uint16_t)(ec_math(m,0,m->dx,16,0,1));
L_1D535: ; /* IDA 0x1D535; sub     cx, dx */
    m->cx=(uint16_t)(ec_math(m,m->cx,m->dx,16,0,1));
L_1D537: ; /* IDA 0x1D537; jle     short loc_1D562 */
    if(!ec_greater(m)) goto L_1D562;
L_1D539: ; /* IDA 0x1D539; add     di, dx */
    m->di=(uint16_t)(ec_math(m,m->di,m->dx,16,0,0));
L_1D53B: ; /* IDA 0x1D53B; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_1D53D: ; /* IDA 0x1D53D; jmp     short loc_1D54C */
    goto L_1D54C;
L_1D53F: ; /* IDA 0x1D53F; mov     si, 28h ; '(' */
    m->si=(uint16_t)(40);
L_1D542: ; /* IDA 0x1D542; sub     si, dx */
    m->si=(uint16_t)(ec_math(m,m->si,m->dx,16,0,1));
L_1D544: ; /* IDA 0x1D544; jle     short loc_1D562 */
    if(!ec_greater(m)) goto L_1D562;
L_1D546: ; /* IDA 0x1D546; cmp     cx, si */
    ec_cmp(m,m->cx,m->si,16);
L_1D548: ; /* IDA 0x1D548; jle     short loc_1D54C */
    if(!ec_greater(m)) goto L_1D54C;
L_1D54A: ; /* IDA 0x1D54A; mov     cx, si */
    m->cx=(uint16_t)(m->si);
L_1D54C: ; /* IDA 0x1D54C; mov     si, cs:word_1D856 */
    m->si=(uint16_t)(dl_read16(m,m->cs,55382));
L_1D551: ; /* IDA 0x1D551; sub     bx, si */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->si,16,0,1));
L_1D553: ; /* IDA 0x1D553; jge     short loc_1D564 */
    if(!ec_less(m)) goto L_1D564;
L_1D555: ; /* IDA 0x1D555; neg     bx */
    m->bx=(uint16_t)(ec_math(m,0,m->bx,16,0,1));
L_1D557: ; /* IDA 0x1D557; sub     ax, bx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->bx,16,0,1));
L_1D559: ; /* IDA 0x1D559; jle     short loc_1D562 */
    if(!ec_greater(m)) goto L_1D562;
L_1D55B: ; /* IDA 0x1D55B; add     di, bp */
    m->di=(uint16_t)(ec_math(m,m->di,m->bp,16,0,0));
L_1D55D: ; /* IDA 0x1D55D; dec     bx */
    m->bx=(uint16_t)(st_dec(m,m->bx,16));
L_1D55E: ; /* IDA 0x1D55E; jnz     short loc_1D55B */
    if(!(m->flags&0x40)) goto L_1D55B;
L_1D560: ; /* IDA 0x1D560; jmp     short loc_1D571 */
    goto L_1D571;
L_1D562: ; /* IDA 0x1D562; jmp     short loc_1D5D1 */
    goto L_1D5D1;
L_1D564: ; /* IDA 0x1D564; mov     si, 17h */
    m->si=(uint16_t)(23);
L_1D567: ; /* IDA 0x1D567; sub     si, bx */
    m->si=(uint16_t)(ec_math(m,m->si,m->bx,16,0,1));
L_1D569: ; /* IDA 0x1D569; jle     short loc_1D562 */
    if(!ec_greater(m)) goto L_1D562;
L_1D56B: ; /* IDA 0x1D56B; cmp     ax, si */
    ec_cmp(m,m->ax,m->si,16);
L_1D56D: ; /* IDA 0x1D56D; jle     short loc_1D571 */
    if(!ec_greater(m)) goto L_1D571;
L_1D56F: ; /* IDA 0x1D56F; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_1D571: ; /* IDA 0x1D571; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D573: ; /* IDA 0x1D573; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D575: ; /* IDA 0x1D575; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D577: ; /* IDA 0x1D577; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_1D579: ; /* IDA 0x1D579; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D57B: ; /* IDA 0x1D57B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D57D: ; /* IDA 0x1D57D; add     si, bx */
    m->si=(uint16_t)(ec_math(m,m->si,m->bx,16,0,0));
L_1D57F: ; /* IDA 0x1D57F; add     si, dx */
    m->si=(uint16_t)(ec_math(m,m->si,m->dx,16,0,0));
L_1D581: ; /* IDA 0x1D581; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_1D583: ; /* IDA 0x1D583; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_1D585: ; /* IDA 0x1D585; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_1D587: ; /* IDA 0x1D587; push    ds */
    ec_push(m,m->ds);
L_1D588: ; /* IDA 0x1D588; mov     ds, cs:word_1D84E */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55374));
L_1D58D: ; /* IDA 0x1D58D; sub     bp, cx */
    m->bp=(uint16_t)(ec_math(m,m->bp,m->cx,16,0,1));
L_1D58F: ; /* IDA 0x1D58F; mov     bh, ch */
    gw_bh(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_1D591: ; /* IDA 0x1D591; mov     dl, al */
    vg_dl(m,(uint8_t)((uint8_t)m->ax));
L_1D593: ; /* IDA 0x1D593; mov     dh, cl */
    vg_dh(m,(uint8_t)((uint8_t)m->cx));
L_1D595: ; /* IDA 0x1D595; push    si */
    ec_push(m,m->si);
L_1D596: ; /* IDA 0x1D596; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_1D599: ; /* IDA 0x1D599; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1D59A: ; /* IDA 0x1D59A; cmp     al, 0FFh */
#if KI_OVERLAY_MUTATION == 9
    ec_cmp(m,(uint8_t)m->ax,254,8);
#else
    ec_cmp(m,(uint8_t)m->ax,255,8);
#endif
L_1D59C: ; /* IDA 0x1D59C; jz      short loc_1D5BE */
    if(m->flags&0x40) goto L_1D5BE;
L_1D59E: ; /* IDA 0x1D59E; mov     bx, [si] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_1D5A0: ; /* IDA 0x1D5A0; test    bl, 10h */
#if KI_OVERLAY_MUTATION == 10
    {uint16_t result=(((uint8_t)m->bx)&(0))&255;ec_logic(m,result,8);}
#else
    {uint16_t result=(((uint8_t)m->bx)&(16))&255;ec_logic(m,result,8);}
#endif
L_1D5A3: ; /* IDA 0x1D5A3; jnz     short loc_1D5BE */
    if(!(m->flags&0x40)) goto L_1D5BE;
L_1D5A5: ; /* IDA 0x1D5A5; mov     bl, bh */
    gw_bl(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_1D5A7: ; /* IDA 0x1D5A7; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1D5A9: ; /* IDA 0x1D5A9; cmp     bl, 5 */
#if KI_OVERLAY_MUTATION == 8
    ec_cmp(m,(uint8_t)m->bx,4,8);
#else
    ec_cmp(m,(uint8_t)m->bx,dl_read8(m,m->cs,0xd5ab),8);
#endif
L_1D5AC: ; /* IDA 0x1D5AC; jnb     short loc_1D5BE */
    if(!(m->flags&1)) goto L_1D5BE;
L_1D5AE: ; /* IDA 0x1D5AE; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_1D5B0: ; /* IDA 0x1D5B0; xchg    al, [bx+si+3] */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->ds,(uint16_t)(m->bx+m->si+3));vg_al(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->bx+m->si+3),a);}
L_1D5B3: ; /* IDA 0x1D5B3; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1D5B4: ; /* IDA 0x1D5B4; mov     [si+1], bl */
    dl_write8(m,m->ds,(uint16_t)(m->si+1),(uint8_t)m->bx);
L_1D5B7: ; /* IDA 0x1D5B7; cmp     al, ah */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
L_1D5B9: ; /* IDA 0x1D5B9; jz      short loc_1D5BE */
    if(m->flags&0x40) goto L_1D5BE;
L_1D5BB: ; /* IDA 0x1D5BB; or      byte ptr [si], 20h */
#if KI_OVERLAY_MUTATION == 11
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(0))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(32))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#endif
L_1D5BE: ; /* IDA 0x1D5BE; add     si, 8 */
    m->si=(uint16_t)(ec_math(m,m->si,8,16,0,0));
L_1D5C1: ; /* IDA 0x1D5C1; loop    loc_1D596 */
    if(--m->cx) goto L_1D596;
L_1D5C3: ; /* IDA 0x1D5C3; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1D5C4: ; /* IDA 0x1D5C4; add     si, 140h */
    m->si=(uint16_t)(ec_math(m,m->si,320,16,0,0));
L_1D5C8: ; /* IDA 0x1D5C8; add     di, bp */
    m->di=(uint16_t)(ec_math(m,m->di,m->bp,16,0,0));
L_1D5CA: ; /* IDA 0x1D5CA; mov     cl, dh */
    vg_cl(m,(uint8_t)((uint8_t)(m->dx>>8)));
L_1D5CC: ; /* IDA 0x1D5CC; dec     dl */
    vg_dl(m,(uint8_t)(st_dec(m,(uint8_t)m->dx,8)));
L_1D5CE: ; /* IDA 0x1D5CE; jnz     short loc_1D595 */
    if(!(m->flags&0x40)) goto L_1D595;
L_1D5D0: ; /* IDA 0x1D5D0; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1D5D1: ; /* IDA 0x1D5D1; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1D5D2: ; /* IDA 0x1D5D2; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1D5D3: ; /* IDA 0x1D5D3; retn */
    return;
}
int ki_overlay_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x1cc9:ov_body_11CC9(m,h);break;
case 0x2533:ov_body_12533(m,h);break;
case 0x2af4:ov_body_12AF4(m,h);break;
case 0x2b2a:ov_body_12B2A(m,h);break;
case 0x2b3c:ov_body_12B3C(m,h);break;
case 0x5d19:ov_body_15D19(m,h);break;
case 0xd51f:ov_body_1D51F(m,h);break;
default:return 0;}return 1;}
