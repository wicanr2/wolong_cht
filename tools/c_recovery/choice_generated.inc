/* Fixed IDA native C choice control; original locators; spec/232, re/112. */
/* Original sub_1036F; IDA 0x1036F. */
static void ch_body_1036F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1036F: ; /* IDA 0x1036F; sub     sp, 18h */
    m->sp=(uint16_t)(ec_math(m,m->sp,24,16,0,1));
L_10372: ; /* IDA 0x10372; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_10374: ; /* IDA 0x10374; call    sub_104B5 */
    ch_call(m,1205,887,h);
L_10377: ; /* IDA 0x10377; cmp     al, ch */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->cx>>8),8);
L_10379: ; /* IDA 0x10379; jnb     short loc_1037D */
    if(!(m->flags&1)) goto L_1037D;
L_1037B: ; /* IDA 0x1037B; mov     ch, al */
    vg_ch(m,(uint8_t)((uint8_t)m->ax));
L_1037D: ; /* IDA 0x1037D; mov     [bp+0], al */
    dl_write8(m,m->ss,(uint16_t)(m->bp),(uint8_t)m->ax);
L_10380: ; /* IDA 0x10380; dec     byte ptr [bp+0] */
    dl_write8(m,m->ss,(uint16_t)(m->bp),st_dec(m,dl_read8(m,m->ss,(uint16_t)(m->bp)),8));
L_10383: ; /* IDA 0x10383; mov     [bp+1], ch */
    dl_write8(m,m->ss,(uint16_t)(m->bp+1),(uint8_t)(m->cx>>8));
L_10386: ; /* IDA 0x10386; dec     byte ptr [bp+1] */
    dl_write8(m,m->ss,(uint16_t)(m->bp+1),st_dec(m,dl_read8(m,m->ss,(uint16_t)(m->bp+1)),8));
L_10389: ; /* IDA 0x10389; mov     cs:word_1036D, bx */
    dl_write16(m,m->cs,877,m->bx);
L_1038E: ; /* IDA 0x1038E; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_10390: ; /* IDA 0x10390; xchg    bl, dh */
    {uint16_t a=(uint8_t)m->bx,b=(uint8_t)(m->dx>>8);gw_bl(m,(uint8_t)(b));vg_dh(m,(uint8_t)(a));}
L_10392: ; /* IDA 0x10392; cmp     si, 0FFFFh */
    ec_cmp(m,m->si,65535,16);
L_10395: ; /* IDA 0x10395; jz      short loc_1039C */
    if(m->flags&0x40) goto L_1039C;
L_10397: ; /* IDA 0x10397; call    sub_103C3 */
    ch_call(m,963,922,h);
L_1039A: ; /* IDA 0x1039A; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1039B: ; /* IDA 0x1039B; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1039C: ; /* IDA 0x1039C; call    sub_103E6 */
    ch_call(m,998,927,h);
L_1039F: ; /* IDA 0x1039F; mov     al, cl */
    vg_al(m,(uint8_t)((uint8_t)m->cx));
L_103A1: ; /* IDA 0x103A1; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_103A3: ; /* IDA 0x103A3; mov     [bp+11h], cl */
    dl_write8(m,m->ss,(uint16_t)(m->bp+17),(uint8_t)m->cx);
L_103A6: ; /* IDA 0x103A6; mov     cl, 4 */
    vg_cl(m,(uint8_t)(4));
L_103A8: ; /* IDA 0x103A8; shl     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,4));
L_103AA: ; /* IDA 0x103AA; mov     [bp+12h], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp+18),m->ax);
L_103AD: ; /* IDA 0x103AD; shl     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,4));
L_103AF: ; /* IDA 0x103AF; shl     bx, cl */
    m->bx=(uint16_t)(dl_shift(m,m->bx,(uint8_t)m->cx,16,4));
L_103B1: ; /* IDA 0x103B1; add     dx, 8 */
    m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,0));
L_103B4: ; /* IDA 0x103B4; add     bx, 8 */
    m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,0));
L_103B7: ; /* IDA 0x103B7; call    sub_10414 */
    ch_call(m,1044,954,h);
L_103BA: ; /* IDA 0x103BA; lahf */
    vg_ah(m,(uint8_t)((m->flags&0xd5)|2));
L_103BB: ; /* IDA 0x103BB; call    sub_104FF */
    ch_call(m,1279,958,h);
L_103BE: ; /* IDA 0x103BE; add     sp, 18h */
    m->sp=(uint16_t)(ec_math(m,m->sp,24,16,0,0));
L_103C1: ; /* IDA 0x103C1; sahf */
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
L_103C2: ; /* IDA 0x103C2; retn */
    return;
}
/* Original sub_103C3; IDA 0x103C3. */
static void ch_body_103C3(KiMachine16 *m,const KiEconomyHooks *h) {
L_103C3: ; /* IDA 0x103C3; push    ax */
    ec_push(m,m->ax);
L_103C4: ; /* IDA 0x103C4; push    bx */
    ec_push(m,m->bx);
L_103C5: ; /* IDA 0x103C5; push    cx */
    ec_push(m,m->cx);
L_103C6: ; /* IDA 0x103C6; push    dx */
    ec_push(m,m->dx);
L_103C7: ; /* IDA 0x103C7; add     cl, 1 */
    vg_cl(m,(uint8_t)(ec_math(m,(uint8_t)m->cx,1,8,0,0)));
L_103CA: ; /* IDA 0x103CA; mov     ch, 3 */
    vg_ch(m,(uint8_t)(3));
L_103CC: ; /* IDA 0x103CC; call    sub_10BCD */
    ch_call(m,3021,975,h);
L_103CF: ; /* IDA 0x103CF; mov     cl, 4 */
    vg_cl(m,(uint8_t)(4));
L_103D1: ; /* IDA 0x103D1; shl     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,4));
L_103D3: ; /* IDA 0x103D3; shl     bx, cl */
    m->bx=(uint16_t)(dl_shift(m,m->bx,(uint8_t)m->cx,16,4));
L_103D5: ; /* IDA 0x103D5; add     dx, 8 */
    m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,0));
L_103D8: ; /* IDA 0x103D8; add     bx, 0Ch */
    m->bx=(uint16_t)(ec_math(m,m->bx,12,16,0,0));
L_103DB: ; /* IDA 0x103DB; mov     ax, 0F01h */
    m->ax=(uint16_t)(3841);
L_103DE: ; /* IDA 0x103DE; call    sub_1F6DC */
    ch_call(m,63196,993,h);
