/* Fixed IDA native C number control; spec/229, re/109. Original locators retained. */
/* Original sub_1062F; IDA 0x1062F. */
static void nr_body_1062F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1062F: ; /* IDA 0x1062F; push    ds */
    ec_push(m,m->ds);
L_10630: ; /* IDA 0x10630; push    es */
    ec_push(m,m->es);
L_10631: ; /* IDA 0x10631; push    cx */
    ec_push(m,m->cx);
L_10632: ; /* IDA 0x10632; push    si */
    ec_push(m,m->si);
L_10633: ; /* IDA 0x10633; push    di */
    ec_push(m,m->di);
L_10634: ; /* IDA 0x10634; push    bp */
    ec_push(m,m->bp);
L_10635: ; /* IDA 0x10635; cld */
    m->flags&=(uint16_t)~0x400u;
L_10636: ; /* IDA 0x10636; mov     cx, 0A0C8h */
    m->cx=(uint16_t)(41160);
L_10639: ; /* IDA 0x10639; mov     es, cx */
    m->es=(uint16_t)(m->cx);
L_1063B: ; /* IDA 0x1063B; mov     ds, cs:word_10D54 */
#if KI_NUMBERS_MUTATION == 3
    m->ds=(uint16_t)(dl_read16(m,m->cs,3404));
#else
    m->ds=(uint16_t)(dl_read16(m,m->cs,3412));
#endif
L_10640: ; /* IDA 0x10640; mov     cl, bl */
    vg_cl(m,(uint8_t)((uint8_t)m->bx));
L_10642: ; /* IDA 0x10642; xor     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_10644: ; /* IDA 0x10644; dec     cx */
#if KI_NUMBERS_MUTATION == 2
    ; /* lost right-edge width adjustment */
#else
    m->cx=(uint16_t)(st_dec(m,m->cx,16));
#endif
L_10645: ; /* IDA 0x10645; add     di, cx */
    m->di=(uint16_t)(ec_math(m,m->di,m->cx,16,0,0));
L_10647: ; /* IDA 0x10647; mov     bp, dx */
    m->bp=(uint16_t)(m->dx);
L_10649: ; /* IDA 0x10649; cmp     dh, 80h */
    ec_cmp(m,(uint8_t)(m->dx>>8),128,8);
L_1064C: ; /* IDA 0x1064C; jb      short loc_1065E */
    if(m->flags&1) goto L_1065E;
L_1064E: ; /* IDA 0x1064E; not     dx */
    m->dx=(uint16_t)(~(m->dx));
L_10650: ; /* IDA 0x10650; not     ax */
    m->ax=(uint16_t)(~(m->ax));
L_10652: ; /* IDA 0x10652; add     ax, 1 */
    m->ax=(uint16_t)(ec_math(m,m->ax,1,16,0,0));
L_10655: ; /* IDA 0x10655; adc     dx, 0 */
#if KI_NUMBERS_MUTATION == 1
    ; /* lost carry for negative multiples of 65536 */
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,0,16,m->flags&1,0));
#endif
L_10658: ; /* IDA 0x10658; dec     bl */
    gw_bl(m,(uint8_t)(st_dec(m,(uint8_t)m->bx,8)));
L_1065A: ; /* IDA 0x1065A; jnz     short loc_1065E */
    if(!(m->flags&0x40)) goto L_1065E;
L_1065C: ; /* IDA 0x1065C; mov     bl, 1 */
    gw_bl(m,(uint8_t)(1));
L_1065E: ; /* IDA 0x1065E; call    sub_10CAC */
    nr_call(m,3244,1633,h);
L_10661: ; /* IDA 0x10661; mov     cx, 0Ah */
    m->cx=(uint16_t)(10);
L_10664: ; /* IDA 0x10664; div     cx */
    dl_div(m,m->cx,16);
L_10666: ; /* IDA 0x10666; call    sub_1069A */
    nr_call(m,1690,1641,h);
L_10669: ; /* IDA 0x10669; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_1066B: ; /* IDA 0x1066B; dec     di */
    m->di=(uint16_t)(st_dec(m,m->di,16));
L_1066C: ; /* IDA 0x1066C; dec     bl */
    gw_bl(m,(uint8_t)(st_dec(m,(uint8_t)m->bx,8)));
