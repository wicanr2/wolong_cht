/* Original native C city/object tick; re/128, spec/248. No scheduler or full main loop. */
/* Original sub_10CDE; IDA 0x10CDE. */
static void ct_body_10CDE(KiMachine16 *m,const KiEconomyHooks *h) {
L_10CDE: ; /* IDA 0x10CDE; push    ax */
    ec_push(m,m->ax);
L_10CDF: ; /* IDA 0x10CDF; mov     ax, 101h */
    m->ax=(uint16_t)(257);
L_10CE2: ; /* IDA 0x10CE2; call    sub_1EB11 */
    ct_call(m,60177,3301,h);
L_10CE5: ; /* IDA 0x10CE5; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10CE6: ; /* IDA 0x10CE6; retn */
    return;
}
/* Original sub_12459; IDA 0x12459. */
static void ct_body_12459(KiMachine16 *m,const KiEconomyHooks *h) {
L_12459: ; /* IDA 0x12459; mov     ds, word ptr ds:0D52h */
    m->ds=(uint16_t)(dl_read16(m,m->ds,3410));
L_1245D: ; /* IDA 0x1245D; mov     si, 2040h */
    m->si=(uint16_t)(8256);
L_12460: ; /* IDA 0x12460; mov     bp, 20h ; ' ' */
#if KI_TICK_MUTATION == 2
    m->bp=31;
#else
    m->bp=(uint16_t)(32);
#endif
L_12463: ; /* IDA 0x12463; cmp     byte ptr [si], 80h */
#if KI_TICK_MUTATION == 1
    ec_cmp(m,dl_read8(m,m->ds,m->si),0x81,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
#endif
L_12466: ; /* IDA 0x12466; jb      short loc_1247F */
    if(m->flags&1) goto L_1247F;
L_12468: ; /* IDA 0x12468; dec     byte ptr [si+0Ch] */
    dl_write8(m,m->ds,(uint16_t)(m->si+12),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->si+12)),8));
L_1246B: ; /* IDA 0x1246B; jnz     short loc_1247F */
    if(!(m->flags&0x40)) goto L_1247F;
L_1246D: ; /* IDA 0x1246D; mov     al, [si+0Dh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+13))));
L_12470: ; /* IDA 0x12470; mov     [si+0Ch], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+12),(uint8_t)m->ax);
L_12473: ; /* IDA 0x12473; or      byte ptr [si], 1 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(1))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_12476: ; /* IDA 0x12476; cmp     si, 2140h */
#if KI_TICK_MUTATION == 3
    ec_cmp(m,m->si,0x2150,16);
#else
    ec_cmp(m,m->si,8512,16);
#endif
L_1247A: ; /* IDA 0x1247A; jb      short loc_1247F */
    if(m->flags&1) goto L_1247F;
L_1247C: ; /* IDA 0x1247C; call    sub_1248A */
    ct_call(m,9354,9343,h);
L_1247F: ; /* IDA 0x1247F; add     si, 10h */
    m->si=(uint16_t)(ec_math(m,m->si,16,16,0,0));
L_12482: ; /* IDA 0x12482; dec     bp */
    m->bp=(uint16_t)(st_dec(m,m->bp,16));
L_12483: ; /* IDA 0x12483; jnz     short loc_12463 */
    if(!(m->flags&0x40)) goto L_12463;
L_12485: ; /* IDA 0x12485; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_12487: ; /* IDA 0x12487; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_12489: ; /* IDA 0x12489; retn */
    return;
}
/* Original sub_1248A; IDA 0x1248A. */
static void ct_body_1248A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1248A: ; /* IDA 0x1248A; mov     dx, [si+8] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+8)));
L_1248D: ; /* IDA 0x1248D; call    sub_124FF */
    ct_call(m,9471,9360,h);
L_12490: ; /* IDA 0x12490; mov     [si+8], dx */
    dl_write16(m,m->ds,(uint16_t)(m->si+8),m->dx);
L_12493: ; /* IDA 0x12493; add     [si+2], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+2),ec_math(m,dl_read16(m,m->ds,(uint16_t)(m->si+2)),m->ax,16,0,0));
L_12496: ; /* IDA 0x12496; mov     dx, [si+0Ah] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+10)));
L_12499: ; /* IDA 0x12499; call    sub_124FF */
    ct_call(m,9471,9372,h);
L_1249C: ; /* IDA 0x1249C; mov     [si+0Ah], dx */
    dl_write16(m,m->ds,(uint16_t)(m->si+10),m->dx);
L_1249F: ; /* IDA 0x1249F; add     [si+4], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+4),ec_math(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),m->ax,16,0,0));
L_124A2: ; /* IDA 0x124A2; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_124A4: ; /* IDA 0x124A4; mov     ax, [si+2] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+2)));
L_124A7: ; /* IDA 0x124A7; cmp     ax, cs:word_10D22 */
    ec_cmp(m,m->ax,dl_read16(m,m->cs,3362),16);
L_124AC: ; /* IDA 0x124AC; jge     short loc_124B0 */
    if(!ec_less(m)) goto L_124B0;
L_124AE: ; /* IDA 0x124AE; mov     dl, 1 */
    vg_dl(m,(uint8_t)(1));
L_124B0: ; /* IDA 0x124B0; cmp     ax, cs:word_10D26 */
    ec_cmp(m,m->ax,dl_read16(m,m->cs,3366),16);
L_124B5: ; /* IDA 0x124B5; jle     short loc_124B9 */
    if(!ec_greater(m)) goto L_124B9;
L_124B7: ; /* IDA 0x124B7; mov     dl, 0FFh */
    vg_dl(m,(uint8_t)(255));
L_124B9: ; /* IDA 0x124B9; cmp     ax, 0FFF0h */
    ec_cmp(m,m->ax,65520,16);
L_124BC: ; /* IDA 0x124BC; jge     short loc_124C3 */
    if(!ec_less(m)) goto L_124C3;
L_124BE: ; /* IDA 0x124BE; mov     ax, 190h */
    m->ax=(uint16_t)(400);
L_124C1: ; /* IDA 0x124C1; jmp     short loc_124CB */
    goto L_124CB;
L_124C3: ; /* IDA 0x124C3; cmp     ax, 190h */
    ec_cmp(m,m->ax,400,16);
L_124C6: ; /* IDA 0x124C6; jle     short loc_124CB */
    if(!ec_greater(m)) goto L_124CB;
L_124C8: ; /* IDA 0x124C8; mov     ax, 0FFF0h */
    m->ax=(uint16_t)(65520);
L_124CB: ; /* IDA 0x124CB; mov     [si+2], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+2),m->ax);
L_124CE: ; /* IDA 0x124CE; mov     ax, [si+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+4)));
L_124D1: ; /* IDA 0x124D1; cmp     ax, cs:word_10D24 */
    ec_cmp(m,m->ax,dl_read16(m,m->cs,3364),16);
L_124D6: ; /* IDA 0x124D6; jge     short loc_124DA */
    if(!ec_less(m)) goto L_124DA;
L_124D8: ; /* IDA 0x124D8; mov     dh, 1 */
    vg_dh(m,(uint8_t)(1));
L_124DA: ; /* IDA 0x124DA; cmp     ax, cs:word_10D28 */
    ec_cmp(m,m->ax,dl_read16(m,m->cs,3368),16);
L_124DF: ; /* IDA 0x124DF; jle     short loc_124E3 */
    if(!ec_greater(m)) goto L_124E3;
L_124E1: ; /* IDA 0x124E1; mov     dh, 0FFh */
    vg_dh(m,(uint8_t)(255));
L_124E3: ; /* IDA 0x124E3; cmp     ax, 0FFF0h */
    ec_cmp(m,m->ax,65520,16);
L_124E6: ; /* IDA 0x124E6; jge     short loc_124ED */
    if(!ec_less(m)) goto L_124ED;
L_124E8: ; /* IDA 0x124E8; mov     ax, 110h */
    m->ax=(uint16_t)(272);
L_124EB: ; /* IDA 0x124EB; jmp     short loc_124F5 */
    goto L_124F5;
L_124ED: ; /* IDA 0x124ED; cmp     ax, 110h */
    ec_cmp(m,m->ax,272,16);
L_124F0: ; /* IDA 0x124F0; jle     short loc_124F5 */
    if(!ec_greater(m)) goto L_124F5;
L_124F2: ; /* IDA 0x124F2; mov     ax, 0FFF0h */
    m->ax=(uint16_t)(65520);
L_124F5: ; /* IDA 0x124F5; mov     [si+4], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+4),m->ax);
L_124F8: ; /* IDA 0x124F8; add     [si+9], dl */
    dl_write8(m,m->ds,(uint16_t)(m->si+9),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->si+9)),(uint8_t)m->dx,8,0,0));
L_124FB: ; /* IDA 0x124FB; add     [si+0Bh], dh */
    dl_write8(m,m->ds,(uint16_t)(m->si+11),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->si+11)),(uint8_t)(m->dx>>8),8,0,0));