L_103E1: ; /* IDA 0x103E1; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_103E2: ; /* IDA 0x103E2; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_103E3: ; /* IDA 0x103E3; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_103E4: ; /* IDA 0x103E4; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_103E5: ; /* IDA 0x103E5; retn */
    return;
}
/* Original sub_103E6; IDA 0x103E6. */
static void ch_body_103E6(KiMachine16 *m,const KiEconomyHooks *h) {
L_103E6: ; /* IDA 0x103E6; push    ax */
    ec_push(m,m->ax);
L_103E7: ; /* IDA 0x103E7; push    bx */
    ec_push(m,m->bx);
L_103E8: ; /* IDA 0x103E8; push    cx */
    ec_push(m,m->cx);
L_103E9: ; /* IDA 0x103E9; push    dx */
    ec_push(m,m->dx);
L_103EA: ; /* IDA 0x103EA; push    cx */
    ec_push(m,m->cx);
L_103EB: ; /* IDA 0x103EB; add     cx, 101h */
    m->cx=(uint16_t)(ec_math(m,m->cx,257,16,0,0));
L_103EF: ; /* IDA 0x103EF; call    sub_10BCD */
    ch_call(m,3021,1010,h);
L_103F2: ; /* IDA 0x103F2; mov     cl, 4 */
    vg_cl(m,(uint8_t)(4));
L_103F4: ; /* IDA 0x103F4; shl     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,4));
L_103F6: ; /* IDA 0x103F6; shl     bx, cl */
    m->bx=(uint16_t)(dl_shift(m,m->bx,(uint8_t)m->cx,16,4));
L_103F8: ; /* IDA 0x103F8; add     dx, 8 */
    m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,0));
L_103FB: ; /* IDA 0x103FB; add     bx, 8 */
    m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,0));
L_103FE: ; /* IDA 0x103FE; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_103FF: ; /* IDA 0x103FF; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_10401: ; /* IDA 0x10401; call    cs:word_1036D */
    ch_call(m,dl_read16(m,m->cs,877),1030,h);
L_10406: ; /* IDA 0x10406; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_10409: ; /* IDA 0x10409; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_1040B: ; /* IDA 0x1040B; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_1040D: ; /* IDA 0x1040D; jnz     short loc_10401 */
    if(!(m->flags&0x40)) goto L_10401;
