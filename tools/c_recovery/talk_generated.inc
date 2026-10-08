/* Fixed IDA native C TALK control; original locators; spec/230, re/110. */
/* Original sub_106F9; IDA 0x106F9. */
static void tx_body_106F9(KiMachine16 *m,const KiEconomyHooks *h) {
L_106F9: ; /* IDA 0x106F9; call    sub_1084A */
    tx_call(m,2122,1788,h);
L_106FC: ; /* IDA 0x106FC; retn */
    return;
}
/* Original sub_1075B; IDA 0x1075B. */
static void tx_body_1075B(KiMachine16 *m,const KiEconomyHooks *h) {
L_1075B: ; /* IDA 0x1075B; push    ds */
    ec_push(m,m->ds);
L_1075C: ; /* IDA 0x1075C; push    si */
    ec_push(m,m->si);
L_1075D: ; /* IDA 0x1075D; push    ax */
    ec_push(m,m->ax);
L_1075E: ; /* IDA 0x1075E; mov     ds, cs:word_10D38 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3384));
L_10763: ; /* IDA 0x10763; cmp     cx, 196h */
    ec_cmp(m,m->cx,406,16);
L_10767: ; /* IDA 0x10767; jb      short loc_1077C */
    if(m->flags&1) goto L_1077C;
L_10769: ; /* IDA 0x10769; sub     cx, 196h */
    m->cx=(uint16_t)(ec_math(m,m->cx,406,16,0,1));
L_1076D: ; /* IDA 0x1076D; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_1076F: ; /* IDA 0x1076F; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_10771: ; /* IDA 0x10771; shl     cx, 1 */
#if KI_TALK_MUTATION == 1
    ; /* wrong variant stride */
#else
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
#endif
L_10773: ; /* IDA 0x10773; add     cl, ah */
#if KI_TALK_MUTATION == 10
    ; /* lost variant byte addition */
#else
    vg_cl(m,(uint8_t)(ec_math(m,(uint8_t)m->cx,(uint8_t)(m->ax>>8),8,0,0)));
#endif
L_10775: ; /* IDA 0x10775; adc     ch, 0 */
    vg_ch(m,(uint8_t)(ec_math(m,(uint8_t)(m->cx>>8),0,8,m->flags&1,0)));
L_10778: ; /* IDA 0x10778; add     cx, 196h */
    m->cx=(uint16_t)(ec_math(m,m->cx,406,16,0,0));
L_1077C: ; /* IDA 0x1077C; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_1077E: ; /* IDA 0x1077E; mov     si, cx */
    m->si=(uint16_t)(m->cx);
L_10780: ; /* IDA 0x10780; mov     si, [si] */
    m->si=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_10782: ; /* IDA 0x10782; mov     cx, 510h */
    m->cx=(uint16_t)(1296);
L_10785: ; /* IDA 0x10785; call    sub_10BCD */
    tx_call(m,3021,1928,h);
L_10788: ; /* IDA 0x10788; push    bx */
    ec_push(m,m->bx);
L_10789: ; /* IDA 0x10789; mov     bh, bl */
    gw_bh(m,(uint8_t)((uint8_t)m->bx));
L_1078B: ; /* IDA 0x1078B; shl     bh, 1 */
    gw_bh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->bx>>8),1,8,4)));
L_1078D: ; /* IDA 0x1078D; shl     bh, 1 */
    gw_bh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->bx>>8),1,8,4)));
L_1078F: ; /* IDA 0x1078F; add     bh, bl */
    gw_bh(m,(uint8_t)(ec_math(m,(uint8_t)(m->bx>>8),(uint8_t)m->bx,8,0,0)));
L_10791: ; /* IDA 0x10791; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_10793: ; /* IDA 0x10793; add     bx, 280h */
    m->bx=(uint16_t)(ec_math(m,m->bx,640,16,0,0));
L_10797: ; /* IDA 0x10797; push    bx */
    ec_push(m,m->bx);
L_10798: ; /* IDA 0x10798; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_1079A: ; /* IDA 0x1079A; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_1079C: ; /* IDA 0x1079C; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1079D: ; /* IDA 0x1079D; call    sub_107D2 */
    tx_call(m,2002,1952,h);
L_107A0: ; /* IDA 0x107A0; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_107A1: ; /* IDA 0x107A1; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_107A3: ; /* IDA 0x107A3; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_107A5: ; /* IDA 0x107A5; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_107A7: ; /* IDA 0x107A7; add     dx, 48h ; 'H' */
    m->dx=(uint16_t)(ec_math(m,m->dx,72,16,0,0));
L_107AA: ; /* IDA 0x107AA; mov     cx, 40B0h */
    m->cx=(uint16_t)(16560);