L_124FE: ; /* IDA 0x124FE; retn */
    return;
}
/* Original sub_124FF; IDA 0x124FF. */
static void ct_body_124FF(KiMachine16 *m,const KiEconomyHooks *h) {
L_124FF: ; /* IDA 0x124FF; call    sub_1ECE0 */
    ct_call(m,60640,9474,h);
L_12502: ; /* IDA 0x12502; and     al, 7 */
#if KI_TICK_MUTATION == 4
    {uint8_t result=(uint8_t)m->ax&3;ec_logic(m,result,8);vg_al(m,result);}
#else
    {uint16_t result=(((uint8_t)m->ax)&(7))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
#endif
L_12504: ; /* IDA 0x12504; cmp     al, 2 */
    ec_cmp(m,(uint8_t)m->ax,2,8);
L_12506: ; /* IDA 0x12506; ja      short loc_1250C */
    if(!(m->flags&(1|0x40))) goto L_1250C;
L_12508: ; /* IDA 0x12508; dec     al */
    vg_al(m,(uint8_t)(st_dec(m,(uint8_t)m->ax,8)));
L_1250A: ; /* IDA 0x1250A; add     dh, al */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),(uint8_t)m->ax,8,0,0)));
L_1250C: ; /* IDA 0x1250C; cmp     dh, 0F1h */
#if KI_TICK_MUTATION == 5
    ec_cmp(m,(uint8_t)(m->dx>>8),0xf2,8);
#else
    ec_cmp(m,(uint8_t)(m->dx>>8),241,8);
#endif
L_1250F: ; /* IDA 0x1250F; jge     short loc_12513 */
    if(!ec_less(m)) goto L_12513;
L_12511: ; /* IDA 0x12511; mov     dh, 0F1h */
#if KI_TICK_MUTATION == 5
    vg_dh(m,0xf2);
#else
    vg_dh(m,(uint8_t)(241));
#endif
L_12513: ; /* IDA 0x12513; cmp     dh, 0Fh */
    ec_cmp(m,(uint8_t)(m->dx>>8),15,8);
L_12516: ; /* IDA 0x12516; jle     short loc_1251A */
    if(!ec_greater(m)) goto L_1251A;
L_12518: ; /* IDA 0x12518; mov     dh, 0Fh */
    vg_dh(m,(uint8_t)(15));
L_1251A: ; /* IDA 0x1251A; add     dl, dh */
    vg_dl(m,(uint8_t)(ec_math(m,(uint8_t)m->dx,(uint8_t)(m->dx>>8),8,0,0)));
L_1251C: ; /* IDA 0x1251C; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_1251E: ; /* IDA 0x1251E; cmp     dl, 0Fh */
#if KI_TICK_MUTATION == 6
    ec_cmp(m,(uint8_t)m->dx,0x10,8);
#else
    ec_cmp(m,(uint8_t)m->dx,15,8);
#endif
L_12521: ; /* IDA 0x12521; jl      short loc_12529 */
    if(ec_less(m)) goto L_12529;
L_12523: ; /* IDA 0x12523; sub     dl, 0Fh */
    vg_dl(m,(uint8_t)(ec_math(m,(uint8_t)m->dx,15,8,0,1)));
L_12526: ; /* IDA 0x12526; inc     ax */
    m->ax=(uint16_t)(ec_inc(m,m->ax,16));
L_12527: ; /* IDA 0x12527; jmp     short locret_12532 */
    goto L_12532;
L_12529: ; /* IDA 0x12529; cmp     dl, 0F1h */
    ec_cmp(m,(uint8_t)m->dx,241,8);
L_1252C: ; /* IDA 0x1252C; jg      short locret_12532 */
    if(ec_greater(m)) goto L_12532;
L_1252E: ; /* IDA 0x1252E; add     dl, 0Fh */
    vg_dl(m,(uint8_t)(ec_math(m,(uint8_t)m->dx,15,8,0,0)));
L_12531: ; /* IDA 0x12531; dec     ax */
    m->ax=(uint16_t)(st_dec(m,m->ax,16));
L_12532: ; /* IDA 0x12532; retn */
    return;
}
/* Original sub_13EFD; IDA 0x13EFD. */
static void ct_body_13EFD(KiMachine16 *m,const KiEconomyHooks *h) {
L_13EFD: ; /* IDA 0x13EFD; mov     si, word_10D1E */
    m->si=(uint16_t)(dl_read16(m,m->ds,3358));
L_13F01: ; /* IDA 0x13F01; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_13F06: ; /* IDA 0x13F06; cmp     byte ptr [si+857h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+2135)),0,8);
L_13F0B: ; /* IDA 0x13F0B; jz      short loc_13F11 */
    if(m->flags&0x40) goto L_13F11;
L_13F0D: ; /* IDA 0x13F0D; dec     byte ptr [si+857h] */
    dl_write8(m,m->ds,(uint16_t)(m->si+2135),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->si+2135)),8));
L_13F11: ; /* IDA 0x13F11; mov     bh, [si+85Ah] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2138))));
L_13F15: ; /* IDA 0x13F15; cmp     bh, [si+841h] */
    ec_cmp(m,(uint8_t)(m->bx>>8),dl_read8(m,m->ds,(uint16_t)(m->si+2113)),8);
L_13F19: ; /* IDA 0x13F19; jz      short loc_13F2C */
    if(m->flags&0x40) goto L_13F2C;