L_1040F: ; /* IDA 0x1040F; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10410: ; /* IDA 0x10410; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10411: ; /* IDA 0x10411; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10412: ; /* IDA 0x10412; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10413: ; /* IDA 0x10413; retn */
    return;
}
/* Original sub_10414; IDA 0x10414. */
static void ch_body_10414(KiMachine16 *m,const KiEconomyHooks *h) {
L_10414: ; /* IDA 0x10414; mov     [bp+0Ah], dx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+10),m->dx);
L_10417: ; /* IDA 0x10417; mov     [bp+0Ch], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+12),m->bx);
L_1041A: ; /* IDA 0x1041A; mov     byte ptr [bp+0Eh], 0 */
    dl_write8(m,m->ss,(uint16_t)(m->bp+14),0);
L_1041E: ; /* IDA 0x1041E; mov     byte ptr [bp+0Fh], 0 */
    dl_write8(m,m->ss,(uint16_t)(m->bp+15),0);
L_10422: ; /* IDA 0x10422; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_10424: ; /* IDA 0x10424; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_10426: ; /* IDA 0x10426; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_10428: ; /* IDA 0x10428; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_1042A: ; /* IDA 0x1042A; mov     bx, dx */
    m->bx=(uint16_t)(m->dx);
L_1042C: ; /* IDA 0x1042C; mov     dx, 50h ; 'P' */
    m->dx=(uint16_t)(80);
L_1042F: ; /* IDA 0x1042F; mul     dx */
    dl_mul(m,m->dx,16);
L_10431: ; /* IDA 0x10431; add     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
L_10433: ; /* IDA 0x10433; mov     [bp+14h], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+20),m->bx);
L_10436: ; /* IDA 0x10436; mov     al, [bp+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+1))));
L_10439: ; /* IDA 0x10439; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1043B: ; /* IDA 0x1043B; mov     dx, 500h */
    m->dx=(uint16_t)(1280);
L_1043E: ; /* IDA 0x1043E; mul     dx */
    dl_mul(m,m->dx,16);
L_10440: ; /* IDA 0x10440; add     bx, ax */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
L_10442: ; /* IDA 0x10442; add     bx, 4B0h */
    m->bx=(uint16_t)(ec_math(m,m->bx,1200,16,0,0));
L_10446: ; /* IDA 0x10446; mov     [bp+16h], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+22),m->bx);
L_10449: ; /* IDA 0x10449; call    sub_1061F */
    ch_call(m,1567,1100,h);
L_1044C: ; /* IDA 0x1044C; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_1044F: ; /* IDA 0x1044F; call    sub_20000 */
    mx_far(m,0x44f,0x454,h);
L_10454: ; /* IDA 0x10454; test    bx, 1 */
    {uint16_t result=((m->bx)&(1))&65535;ec_logic(m,result,16);}
L_10458: ; /* IDA 0x10458; jnz     short loc_10491 */
    if(!(m->flags&0x40)) goto L_10491;
L_1045A: ; /* IDA 0x1045A; test    bx, 2 */
    {uint16_t result=((m->bx)&(2))&65535;ec_logic(m,result,16);}
L_1045E: ; /* IDA 0x1045E; jnz     short loc_104A4 */
    if(!(m->flags&0x40)) goto L_104A4;
L_10460: ; /* IDA 0x10460; cmp     dx, 83h */
    ec_cmp(m,m->dx,131,16);
L_10464: ; /* IDA 0x10464; ja      short loc_1047F */
#if KI_CHOICE_MUTATION == 3
    if(!(m->flags&1)) goto L_1047F;
#else
    if(!(m->flags&(1|0x40))) goto L_1047F;
#endif
L_10466: ; /* IDA 0x10466; cmp     dx, 64h ; 'd' */
    ec_cmp(m,m->dx,100,16);
L_10469: ; /* IDA 0x10469; jnb     short loc_1044C */
    if(!(m->flags&1)) goto L_1044C;
L_1046B: ; /* IDA 0x1046B; ja      short loc_1047F */
    if(!(m->flags&(1|0x40))) goto L_1047F;
L_1046D: ; /* IDA 0x1046D; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_10470: ; /* IDA 0x10470; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_10472: ; /* IDA 0x10472; mov     dx, 83h */
    m->dx=(uint16_t)(131);
L_10475: ; /* IDA 0x10475; call    sub_20000 */
    mx_far(m,0x475,0x47a,h);
L_1047A: ; /* IDA 0x1047A; call    sub_1054D */
    ch_call(m,1357,1149,h);
L_1047D: ; /* IDA 0x1047D; jmp     short loc_1044C */
    goto L_1044C;
L_1047F: ; /* IDA 0x1047F; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_10482: ; /* IDA 0x10482; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_10484: ; /* IDA 0x10484; mov     dx, 64h ; 'd' */
    m->dx=(uint16_t)(100);
L_10487: ; /* IDA 0x10487; call    sub_20000 */
    mx_far(m,0x487,0x48c,h);
L_1048C: ; /* IDA 0x1048C; call    sub_1059B */
    ch_call(m,1435,1167,h);
L_1048F: ; /* IDA 0x1048F; jmp     short loc_1044C */
    goto L_1044C;
L_10491: ; /* IDA 0x10491; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_10494: ; /* IDA 0x10494; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_10496: ; /* IDA 0x10496; call    sub_20000 */
    mx_far(m,0x496,0x49b,h);
L_1049B: ; /* IDA 0x1049B; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_1049D: ; /* IDA 0x1049D; jz      short loc_1044C */
    if(m->flags&0x40) goto L_1044C;
L_1049F: ; /* IDA 0x1049F; mov     al, [bp+0Eh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+14))));
L_104A2: ; /* IDA 0x104A2; clc */
#if KI_CHOICE_MUTATION == 5
    m->flags|=1;
#else
    m->flags&=(uint16_t)~1u;
#endif
L_104A3: ; /* IDA 0x104A3; retn */
    return;
L_104A4: ; /* IDA 0x104A4; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_104A7: ; /* IDA 0x104A7; mov     bx, 1 */
    m->bx=(uint16_t)(1);
L_104AA: ; /* IDA 0x104AA; call    sub_20000 */
    mx_far(m,0x4aa,0x4af,h);
L_104AF: ; /* IDA 0x104AF; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_104B1: ; /* IDA 0x104B1; jz      short loc_1044C */
    if(m->flags&0x40) goto L_1044C;
L_104B3: ; /* IDA 0x104B3; stc */
#if KI_CHOICE_MUTATION == 6
    m->flags&=(uint16_t)~1u;
#else
    m->flags|=1;
#endif
L_104B4: ; /* IDA 0x104B4; retn */
    return;
}
/* Original sub_104B5; IDA 0x104B5. */
static void ch_body_104B5(KiMachine16 *m,const KiEconomyHooks *h) {
L_104B5: ; /* IDA 0x104B5; push    ax */
    ec_push(m,m->ax);
L_104B6: ; /* IDA 0x104B6; push    bx */
    ec_push(m,m->bx);
L_104B7: ; /* IDA 0x104B7; push    cx */
    ec_push(m,m->cx);
L_104B8: ; /* IDA 0x104B8; push    dx */
    ec_push(m,m->dx);
L_104B9: ; /* IDA 0x104B9; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_104BC: ; /* IDA 0x104BC; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_104BE: ; /* IDA 0x104BE; call    sub_20000 */
    mx_far(m,0x4be,0x4c3,h);
L_104C3: ; /* IDA 0x104C3; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_104C6: ; /* IDA 0x104C6; mov     bx, 1 */
    m->bx=(uint16_t)(1);
L_104C9: ; /* IDA 0x104C9; call    sub_20000 */
    mx_far(m,0x4c9,0x4ce,h);
L_104CE: ; /* IDA 0x104CE; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_104D1: ; /* IDA 0x104D1; call    sub_20000 */
    mx_far(m,0x4d1,0x4d6,h);
L_104D6: ; /* IDA 0x104D6; mov     [bp+6], dx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+6),m->dx);
L_104D9: ; /* IDA 0x104D9; mov     [bp+8], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+8),m->bx);
L_104DC: ; /* IDA 0x104DC; mov     [bp+10h], al */
    dl_write8(m,m->ss,(uint16_t)(m->bp+16),(uint8_t)m->ax);
L_104DF: ; /* IDA 0x104DF; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_104E2: ; /* IDA 0x104E2; call    sub_20000 */
    mx_far(m,0x4e2,0x4e7,h);
L_104E7: ; /* IDA 0x104E7; mov     [bp+2], cx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->cx);
L_104EA: ; /* IDA 0x104EA; mov     [bp+4], dx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+4),m->dx);
L_104ED: ; /* IDA 0x104ED; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_104F0: ; /* IDA 0x104F0; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_104F2: ; /* IDA 0x104F2; mov     dx, 64h ; 'd' */
    m->dx=(uint16_t)(100);
L_104F5: ; /* IDA 0x104F5; call    sub_20000 */
    mx_far(m,0x4f5,0x4fa,h);
L_104FA: ; /* IDA 0x104FA; pop     dx */
#if KI_CHOICE_MUTATION == 8
    ; /* lost old DX restore */
#else
    m->dx=(uint16_t)(ec_pop(m));
#endif
L_104FB: ; /* IDA 0x104FB; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_104FC: ; /* IDA 0x104FC; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_104FD: ; /* IDA 0x104FD; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_104FE: ; /* IDA 0x104FE; retn */
    return;
}
/* Original sub_104FF; IDA 0x104FF. */
static void ch_body_104FF(KiMachine16 *m,const KiEconomyHooks *h) {
L_104FF: ; /* IDA 0x104FF; push    ax */
    ec_push(m,m->ax);
L_10500: ; /* IDA 0x10500; push    bx */
    ec_push(m,m->bx);
L_10501: ; /* IDA 0x10501; push    cx */
    ec_push(m,m->cx);
L_10502: ; /* IDA 0x10502; push    dx */
    ec_push(m,m->dx);
L_10503: ; /* IDA 0x10503; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_10506: ; /* IDA 0x10506; mov     cx, [bp+2] */
    m->cx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_10509: ; /* IDA 0x10509; mov     dx, [bp+4] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+4)));
L_1050C: ; /* IDA 0x1050C; call    sub_20000 */
    mx_far(m,0x50c,0x511,h);
L_10511: ; /* IDA 0x10511; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_10514: ; /* IDA 0x10514; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_10516: ; /* IDA 0x10516; call    sub_20000 */
    mx_far(m,0x516,0x51b,h);
L_1051B: ; /* IDA 0x1051B; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_1051E: ; /* IDA 0x1051E; mov     bx, 1 */
    m->bx=(uint16_t)(1);
L_10521: ; /* IDA 0x10521; call    sub_20000 */
    mx_far(m,0x521,0x526,h);
L_10526: ; /* IDA 0x10526; cmp     byte ptr [bp+10h], 0 */
    ec_cmp(m,dl_read8(m,m->ss,(uint16_t)(m->bp+16)),0,8);
L_1052A: ; /* IDA 0x1052A; jl      short loc_10548 */
    if(ec_less(m)) goto L_10548;
L_1052C: ; /* IDA 0x1052C; jz      short loc_1053E */
    if(m->flags&0x40) goto L_1053E;
L_1052E: ; /* IDA 0x1052E; mov     dx, [bp+6] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+6)));
L_10531: ; /* IDA 0x10531; mov     bx, [bp+8] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+8)));
L_10534: ; /* IDA 0x10534; mov     ax, 9 */
    m->ax=(uint16_t)(9);
L_10537: ; /* IDA 0x10537; call    sub_20000 */
    mx_far(m,0x537,0x53c,h);
L_1053C: ; /* IDA 0x1053C; jmp     short loc_10548 */
    goto L_10548;
L_1053E: ; /* IDA 0x1053E; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_10541: ; /* IDA 0x10541; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_10543: ; /* IDA 0x10543; call    sub_20000 */
    mx_far(m,0x543,0x548,h);
L_10548: ; /* IDA 0x10548; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10549: ; /* IDA 0x10549; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1054A: ; /* IDA 0x1054A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1054B: ; /* IDA 0x1054B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1054C: ; /* IDA 0x1054C; retn */
    return;
}
/* Original sub_1054D; IDA 0x1054D. */
static void ch_body_1054D(KiMachine16 *m,const KiEconomyHooks *h) {
L_1054D: ; /* IDA 0x1054D; cmp     byte ptr [bp+0Eh], 0 */
    ec_cmp(m,dl_read8(m,m->ss,(uint16_t)(m->bp+14)),0,8);
L_10551: ; /* IDA 0x10551; jz      short locret_1059A */
    if(m->flags&0x40) goto L_1059A;
L_10553: ; /* IDA 0x10553; call    sub_1061F */
    ch_call(m,1567,1366,h);
L_10556: ; /* IDA 0x10556; cmp     byte ptr [bp+0Fh], 0 */
    ec_cmp(m,dl_read8(m,m->ss,(uint16_t)(m->bp+15)),0,8);
L_1055A: ; /* IDA 0x1055A; jz      short loc_10568 */
    if(m->flags&0x40) goto L_10568;
L_1055C: ; /* IDA 0x1055C; sub     word ptr [bp+0Ch], 10h */
    dl_write16(m,m->ss,(uint16_t)(m->bp+12),ec_math(m,dl_read16(m,m->ss,(uint16_t)(m->bp+12)),16,16,0,1));
L_10560: ; /* IDA 0x10560; dec     byte ptr [bp+0Eh] */
    dl_write8(m,m->ss,(uint16_t)(m->bp+14),st_dec(m,dl_read8(m,m->ss,(uint16_t)(m->bp+14)),8));
L_10563: ; /* IDA 0x10563; dec     byte ptr [bp+0Fh] */
    dl_write8(m,m->ss,(uint16_t)(m->bp+15),st_dec(m,dl_read8(m,m->ss,(uint16_t)(m->bp+15)),8));
L_10566: ; /* IDA 0x10566; jmp     short loc_10597 */
    goto L_10597;
L_10568: ; /* IDA 0x10568; mov     dx, [bp+16h] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+22)));
L_1056B: ; /* IDA 0x1056B; mov     bx, dx */
    m->bx=(uint16_t)(m->dx);
L_1056D: ; /* IDA 0x1056D; sub     bx, 500h */
    m->bx=(uint16_t)(ec_math(m,m->bx,1280,16,0,1));
L_10571: ; /* IDA 0x10571; mov     al, [bp+11h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+17))));
L_10574: ; /* IDA 0x10574; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_10576: ; /* IDA 0x10576; mov     cx, [bp+14h] */
    m->cx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+20)));