L_1066E: ; /* IDA 0x1066E; jz      short loc_10688 */
    if(m->flags&0x40) goto L_10688;
L_10670: ; /* IDA 0x10670; and     ax, ax */
    {uint16_t result=((m->ax)&(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_10672: ; /* IDA 0x10672; jnz     short loc_10661 */
    if(!(m->flags&0x40)) goto L_10661;
L_10674: ; /* IDA 0x10674; mov     dh, bh */
#if KI_NUMBERS_MUTATION == 8
    vg_dh(m,(uint8_t)((uint8_t)(m->ax>>8)));
#else
    vg_dh(m,(uint8_t)((uint8_t)(m->bx>>8)));
#endif
L_10676: ; /* IDA 0x10676; shr     dh, 1 */
    vg_dh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->dx>>8),1,8,5)));
L_10678: ; /* IDA 0x10678; shr     dh, 1 */
    vg_dh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->dx>>8),1,8,5)));
L_1067A: ; /* IDA 0x1067A; shr     dh, 1 */
    vg_dh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->dx>>8),1,8,5)));
L_1067C: ; /* IDA 0x1067C; shr     dh, 1 */
    vg_dh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->dx>>8),1,8,5)));
L_1067E: ; /* IDA 0x1067E; mov     ah, dh */
    vg_ah(m,(uint8_t)((uint8_t)(m->dx>>8)));
L_10680: ; /* IDA 0x10680; call    sub_106DE */
    nr_call(m,1758,1667,h);
L_10683: ; /* IDA 0x10683; dec     di */
    m->di=(uint16_t)(st_dec(m,m->di,16));
L_10684: ; /* IDA 0x10684; dec     bl */
    gw_bl(m,(uint8_t)(st_dec(m,(uint8_t)m->bx,8)));
L_10686: ; /* IDA 0x10686; jnz     short loc_1067E */
    if(!(m->flags&0x40)) goto L_1067E;
L_10688: ; /* IDA 0x10688; cmp     bp, 0 */
    ec_cmp(m,m->bp,0,16);
L_1068B: ; /* IDA 0x1068B; jge     short loc_10693 */
    if(!ec_less(m)) goto L_10693;
L_1068D: ; /* IDA 0x1068D; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_10690: ; /* IDA 0x10690; call    sub_1069A */
    nr_call(m,1690,1683,h);
L_10693: ; /* IDA 0x10693; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_10694: ; /* IDA 0x10694; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_10695: ; /* IDA 0x10695; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_10696: ; /* IDA 0x10696; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10697: ; /* IDA 0x10697; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_10698: ; /* IDA 0x10698; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_10699: ; /* IDA 0x10699; retn */
    return;
}
/* Original sub_1069A; IDA 0x1069A. */
static void nr_body_1069A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1069A: ; /* IDA 0x1069A; push    ax */
    ec_push(m,m->ax);
L_1069B: ; /* IDA 0x1069B; mov     si, dx */
    m->si=(uint16_t)(m->dx);