L_107AD: ; /* IDA 0x107AD; mov     ah, 0 */
    vg_ah(m,(uint8_t)(0));
L_107AF: ; /* IDA 0x107AF; call    sub_10AD9 */
    tx_call(m,2777,1970,h);
L_107B2: ; /* IDA 0x107B2; add     dx, 8 */
    m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,0));
L_107B5: ; /* IDA 0x107B5; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_107B6: ; /* IDA 0x107B6; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_107B8: ; /* IDA 0x107B8; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_107BA: ; /* IDA 0x107BA; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_107BC: ; /* IDA 0x107BC; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_107BE: ; /* IDA 0x107BE; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_107C1: ; /* IDA 0x107C1; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_107C2: ; /* IDA 0x107C2; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_107C4: ; /* IDA 0x107C4; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_107C6: ; /* IDA 0x107C6; push    ax */
    ec_push(m,m->ax);
L_107C7: ; /* IDA 0x107C7; push    bx */
    ec_push(m,m->bx);
L_107C8: ; /* IDA 0x107C8; push    dx */
    ec_push(m,m->dx);
L_107C9: ; /* IDA 0x107C9; call    sub_106F9 */
    tx_call(m,1785,1996,h);
L_107CC: ; /* IDA 0x107CC; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_107CD: ; /* IDA 0x107CD; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_107CE: ; /* IDA 0x107CE; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_107CF: ; /* IDA 0x107CF; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_107D0: ; /* IDA 0x107D0; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_107D1: ; /* IDA 0x107D1; retn */
    return;
}
/* Original sub_107D2; IDA 0x107D2. */
static void tx_body_107D2(KiMachine16 *m,const KiEconomyHooks *h) {
L_107D2: ; /* IDA 0x107D2; push    ds */
    ec_push(m,m->ds);
L_107D3: ; /* IDA 0x107D3; push    ax */
    ec_push(m,m->ax);
L_107D4: ; /* IDA 0x107D4; push    bx */
    ec_push(m,m->bx);
L_107D5: ; /* IDA 0x107D5; push    cx */
    ec_push(m,m->cx);
L_107D6: ; /* IDA 0x107D6; push    dx */
    ec_push(m,m->dx);
L_107D7: ; /* IDA 0x107D7; push    si */
    ec_push(m,m->si);
L_107D8: ; /* IDA 0x107D8; push    di */
    ec_push(m,m->di);
L_107D9: ; /* IDA 0x107D9; push    bx */
    ec_push(m,m->bx);
L_107DA: ; /* IDA 0x107DA; mov     bx, cs */
    m->bx=(uint16_t)(m->cs);
L_107DC: ; /* IDA 0x107DC; mov     ds, bx */
    m->ds=(uint16_t)(m->bx);
L_107DE: ; /* IDA 0x107DE; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_107E0: ; /* IDA 0x107E0; cmp     al, [bx+846h] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+2118)),8);
L_107E4: ; /* IDA 0x107E4; jz      short loc_10828 */
    if(m->flags&0x40) goto L_10828;
L_107E6: ; /* IDA 0x107E6; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_107E7: ; /* IDA 0x107E7; cmp     bl, 4 */
    ec_cmp(m,(uint8_t)m->bx,4,8);
L_107EA: ; /* IDA 0x107EA; jb      short loc_107E0 */
    if(m->flags&1) goto L_107E0;
L_107EC: ; /* IDA 0x107EC; mov     bl, byte_10845 */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,2117)));
L_107F0: ; /* IDA 0x107F0; push    bx */
    ec_push(m,m->bx);
L_107F1: ; /* IDA 0x107F1; mov     [bx+846h], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+2118),(uint8_t)m->ax);
L_107F5: ; /* IDA 0x107F5; xchg    bh, bl */
    {uint16_t a=(uint8_t)(m->bx>>8),b=(uint8_t)m->bx;gw_bh(m,(uint8_t)(b));gw_bl(m,(uint8_t)(a));}
L_107F7: ; /* IDA 0x107F7; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_107F9: ; /* IDA 0x107F9; inc     bh */
    gw_bh(m,(uint8_t)(ec_inc(m,(uint8_t)(m->bx>>8),8)));