L_10579: ; /* IDA 0x10579; call    nullsub_5 */
    ch_call(m,1517,1404,h);
L_1057C: ; /* IDA 0x1057C; sub     dx, 50h ; 'P' */
    m->dx=(uint16_t)(ec_math(m,m->dx,80,16,0,1));
L_1057F: ; /* IDA 0x1057F; sub     bx, 50h ; 'P' */
    m->bx=(uint16_t)(ec_math(m,m->bx,80,16,0,1));
L_10582: ; /* IDA 0x10582; cmp     bx, cx */
    ec_cmp(m,m->bx,m->cx,16);
L_10584: ; /* IDA 0x10584; jnb     short loc_10579 */
    if(!(m->flags&1)) goto L_10579;
L_10586: ; /* IDA 0x10586; dec     byte ptr [bp+0Eh] */
#if KI_CHOICE_MUTATION == 10
    ; /* lost absolute selected-row decrement */
#else
    dl_write8(m,m->ss,(uint16_t)(m->bp+14),st_dec(m,dl_read8(m,m->ss,(uint16_t)(m->bp+14)),8));
#endif
L_10589: ; /* IDA 0x10589; mov     dx, [bp+0Ah] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+10)));
L_1058C: ; /* IDA 0x1058C; mov     bx, [bp+0Ch] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+12)));
L_1058F: ; /* IDA 0x1058F; mov     al, [bp+0Eh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+14))));
L_10592: ; /* IDA 0x10592; call    cs:word_1036D */
    ch_call(m,dl_read16(m,m->cs,877),1431,h);
L_10597: ; /* IDA 0x10597; call    sub_1061F */
    ch_call(m,1567,1434,h);
L_1059A: ; /* IDA 0x1059A; retn */
    return;
}
/* Original sub_1059B; IDA 0x1059B. */
static void ch_body_1059B(KiMachine16 *m,const KiEconomyHooks *h) {
L_1059B: ; /* IDA 0x1059B; mov     al, [bp+0Eh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+14))));
L_1059E: ; /* IDA 0x1059E; cmp     al, [bp+0] */
#if KI_CHOICE_MUTATION == 4
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ss,(uint16_t)(m->bp+1)),8);
#else
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ss,(uint16_t)(m->bp)),8);
#endif
L_105A1: ; /* IDA 0x105A1; jnb     short locret_105EC */
    if(!(m->flags&1)) goto L_105EC;
L_105A3: ; /* IDA 0x105A3; call    sub_1061F */
    ch_call(m,1567,1446,h);
L_105A6: ; /* IDA 0x105A6; mov     al, [bp+0Fh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+15))));
L_105A9: ; /* IDA 0x105A9; cmp     al, [bp+1] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ss,(uint16_t)(m->bp+1)),8);
L_105AC: ; /* IDA 0x105AC; jz      short loc_105BA */
    if(m->flags&0x40) goto L_105BA;
L_105AE: ; /* IDA 0x105AE; add     word ptr [bp+0Ch], 10h */
    dl_write16(m,m->ss,(uint16_t)(m->bp+12),ec_math(m,dl_read16(m,m->ss,(uint16_t)(m->bp+12)),16,16,0,0));
L_105B2: ; /* IDA 0x105B2; inc     byte ptr [bp+0Eh] */
    dl_write8(m,m->ss,(uint16_t)(m->bp+14),ec_inc(m,dl_read8(m,m->ss,(uint16_t)(m->bp+14)),8));