L_13F1B: ; /* IDA 0x13F1B; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_13F1D: ; /* IDA 0x13F1D; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13F1F: ; /* IDA 0x13F1F; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13F21: ; /* IDA 0x13F21; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_13F23: ; /* IDA 0x13F23; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_13F25: ; /* IDA 0x13F25; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_13F27: ; /* IDA 0x13F27; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_13F29: ; /* IDA 0x13F29; mov     [bx+17h], ah */
#if KI_TICK_MUTATION == 8
    dl_write8(m,m->ds,(uint16_t)(m->bx+0x16),(uint8_t)(m->ax>>8));
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),(uint8_t)(m->ax>>8));
#endif
L_13F2C: ; /* IDA 0x13F2C; mov     ax, [si+84Ah] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+2122)));
L_13F30: ; /* IDA 0x13F30; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_13F32: ; /* IDA 0x13F32; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_13F34: ; /* IDA 0x13F34; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_13F36: ; /* IDA 0x13F36; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_13F38: ; /* IDA 0x13F38; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_13F3A: ; /* IDA 0x13F3A; add     ax, bx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->bx,16,0,0));
L_13F3C: ; /* IDA 0x13F3C; add     ax, cs:word_19872 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,39026),16,0,0));
L_13F41: ; /* IDA 0x13F41; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_13F43: ; /* IDA 0x13F43; mov     bx, [si+848h] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+2120)));
L_13F47: ; /* IDA 0x13F47; mov     al, es:[bx] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->bx))));
L_13F4A: ; /* IDA 0x13F4A; and     al, 7Fh */
    {uint16_t result=(((uint8_t)m->ax)&(127))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_13F4C: ; /* IDA 0x13F4C; mov     [si+858h], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+2136),(uint8_t)m->ax);
L_13F50: ; /* IDA 0x13F50; cmp     byte ptr [si+841h], 18h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+2113)),24,8);
L_13F55: ; /* IDA 0x13F55; jz      short loc_13F5A */
    if(m->flags&0x40) goto L_13F5A;
L_13F57: ; /* IDA 0x13F57; call    sub_13F74 */
    ct_call(m,16244,16218,h);
L_13F5A: ; /* IDA 0x13F5A; call    sub_14194 */
    ct_call(m,16788,16221,h);
L_13F5D: ; /* IDA 0x13F5D; call    sub_14269 */
    ct_call(m,17001,16224,h);
L_13F60: ; /* IDA 0x13F60; add     si, 20h ; ' ' */
    m->si=(uint16_t)(ec_math(m,m->si,32,16,0,0));
L_13F63: ; /* IDA 0x13F63; cmp     si, 1800h */
#if KI_TICK_MUTATION == 7
    ec_cmp(m,m->si,0x17e0,16);
#else
    ec_cmp(m,m->si,6144,16);
#endif
L_13F67: ; /* IDA 0x13F67; jb      short loc_13F6B */
    if(m->flags&1) goto L_13F6B;
L_13F69: ; /* IDA 0x13F69; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_13F6B: ; /* IDA 0x13F6B; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_13F6D: ; /* IDA 0x13F6D; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_13F6F: ; /* IDA 0x13F6F; mov     word_10D1E, si */
    dl_write16(m,m->ds,3358,m->si);
L_13F73: ; /* IDA 0x13F73; retn */
    return;
}
/* Original sub_13F74; IDA 0x13F74. */
static void ct_body_13F74(KiMachine16 *m,const KiEconomyHooks *h) {
L_13F74: ; /* IDA 0x13F74; push    es */
    ec_push(m,m->es);
L_13F75: ; /* IDA 0x13F75; push    ax */
    ec_push(m,m->ax);
L_13F76: ; /* IDA 0x13F76; push    bx */
    ec_push(m,m->bx);
L_13F77: ; /* IDA 0x13F77; push    cx */
    ec_push(m,m->cx);
L_13F78: ; /* IDA 0x13F78; push    dx */
    ec_push(m,m->dx);
L_13F79: ; /* IDA 0x13F79; push    di */
    ec_push(m,m->di);
L_13F7A: ; /* IDA 0x13F7A; push    bp */
    ec_push(m,m->bp);
L_13F7B: ; /* IDA 0x13F7B; sub     sp, 10h */
    m->sp=(uint16_t)(ec_math(m,m->sp,16,16,0,1));
L_13F7E: ; /* IDA 0x13F7E; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_13F80: ; /* IDA 0x13F80; mov     ax, ss */
    m->ax=(uint16_t)(m->ss);
L_13F82: ; /* IDA 0x13F82; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_13F84: ; /* IDA 0x13F84; mov     di, bp */
    m->di=(uint16_t)(m->bp);
L_13F86: ; /* IDA 0x13F86; mov     cx, 8 */
    m->cx=(uint16_t)(8);
L_13F89: ; /* IDA 0x13F89; mov     ax, 0FFFFh */
    m->ax=(uint16_t)(65535);
L_13F8C: ; /* IDA 0x13F8C; cld */
    m->flags&=(uint16_t)~0x400u;
L_13F8D: ; /* IDA 0x13F8D; rep stosw */
    mc_string(m,2,1);
L_13F8F: ; /* IDA 0x13F8F; call    sub_13FA9 */
    ct_call(m,16297,16274,h);
L_13F92: ; /* IDA 0x13F92; mov     [si+854h], ch */
    dl_write8(m,m->ds,(uint16_t)(m->si+2132),(uint8_t)(m->cx>>8));
L_13F96: ; /* IDA 0x13F96; call    sub_14028 */
    ct_call(m,16424,16281,h);
L_13F99: ; /* IDA 0x13F99; jnb     short loc_13F9E */
    if(!(m->flags&1)) goto L_13F9E;
L_13F9B: ; /* IDA 0x13F9B; call    sub_14057 */
    ct_call(m,16471,16286,h);
L_13F9E: ; /* IDA 0x13F9E; add     sp, 10h */
    m->sp=(uint16_t)(ec_math(m,m->sp,16,16,0,0));
L_13FA1: ; /* IDA 0x13FA1; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_13FA2: ; /* IDA 0x13FA2; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_13FA3: ; /* IDA 0x13FA3; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_13FA4: ; /* IDA 0x13FA4; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_13FA5: ; /* IDA 0x13FA5; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_13FA6: ; /* IDA 0x13FA6; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_13FA7: ; /* IDA 0x13FA7; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_13FA8: ; /* IDA 0x13FA8; retn */
    return;
}
/* Original sub_13FA9; IDA 0x13FA9. */
static void ct_body_13FA9(KiMachine16 *m,const KiEconomyHooks *h) {
L_13FA9: ; /* IDA 0x13FA9; xor     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_13FAB: ; /* IDA 0x13FAB; cmp     byte ptr [si+85Bh], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+2139)),0,8);
L_13FB0: ; /* IDA 0x13FB0; jnz     short loc_13FB3 */
    if(!(m->flags&0x40)) goto L_13FB3;
L_13FB2: ; /* IDA 0x13FB2; retn */
    return;
L_13FB3: ; /* IDA 0x13FB3; mov     dl, [si+841h] */
    vg_dl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2113))));
L_13FB7: ; /* IDA 0x13FB7; mov     bh, dl */
    gw_bh(m,(uint8_t)((uint8_t)m->dx));
L_13FB9: ; /* IDA 0x13FB9; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_13FBB: ; /* IDA 0x13FBB; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13FBD: ; /* IDA 0x13FBD; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13FBF: ; /* IDA 0x13FBF; mov     dh, [bx+19h] */
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+25))));
L_13FC2: ; /* IDA 0x13FC2; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13FC4: ; /* IDA 0x13FC4; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13FC6: ; /* IDA 0x13FC6; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_13FC8: ; /* IDA 0x13FC8; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13FCA: ; /* IDA 0x13FCA; add     di, bx */
    m->di=(uint16_t)(ec_math(m,m->di,m->bx,16,0,0));
L_13FCC: ; /* IDA 0x13FCC; push    si */
    ec_push(m,m->si);
L_13FCD: ; /* IDA 0x13FCD; push    bp */
    ec_push(m,m->bp);
L_13FCE: ; /* IDA 0x13FCE; add     si, 85Ch */
    m->si=(uint16_t)(ec_math(m,m->si,2140,16,0,0));
L_13FD2: ; /* IDA 0x13FD2; mov     cl, 4 */
    vg_cl(m,(uint8_t)(4));
L_13FD4: ; /* IDA 0x13FD4; mov     bh, [si] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
L_13FD6: ; /* IDA 0x13FD6; cmp     bh, 0FFh */
    ec_cmp(m,(uint8_t)(m->bx>>8),255,8);
L_13FD9: ; /* IDA 0x13FD9; jz      short loc_14025 */
    if(m->flags&0x40) goto L_14025;
L_13FDB: ; /* IDA 0x13FDB; push    di */
    ec_push(m,m->di);
L_13FDC: ; /* IDA 0x13FDC; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_13FDE: ; /* IDA 0x13FDE; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13FE0: ; /* IDA 0x13FE0; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13FE2: ; /* IDA 0x13FE2; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_13FE4: ; /* IDA 0x13FE4; mov     al, [bx+841h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+2113))));
L_13FE8: ; /* IDA 0x13FE8; cmp     dl, al */
    ec_cmp(m,(uint8_t)m->dx,(uint8_t)m->ax,8);
L_13FEA: ; /* IDA 0x13FEA; jz      short loc_1401F */
    if(m->flags&0x40) goto L_1401F;
L_13FEC: ; /* IDA 0x13FEC; cmp     al, 18h */
#if KI_TICK_MUTATION == 9
    ec_cmp(m,(uint8_t)m->ax,0x17,8);
#else
    ec_cmp(m,(uint8_t)m->ax,24,8);
#endif
L_13FEE: ; /* IDA 0x13FEE; jz      short loc_14003 */
    if(m->flags&0x40) goto L_14003;
L_13FF0: ; /* IDA 0x13FF0; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_13FF2: ; /* IDA 0x13FF2; add     di, ax */
    m->di=(uint16_t)(ec_math(m,m->di,m->ax,16,0,0));
L_13FF4: ; /* IDA 0x13FF4; cmp     byte ptr [di+600h], 80h */
#if KI_TICK_MUTATION == 10
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+0x600)),0x81,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+1536)),128,8);
#endif
L_13FF9: ; /* IDA 0x13FF9; jnb     short loc_1401F */
    if(!(m->flags&1)) goto L_1401F;
L_13FFB: ; /* IDA 0x13FFB; mov     byte ptr [bp+0], 0FEh */
    dl_write8(m,m->ss,(uint16_t)(m->bp),254);
L_13FFF: ; /* IDA 0x13FFF; add     ch, [bx+858h] */
    vg_ch(m,(uint8_t)(ec_math(m,(uint8_t)(m->cx>>8),dl_read8(m,m->ds,(uint16_t)(m->bx+2136)),8,0,0)));
L_14003: ; /* IDA 0x14003; cmp     dh, al */
    ec_cmp(m,(uint8_t)(m->dx>>8),(uint8_t)m->ax,8);
L_14005: ; /* IDA 0x14005; jnz     short loc_1401F */
    if(!(m->flags&0x40)) goto L_1401F;
L_14007: ; /* IDA 0x14007; mov     al, [si] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
L_14009: ; /* IDA 0x14009; mov     [bp+0], al */
    dl_write8(m,m->ss,(uint16_t)(m->bp),(uint8_t)m->ax);
L_1400C: ; /* IDA 0x1400C; mov     al, [bx+841h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+2113))));
L_14010: ; /* IDA 0x14010; mov     [bp+1], al */
    dl_write8(m,m->ss,(uint16_t)(m->bp+1),(uint8_t)m->ax);
L_14013: ; /* IDA 0x14013; mov     al, [bx+858h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+2136))));
L_14017: ; /* IDA 0x14017; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_14019: ; /* IDA 0x14019; mov     [bp+2], al */
    dl_write8(m,m->ss,(uint16_t)(m->bp+2),(uint8_t)m->ax);
L_1401C: ; /* IDA 0x1401C; add     bp, 4 */
    m->bp=(uint16_t)(ec_math(m,m->bp,4,16,0,0));
L_1401F: ; /* IDA 0x1401F; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14020: ; /* IDA 0x14020; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_14021: ; /* IDA 0x14021; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_14023: ; /* IDA 0x14023; jnz     short loc_13FD4 */
    if(!(m->flags&0x40)) goto L_13FD4;