L_107FB: ; /* IDA 0x107FB; and     bh, 3 */
#if KI_TALK_MUTATION == 6
    {uint16_t result=(((uint8_t)(m->bx>>8))&(1))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
#else
    {uint16_t result=(((uint8_t)(m->bx>>8))&(3))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
#endif
L_107FE: ; /* IDA 0x107FE; mov     byte_10845, bh */
    dl_write8(m,m->ds,2117,(uint8_t)(m->bx>>8));
L_10802: ; /* IDA 0x10802; mov     dx, 0D79h */
    m->dx=(uint16_t)(3449);
L_10805: ; /* IDA 0x10805; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_10807: ; /* IDA 0x10807; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_10809: ; /* IDA 0x10809; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_1080B: ; /* IDA 0x1080B; mov     bx, word_10D40 */
    m->bx=(uint16_t)(dl_read16(m,m->ds,3392));
L_1080F: ; /* IDA 0x1080F; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_10811: ; /* IDA 0x10811; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_10813: ; /* IDA 0x10813; mov     al, cl */
    vg_al(m,(uint8_t)((uint8_t)m->cx));
L_10815: ; /* IDA 0x10815; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_10817: ; /* IDA 0x10817; rcl     cx, 1 */
    m->cx=(uint16_t)(tx_rcl(m,m->cx,16));
L_10819: ; /* IDA 0x10819; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1081B: ; /* IDA 0x1081B; rcl     cx, 1 */
    m->cx=(uint16_t)(tx_rcl(m,m->cx,16));
L_1081D: ; /* IDA 0x1081D; shl     ax, 1 */
#if KI_TALK_MUTATION == 7
    ; /* wrong portrait offset stride */
#else
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
#endif
L_1081F: ; /* IDA 0x1081F; rcl     cx, 1 */
    m->cx=(uint16_t)(tx_rcl(m,m->cx,16));
L_10821: ; /* IDA 0x10821; mov     di, 800h */
#if KI_TALK_MUTATION == 8
    m->di=(uint16_t)(1024);
#else
    m->di=(uint16_t)(2048);
#endif
L_10824: ; /* IDA 0x10824; call    sub_1E38C */
    tx_call(m,58252,2087,h);
L_10827: ; /* IDA 0x10827; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10828: ; /* IDA 0x10828; xchg    bl, bh */
    {uint16_t a=(uint8_t)m->bx,b=(uint8_t)(m->bx>>8);gw_bl(m,(uint8_t)(b));gw_bh(m,(uint8_t)(a));}
L_1082A: ; /* IDA 0x1082A; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1082C: ; /* IDA 0x1082C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1082E: ; /* IDA 0x1082E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10830: ; /* IDA 0x10830; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_10832: ; /* IDA 0x10832; mov     ds, word_10D40 */
    m->ds=(uint16_t)(dl_read16(m,m->ds,3392));
L_10836: ; /* IDA 0x10836; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10837: ; /* IDA 0x10837; mov     ax, 4004h */
    m->ax=(uint16_t)(16388);
L_1083A: ; /* IDA 0x1083A; call    sub_1FA37 */
    tx_call(m,64055,2109,h);
L_1083D: ; /* IDA 0x1083D; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1083E: ; /* IDA 0x1083E; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1083F: ; /* IDA 0x1083F; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10840: ; /* IDA 0x10840; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10841: ; /* IDA 0x10841; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10842: ; /* IDA 0x10842; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10843: ; /* IDA 0x10843; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_10844: ; /* IDA 0x10844; retn */
    return;
}
/* Original sub_1084A; IDA 0x1084A. */
static void tx_body_1084A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1084A: ; /* IDA 0x1084A; push    es */
    ec_push(m,m->es);
L_1084B: ; /* IDA 0x1084B; push    cx */
    ec_push(m,m->cx);
L_1084C: ; /* IDA 0x1084C; push    bp */
    ec_push(m,m->bp);
L_1084D: ; /* IDA 0x1084D; mov     ch, ah */
    vg_ch(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1084F: ; /* IDA 0x1084F; push    dx */
    ec_push(m,m->dx);
L_10850: ; /* IDA 0x10850; mov     bp, 10h */
    m->bp=(uint16_t)(16);
L_10853: ; /* IDA 0x10853; mov     ax, [si] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_10855: ; /* IDA 0x10855; xchg    al, ah */
    {uint16_t a=(uint8_t)m->ax,b=(uint8_t)(m->ax>>8);vg_al(m,(uint8_t)(b));vg_ah(m,(uint8_t)(a));}
L_10857: ; /* IDA 0x10857; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_10858: ; /* IDA 0x10858; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_10859: ; /* IDA 0x10859; and     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))&((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1085B: ; /* IDA 0x1085B; jz      short loc_1087B */
    if(m->flags&0x40) goto L_1087B;
L_1085D: ; /* IDA 0x1085D; cmp     ah, 5Ch ; '\' */
    ec_cmp(m,(uint8_t)(m->ax>>8),92,8);
L_10860: ; /* IDA 0x10860; jz      short loc_1088A */
    if(m->flags&0x40) goto L_1088A;
L_10862: ; /* IDA 0x10862; call    sub_1F878 */
    tx_call(m,63608,2149,h);
L_10865: ; /* IDA 0x10865; jnb     short loc_1086B */
    if(!(m->flags&1)) goto L_1086B;
L_10867: ; /* IDA 0x10867; mov     bp, 8 */
    m->bp=(uint16_t)(8);
L_1086A: ; /* IDA 0x1086A; dec     si */
#if KI_TALK_MUTATION == 2
    ; /* lost half-width rewind */
#else
    m->si=(uint16_t)(st_dec(m,m->si,16));
#endif
L_1086B: ; /* IDA 0x1086B; push    cx */
    ec_push(m,m->cx);
L_1086C: ; /* IDA 0x1086C; xchg    ax, cx */
    {uint16_t a=m->ax,b=m->cx;m->ax=(uint16_t)(b);m->cx=(uint16_t)(a);}
L_1086D: ; /* IDA 0x1086D; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1086F: ; /* IDA 0x1086F; push    bx */
    ec_push(m,m->bx);
L_10870: ; /* IDA 0x10870; push    dx */
    ec_push(m,m->dx);
L_10871: ; /* IDA 0x10871; call    loc_1F75E */
    tx_call(m,63326,2164,h);
L_10874: ; /* IDA 0x10874; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10875: ; /* IDA 0x10875; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10876: ; /* IDA 0x10876; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10877: ; /* IDA 0x10877; add     dx, bp */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bp,16,0,0));
L_10879: ; /* IDA 0x10879; jmp     short loc_10850 */
    goto L_10850;