L_105B5: ; /* IDA 0x105B5; inc     byte ptr [bp+0Fh] */
    dl_write8(m,m->ss,(uint16_t)(m->bp+15),ec_inc(m,dl_read8(m,m->ss,(uint16_t)(m->bp+15)),8));
L_105B8: ; /* IDA 0x105B8; jmp     short loc_105E9 */
    goto L_105E9;
L_105BA: ; /* IDA 0x105BA; mov     dx, [bp+14h] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+20)));
L_105BD: ; /* IDA 0x105BD; mov     bx, dx */
    m->bx=(uint16_t)(m->dx);
L_105BF: ; /* IDA 0x105BF; add     bx, 500h */
    m->bx=(uint16_t)(ec_math(m,m->bx,1280,16,0,0));
L_105C3: ; /* IDA 0x105C3; mov     al, [bp+11h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+17))));
L_105C6: ; /* IDA 0x105C6; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_105C8: ; /* IDA 0x105C8; mov     cx, [bp+16h] */
    m->cx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+22)));
L_105CB: ; /* IDA 0x105CB; call    nullsub_5 */
    ch_call(m,1517,1486,h);
L_105CE: ; /* IDA 0x105CE; add     dx, 50h ; 'P' */
    m->dx=(uint16_t)(ec_math(m,m->dx,80,16,0,0));
L_105D1: ; /* IDA 0x105D1; add     bx, 50h ; 'P' */
    m->bx=(uint16_t)(ec_math(m,m->bx,80,16,0,0));
L_105D4: ; /* IDA 0x105D4; cmp     bx, cx */
    ec_cmp(m,m->bx,m->cx,16);
L_105D6: ; /* IDA 0x105D6; jbe     short loc_105CB */
    if(m->flags&(1|0x40)) goto L_105CB;
L_105D8: ; /* IDA 0x105D8; inc     byte ptr [bp+0Eh] */
#if KI_CHOICE_MUTATION == 11
    ; /* lost absolute selected-row increment */
#else
    dl_write8(m,m->ss,(uint16_t)(m->bp+14),ec_inc(m,dl_read8(m,m->ss,(uint16_t)(m->bp+14)),8));
#endif
L_105DB: ; /* IDA 0x105DB; mov     dx, [bp+0Ah] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+10)));
L_105DE: ; /* IDA 0x105DE; mov     bx, [bp+0Ch] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+12)));
L_105E1: ; /* IDA 0x105E1; mov     al, [bp+0Eh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+14))));
L_105E4: ; /* IDA 0x105E4; call    cs:word_1036D */
    ch_call(m,dl_read16(m,m->cs,877),1513,h);
L_105E9: ; /* IDA 0x105E9; call    sub_1061F */
    ch_call(m,1567,1516,h);
L_105EC: ; /* IDA 0x105EC; retn */
    return;
}
/* Original nullsub_5; IDA 0x105ED. */
static void ch_body_105ED(KiMachine16 *m,const KiEconomyHooks *h) {
L_105ED: ; /* IDA 0x105ED; retn */
    return;
}
/* Original sub_1061F; IDA 0x1061F. */
static void ch_body_1061F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1061F: ; /* IDA 0x1061F; mov     dx, [bp+0Ah] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+10)));
L_10622: ; /* IDA 0x10622; mov     bx, [bp+0Ch] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+12)));
L_10625: ; /* IDA 0x10625; mov     si, [bp+12h] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+18)));
L_10628: ; /* IDA 0x10628; mov     di, 10h */
    m->di=(uint16_t)(16);
L_1062B: ; /* IDA 0x1062B; call    sub_10B46 */
    ch_call(m,2886,1582,h);