L_14025: ; /* IDA 0x14025; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_14026: ; /* IDA 0x14026; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14027: ; /* IDA 0x14027; retn */
    return;
}
/* Original sub_14028; IDA 0x14028. */
static void ct_body_14028(KiMachine16 *m,const KiEconomyHooks *h) {
L_14028: ; /* IDA 0x14028; and     byte ptr [si+840h], 3Fh */
#if KI_TICK_MUTATION == 11
    {uint8_t result=dl_read8(m,m->ds,(uint16_t)(m->si+0x840))&0x7f;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si+0x840),result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si+2112)))&(63))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si+2112),result);}
#endif
L_1402D: ; /* IDA 0x1402D; mov     al, 80h */
    vg_al(m,(uint8_t)(128));
L_1402F: ; /* IDA 0x1402F; cmp     byte ptr [bp+0], 0FEh */
    ec_cmp(m,dl_read8(m,m->ss,(uint16_t)(m->bp)),254,8);
L_14033: ; /* IDA 0x14033; jb      short loc_1403E */
    if(m->flags&1) goto L_1403E;
L_14035: ; /* IDA 0x14035; jz      short loc_14040 */
    if(m->flags&0x40) goto L_14040;
L_14037: ; /* IDA 0x14037; mov     byte ptr [si+857h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+2135),0);
L_1403C: ; /* IDA 0x1403C; jmp     short loc_14053 */
    goto L_14053;
L_1403E: ; /* IDA 0x1403E; or      al, 40h */
    {uint16_t result=(((uint8_t)m->ax)|(64))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_14040: ; /* IDA 0x14040; or      [si+840h], al */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si+2112)))|((uint8_t)m->ax))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si+2112),result);}
L_14044: ; /* IDA 0x14044; cmp     byte ptr [si+858h], 1 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+2136)),1,8);
L_14049: ; /* IDA 0x14049; jnb     short loc_14055 */
    if(!(m->flags&1)) goto L_14055;
L_1404B: ; /* IDA 0x1404B; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_1404D: ; /* IDA 0x1404D; call    sub_140C9 */
    ct_call(m,16585,16464,h);
L_14050: ; /* IDA 0x14050; call    sub_140B3 */
    ct_call(m,16563,16467,h);
L_14053: ; /* IDA 0x14053; clc */
    m->flags&=(uint16_t)~1u;
L_14054: ; /* IDA 0x14054; retn */
    return;
L_14055: ; /* IDA 0x14055; stc */
    m->flags|=1;
L_14056: ; /* IDA 0x14056; retn */
    return;
}
/* Original sub_14057; IDA 0x14057. */
static void ct_body_14057(KiMachine16 *m,const KiEconomyHooks *h) {
L_14057: ; /* IDA 0x14057; cmp     byte ptr [bp+0], 0FEh */
    ec_cmp(m,dl_read8(m,m->ss,(uint16_t)(m->bp)),254,8);
L_1405B: ; /* IDA 0x1405B; jnb     short locret_140B2 */
    if(!(m->flags&1)) goto L_140B2;
L_1405D: ; /* IDA 0x1405D; call    sub_1ECE0 */
    ct_call(m,60640,16480,h);
L_14060: ; /* IDA 0x14060; and     al, 3 */
    {uint16_t result=(((uint8_t)m->ax)&(3))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_14062: ; /* IDA 0x14062; mov     di, bp */
    m->di=(uint16_t)(m->bp);
L_14064: ; /* IDA 0x14064; cmp     byte ptr ss:[di], 0FEh */
    ec_cmp(m,dl_read8(m,m->ss,(uint16_t)(m->di)),254,8);
L_14068: ; /* IDA 0x14068; jnb     short loc_14062 */
    if(!(m->flags&1)) goto L_14062;
L_1406A: ; /* IDA 0x1406A; dec     al */
    vg_al(m,(uint8_t)(st_dec(m,(uint8_t)m->ax,8)));
L_1406C: ; /* IDA 0x1406C; jz      short loc_14073 */
    if(m->flags&0x40) goto L_14073;
L_1406E: ; /* IDA 0x1406E; add     di, 4 */
    m->di=(uint16_t)(ec_math(m,m->di,4,16,0,0));
L_14071: ; /* IDA 0x14071; jmp     short loc_14064 */
    goto L_14064;
L_14073: ; /* IDA 0x14073; cmp     byte ptr [si+858h], 1 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+2136)),1,8);
L_14078: ; /* IDA 0x14078; ja      short loc_14099 */
    if(!(m->flags&(1|0x40))) goto L_14099;
L_1407A: ; /* IDA 0x1407A; mov     al, [si+854h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2132))));
L_1407E: ; /* IDA 0x1407E; add     al, 2 */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,2,8,0,0)));
L_14080: ; /* IDA 0x14080; sub     al, [si+858h] */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+2136)),8,0,1)));
L_14084: ; /* IDA 0x14084; jbe     short locret_140B2 */
    if(m->flags&(1|0x40)) goto L_140B2;
L_14086: ; /* IDA 0x14086; mov     cl, cs:byte_10CFF */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_1408B: ; /* IDA 0x1408B; cmp     cl, [si+841h] */
    ec_cmp(m,(uint8_t)m->cx,dl_read8(m,m->ds,(uint16_t)(m->si+2113)),8);
L_1408F: ; /* IDA 0x1408F; jz      short locret_140B2 */
    if(m->flags&0x40) goto L_140B2;
L_14091: ; /* IDA 0x14091; call    sub_140C9 */
    ct_call(m,16585,16532,h);
L_14094: ; /* IDA 0x14094; call    sub_140B3 */
    ct_call(m,16563,16535,h);
L_14097: ; /* IDA 0x14097; jmp     short locret_140B2 */
    goto L_140B2;
L_14099: ; /* IDA 0x14099; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1409B: ; /* IDA 0x1409B; jnz     short loc_1409F */
    if(!(m->flags&0x40)) goto L_1409F;
L_1409D: ; /* IDA 0x1409D; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_1409F: ; /* IDA 0x1409F; mov     cl, ss:[di] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->di))));
L_140A2: ; /* IDA 0x140A2; mov     ch, 0 */
    vg_ch(m,(uint8_t)(0));
L_140A4: ; /* IDA 0x140A4; mov     dl, al */
    vg_dl(m,(uint8_t)((uint8_t)m->ax));
L_140A6: ; /* IDA 0x140A6; mov     dh, [si+858h] */
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2136))));
L_140AA: ; /* IDA 0x140AA; call    sub_14155 */
    ct_call(m,16725,16557,h);
L_140AD: ; /* IDA 0x140AD; mov     byte ptr [si+857h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+2135),0);
L_140B2: ; /* IDA 0x140B2; retn */
    return;
}
/* Original sub_140B3; IDA 0x140B3. */
static void ct_body_140B3(KiMachine16 *m,const KiEconomyHooks *h) {
L_140B3: ; /* IDA 0x140B3; mov     bh, [si+841h] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2113))));
L_140B7: ; /* IDA 0x140B7; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_140B9: ; /* IDA 0x140B9; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_140BB: ; /* IDA 0x140BB; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_140BD: ; /* IDA 0x140BD; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_140BF: ; /* IDA 0x140BF; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_140C1: ; /* IDA 0x140C1; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_140C3: ; /* IDA 0x140C3; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_140C5: ; /* IDA 0x140C5; mov     [bx+16h], ah */
    dl_write8(m,m->ds,(uint16_t)(m->bx+22),(uint8_t)(m->ax>>8));
L_140C8: ; /* IDA 0x140C8; retn */
    return;
}
/* Original sub_140C9; IDA 0x140C9. */
static void ct_body_140C9(KiMachine16 *m,const KiEconomyHooks *h) {
L_140C9: ; /* IDA 0x140C9; cmp     byte ptr [si+857h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+2135)),0,8);
L_140CE: ; /* IDA 0x140CE; jz      short loc_140D1 */
    if(m->flags&0x40) goto L_140D1;
L_140D0: ; /* IDA 0x140D0; retn */
    return;
L_140D1: ; /* IDA 0x140D1; push    si */
    ec_push(m,m->si);
L_140D2: ; /* IDA 0x140D2; mov     cl, cs:byte_10CFF */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_140D7: ; /* IDA 0x140D7; cmp     cl, [si+841h] */
    ec_cmp(m,(uint8_t)m->cx,dl_read8(m,m->ds,(uint16_t)(m->si+2113)),8);
L_140DB: ; /* IDA 0x140DB; jnz     short loc_140FF */
    if(!(m->flags&0x40)) goto L_140FF;
L_140DD: ; /* IDA 0x140DD; mov     di, si */
    m->di=(uint16_t)(m->si);
L_140DF: ; /* IDA 0x140DF; add     di, 840h */
    m->di=(uint16_t)(ec_math(m,m->di,2112,16,0,0));
L_140E3: ; /* IDA 0x140E3; push    di */
    ec_push(m,m->di);
L_140E4: ; /* IDA 0x140E4; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_140E6: ; /* IDA 0x140E6; call    sub_10CDE */
    ct_call(m,3294,16617,h);
L_140E9: ; /* IDA 0x140E9; mov     cx, 26h ; '&' */
    m->cx=(uint16_t)(38);
L_140EC: ; /* IDA 0x140EC; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_140EE: ; /* IDA 0x140EE; call    sub_18810 */
    ct_call(m,34832,16625,h);
L_140F1: ; /* IDA 0x140F1; add     sp, 2 */
    m->sp=(uint16_t)(ec_math(m,m->sp,2,16,0,0));
L_140F4: ; /* IDA 0x140F4; mov     di, si */
    m->di=(uint16_t)(m->si);
L_140F6: ; /* IDA 0x140F6; call    sub_1ECE0 */
    ct_call(m,60640,16633,h);
L_140F9: ; /* IDA 0x140F9; and     al, 0Fh */
    {uint16_t result=(((uint8_t)m->ax)&(15))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_140FB: ; /* IDA 0x140FB; add     al, 18h */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,24,8,0,0)));