L_1087B: ; /* IDA 0x1087B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1087C: ; /* IDA 0x1087C; and     ax, ax */
    {uint16_t result=((m->ax)&(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_1087E: ; /* IDA 0x1087E; jz      short loc_10886 */
    if(m->flags&0x40) goto L_10886;
L_10880: ; /* IDA 0x10880; dec     si */
    m->si=(uint16_t)(st_dec(m,m->si,16));
L_10881: ; /* IDA 0x10881; add     bx, 10h */
#if KI_TALK_MUTATION == 3
    m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,0));
#else
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
#endif
L_10884: ; /* IDA 0x10884; jmp     short loc_1084F */
    goto L_1084F;
L_10886: ; /* IDA 0x10886; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_10887: ; /* IDA 0x10887; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10888: ; /* IDA 0x10888; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_10889: ; /* IDA 0x10889; retn */
    return;
L_1088A: ; /* IDA 0x1088A; sub     al, 31h ; '1' */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,49,8,0,1)));
L_1088C: ; /* IDA 0x1088C; push    ds */
    ec_push(m,m->ds);
L_1088D: ; /* IDA 0x1088D; push    cx */
    ec_push(m,m->cx);
L_1088E: ; /* IDA 0x1088E; push    si */
    ec_push(m,m->si);
L_1088F: ; /* IDA 0x1088F; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_10891: ; /* IDA 0x10891; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_10893: ; /* IDA 0x10893; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_10895: ; /* IDA 0x10895; mov     ax, cs:[si+8A4h] */
    m->ax=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->si+2212)));
L_1089A: ; /* IDA 0x1089A; call    ax */
    tx_call(m,m->ax,2204,h);