L_1062E: ; /* IDA 0x1062E; retn */
    return;
}
/* Original sub_10B46; IDA 0x10B46. */
static void ch_body_10B46(KiMachine16 *m,const KiEconomyHooks *h) {
L_10B46: ; /* IDA 0x10B46; push    ds */
    ec_push(m,m->ds);
L_10B47: ; /* IDA 0x10B47; push    ax */
    ec_push(m,m->ax);
L_10B48: ; /* IDA 0x10B48; push    cx */
    ec_push(m,m->cx);
L_10B49: ; /* IDA 0x10B49; push    bp */
    ec_push(m,m->bp);
L_10B4A: ; /* IDA 0x10B4A; mov     cl, dl */
    vg_cl(m,(uint8_t)((uint8_t)m->dx));
L_10B4C: ; /* IDA 0x10B4C; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_10B4E: ; /* IDA 0x10B4E; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_10B50: ; /* IDA 0x10B50; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_10B52: ; /* IDA 0x10B52; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10B54: ; /* IDA 0x10B54; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10B56: ; /* IDA 0x10B56; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10B58: ; /* IDA 0x10B58; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10B5A: ; /* IDA 0x10B5A; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_10B5C: ; /* IDA 0x10B5C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10B5E: ; /* IDA 0x10B5E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10B60: ; /* IDA 0x10B60; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_10B62: ; /* IDA 0x10B62; and     cx, 7 */
    {uint16_t result=((m->cx)&(7))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_10B65: ; /* IDA 0x10B65; mov     dx, 80FFh */
    m->dx=(uint16_t)(33023);
L_10B68: ; /* IDA 0x10B68; shr     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
L_10B6A: ; /* IDA 0x10B6A; add     cx, si */
    m->cx=(uint16_t)(ec_math(m,m->cx,m->si,16,0,0));
L_10B6C: ; /* IDA 0x10B6C; dec     cx */
    m->cx=(uint16_t)(st_dec(m,m->cx,16));
L_10B6D: ; /* IDA 0x10B6D; mov     bp, cx */
    m->bp=(uint16_t)(m->cx);
L_10B6F: ; /* IDA 0x10B6F; shr     bp, 1 */
    m->bp=(uint16_t)(dl_shift(m,m->bp,1,16,5));
L_10B71: ; /* IDA 0x10B71; shr     bp, 1 */
    m->bp=(uint16_t)(dl_shift(m,m->bp,1,16,5));
L_10B73: ; /* IDA 0x10B73; shr     bp, 1 */
    m->bp=(uint16_t)(dl_shift(m,m->bp,1,16,5));
L_10B75: ; /* IDA 0x10B75; inc     bp */
    m->bp=(uint16_t)(ec_inc(m,m->bp,16));
L_10B76: ; /* IDA 0x10B76; and     cx, 7 */
    {uint16_t result=((m->cx)&(7))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_10B79: ; /* IDA 0x10B79; xor     cl, 7 */
    {uint16_t result=(((uint8_t)m->cx)^(7))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_10B7C: ; /* IDA 0x10B7C; mov     ax, 0FEFFh */
    m->ax=(uint16_t)(65279);
L_10B7F: ; /* IDA 0x10B7F; shl     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,4));
L_10B81: ; /* IDA 0x10B81; not     ah */
    vg_ah(m,(uint8_t)(~((uint8_t)(m->ax>>8))));
L_10B83: ; /* IDA 0x10B83; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_10B85: ; /* IDA 0x10B85; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_10B87: ; /* IDA 0x10B87; push    ax */
    ec_push(m,m->ax);
L_10B88: ; /* IDA 0x10B88; push    dx */
    ec_push(m,m->dx);
L_10B89: ; /* IDA 0x10B89; call    sub_10CAC */
    ch_call(m,3244,2956,h);
L_10B8C: ; /* IDA 0x10B8C; mov     dx, 3CFh */
    m->dx=(uint16_t)(975);
L_10B8F: ; /* IDA 0x10B8F; mov     al, 0Ch */
    vg_al(m,(uint8_t)(12));
L_10B91: ; /* IDA 0x10B91; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_10B92: ; /* IDA 0x10B92; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_10B93: ; /* IDA 0x10B93; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_10B95: ; /* IDA 0x10B95; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_10B96: ; /* IDA 0x10B96; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_10B97: ; /* IDA 0x10B97; mov     al, 18h */
    vg_al(m,(uint8_t)(24));
L_10B99: ; /* IDA 0x10B99; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_10B9A: ; /* IDA 0x10B9A; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10B9B: ; /* IDA 0x10B9B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10B9C: ; /* IDA 0x10B9C; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_10B9F: ; /* IDA 0x10B9F; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_10BA1: ; /* IDA 0x10BA1; call    sub_10BAF */
    ch_call(m,2991,2980,h);
L_10BA4: ; /* IDA 0x10BA4; add     si, 50h ; 'P' */
    m->si=(uint16_t)(ec_math(m,m->si,80,16,0,0));
L_10BA7: ; /* IDA 0x10BA7; dec     di */
    m->di=(uint16_t)(st_dec(m,m->di,16));
L_10BA8: ; /* IDA 0x10BA8; jnz     short loc_10BA1 */
    if(!(m->flags&0x40)) goto L_10BA1;
L_10BAA: ; /* IDA 0x10BAA; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_10BAB: ; /* IDA 0x10BAB; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10BAC: ; /* IDA 0x10BAC; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10BAD: ; /* IDA 0x10BAD; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_10BAE: ; /* IDA 0x10BAE; retn */
    return;
}
/* Original sub_10BAF; IDA 0x10BAF. */
static void ch_body_10BAF(KiMachine16 *m,const KiEconomyHooks *h) {
L_10BAF: ; /* IDA 0x10BAF; mov     cx, bp */
    m->cx=(uint16_t)(m->bp);
L_10BB1: ; /* IDA 0x10BB1; push    si */
    ec_push(m,m->si);
L_10BB2: ; /* IDA 0x10BB2; mov     ah, [si] */
#if KI_CHOICE_MUTATION == 7
    ; /* lost XOR latch read */
#else
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
#endif
L_10BB4: ; /* IDA 0x10BB4; mov     [si], dl */
    dl_write8(m,m->ds,(uint16_t)(m->si),(uint8_t)m->dx);
L_10BB6: ; /* IDA 0x10BB6; dec     cx */
    m->cx=(uint16_t)(st_dec(m,m->cx,16));
L_10BB7: ; /* IDA 0x10BB7; jz      short loc_10BCB */
    if(m->flags&0x40) goto L_10BCB;
L_10BB9: ; /* IDA 0x10BB9; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_10BBA: ; /* IDA 0x10BBA; dec     cx */
    m->cx=(uint16_t)(st_dec(m,m->cx,16));
L_10BBB: ; /* IDA 0x10BBB; jz      short loc_10BC7 */
    if(m->flags&0x40) goto L_10BC7;
L_10BBD: ; /* IDA 0x10BBD; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_10BBF: ; /* IDA 0x10BBF; mov     ah, [si] */
#if KI_CHOICE_MUTATION == 7
    ; /* lost XOR latch read */
#else
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
#endif
L_10BC1: ; /* IDA 0x10BC1; mov     [si], al */
    dl_write8(m,m->ds,(uint16_t)(m->si),(uint8_t)m->ax);
L_10BC3: ; /* IDA 0x10BC3; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_10BC4: ; /* IDA 0x10BC4; dec     cx */
    m->cx=(uint16_t)(st_dec(m,m->cx,16));
L_10BC5: ; /* IDA 0x10BC5; jnz     short loc_10BBF */
    if(!(m->flags&0x40)) goto L_10BBF;
L_10BC7: ; /* IDA 0x10BC7; mov     ah, [si] */
#if KI_CHOICE_MUTATION == 7
    ; /* lost XOR latch read */
#else
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
#endif
L_10BC9: ; /* IDA 0x10BC9; mov     [si], bl */
    dl_write8(m,m->ds,(uint16_t)(m->si),(uint8_t)m->bx);
L_10BCB: ; /* IDA 0x10BCB; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_10BCC: ; /* IDA 0x10BCC; retn */
    return;
}
/* Original sub_19479; IDA 0x19479. */
static void ch_body_19479(KiMachine16 *m,const KiEconomyHooks *h) {
L_19479: ; /* IDA 0x19479; push    bp */
    ec_push(m,m->bp);
L_1947A: ; /* IDA 0x1947A; sub     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,1));
L_1947D: ; /* IDA 0x1947D; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_1947F: ; /* IDA 0x1947F; sub     dh, 2 */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),2,8,0,1)));
L_19482: ; /* IDA 0x19482; push    ax */
    ec_push(m,m->ax);
L_19483: ; /* IDA 0x19483; push    bx */
    ec_push(m,m->bx);
L_19484: ; /* IDA 0x19484; push    cx */
    ec_push(m,m->cx);
L_19485: ; /* IDA 0x19485; push    dx */
    ec_push(m,m->dx);
L_19486: ; /* IDA 0x19486; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_19488: ; /* IDA 0x19488; mov     ah, bl */
    vg_ah(m,(uint8_t)((uint8_t)m->bx));
L_1948A: ; /* IDA 0x1948A; xchg    dh, bl */
    {uint16_t a=(uint8_t)(m->dx>>8),b=(uint8_t)m->bx;vg_dh(m,(uint8_t)(b));gw_bl(m,(uint8_t)(a));}
L_1948C: ; /* IDA 0x1948C; inc     cl */
    vg_cl(m,(uint8_t)(ec_inc(m,(uint8_t)m->cx,8)));
L_1948E: ; /* IDA 0x1948E; inc     ch */
    vg_ch(m,(uint8_t)(ec_inc(m,(uint8_t)(m->cx>>8),8)));