L_1069D: ; /* IDA 0x1069D; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_1069F: ; /* IDA 0x1069F; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_106A1: ; /* IDA 0x106A1; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_106A3: ; /* IDA 0x106A3; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_106A5: ; /* IDA 0x106A5; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_106A8: ; /* IDA 0x106A8; mov     al, bh */
    vg_al(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_106AA: ; /* IDA 0x106AA; and     al, 0Fh */
#if KI_NUMBERS_MUTATION == 4
    {uint16_t result=(((uint8_t)m->ax)&(0))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
#else
    {uint16_t result=(((uint8_t)m->ax)&(15))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
#endif
L_106AC: ; /* IDA 0x106AC; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_106AD: ; /* IDA 0x106AD; mov     cl, 10h */
    vg_cl(m,(uint8_t)(16));
L_106AF: ; /* IDA 0x106AF; mov     al, es:[di] */
#if KI_NUMBERS_MUTATION == 5
    ; /* lost latch read */
#else
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
#endif
L_106B2: ; /* IDA 0x106B2; movsb */
    vg_movsb(m);
L_106B3: ; /* IDA 0x106B3; add     di, 4Fh ; 'O' */
#if KI_NUMBERS_MUTATION == 7
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
#else
    m->di=(uint16_t)(ec_math(m,m->di,79,16,0,0));
#endif
L_106B6: ; /* IDA 0x106B6; loop    loc_106AF */
    if(--m->cx) goto L_106AF;
L_106B8: ; /* IDA 0x106B8; sub     si, 10h */
    m->si=(uint16_t)(ec_math(m,m->si,16,16,0,1));
L_106BB: ; /* IDA 0x106BB; sub     di, 500h */
    m->di=(uint16_t)(ec_math(m,m->di,1280,16,0,1));
L_106BF: ; /* IDA 0x106BF; mov     al, bh */
    vg_al(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_106C1: ; /* IDA 0x106C1; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_106C3: ; /* IDA 0x106C3; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_106C5: ; /* IDA 0x106C5; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_106C7: ; /* IDA 0x106C7; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_106C9: ; /* IDA 0x106C9; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_106CA: ; /* IDA 0x106CA; mov     cl, 10h */
    vg_cl(m,(uint8_t)(16));
L_106CC: ; /* IDA 0x106CC; mov     al, es:[di] */
#if KI_NUMBERS_MUTATION == 5
    ; /* lost latch read */
#else
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
#endif
L_106CF: ; /* IDA 0x106CF; lodsb */
    vg_al(m,dl_read8(m,m->ds,m->si));m->si=(uint16_t)(m->si+((m->flags&0x400)?-1:1));
L_106D0: ; /* IDA 0x106D0; not     al */
    vg_al(m,(uint8_t)(~((uint8_t)m->ax)));
L_106D2: ; /* IDA 0x106D2; stosb */
    dl_write8(m,m->es,m->di,(uint8_t)m->ax);m->di=(uint16_t)(m->di+((m->flags&0x400)?-1:1));
L_106D3: ; /* IDA 0x106D3; add     di, 4Fh ; 'O' */
#if KI_NUMBERS_MUTATION == 7
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
#else
    m->di=(uint16_t)(ec_math(m,m->di,79,16,0,0));
#endif
L_106D6: ; /* IDA 0x106D6; loop    loc_106CC */
    if(--m->cx) goto L_106CC;
L_106D8: ; /* IDA 0x106D8; sub     di, 500h */
    m->di=(uint16_t)(ec_math(m,m->di,1280,16,0,1));
L_106DC: ; /* IDA 0x106DC; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_106DD: ; /* IDA 0x106DD; retn */
    return;
}
/* Original sub_106DE; IDA 0x106DE. */
static void nr_body_106DE(KiMachine16 *m,const KiEconomyHooks *h) {
L_106DE: ; /* IDA 0x106DE; push    dx */
    ec_push(m,m->dx);
L_106DF: ; /* IDA 0x106DF; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_106E2: ; /* IDA 0x106E2; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_106E4: ; /* IDA 0x106E4; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_106E5: ; /* IDA 0x106E5; mov     cl, 10h */
    vg_cl(m,(uint8_t)(16));
L_106E7: ; /* IDA 0x106E7; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_106E9: ; /* IDA 0x106E9; stosb */
#if KI_NUMBERS_MUTATION == 6
    dl_write8(m,m->es,m->di,(uint8_t)m->ax);m->di=(uint16_t)(m->di+1);
#else
    dl_write8(m,m->es,m->di,(uint8_t)m->ax);m->di=(uint16_t)(m->di+((m->flags&0x400)?-1:1));
#endif
L_106EA: ; /* IDA 0x106EA; add     di, 4Fh ; 'O' */
    m->di=(uint16_t)(ec_math(m,m->di,79,16,0,0));
L_106ED: ; /* IDA 0x106ED; loop    loc_106E9 */
    if(--m->cx) goto L_106E9;
L_106EF: ; /* IDA 0x106EF; sub     di, 500h */
    m->di=(uint16_t)(ec_math(m,m->di,1280,16,0,1));
L_106F3: ; /* IDA 0x106F3; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_106F4: ; /* IDA 0x106F4; retn */
    return;
}
/* Original ; IDA 0x10984. */
static void nr_body_10984(KiMachine16 *m,const KiEconomyHooks *h) {
L_10984: ; /* IDA 0x10984; mov     ax, ss:[di] */
    m->ax=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->di)));
L_10987: ; /* IDA 0x10987; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_10988: ; /* IDA 0x10988; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_10989: ; /* IDA 0x10989; push    bx */
    ec_push(m,m->bx);
L_1098A: ; /* IDA 0x1098A; push    dx */
    ec_push(m,m->dx);
L_1098B: ; /* IDA 0x1098B; push    di */
    ec_push(m,m->di);
L_1098C: ; /* IDA 0x1098C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1098E: ; /* IDA 0x1098E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10990: ; /* IDA 0x10990; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10992: ; /* IDA 0x10992; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10994: ; /* IDA 0x10994; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_10996: ; /* IDA 0x10996; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10998: ; /* IDA 0x10998; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1099A: ; /* IDA 0x1099A; add     di, bx */
    m->di=(uint16_t)(ec_math(m,m->di,m->bx,16,0,0));
L_1099C: ; /* IDA 0x1099C; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1099E: ; /* IDA 0x1099E; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_109A0: ; /* IDA 0x109A0; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_109A2: ; /* IDA 0x109A2; add     di, dx */
    m->di=(uint16_t)(ec_math(m,m->di,m->dx,16,0,0));
L_109A4: ; /* IDA 0x109A4; cwd */
    m->dx=(m->ax&0x8000)?0xffff:0;
L_109A5: ; /* IDA 0x109A5; mov     bx, 0F06h */
    m->bx=(uint16_t)(3846);
L_109A8: ; /* IDA 0x109A8; call    sub_1062F */
    nr_call(m,1583,2475,h);
L_109AB: ; /* IDA 0x109AB; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_109AC: ; /* IDA 0x109AC; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_109AD: ; /* IDA 0x109AD; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_109AE: ; /* IDA 0x109AE; retn */
    return;
}
/* Original sub_11E17; IDA 0x11E17. */
static void nr_body_11E17(KiMachine16 *m,const KiEconomyHooks *h) {
L_11E17: ; /* IDA 0x11E17; mov     ax, ds:0CF6h */
    m->ax=(uint16_t)(dl_read16(m,m->ds,3318));
L_11E1A: ; /* IDA 0x11E1A; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_11E1C: ; /* IDA 0x11E1C; mov     bx, 903h */
    m->bx=(uint16_t)(2307);
L_11E1F: ; /* IDA 0x11E1F; mov     di, 2BBh */
    m->di=(uint16_t)(699);
L_11E22: ; /* IDA 0x11E22; call    sub_1062F */
    nr_call(m,1583,7717,h);
L_11E25: ; /* IDA 0x11E25; mov     al, ds:0CF4h */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,3316)));
L_11E28: ; /* IDA 0x11E28; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_11E2A: ; /* IDA 0x11E2A; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_11E2C: ; /* IDA 0x11E2C; mov     bx, 902h */
    m->bx=(uint16_t)(2306);
L_11E2F: ; /* IDA 0x11E2F; mov     di, 2C0h */
    m->di=(uint16_t)(704);
L_11E32: ; /* IDA 0x11E32; call    sub_1062F */
    nr_call(m,1583,7733,h);
L_11E35: ; /* IDA 0x11E35; mov     al, ds:0CF0h */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,3312)));
L_11E38: ; /* IDA 0x11E38; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_11E3A: ; /* IDA 0x11E3A; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_11E3C: ; /* IDA 0x11E3C; mov     bx, 902h */
    m->bx=(uint16_t)(2306);
L_11E3F: ; /* IDA 0x11E3F; mov     di, 2C4h */
    m->di=(uint16_t)(708);
L_11E42: ; /* IDA 0x11E42; call    sub_1062F */
    nr_call(m,1583,7749,h);
L_11E45: ; /* IDA 0x11E45; retn */
    return;
}
int ki_numbers_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x62f:nr_body_1062F(m,h);break;
case 0x69a:nr_body_1069A(m,h);break;
case 0x6de:nr_body_106DE(m,h);break;
case 0x984:nr_body_10984(m,h);break;
case 0x1e17:nr_body_11E17(m,h);break;
default:return 0;}return 1;}