L_1089C: ; /* IDA 0x1089C; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1089D: ; /* IDA 0x1089D; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1089E: ; /* IDA 0x1089E; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1089F: ; /* IDA 0x1089F; add     dx, 30h ; '0' */
#if KI_TALK_MUTATION == 4
    m->dx=(uint16_t)(ec_math(m,m->dx,32,16,0,0));
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,48,16,0,0));
#endif
L_108A2: ; /* IDA 0x108A2; jmp     short loc_10850 */
    goto L_10850;
}
/* Original sub_18853; IDA 0x18853. */
static void tx_body_18853(KiMachine16 *m,const KiEconomyHooks *h) {
L_18853: ; /* IDA 0x18853; push    ax */
    ec_push(m,m->ax);
L_18854: ; /* IDA 0x18854; push    bx */
    ec_push(m,m->bx);
L_18855: ; /* IDA 0x18855; push    dx */
    ec_push(m,m->dx);
L_18856: ; /* IDA 0x18856; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_18858: ; /* IDA 0x18858; mov     dx, 0 */
    m->dx=(uint16_t)(0);
L_1885B: ; /* IDA 0x1885B; mov     bx, 12h */
    m->bx=(uint16_t)(18);
L_1885E: ; /* IDA 0x1885E; push    ax */
    ec_push(m,m->ax);
L_1885F: ; /* IDA 0x1885F; push    cx */
    ec_push(m,m->cx);
L_18860: ; /* IDA 0x18860; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_18862: ; /* IDA 0x18862; cmp     cx, 0FFFFh */
    ec_cmp(m,m->cx,65535,16);
L_18865: ; /* IDA 0x18865; jz      short loc_18869 */
    if(m->flags&0x40) goto L_18869;
L_18867: ; /* IDA 0x18867; mov     al, 1Eh */
    vg_al(m,(uint8_t)(30));
L_18869: ; /* IDA 0x18869; mov     cx, 510h */
    m->cx=(uint16_t)(1296);
L_1886C: ; /* IDA 0x1886C; call    sub_189A4 */
    tx_call(m,35236,34927,h);
L_1886F: ; /* IDA 0x1886F; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_18870: ; /* IDA 0x18870; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_18871: ; /* IDA 0x18871; cmp     cx, 0FFFFh */
    ec_cmp(m,m->cx,65535,16);
L_18874: ; /* IDA 0x18874; jz      short loc_1887F */
    if(m->flags&0x40) goto L_1887F;
L_18876: ; /* IDA 0x18876; mov     dx, 0 */
    m->dx=(uint16_t)(0);
L_18879: ; /* IDA 0x18879; mov     bx, 14h */
    m->bx=(uint16_t)(20);
L_1887C: ; /* IDA 0x1887C; call    sub_1075B */
    tx_call(m,1883,34943,h);
L_1887F: ; /* IDA 0x1887F; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_18880: ; /* IDA 0x18880; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_18881: ; /* IDA 0x18881; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_18882: ; /* IDA 0x18882; retn */
    return;
}
/* Original sub_189A4; IDA 0x189A4. */
static void tx_body_189A4(KiMachine16 *m,const KiEconomyHooks *h) {
L_189A4: ; /* IDA 0x189A4; push    ax */
    ec_push(m,m->ax);
L_189A5: ; /* IDA 0x189A5; push    bx */
    ec_push(m,m->bx);
L_189A6: ; /* IDA 0x189A6; push    cx */
    ec_push(m,m->cx);
L_189A7: ; /* IDA 0x189A7; push    dx */
    ec_push(m,m->dx);
L_189A8: ; /* IDA 0x189A8; push    si */
    ec_push(m,m->si);
L_189A9: ; /* IDA 0x189A9; push    di */
    ec_push(m,m->di);
L_189AA: ; /* IDA 0x189AA; push    bp */
    ec_push(m,m->bp);
L_189AB: ; /* IDA 0x189AB; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_189AD: ; /* IDA 0x189AD; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_189AF: ; /* IDA 0x189AF; mov     bp, cx */
    m->bp=(uint16_t)(m->cx);
L_189B1: ; /* IDA 0x189B1; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_189B3: ; /* IDA 0x189B3; push    ax */
    ec_push(m,m->ax);
L_189B4: ; /* IDA 0x189B4; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_189B6: ; /* IDA 0x189B6; jnz     short loc_189BA */
    if(!(m->flags&0x40)) goto L_189BA;
L_189B8: ; /* IDA 0x189B8; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_189BA: ; /* IDA 0x189BA; call    sub_1D5D4 */
    tx_call(m,54740,35261,h);
L_189BD: ; /* IDA 0x189BD; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_189BE: ; /* IDA 0x189BE; add     si, 2 */
    m->si=(uint16_t)(ec_math(m,m->si,2,16,0,0));
L_189C1: ; /* IDA 0x189C1; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_189C3: ; /* IDA 0x189C3; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_189C5: ; /* IDA 0x189C5; mov     cx, bp */
    m->cx=(uint16_t)(m->bp);
L_189C7: ; /* IDA 0x189C7; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_189C9: ; /* IDA 0x189C9; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_189CB: ; /* IDA 0x189CB; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_189CD: ; /* IDA 0x189CD; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_189CF: ; /* IDA 0x189CF; jz      short loc_189D3 */
    if(m->flags&0x40) goto L_189D3;
L_189D1: ; /* IDA 0x189D1; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_189D3: ; /* IDA 0x189D3; call    sub_189DE */
    tx_call(m,35294,35286,h);
L_189D6: ; /* IDA 0x189D6; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_189D7: ; /* IDA 0x189D7; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_189D8: ; /* IDA 0x189D8; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_189D9: ; /* IDA 0x189D9; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_189DA: ; /* IDA 0x189DA; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_189DB: ; /* IDA 0x189DB; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_189DC: ; /* IDA 0x189DC; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_189DD: ; /* IDA 0x189DD; retn */
    return;
}
/* Original sub_1E38C; IDA 0x1E38C. */
static void tx_body_1E38C(KiMachine16 *m,const KiEconomyHooks *h) {
L_1E38C: ; /* IDA 0x1E38C; push    ds */
    ec_push(m,m->ds);
L_1E38D: ; /* IDA 0x1E38D; push    ax */
    ec_push(m,m->ax);
L_1E38E: ; /* IDA 0x1E38E; push    bx */
    ec_push(m,m->bx);
L_1E38F: ; /* IDA 0x1E38F; push    cx */
    ec_push(m,m->cx);
L_1E390: ; /* IDA 0x1E390; push    dx */
    ec_push(m,m->dx);
L_1E391: ; /* IDA 0x1E391; push    si */
    ec_push(m,m->si);
L_1E392: ; /* IDA 0x1E392; push    di */
    ec_push(m,m->di);
L_1E393: ; /* IDA 0x1E393; call    sub_1F4DF */
    tx_call(m,62687,58262,h);
L_1E396: ; /* IDA 0x1E396; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1E397: ; /* IDA 0x1E397; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1E398: ; /* IDA 0x1E398; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1E399: ; /* IDA 0x1E399; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1E39A: ; /* IDA 0x1E39A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1E39B: ; /* IDA 0x1E39B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1E39C: ; /* IDA 0x1E39C; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1E39D: ; /* IDA 0x1E39D; jnb     short locret_1E3A5 */
    if(!(m->flags&1)) goto L_1E3A5;
L_1E39F: ; /* IDA 0x1E39F; mov     ah, 0 */
    vg_ah(m,(uint8_t)(0));
L_1E3A1: ; /* IDA 0x1E3A1; int     50h; through 57 - IRQ0-IRQ7 relocated by DESQview */
    gl_interrupt(m,80,58275);
L_1E3A3: ; /* IDA 0x1E3A3; jmp     short sub_1E38C */
    goto L_1E38C;
L_1E3A5: ; /* IDA 0x1E3A5; retn */
    return;
}
/* Original sub_1F4DF; IDA 0x1F4DF. */
static void tx_body_1F4DF(KiMachine16 *m,const KiEconomyHooks *h) {
L_1F4DF: ; /* IDA 0x1F4DF; push    bp */
    ec_push(m,m->bp);
L_1F4E0: ; /* IDA 0x1F4E0; sub     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,1));
L_1F4E3: ; /* IDA 0x1F4E3; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_1F4E5: ; /* IDA 0x1F4E5; mov     [bp+2], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->ax);
L_1F4E8: ; /* IDA 0x1F4E8; mov     [bp+4], cx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+4),m->cx);
L_1F4EB: ; /* IDA 0x1F4EB; mov     ax, 3D00h */
    m->ax=(uint16_t)(15616);
L_1F4EE: ; /* IDA 0x1F4EE; int     21h; DOS - 2+ - OPEN DISK FILE WITH HANDLE */
    gl_interrupt(m,33,62704);
L_1F4F0: ; /* IDA 0x1F4F0; jb      short loc_1F51F */
    if(m->flags&1) goto L_1F51F;
L_1F4F2: ; /* IDA 0x1F4F2; push    ds */
    ec_push(m,m->ds);
L_1F4F3: ; /* IDA 0x1F4F3; mov     ds, bx */
    m->ds=(uint16_t)(m->bx);
L_1F4F5: ; /* IDA 0x1F4F5; mov     [bp+0], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->ax);
L_1F4F8: ; /* IDA 0x1F4F8; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_1F4FA: ; /* IDA 0x1F4FA; mov     ax, 4200h */
    m->ax=(uint16_t)(16896);
L_1F4FD: ; /* IDA 0x1F4FD; mov     cx, [bp+4] */
    m->cx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+4)));