L_19490: ; /* IDA 0x19490; mov     [bp+0], dx */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->dx);
L_19493: ; /* IDA 0x19493; mov     [bp+2], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->bx);
L_19496: ; /* IDA 0x19496; mov     [bp+4], cx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+4),m->cx);
L_19499: ; /* IDA 0x19499; call    sub_1D5D4 */
    ch_call(m,54740,38044,h);
L_1949C: ; /* IDA 0x1949C; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1949D: ; /* IDA 0x1949D; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1949E: ; /* IDA 0x1949E; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1949F: ; /* IDA 0x1949F; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_194A0: ; /* IDA 0x194A0; push    bp */
    ec_push(m,m->bp);
L_194A1: ; /* IDA 0x194A1; add     dh, 2 */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),2,8,0,0)));
L_194A4: ; /* IDA 0x194A4; call    sub_1036F */
    ch_call(m,879,38055,h);
L_194A7: ; /* IDA 0x194A7; lahf */
    vg_ah(m,(uint8_t)((m->flags&0xd5)|2));
L_194A8: ; /* IDA 0x194A8; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_194A9: ; /* IDA 0x194A9; push    ax */
    ec_push(m,m->ax);
L_194AA: ; /* IDA 0x194AA; mov     dx, [bp+0] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_194AD: ; /* IDA 0x194AD; mov     bx, [bp+2] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_194B0: ; /* IDA 0x194B0; mov     cx, [bp+4] */
    m->cx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+4)));
L_194B3: ; /* IDA 0x194B3; mov     ah, 1 */
    vg_ah(m,(uint8_t)(1));
L_194B5: ; /* IDA 0x194B5; call    sub_1D5D4 */
    ch_call(m,54740,38072,h);
L_194B8: ; /* IDA 0x194B8; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_194B9: ; /* IDA 0x194B9; add     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,0));
L_194BC: ; /* IDA 0x194BC; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_194BD: ; /* IDA 0x194BD; sahf */
#if KI_CHOICE_MUTATION == 9
    ; /* lost selector carry restore */
#else
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
#endif
L_194BE: ; /* IDA 0x194BE; retn */
    return;
}
/* Original sub_194BF; IDA 0x194BF. */
static void ch_body_194BF(KiMachine16 *m,const KiEconomyHooks *h) {
L_194BF: ; /* IDA 0x194BF; push    bp */
    ec_push(m,m->bp);
L_194C0: ; /* IDA 0x194C0; sub     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,1));
L_194C3: ; /* IDA 0x194C3; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_194C5: ; /* IDA 0x194C5; sub     dh, 2 */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),2,8,0,1)));
L_194C8: ; /* IDA 0x194C8; push    ax */
    ec_push(m,m->ax);
L_194C9: ; /* IDA 0x194C9; push    bx */
    ec_push(m,m->bx);
L_194CA: ; /* IDA 0x194CA; push    cx */
    ec_push(m,m->cx);
L_194CB: ; /* IDA 0x194CB; push    dx */
    ec_push(m,m->dx);
L_194CC: ; /* IDA 0x194CC; mov     ax, cx */
    m->ax=(uint16_t)(m->cx);
L_194CE: ; /* IDA 0x194CE; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_194D0: ; /* IDA 0x194D0; xchg    dh, bl */
    {uint16_t a=(uint8_t)(m->dx>>8),b=(uint8_t)m->bx;vg_dh(m,(uint8_t)(b));gw_bl(m,(uint8_t)(a));}
L_194D2: ; /* IDA 0x194D2; add     bx, 2 */
    m->bx=(uint16_t)(ec_math(m,m->bx,2,16,0,0));
L_194D5: ; /* IDA 0x194D5; mov     cl, 4 */
    vg_cl(m,(uint8_t)(4));
L_194D7: ; /* IDA 0x194D7; shl     bx, cl */
    m->bx=(uint16_t)(dl_shift(m,m->bx,(uint8_t)m->cx,16,4));
L_194D9: ; /* IDA 0x194D9; shl     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,4));
L_194DB: ; /* IDA 0x194DB; add     ax, 101h */
    m->ax=(uint16_t)(ec_math(m,m->ax,257,16,0,0));
L_194DE: ; /* IDA 0x194DE; shl     ah, cl */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),(uint8_t)m->cx,8,4)));
L_194E0: ; /* IDA 0x194E0; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_194E2: ; /* IDA 0x194E2; mov     [bp+0], dx */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->dx);
L_194E5: ; /* IDA 0x194E5; mov     [bp+2], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->bx);
L_194E8: ; /* IDA 0x194E8; mov     [bp+4], cx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+4),m->cx);
L_194EB: ; /* IDA 0x194EB; call    sub_19796 */
    ch_call(m,38806,38126,h);
L_194EE: ; /* IDA 0x194EE; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_194EF: ; /* IDA 0x194EF; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_194F0: ; /* IDA 0x194F0; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_194F1: ; /* IDA 0x194F1; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_194F2: ; /* IDA 0x194F2; push    bp */
    ec_push(m,m->bp);
L_194F3: ; /* IDA 0x194F3; add     dh, 2 */
    vg_dh(m,(uint8_t)(ec_math(m,(uint8_t)(m->dx>>8),2,8,0,0)));
L_194F6: ; /* IDA 0x194F6; call    sub_1036F */
    ch_call(m,879,38137,h);
L_194F9: ; /* IDA 0x194F9; lahf */
    vg_ah(m,(uint8_t)((m->flags&0xd5)|2));
L_194FA: ; /* IDA 0x194FA; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_194FB: ; /* IDA 0x194FB; push    ax */
    ec_push(m,m->ax);
L_194FC: ; /* IDA 0x194FC; mov     dx, [bp+0] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_194FF: ; /* IDA 0x194FF; mov     bx, [bp+2] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_19502: ; /* IDA 0x19502; mov     cx, [bp+4] */
    m->cx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+4)));
L_19505: ; /* IDA 0x19505; call    sub_197C3 */
    ch_call(m,38851,38152,h);
L_19508: ; /* IDA 0x19508; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_19509: ; /* IDA 0x19509; add     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,0));
L_1950C: ; /* IDA 0x1950C; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1950D: ; /* IDA 0x1950D; sahf */
#if KI_CHOICE_MUTATION == 9
    ; /* lost selector carry restore */