L_140FD: ; /* IDA 0x140FD; jmp     short loc_1414F */
    goto L_1414F;
L_140FF: ; /* IDA 0x140FF; mov     di, si */
    m->di=(uint16_t)(m->si);
L_14101: ; /* IDA 0x14101; mov     cx, si */
    m->cx=(uint16_t)(m->si);
L_14103: ; /* IDA 0x14103; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_14105: ; /* IDA 0x14105; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_14107: ; /* IDA 0x14107; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_14109: ; /* IDA 0x14109; mov     cl, al */
    vg_cl(m,(uint8_t)((uint8_t)m->ax));
L_1410B: ; /* IDA 0x1410B; mov     bh, [si+841h] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2113))));
L_1410F: ; /* IDA 0x1410F; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_14111: ; /* IDA 0x14111; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14113: ; /* IDA 0x14113; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14115: ; /* IDA 0x14115; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_14117: ; /* IDA 0x14117; call    sub_14575 */
    ct_call(m,17781,16666,h);
L_1411A: ; /* IDA 0x1411A; jb      short loc_14153 */
    if(m->flags&1) goto L_14153;
L_1411C: ; /* IDA 0x1411C; mov     bh, [si+3] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+3))));
L_1411F: ; /* IDA 0x1411F; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_14121: ; /* IDA 0x14121; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14123: ; /* IDA 0x14123; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14125: ; /* IDA 0x14125; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14127: ; /* IDA 0x14127; mov     ax, [di+848h] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+2120)));
L_1412B: ; /* IDA 0x1412B; sub     ax, [bx+848h] */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->ds,(uint16_t)(m->bx+2120)),16,0,1));
L_1412F: ; /* IDA 0x1412F; jnb     short loc_14133 */
    if(!(m->flags&1)) goto L_14133;
L_14131: ; /* IDA 0x14131; neg     ax */
    m->ax=(uint16_t)(ec_math(m,0,m->ax,16,0,1));
L_14133: ; /* IDA 0x14133; mov     dx, [di+84Ah] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+2122)));
L_14137: ; /* IDA 0x14137; sub     dx, [bx+848h] */
#if KI_TICK_MUTATION == 15
    m->dx=ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->bx+0x84a)),16,0,1);
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->bx+2120)),16,0,1));
#endif
L_1413B: ; /* IDA 0x1413B; jnb     short loc_1413F */
    if(!(m->flags&1)) goto L_1413F;
L_1413D: ; /* IDA 0x1413D; neg     dx */
    m->dx=(uint16_t)(ec_math(m,0,m->dx,16,0,1));
L_1413F: ; /* IDA 0x1413F; add     ax, dx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->dx,16,0,0));
L_14141: ; /* IDA 0x14141; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14143: ; /* IDA 0x14143; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14145: ; /* IDA 0x14145; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14147: ; /* IDA 0x14147; cmp     ax, 1Eh */
    ec_cmp(m,m->ax,30,16);
L_1414A: ; /* IDA 0x1414A; jbe     short loc_1414F */
    if(m->flags&(1|0x40)) goto L_1414F;
L_1414C: ; /* IDA 0x1414C; mov     ax, 1Eh */
    m->ax=(uint16_t)(30);
L_1414F: ; /* IDA 0x1414F; mov     [di+857h], al */
    dl_write8(m,m->ds,(uint16_t)(m->di+2135),(uint8_t)m->ax);
L_14153: ; /* IDA 0x14153; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14154: ; /* IDA 0x14154; retn */
    return;
}
/* Original sub_14155; IDA 0x14155. */
static void ct_body_14155(KiMachine16 *m,const KiEconomyHooks *h) {
L_14155: ; /* IDA 0x14155; sub     dh, dl */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),(uint8_t)m->dx,8,0,1)));
L_14157: ; /* IDA 0x14157; mov     di, 2240h */
    m->di=(uint16_t)(8768);
L_1415A: ; /* IDA 0x1415A; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_1415C: ; /* IDA 0x1415C; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1415E: ; /* IDA 0x1415E; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14160: ; /* IDA 0x14160; cmp     bx, [di+0Eh] */
    ec_cmp(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->di+14)),16);
L_14163: ; /* IDA 0x14163; jnz     short loc_1418E */
    if(!(m->flags&0x40)) goto L_1418E;
L_14165: ; /* IDA 0x14165; cmp     byte ptr [di], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di)),128,8);
L_14168: ; /* IDA 0x14168; jb      short loc_1418E */
    if(m->flags&1) goto L_1418E;
L_1416A: ; /* IDA 0x1416A; and     dh, dh */
    {uint16_t result=(((uint8_t)(m->dx>>8))&((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);vg_dh(m,(uint8_t)(result));}
L_1416C: ; /* IDA 0x1416C; jz      short loc_14179 */
    if(m->flags&0x40) goto L_14179;
L_1416E: ; /* IDA 0x1416E; call    sub_1ECE0 */
    ct_call(m,60640,16753,h);
L_14171: ; /* IDA 0x14171; cmp     al, 40h ; '@' */
    ec_cmp(m,(uint8_t)m->ax,64,8);
L_14173: ; /* IDA 0x14173; jnb     short loc_14179 */
    if(!(m->flags&1)) goto L_14179;
L_14175: ; /* IDA 0x14175; dec     dh */
    vg_dh(m,(uint8_t)(st_dec(m,(uint8_t)(m->dx>>8),8)));
L_14177: ; /* IDA 0x14177; jmp     short loc_1418E */
    goto L_1418E;
L_14179: ; /* IDA 0x14179; test    byte ptr [di], 4 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->di)))&(4))&255;ec_logic(m,result,8);}
L_1417C: ; /* IDA 0x1417C; jz      short loc_1418A */
    if(m->flags&0x40) goto L_1418A;
L_1417E: ; /* IDA 0x1417E; cmp     byte ptr [di+23h], 8 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+35)),8,8);
L_14182: ; /* IDA 0x14182; jnb     short loc_1418A */
    if(!(m->flags&1)) goto L_1418A;
L_14184: ; /* IDA 0x14184; mov     [di+20h], cl */
    dl_write8(m,m->ds,(uint16_t)(m->di+32),(uint8_t)m->cx);
L_14187: ; /* IDA 0x14187; mov     [di+23h], ch */
    dl_write8(m,m->ds,(uint16_t)(m->di+35),(uint8_t)(m->cx>>8));
L_1418A: ; /* IDA 0x1418A; dec     dl */
    vg_dl(m,(uint8_t)(st_dec(m,(uint8_t)m->dx,8)));
L_1418C: ; /* IDA 0x1418C; jz      short locret_14193 */
    if(m->flags&0x40) goto L_14193;
L_1418E: ; /* IDA 0x1418E; add     di, 40h ; '@' */
    m->di=(uint16_t)(ec_math(m,m->di,64,16,0,0));
L_14191: ; /* IDA 0x14191; jmp     short loc_14160 */
    goto L_14160;
L_14193: ; /* IDA 0x14193; retn */
    return;
}
/* Original sub_14194; IDA 0x14194. */
static void ct_body_14194(KiMachine16 *m,const KiEconomyHooks *h) {
L_14194: ; /* IDA 0x14194; mov     cl, 8 */
#if KI_TICK_MUTATION == 12
    vg_cl(m,7);
#else
    vg_cl(m,(uint8_t)(8));
#endif
L_14196: ; /* IDA 0x14196; mov     dl, 4 */
    vg_dl(m,(uint8_t)(4));
L_14198: ; /* IDA 0x14198; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_1419C: ; /* IDA 0x1419C; cmp     al, [si+841h] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+2113)),8);
L_141A0: ; /* IDA 0x141A0; jnz     short loc_141CC */
    if(!(m->flags&0x40)) goto L_141CC;
L_141A2: ; /* IDA 0x141A2; mov     cl, 5 */
    vg_cl(m,(uint8_t)(5));
L_141A4: ; /* IDA 0x141A4; mov     dl, 1 */
    vg_dl(m,(uint8_t)(1));
L_141A6: ; /* IDA 0x141A6; mov     bh, [si+859h] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2137))));
L_141AA: ; /* IDA 0x141AA; cmp     bh, 0FFh */
    ec_cmp(m,(uint8_t)(m->bx>>8),255,8);
L_141AD: ; /* IDA 0x141AD; jz      short loc_141CC */
    if(m->flags&0x40) goto L_141CC;
L_141AF: ; /* IDA 0x141AF; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_141B1: ; /* IDA 0x141B1; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_141B3: ; /* IDA 0x141B3; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_141B5: ; /* IDA 0x141B5; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_141B7: ; /* IDA 0x141B7; cmp     byte ptr [bx+425Ah], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+16986)),0,8);
L_141BC: ; /* IDA 0x141BC; jz      short loc_141CC */
    if(m->flags&0x40) goto L_141CC;
L_141BE: ; /* IDA 0x141BE; dec     byte ptr [bx+425Ah] */
    dl_write8(m,m->ds,(uint16_t)(m->bx+16986),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->bx+16986)),8));