L_1F500: ; /* IDA 0x1F500; mov     dx, [bp+2] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_1F503: ; /* IDA 0x1F503; int     21h; DOS - 2+ - MOVE FILE READ/WRITE POINTER (LSEEK) */
    gl_interrupt(m,33,62725);
L_1F505: ; /* IDA 0x1F505; jb      short loc_1F513 */
    if(m->flags&1) goto L_1F513;
L_1F507: ; /* IDA 0x1F507; mov     ah, 3Fh ; '?' */
    vg_ah(m,(uint8_t)(63));
L_1F509: ; /* IDA 0x1F509; mov     dx, si */
    m->dx=(uint16_t)(m->si);
L_1F50B: ; /* IDA 0x1F50B; mov     cx, di */
    m->cx=(uint16_t)(m->di);
L_1F50D: ; /* IDA 0x1F50D; mov     bx, [bp+0] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1F510: ; /* IDA 0x1F510; int     21h; DOS - 2+ - READ FROM FILE WITH HANDLE */
    gl_interrupt(m,33,62738);
L_1F512: ; /* IDA 0x1F512; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1F513: ; /* IDA 0x1F513; lahf */
    vg_ah(m,(uint8_t)((m->flags&0xd5)|2));
L_1F514: ; /* IDA 0x1F514; mov     ah, 3Eh ; '>' */
    vg_ah(m,(uint8_t)(62));
L_1F516: ; /* IDA 0x1F516; mov     bx, [bp+0] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1F519: ; /* IDA 0x1F519; int     21h; DOS - 2+ - CLOSE A FILE WITH HANDLE */
    gl_interrupt(m,33,62747);