#else
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
#endif
L_1950E: ; /* IDA 0x1950E; retn */
    return;
}
/* Original loc_19409; IDA 0x19409. */
static void ch_body_19409(KiMachine16 *m,const KiEconomyHooks *h) {
L_19409: ; /* IDA 0x19409; push    ds */
    ec_push(m,m->ds);
L_1940A: ; /* IDA 0x1940A; push    bx */
    ec_push(m,m->bx);
L_1940B: ; /* IDA 0x1940B; push    cx */
    ec_push(m,m->cx);
L_1940C: ; /* IDA 0x1940C; push    dx */
    ec_push(m,m->dx);
L_1940D: ; /* IDA 0x1940D; push    si */
    ec_push(m,m->si);
L_1940E: ; /* IDA 0x1940E; push    di */
    ec_push(m,m->di);
L_1940F: ; /* IDA 0x1940F; mov     ds, cs:word_10D38 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3384));
L_19414: ; /* IDA 0x19414; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_19416: ; /* IDA 0x19416; mov     si, cx */
    m->si=(uint16_t)(m->cx);
L_19418: ; /* IDA 0x19418; mov     si, [si] */
    m->si=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_1941A: ; /* IDA 0x1941A; push    ax */
    ec_push(m,m->ax);
L_1941B: ; /* IDA 0x1941B; mov     ch, al */
    vg_ch(m,(uint8_t)((uint8_t)m->ax));
L_1941D: ; /* IDA 0x1941D; mov     di, si */
    m->di=(uint16_t)(m->si);
L_1941F: ; /* IDA 0x1941F; xor     cl, cl */
    {uint16_t result=(((uint8_t)m->cx)^((uint8_t)m->cx))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_19421: ; /* IDA 0x19421; mov     ax, [di] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di)));
L_19423: ; /* IDA 0x19423; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_19424: ; /* IDA 0x19424; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_19425: ; /* IDA 0x19425; inc     cl */
    vg_cl(m,(uint8_t)(ec_inc(m,(uint8_t)m->cx,8)));
L_19427: ; /* IDA 0x19427; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_19429: ; /* IDA 0x19429; jnz     short loc_19421 */
    if(!(m->flags&0x40)) goto L_19421;
L_1942B: ; /* IDA 0x1942B; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_1942D: ; /* IDA 0x1942D; mov     al, cl */
    vg_al(m,(uint8_t)((uint8_t)m->cx));
L_1942F: ; /* IDA 0x1942F; shl     al, 1 */
#if KI_CHOICE_MUTATION == 1
    ; /* lost row stride doubling */
#else
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
#endif
L_19431: ; /* IDA 0x19431; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_19433: ; /* IDA 0x19433; mov     cs:byte_19465, al */
    dl_write8(m,m->cs,37989,(uint8_t)m->ax);
L_19437: ; /* IDA 0x19437; mov     cs:word_19467, si */
    dl_write16(m,m->cs,37991,m->si);
L_1943C: ; /* IDA 0x1943C; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1943D: ; /* IDA 0x1943D; mov     bx, 945Ah */
    m->bx=(uint16_t)(37978);
L_19440: ; /* IDA 0x19440; mov     dl, 9 */
#if KI_CHOICE_MUTATION == 12
    vg_dl(m,9);
#else
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,0x9441)));
#endif
L_19442: ; /* IDA 0x19442; mov     dh, 0Fh */
    vg_dh(m,(uint8_t)(dl_read8(m,m->cs,0x9443)));
L_19444: ; /* IDA 0x19444; mov     si, 0FFFFh */
    m->si=(uint16_t)(65535);
L_19447: ; /* IDA 0x19447; and     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))&((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_19449: ; /* IDA 0x19449; jnz     short loc_19450 */
    if(!(m->flags&0x40)) goto L_19450;
L_1944B: ; /* IDA 0x1944B; call    sub_19479 */
    ch_call(m,38009,37966,h);
L_1944E: ; /* IDA 0x1944E; jmp     short loc_19453 */
    goto L_19453;
L_19450: ; /* IDA 0x19450; call    sub_194BF */
    ch_call(m,38079,37971,h);
L_19453: ; /* IDA 0x19453; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_19454: ; /* IDA 0x19454; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_19455: ; /* IDA 0x19455; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_19456: ; /* IDA 0x19456; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_19457: ; /* IDA 0x19457; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_19458: ; /* IDA 0x19458; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_19459: ; /* IDA 0x19459; retn */
    return;
}
/* Original (raw code); IDA 0x1945A. */
static void ch_body_1945A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1945A: ; /* IDA 0x1945A; push    ds */
    ec_push(m,m->ds);
L_1945B: ; /* IDA 0x1945B; push    ax */
    ec_push(m,m->ax);
L_1945C: ; /* IDA 0x1945C; push    bx */
    ec_push(m,m->bx);
L_1945D: ; /* IDA 0x1945D; push    dx */
    ec_push(m,m->dx);
L_1945E: ; /* IDA 0x1945E; push    si */
    ec_push(m,m->si);
L_1945F: ; /* IDA 0x1945F; mov     ds, cs:word_10D38 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3384));
L_19464: ; /* IDA 0x19464; mov     ah, 0Ch */
    vg_ah(m,(uint8_t)(dl_read8(m,m->cs,0x9465)));
L_19466: ; /* IDA 0x19466; mov     si, 1234h */
#if KI_CHOICE_MUTATION == 2
    m->si=0x1234;
#else
    m->si=dl_read16(m,m->cs,0x9467);
#endif
L_19469: ; /* IDA 0x19469; mul     ah */
    dl_mul(m,(uint8_t)(m->ax>>8),8);
L_1946B: ; /* IDA 0x1946B; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_1946D: ; /* IDA 0x1946D; mov     ax, 0F01h */
    m->ax=(uint16_t)(3841);
L_19470: ; /* IDA 0x19470; call    sub_1F6DC */
    ch_call(m,63196,38003,h);
L_19473: ; /* IDA 0x19473; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_19474: ; /* IDA 0x19474; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_19475: ; /* IDA 0x19475; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_19476: ; /* IDA 0x19476; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_19477: ; /* IDA 0x19477; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_19478: ; /* IDA 0x19478; retn */
    return;
}
int ki_choice_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x36f:ch_body_1036F(m,h);break;
case 0x3c3:ch_body_103C3(m,h);break;
case 0x3e6:ch_body_103E6(m,h);break;
case 0x414:ch_body_10414(m,h);break;
case 0x4b5:ch_body_104B5(m,h);break;
case 0x4ff:ch_body_104FF(m,h);break;
case 0x54d:ch_body_1054D(m,h);break;
case 0x59b:ch_body_1059B(m,h);break;
case 0x5ed:ch_body_105ED(m,h);break;
case 0x61f:ch_body_1061F(m,h);break;
case 0xb46:ch_body_10B46(m,h);break;
case 0xbaf:ch_body_10BAF(m,h);break;
case 0x9479:ch_body_19479(m,h);break;
case 0x94bf:ch_body_194BF(m,h);break;
case 0x9409:ch_body_19409(m,h);break;
case 0x945a:ch_body_1945A(m,h);break;
default:return 0;}return 1;}