L_141C2: ; /* IDA 0x141C2; add     cl, [bx+4253h] */
    vg_cl(m,(uint8_t)(ec_math(m,(uint8_t)m->cx,dl_read8(m,m->ds,(uint16_t)(m->bx+16979)),8,0,0)));
L_141C6: ; /* IDA 0x141C6; add     dl, [bx+4251h] */
    vg_dl(m,(uint8_t)(ec_math(m,(uint8_t)m->dx,dl_read8(m,m->ds,(uint16_t)(m->bx+16977)),8,0,0)));
L_141CA: ; /* IDA 0x141CA; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_141CC: ; /* IDA 0x141CC; mov     ch, cl */
    vg_ch(m,(uint8_t)((uint8_t)m->cx));
L_141CE: ; /* IDA 0x141CE; sub     ch, 0Fh */
    vg_ch(m,(uint8_t)(ec_math(m,(uint8_t)(m->cx>>8),15,8,0,1)));
L_141D1: ; /* IDA 0x141D1; ja      short loc_141D5 */
    if(!(m->flags&(1|0x40))) goto L_141D5;
L_141D3: ; /* IDA 0x141D3; mov     ch, 1 */
    vg_ch(m,(uint8_t)(1));
L_141D5: ; /* IDA 0x141D5; call    sub_1ECE0 */
    ct_call(m,60640,16856,h);
L_141D8: ; /* IDA 0x141D8; and     al, 0Fh */
    {uint16_t result=(((uint8_t)m->ax)&(15))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_141DA: ; /* IDA 0x141DA; cmp     cl, al */
    ec_cmp(m,(uint8_t)m->cx,(uint8_t)m->ax,8);
L_141DC: ; /* IDA 0x141DC; jb      short loc_141EE */
    if(m->flags&1) goto L_141EE;
L_141DE: ; /* IDA 0x141DE; mov     al, [si+850h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2128))));
L_141E2: ; /* IDA 0x141E2; add     al, ch */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,(uint8_t)(m->cx>>8),8,0,0)));
L_141E4: ; /* IDA 0x141E4; cmp     al, 0C8h */
#if KI_TICK_MUTATION == 13
    ec_cmp(m,(uint8_t)m->ax,0xc7,8);
#else
    ec_cmp(m,(uint8_t)m->ax,200,8);
#endif
L_141E6: ; /* IDA 0x141E6; jbe     short loc_141EA */
    if(m->flags&(1|0x40)) goto L_141EA;
L_141E8: ; /* IDA 0x141E8; mov     al, 0C8h */
#if KI_TICK_MUTATION == 13
    vg_al(m,0xc7);
#else
    vg_al(m,(uint8_t)(200));
#endif
L_141EA: ; /* IDA 0x141EA; mov     [si+850h], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+2128),(uint8_t)m->ax);
L_141EE: ; /* IDA 0x141EE; call    sub_1ECE0 */
    ct_call(m,60640,16881,h);
L_141F1: ; /* IDA 0x141F1; and     al, 0Fh */
    {uint16_t result=(((uint8_t)m->ax)&(15))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_141F3: ; /* IDA 0x141F3; cmp     cl, al */
    ec_cmp(m,(uint8_t)m->cx,(uint8_t)m->ax,8);
L_141F5: ; /* IDA 0x141F5; jb      short loc_1420B */
    if(m->flags&1) goto L_1420B;
L_141F7: ; /* IDA 0x141F7; mov     al, [si+851h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2129))));
L_141FB: ; /* IDA 0x141FB; shr     ch, 1 */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),1,8,5)));
L_141FD: ; /* IDA 0x141FD; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_141FF: ; /* IDA 0x141FF; add     al, ch */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,(uint8_t)(m->cx>>8),8,0,0)));
L_14201: ; /* IDA 0x14201; cmp     al, 0C8h */
    ec_cmp(m,(uint8_t)m->ax,200,8);
L_14203: ; /* IDA 0x14203; jbe     short loc_14207 */
    if(m->flags&(1|0x40)) goto L_14207;
L_14205: ; /* IDA 0x14205; mov     al, 0C8h */
    vg_al(m,(uint8_t)(200));
L_14207: ; /* IDA 0x14207; mov     [si+851h], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+2129),(uint8_t)m->ax);
L_1420B: ; /* IDA 0x1420B; mov     cx, [si+852h] */
    m->cx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+2130)));
L_1420F: ; /* IDA 0x1420F; cmp     ch, cl */
    ec_cmp(m,(uint8_t)(m->cx>>8),(uint8_t)m->cx,8);
L_14211: ; /* IDA 0x14211; jnb     short loc_1422F */
    if(!(m->flags&1)) goto L_1422F;
L_14213: ; /* IDA 0x14213; call    sub_1ECE0 */
    ct_call(m,60640,16918,h);
L_14216: ; /* IDA 0x14216; cmp     al, 18h */
    ec_cmp(m,(uint8_t)m->ax,24,8);
L_14218: ; /* IDA 0x14218; jnb     short locret_14235 */
    if(!(m->flags&1)) goto L_14235;
L_1421A: ; /* IDA 0x1421A; sub     [si+850h], dl */
    dl_write8(m,m->ds,(uint16_t)(m->si+2128),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->si+2128)),(uint8_t)m->dx,8,0,1));
L_1421E: ; /* IDA 0x1421E; jnb     short loc_14225 */
    if(!(m->flags&1)) goto L_14225;
L_14220: ; /* IDA 0x14220; mov     byte ptr [si+850h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+2128),0);
L_14225: ; /* IDA 0x14225; add     ch, dl */
    vg_ch(m,(uint8_t)(ec_math(m,(uint8_t)(m->cx>>8),(uint8_t)m->dx,8,0,0)));
L_14227: ; /* IDA 0x14227; jnb     short loc_1422B */
    if(!(m->flags&1)) goto L_1422B;
L_14229: ; /* IDA 0x14229; mov     ch, 0FFh */
    vg_ch(m,(uint8_t)(255));
L_1422B: ; /* IDA 0x1422B; cmp     ch, cl */
    ec_cmp(m,(uint8_t)(m->cx>>8),(uint8_t)m->cx,8);
L_1422D: ; /* IDA 0x1422D; jbe     short loc_14231 */
    if(m->flags&(1|0x40)) goto L_14231;
L_1422F: ; /* IDA 0x1422F; mov     ch, cl */
    vg_ch(m,(uint8_t)((uint8_t)m->cx));
L_14231: ; /* IDA 0x14231; mov     [si+853h], ch */
    dl_write8(m,m->ds,(uint16_t)(m->si+2131),(uint8_t)(m->cx>>8));
L_14235: ; /* IDA 0x14235; retn */
    return;
}
/* Original sub_14269; IDA 0x14269. */
static void ct_body_14269(KiMachine16 *m,const KiEconomyHooks *h) {
L_14269: ; /* IDA 0x14269; push    ax */
    ec_push(m,m->ax);
L_1426A: ; /* IDA 0x1426A; mov     ah, [si+855h] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2133))));
L_1426E: ; /* IDA 0x1426E; mov     al, [si+851h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2129))));
L_14272: ; /* IDA 0x14272; sub     al, ah */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8,0,1)));
L_14274: ; /* IDA 0x14274; jb      short loc_1427C */
    if(m->flags&1) goto L_1427C;
L_14276: ; /* IDA 0x14276; mov     [si+851h], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+2129),(uint8_t)m->ax);
L_1427A: ; /* IDA 0x1427A; jmp     short loc_142A9 */
    goto L_142A9;
L_1427C: ; /* IDA 0x1427C; mov     byte ptr [si+851h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+2129),0);
L_14281: ; /* IDA 0x14281; neg     al */
    vg_al(m,(uint8_t)(ec_math(m,0,(uint8_t)m->ax,8,0,1)));