L_1F51B: ; /* IDA 0x1F51B; sahf */
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
L_1F51C: ; /* IDA 0x1F51C; mov     ax, 0FFFFh */
    m->ax=(uint16_t)(65535);
L_1F51F: ; /* IDA 0x1F51F; mov     dx, ax */
    m->dx=(uint16_t)(m->ax);
L_1F521: ; /* IDA 0x1F521; lahf */
    vg_ah(m,(uint8_t)((m->flags&0xd5)|2));
L_1F522: ; /* IDA 0x1F522; add     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,0));
L_1F525: ; /* IDA 0x1F525; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1F526: ; /* IDA 0x1F526; sahf */
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
L_1F527: ; /* IDA 0x1F527; retn */
    return;
}
/* Original (unnamed code); IDA 0x108B2. */
static void tx_body_108B2(KiMachine16 *m,const KiEconomyHooks *h) {
L_108B2: ; /* IDA 0x108B2; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_108B7: ; /* IDA 0x108B7; mov     ax, ss:[di] */
    m->ax=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->di)));
L_108BA: ; /* IDA 0x108BA; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_108BB: ; /* IDA 0x108BB; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_108BC: ; /* IDA 0x108BC; cmp     ah, 0FFh */
    ec_cmp(m,(uint8_t)(m->ax>>8),255,8);
L_108BF: ; /* IDA 0x108BF; jnz     short loc_108CE */
    if(!(m->flags&0x40)) goto L_108CE;
L_108C1: ; /* IDA 0x108C1; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_108C3: ; /* IDA 0x108C3; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_108C5: ; /* IDA 0x108C5; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_108C7: ; /* IDA 0x108C7; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_108C9: ; /* IDA 0x108C9; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_108CB: ; /* IDA 0x108CB; add     ax, 4240h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16960,16,0,0));
L_108CE: ; /* IDA 0x108CE; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_108D0: ; /* IDA 0x108D0; add     si, 8 */
#if KI_TALK_MUTATION == 5
    m->si=(uint16_t)(ec_math(m,m->si,2,16,0,0));
#else
    m->si=(uint16_t)(ec_math(m,m->si,8,16,0,0));
#endif
L_108D3: ; /* IDA 0x108D3; mov     ah, 9 */
    vg_ah(m,(uint8_t)(9));
L_108D5: ; /* IDA 0x108D5; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_108D7: ; /* IDA 0x108D7; call    loc_10701 */
    tx_call(m,1793,2266,h);
L_108DA: ; /* IDA 0x108DA; retn */
    return;
}
/* Original (unnamed code); IDA 0x108DB. */
static void tx_body_108DB(KiMachine16 *m,const KiEconomyHooks *h) {
L_108DB: ; /* IDA 0x108DB; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_108E0: ; /* IDA 0x108E0; mov     ax, ss:[di] */
    m->ax=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->di)));
L_108E3: ; /* IDA 0x108E3; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_108E4: ; /* IDA 0x108E4; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_108E5: ; /* IDA 0x108E5; cmp     ah, 0FFh */
    ec_cmp(m,(uint8_t)(m->ax>>8),255,8);
L_108E8: ; /* IDA 0x108E8; jnz     short loc_108F7 */
    if(!(m->flags&0x40)) goto L_108F7;
L_108EA: ; /* IDA 0x108EA; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_108EC: ; /* IDA 0x108EC; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_108EE: ; /* IDA 0x108EE; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_108F0: ; /* IDA 0x108F0; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_108F2: ; /* IDA 0x108F2; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_108F4: ; /* IDA 0x108F4; add     ax, 840h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2112,16,0,0));
L_108F7: ; /* IDA 0x108F7; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_108F9: ; /* IDA 0x108F9; add     si, 2 */
    m->si=(uint16_t)(ec_math(m,m->si,2,16,0,0));
L_108FC: ; /* IDA 0x108FC; mov     ah, 0Bh */
    vg_ah(m,(uint8_t)(11));
L_108FE: ; /* IDA 0x108FE; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_10900: ; /* IDA 0x10900; call    loc_10701 */
    tx_call(m,1793,2307,h);
L_10903: ; /* IDA 0x10903; retn */
    return;
}
/* Original (unnamed code); IDA 0x10904. */
static void tx_body_10904(KiMachine16 *m,const KiEconomyHooks *h) {
L_10904: ; /* IDA 0x10904; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_10909: ; /* IDA 0x10909; mov     ax, ss:[di] */
    m->ax=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->di)));
L_1090C: ; /* IDA 0x1090C; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1090D: ; /* IDA 0x1090D; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1090E: ; /* IDA 0x1090E; cmp     ah, 0FFh */
    ec_cmp(m,(uint8_t)(m->ax>>8),255,8);