L_14283: ; /* IDA 0x14283; push    ax */
    ec_push(m,m->ax);
L_14284: ; /* IDA 0x14284; sub     [si+850h], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+2128),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->si+2128)),(uint8_t)m->ax,8,0,1));
L_14288: ; /* IDA 0x14288; jnb     short loc_1428F */
    if(!(m->flags&1)) goto L_1428F;
L_1428A: ; /* IDA 0x1428A; mov     byte ptr [si+850h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+2128),0);
L_1428F: ; /* IDA 0x1428F; mul     byte ptr [si+84Fh] */
#if KI_TICK_MUTATION == 14
    dl_mul(m,dl_read8(m,m->ds,(uint16_t)(m->si+0x84e)),8);
#else
    dl_mul(m,dl_read8(m,m->ds,(uint16_t)(m->si+2127)),8);
#endif
L_14293: ; /* IDA 0x14293; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14295: ; /* IDA 0x14295; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14297: ; /* IDA 0x14297; sub     [si+84Eh], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+2126),ec_math(m,dl_read16(m,m->ds,(uint16_t)(m->si+2126)),m->ax,16,0,1));
L_1429B: ; /* IDA 0x1429B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1429C: ; /* IDA 0x1429C; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_1429E: ; /* IDA 0x1429E; sub     [si+853h], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+2131),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->si+2131)),(uint8_t)m->ax,8,0,1));
L_142A2: ; /* IDA 0x142A2; jnb     short loc_142A9 */
    if(!(m->flags&1)) goto L_142A9;
L_142A4: ; /* IDA 0x142A4; mov     byte ptr [si+853h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+2131),0);
L_142A9: ; /* IDA 0x142A9; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_142AA: ; /* IDA 0x142AA; retn */
    return;
}
/* Original sub_14575; IDA 0x14575. */
static void ct_body_14575(KiMachine16 *m,const KiEconomyHooks *h) {
L_14575: ; /* IDA 0x14575; push    ax */
    ec_push(m,m->ax);
L_14576: ; /* IDA 0x14576; push    bx */
    ec_push(m,m->bx);
L_14577: ; /* IDA 0x14577; mov     bl, 5 */
    gw_bl(m,(uint8_t)(5));
L_14579: ; /* IDA 0x14579; mov     ax, [si+21h] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+33)));
L_1457C: ; /* IDA 0x1457C; cmp     ax, 0A0h */
    ec_cmp(m,m->ax,160,16);
L_1457F: ; /* IDA 0x1457F; jle     short loc_1458C */
    if(!ec_greater(m)) goto L_1458C;
L_14581: ; /* IDA 0x14581; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14583: ; /* IDA 0x14583; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14585: ; /* IDA 0x14585; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14587: ; /* IDA 0x14587; add     bl, 2 */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,2,8,0,0)));
L_1458A: ; /* IDA 0x1458A; mov     bl, ah */
    gw_bl(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1458C: ; /* IDA 0x1458C; sub     bl, [si+14h] */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,dl_read8(m,m->ds,(uint16_t)(m->si+20)),8,0,1)));
L_1458F: ; /* IDA 0x1458F; ja      short loc_14594 */
    if(!(m->flags&(1|0x40))) goto L_14594;
L_14591: ; /* IDA 0x14591; stc */
    m->flags|=1;
L_14592: ; /* IDA 0x14592; jmp     short loc_145BE */
    goto L_145BE;
L_14594: ; /* IDA 0x14594; cmp     cl, bl */
    ec_cmp(m,(uint8_t)m->cx,(uint8_t)m->bx,8);
L_14596: ; /* IDA 0x14596; jbe     short loc_1459A */
    if(m->flags&(1|0x40)) goto L_1459A;
L_14598: ; /* IDA 0x14598; mov     cl, bl */
    vg_cl(m,(uint8_t)((uint8_t)m->bx));
L_1459A: ; /* IDA 0x1459A; push    di */
    ec_push(m,m->di);
L_1459B: ; /* IDA 0x1459B; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_1459D: ; /* IDA 0x1459D; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1459F: ; /* IDA 0x1459F; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_145A1: ; /* IDA 0x145A1; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_145A3: ; /* IDA 0x145A3; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_145A5: ; /* IDA 0x145A5; call    sub_145C1 */
    ct_call(m,17857,17832,h);
L_145A8: ; /* IDA 0x145A8; jb      short loc_145B7 */
    if(m->flags&1) goto L_145B7;
L_145AA: ; /* IDA 0x145AA; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_145AC: ; /* IDA 0x145AC; mov     [di+20h], ch */
    dl_write8(m,m->ds,(uint16_t)(m->di+32),(uint8_t)(m->cx>>8));
L_145AF: ; /* IDA 0x145AF; mov     byte ptr [di+23h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->di+35),0);
L_145B3: ; /* IDA 0x145B3; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_145B5: ; /* IDA 0x145B5; jnz     short loc_145A5 */
    if(!(m->flags&0x40)) goto L_145A5;
L_145B7: ; /* IDA 0x145B7; clc */
    m->flags&=(uint16_t)~1u;
L_145B8: ; /* IDA 0x145B8; and     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))&((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_145BA: ; /* IDA 0x145BA; jnz     short loc_145BD */
    if(!(m->flags&0x40)) goto L_145BD;
L_145BC: ; /* IDA 0x145BC; stc */
    m->flags|=1;
L_145BD: ; /* IDA 0x145BD; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_145BE: ; /* IDA 0x145BE; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_145BF: ; /* IDA 0x145BF; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_145C0: ; /* IDA 0x145C0; retn */
    return;
}
/* Original sub_145C1; IDA 0x145C1. */
static void ct_body_145C1(KiMachine16 *m,const KiEconomyHooks *h) {
L_145C1: ; /* IDA 0x145C1; push    ax */
    ec_push(m,m->ax);
L_145C2: ; /* IDA 0x145C2; push    bx */
    ec_push(m,m->bx);
L_145C3: ; /* IDA 0x145C3; push    cx */
    ec_push(m,m->cx);
L_145C4: ; /* IDA 0x145C4; push    dx */
    ec_push(m,m->dx);
L_145C5: ; /* IDA 0x145C5; push    si */
    ec_push(m,m->si);
L_145C6: ; /* IDA 0x145C6; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_145C8: ; /* IDA 0x145C8; mov     di, 4240h */
    m->di=(uint16_t)(16960);
L_145CB: ; /* IDA 0x145CB; mov     cx, 7Fh */
    m->cx=(uint16_t)(127);
L_145CE: ; /* IDA 0x145CE; cmp     al, [di+1Ch] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->di+28)),8);
L_145D1: ; /* IDA 0x145D1; jnz     short loc_145E3 */
    if(!(m->flags&0x40)) goto L_145E3;
L_145D3: ; /* IDA 0x145D3; cmp     byte ptr [di+17h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+23)),0,8);
L_145D7: ; /* IDA 0x145D7; jnz     short loc_145E3 */
    if(!(m->flags&0x40)) goto L_145E3;
L_145D9: ; /* IDA 0x145D9; cmp     ah, [di+11h] */
    ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->ds,(uint16_t)(m->di+17)),8);
L_145DC: ; /* IDA 0x145DC; jnb     short loc_145E3 */
    if(!(m->flags&1)) goto L_145E3;
L_145DE: ; /* IDA 0x145DE; mov     ah, [di+11h] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+17))));
L_145E1: ; /* IDA 0x145E1; mov     si, di */
    m->si=(uint16_t)(m->di);
L_145E3: ; /* IDA 0x145E3; add     di, 20h ; ' ' */
    m->di=(uint16_t)(ec_math(m,m->di,32,16,0,0));
L_145E6: ; /* IDA 0x145E6; loop    loc_145CE */
    if(--m->cx) goto L_145CE;
L_145E8: ; /* IDA 0x145E8; and     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))&((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_145EA: ; /* IDA 0x145EA; jz      short loc_145F1 */
    if(m->flags&0x40) goto L_145F1;
L_145EC: ; /* IDA 0x145EC; call    sub_16E8F */
    ct_call(m,28303,17903,h);
L_145EF: ; /* IDA 0x145EF; jmp     short loc_145F2 */
    goto L_145F2;
L_145F1: ; /* IDA 0x145F1; stc */
    m->flags|=1;
L_145F2: ; /* IDA 0x145F2; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_145F3: ; /* IDA 0x145F3; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_145F4: ; /* IDA 0x145F4; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_145F5: ; /* IDA 0x145F5; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_145F6: ; /* IDA 0x145F6; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_145F7: ; /* IDA 0x145F7; retn */
    return;
}
/* Original sub_1EB11; IDA 0x1EB11. */
static void ct_body_1EB11(KiMachine16 *m,const KiEconomyHooks *h) {
L_1EB11: ; /* IDA 0x1EB11; push    cx */
    ec_push(m,m->cx);
L_1EB12: ; /* IDA 0x1EB12; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_1EB14: ; /* IDA 0x1EB14; mov     ah, ch */
    vg_ah(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_1EB16: ; /* IDA 0x1EB16; call    sub_1EB5E */
    ct_call(m,60254,60185,h);
L_1EB19: ; /* IDA 0x1EB19; jnb     short loc_1EB16 */
    if(!(m->flags&1)) goto L_1EB16;
L_1EB1B: ; /* IDA 0x1EB1B; in      al, 61h; PC/XT PPI port B bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(97)));
L_1EB1D: ; /* IDA 0x1EB1D; or      al, 3 */
#if KI_TICK_MUTATION == 16
    {uint8_t result=(uint8_t)m->ax|1;ec_logic(m,result,8);vg_al(m,result);}
#else
    {uint16_t result=(((uint8_t)m->ax)|(3))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
#endif
L_1EB1F: ; /* IDA 0x1EB1F; out     61h, al; PC/XT PPI port B bits: */
    dl_out(m,97,(uint8_t)m->ax,8);
L_1EB21: ; /* IDA 0x1EB21; call    sub_1EB5E */
    ct_call(m,60254,60196,h);
L_1EB24: ; /* IDA 0x1EB24; jb      short loc_1EB21 */
    if(m->flags&1) goto L_1EB21;
L_1EB26: ; /* IDA 0x1EB26; call    sub_1EB5E */
    ct_call(m,60254,60201,h);
L_1EB29: ; /* IDA 0x1EB29; jnb     short loc_1EB26 */
    if(!(m->flags&1)) goto L_1EB26;
L_1EB2B: ; /* IDA 0x1EB2B; call    sub_1EB5E */
    ct_call(m,60254,60206,h);
L_1EB2E: ; /* IDA 0x1EB2E; jb      short loc_1EB2B */
    if(m->flags&1) goto L_1EB2B;
L_1EB30: ; /* IDA 0x1EB30; call    sub_1EB5E */
    ct_call(m,60254,60211,h);
L_1EB33: ; /* IDA 0x1EB33; jnb     short loc_1EB30 */
    if(!(m->flags&1)) goto L_1EB30;
L_1EB35: ; /* IDA 0x1EB35; in      al, 61h; PC/XT PPI port B bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(97)));
L_1EB37: ; /* IDA 0x1EB37; and     al, 0FCh */
    {uint16_t result=(((uint8_t)m->ax)&(252))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1EB39: ; /* IDA 0x1EB39; out     61h, al; PC/XT PPI port B bits: */
    dl_out(m,97,(uint8_t)m->ax,8);
L_1EB3B: ; /* IDA 0x1EB3B; cmp     cl, 1 */
    ec_cmp(m,(uint8_t)m->cx,1,8);
L_1EB3E: ; /* IDA 0x1EB3E; jz      short loc_1EB5C */
    if(m->flags&0x40) goto L_1EB5C;
L_1EB40: ; /* IDA 0x1EB40; call    sub_1EB5E */
    ct_call(m,60254,60227,h);
L_1EB43: ; /* IDA 0x1EB43; jb      short loc_1EB40 */
    if(m->flags&1) goto L_1EB40;
L_1EB45: ; /* IDA 0x1EB45; call    sub_1EB5E */
    ct_call(m,60254,60232,h);
L_1EB48: ; /* IDA 0x1EB48; jnb     short loc_1EB45 */
    if(!(m->flags&1)) goto L_1EB45;
L_1EB4A: ; /* IDA 0x1EB4A; call    sub_1EB5E */
    ct_call(m,60254,60237,h);
L_1EB4D: ; /* IDA 0x1EB4D; jb      short loc_1EB4A */
    if(m->flags&1) goto L_1EB4A;
L_1EB4F: ; /* IDA 0x1EB4F; call    sub_1EB5E */
    ct_call(m,60254,60242,h);
L_1EB52: ; /* IDA 0x1EB52; jnb     short loc_1EB4F */
    if(!(m->flags&1)) goto L_1EB4F;
L_1EB54: ; /* IDA 0x1EB54; dec     ah */
    vg_ah(m,(uint8_t)(st_dec(m,(uint8_t)(m->ax>>8),8)));
L_1EB56: ; /* IDA 0x1EB56; jnz     short loc_1EB40 */
    if(!(m->flags&0x40)) goto L_1EB40;
L_1EB58: ; /* IDA 0x1EB58; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_1EB5A: ; /* IDA 0x1EB5A; jnz     short loc_1EB14 */
    if(!(m->flags&0x40)) goto L_1EB14;
L_1EB5C: ; /* IDA 0x1EB5C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1EB5D: ; /* IDA 0x1EB5D; retn */
    return;
}
/* Original sub_1EB5E; IDA 0x1EB5E. */
static void ct_body_1EB5E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1EB5E: ; /* IDA 0x1EB5E; push    dx */
    ec_push(m,m->dx);
L_1EB5F: ; /* IDA 0x1EB5F; mov     dx, 3DAh */
    m->dx=(uint16_t)(986);
L_1EB62: ; /* IDA 0x1EB62; in      al, dx; Video status bits: */
    vg_al(m,(uint8_t)(wolong_main_in8(m->dx)));
L_1EB63: ; /* IDA 0x1EB63; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1EB64: ; /* IDA 0x1EB64; test    al, 8 */
    {uint16_t result=(((uint8_t)m->ax)&(8))&255;ec_logic(m,result,8);}
L_1EB66: ; /* IDA 0x1EB66; jz      short loc_1EB6A */
    if(m->flags&0x40) goto L_1EB6A;
L_1EB68: ; /* IDA 0x1EB68; stc */
    m->flags|=1;
L_1EB69: ; /* IDA 0x1EB69; retn */
    return;
L_1EB6A: ; /* IDA 0x1EB6A; clc */
    m->flags&=(uint16_t)~1u;
L_1EB6B: ; /* IDA 0x1EB6B; retn */
    return;
}
int ki_tick_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0xcde:ct_body_10CDE(m,h);break;
case 0x2459:ct_body_12459(m,h);break;
case 0x248a:ct_body_1248A(m,h);break;
case 0x24ff:ct_body_124FF(m,h);break;
case 0x3efd:ct_body_13EFD(m,h);break;
case 0x3f74:ct_body_13F74(m,h);break;
case 0x3fa9:ct_body_13FA9(m,h);break;
case 0x4028:ct_body_14028(m,h);break;
case 0x4057:ct_body_14057(m,h);break;
case 0x40b3:ct_body_140B3(m,h);break;
case 0x40c9:ct_body_140C9(m,h);break;
case 0x4155:ct_body_14155(m,h);break;
case 0x4194:ct_body_14194(m,h);break;
case 0x4269:ct_body_14269(m,h);break;
case 0x4575:ct_body_14575(m,h);break;
case 0x45c1:ct_body_145C1(m,h);break;
case 0xeb11:ct_body_1EB11(m,h);break;
case 0xeb5e:ct_body_1EB5E(m,h);break;
default:return 0;}return 1;}
static int ct_known(uint16_t t) {switch(t) {case 0xcde: case 0x2459: case 0x248a: case 0x24ff: case 0x3efd: case 0x3f74: case 0x3fa9: case 0x4028: case 0x4057: case 0x40b3: case 0x40c9: case 0x4155: case 0x4194: case 0x4269: case 0x4575: case 0x45c1: case 0xeb11: case 0xeb5e: return 1;default:return 0;}}
void sub_10CDE(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0xcde,h);m->ip=ec_pop(m);}
void sub_12459(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x2459,h);m->ip=ec_pop(m);}
void sub_1248A(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x248a,h);m->ip=ec_pop(m);}
void sub_124FF(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x24ff,h);m->ip=ec_pop(m);}
void sub_13EFD(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x3efd,h);m->ip=ec_pop(m);}
void sub_13F74(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x3f74,h);m->ip=ec_pop(m);}
void sub_13FA9(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x3fa9,h);m->ip=ec_pop(m);}
void sub_14028(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x4028,h);m->ip=ec_pop(m);}
void sub_14057(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x4057,h);m->ip=ec_pop(m);}
void sub_140B3(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x40b3,h);m->ip=ec_pop(m);}
void sub_140C9(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x40c9,h);m->ip=ec_pop(m);}
void sub_14155(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x4155,h);m->ip=ec_pop(m);}
void sub_14194(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x4194,h);m->ip=ec_pop(m);}
void sub_14269(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x4269,h);m->ip=ec_pop(m);}
void sub_14575(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x4575,h);m->ip=ec_pop(m);}
void sub_145C1(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0x45c1,h);m->ip=ec_pop(m);}
void sub_1EB11(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0xeb11,h);m->ip=ec_pop(m);}
void sub_1EB5E(KiMachine16 *m,const KiEconomyHooks *h) {ki_tick_body(m,0xeb5e,h);m->ip=ec_pop(m);}