L_10911: ; /* IDA 0x10911; jnz     short loc_1091E */
    if(!(m->flags&0x40)) goto L_1091E;
L_10913: ; /* IDA 0x10913; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_10915: ; /* IDA 0x10915; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_10917: ; /* IDA 0x10917; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_10919: ; /* IDA 0x10919; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1091B: ; /* IDA 0x1091B; add     ax, 0 */
    m->ax=(uint16_t)(ec_math(m,m->ax,0,16,0,0));
L_1091E: ; /* IDA 0x1091E; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_10920: ; /* IDA 0x10920; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_10923: ; /* IDA 0x10923; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_10925: ; /* IDA 0x10925; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_10927: ; /* IDA 0x10927; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_10929: ; /* IDA 0x10929; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1092B: ; /* IDA 0x1092B; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1092D: ; /* IDA 0x1092D; add     si, 4242h */
    m->si=(uint16_t)(ec_math(m,m->si,16962,16,0,0));
L_10931: ; /* IDA 0x10931; mov     ah, 0Ch */
    vg_ah(m,(uint8_t)(12));
L_10933: ; /* IDA 0x10933; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_10935: ; /* IDA 0x10935; call    loc_10701 */
    tx_call(m,1793,2360,h);
L_10938: ; /* IDA 0x10938; retn */
    return;
}
/* Original (unnamed code); IDA 0x10939. */
static void tx_body_10939(KiMachine16 *m,const KiEconomyHooks *h) {
L_10939: ; /* IDA 0x10939; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_1093E: ; /* IDA 0x1093E; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_10943: ; /* IDA 0x10943; mov     ah, [si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_10946: ; /* IDA 0x10946; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_10948: ; /* IDA 0x10948; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1094A: ; /* IDA 0x1094A; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1094C: ; /* IDA 0x1094C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1094E: ; /* IDA 0x1094E; add     ax, 4248h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16968,16,0,0));
L_10951: ; /* IDA 0x10951; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_10953: ; /* IDA 0x10953; mov     ah, 9 */
    vg_ah(m,(uint8_t)(9));
L_10955: ; /* IDA 0x10955; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_10957: ; /* IDA 0x10957; call    loc_10701 */
    tx_call(m,1793,2394,h);
L_1095A: ; /* IDA 0x1095A; retn */
    return;
}
/* Original (unnamed code); IDA 0x1095B. */
static void tx_body_1095B(KiMachine16 *m,const KiEconomyHooks *h) {
L_1095B: ; /* IDA 0x1095B; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_10960: ; /* IDA 0x10960; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_10965: ; /* IDA 0x10965; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_10968: ; /* IDA 0x10968; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1096A: ; /* IDA 0x1096A; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1096C: ; /* IDA 0x1096C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1096E: ; /* IDA 0x1096E; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_10970: ; /* IDA 0x10970; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_10972: ; /* IDA 0x10972; add     si, 4242h */
    m->si=(uint16_t)(ec_math(m,m->si,16962,16,0,0));
L_10976: ; /* IDA 0x10976; mov     ah, 0Ch */
    vg_ah(m,(uint8_t)(12));
L_10978: ; /* IDA 0x10978; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1097A: ; /* IDA 0x1097A; call    loc_10701 */
    tx_call(m,1793,2429,h);
L_1097D: ; /* IDA 0x1097D; retn */
    return;
}
/* Original (unnamed code); IDA 0x1097E. */
static void tx_body_1097E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1097E: ; /* IDA 0x1097E; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1097F: ; /* IDA 0x1097F; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_10980: ; /* IDA 0x10980; sub     dx, 30h ; '0' */
#if KI_TALK_MUTATION == 9
    ; /* lost no-output marker cursor cancellation */
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,48,16,0,1));
#endif
L_10983: ; /* IDA 0x10983; retn */
    return;
}
int ki_talk_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x6f9:tx_body_106F9(m,h);break;
case 0x75b:tx_body_1075B(m,h);break;
case 0x7d2:tx_body_107D2(m,h);break;
case 0x84a:tx_body_1084A(m,h);break;
case 0x8853:tx_body_18853(m,h);break;
case 0x89a4:tx_body_189A4(m,h);break;
case 0xe38c:tx_body_1E38C(m,h);break;
case 0xf4df:tx_body_1F4DF(m,h);break;
case 0x8b2:tx_body_108B2(m,h);break;
case 0x8db:tx_body_108DB(m,h);break;
case 0x904:tx_body_10904(m,h);break;
case 0x939:tx_body_10939(m,h);break;
case 0x95b:tx_body_1095B(m,h);break;
case 0x97e:tx_body_1097E(m,h);break;
default:return 0;}return 1;}
