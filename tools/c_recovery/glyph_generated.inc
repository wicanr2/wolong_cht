/* Generated native C glyph control; locator=proven, semantic evidence re/108/spec/228.
 * Original names/addresses/operands retained. No opcode dispatch. */
/* Original sub_1F720; IDA 0x1F720; raw control; source re/108 */
static void gl_body_1F720(KiMachine16 *m,const KiEconomyHooks *h) {
L_1F720: ; /* IDA 0x1f720; push    es */
    ec_push(m,m->es);
L_1F721: ; /* IDA 0x1f721; push    ax */
    ec_push(m,m->ax);
L_1F722: ; /* IDA 0x1f722; push    bx */
    ec_push(m,m->bx);
L_1F723: ; /* IDA 0x1f723; push    dx */
    ec_push(m,m->dx);
L_1F724: ; /* IDA 0x1f724; push    bp */
    ec_push(m,m->bp);
L_1F725: ; /* IDA 0x1f725; mov     ax, 5000h */
    m->ax=(uint16_t)(20480);
L_1F728: ; /* IDA 0x1f728; mov     bx, 0 */
    m->bx=(uint16_t)(0);
L_1F72B: ; /* IDA 0x1f72b; mov     dx, 810h */
    m->dx=(uint16_t)(2064);
L_1F72E: ; /* IDA 0x1f72e; xor     bp, bp */
    {uint16_t result=((m->bp)^(m->bp))&65535;ec_logic(m,result,16);m->bp=(uint16_t)(result);}
L_1F730: ; /* IDA 0x1f730; int     15h */
    gl_interrupt(m,21,63282);
L_1F732: ; /* IDA 0x1f732; jb      short loc_1F758 */
    if(m->flags&1) goto L_1F758;
L_1F734: ; /* IDA 0x1f734; mov     cs:word_1F791, bx */
    dl_write16(m,m->cs,63377,m->bx);
L_1F739: ; /* IDA 0x1f739; mov     cs:word_1F793, es */
    dl_write16(m,m->cs,63379,m->es);
L_1F73E: ; /* IDA 0x1f73e; mov     ax, 5000h */
    m->ax=(uint16_t)(20480);
L_1F741: ; /* IDA 0x1f741; mov     bx, 100h */
    m->bx=(uint16_t)(256);
L_1F744: ; /* IDA 0x1f744; mov     dx, 1010h */
    m->dx=(uint16_t)(4112);
L_1F747: ; /* IDA 0x1f747; xor     bp, bp */
    {uint16_t result=((m->bp)^(m->bp))&65535;ec_logic(m,result,16);m->bp=(uint16_t)(result);}
L_1F749: ; /* IDA 0x1f749; int     15h */
    gl_interrupt(m,21,63307);
L_1F74B: ; /* IDA 0x1f74b; jb      short loc_1F758 */
    if(m->flags&1) goto L_1F758;
L_1F74D: ; /* IDA 0x1f74d; mov     cs:word_1F78A, bx */
    dl_write16(m,m->cs,63370,m->bx);
L_1F752: ; /* IDA 0x1f752; mov     cs:word_1F78C, es */
    dl_write16(m,m->cs,63372,m->es);
L_1F757: ; /* IDA 0x1f757; clc */
    m->flags&=(uint16_t)~1u;
L_1F758: ; /* IDA 0x1f758; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1F759: ; /* IDA 0x1f759; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1F75A: ; /* IDA 0x1f75a; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1F75B: ; /* IDA 0x1f75b; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1F75C: ; /* IDA 0x1f75c; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1F75D: ; /* IDA 0x1f75d; retn */
    return;
}
/* Original sub_1F7A4; IDA 0x1F7A4; raw control; source re/108 */
static void gl_body_1F7A4(KiMachine16 *m,const KiEconomyHooks *h) {
L_1F7A4: ; /* IDA 0x1f7a4; test    byte ptr cs:word_1F7A1, 1 */
    {uint16_t result=((dl_read8(m,m->cs,63393))&(1))&255;ec_logic(m,result,8);}
L_1F7AA: ; /* IDA 0x1f7aa; jz      short $+2 */
    if(m->flags&0x40) goto L_1F7AC;
L_1F7AC: ; /* IDA 0x1f7ac; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_1F7AF: ; /* IDA 0x1f7af; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1F7B1: ; /* IDA 0x1f7b1; mov     al, dl */
    vg_al(m,(uint8_t)((uint8_t)m->dx));
L_1F7B3: ; /* IDA 0x1f7b3; and     al, 7 */
    {uint16_t result=(((uint8_t)m->ax)&(7))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F7B5: ; /* IDA 0x1f7b5; mov     cl, 3 */
    vg_cl(m,(uint8_t)(3));
L_1F7B7: ; /* IDA 0x1f7b7; shr     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
L_1F7B9: ; /* IDA 0x1f7b9; inc     cl */
    vg_cl(m,(uint8_t)(ec_inc(m,(uint8_t)m->cx,8)));
L_1F7BB: ; /* IDA 0x1f7bb; shl     bx, cl */
    m->bx=(uint16_t)(dl_shift(m,m->bx,(uint8_t)m->cx,16,4));
L_1F7BD: ; /* IDA 0x1f7bd; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_1F7BF: ; /* IDA 0x1f7bf; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1F7C1: ; /* IDA 0x1f7c1; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1F7C3: ; /* IDA 0x1f7c3; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_1F7C5: ; /* IDA 0x1f7c5; mov     cl, al */
    vg_cl(m,(uint8_t)((uint8_t)m->ax));
L_1F7C7: ; /* IDA 0x1f7c7; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1F7CA: ; /* IDA 0x1f7ca; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_1F7CC: ; /* IDA 0x1f7cc; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F7CD: ; /* IDA 0x1f7cd; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F7CE: ; /* IDA 0x1f7ce; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1F7D0: ; /* IDA 0x1f7d0; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F7D1: ; /* IDA 0x1f7d1; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1F7D2: ; /* IDA 0x1f7d2; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1F7D4: ; /* IDA 0x1f7d4; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F7D5: ; /* IDA 0x1f7d5; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F7D6: ; /* IDA 0x1f7d6; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F7D8: ; /* IDA 0x1f7d8; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F7D9: ; /* IDA 0x1f7d9; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1F7DA: ; /* IDA 0x1f7da; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_1F7DC: ; /* IDA 0x1f7dc; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F7DD: ; /* IDA 0x1f7dd; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F7DE: ; /* IDA 0x1f7de; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1F7E0: ; /* IDA 0x1f7e0; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F7E1: ; /* IDA 0x1f7e1; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1F7E4: ; /* IDA 0x1f7e4; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F7E6: ; /* IDA 0x1f7e6; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F7E7: ; /* IDA 0x1f7e7; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F7E8: ; /* IDA 0x1f7e8; mov     al, byte ptr cs:word_1F7A1+1 */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,63394)));
L_1F7EC: ; /* IDA 0x1f7ec; and     al, 0Fh */
    {uint16_t result=(((uint8_t)m->ax)&(15))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F7EE: ; /* IDA 0x1f7ee; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F7EF: ; /* IDA 0x1f7ef; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_1F7F1: ; /* IDA 0x1f7f1; mov     si, bp */
    m->si=(uint16_t)(m->bp);
L_1F7F3: ; /* IDA 0x1f7f3; mov     ch, 10h */
    vg_ch(m,(uint8_t)(16));
L_1F7F5: ; /* IDA 0x1f7f5; mov     ax, es:[si] */
    m->ax=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->si)));
L_1F7F8: ; /* IDA 0x1f7f8; xchg    al, ah */
#if KI_GLYPH_MUTATION == 3
    ; /* lost byte exchange */
#else
    {uint16_t a=(uint8_t)m->ax,b=(uint8_t)(m->ax>>8);vg_al(m,(uint8_t)(b));vg_ah(m,(uint8_t)(a));}
#endif
L_1F7FA: ; /* IDA 0x1f7fa; cmp     cs:byte_1F7A3, 0 */
    ec_cmp(m,dl_read8(m,m->cs,63395),0,8);
L_1F800: ; /* IDA 0x1f800; jnz     short loc_1F805 */
    if(!(m->flags&0x40)) goto L_1F805;
L_1F802: ; /* IDA 0x1f802; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F804: ; /* IDA 0x1f804; dec     si */
#if KI_GLYPH_MUTATION == 4
    ; /* lost narrow stride */
#else
    m->si=(uint16_t)(st_dec(m,m->si,16));
#endif
L_1F805: ; /* IDA 0x1f805; mov     dh, al */
    vg_dh(m,(uint8_t)((uint8_t)m->ax));
L_1F807: ; /* IDA 0x1f807; xor     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_1F809: ; /* IDA 0x1f809; shr     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,5));
L_1F80B: ; /* IDA 0x1f80b; shr     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
L_1F80D: ; /* IDA 0x1f80d; mov     dh, [di] */
#if KI_GLYPH_MUTATION == 5
    ; /* lost latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
#endif
L_1F80F: ; /* IDA 0x1f80f; mov     [di], ah */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)(m->ax>>8));
L_1F811: ; /* IDA 0x1f811; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1F812: ; /* IDA 0x1f812; mov     dh, [di] */
#if KI_GLYPH_MUTATION == 5
    ; /* lost latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
#endif
L_1F814: ; /* IDA 0x1f814; mov     [di], al */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)m->ax);
L_1F816: ; /* IDA 0x1f816; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1F817: ; /* IDA 0x1f817; mov     dh, [di] */
#if KI_GLYPH_MUTATION == 5
    ; /* lost latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
#endif
L_1F819: ; /* IDA 0x1f819; mov     [di], dl */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)m->dx);
L_1F81B: ; /* IDA 0x1f81b; add     di, 4Eh ; 'N' */
#if KI_GLYPH_MUTATION == 6
    m->di=(uint16_t)(ec_math(m,m->di,77,16,0,0));
#else
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
#endif
L_1F81E: ; /* IDA 0x1f81e; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_1F81F: ; /* IDA 0x1f81f; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_1F820: ; /* IDA 0x1f820; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_1F822: ; /* IDA 0x1f822; jnz     short loc_1F7F5 */
    if(!(m->flags&0x40)) goto L_1F7F5;
L_1F824: ; /* IDA 0x1f824; test    byte ptr cs:word_1F7A1, 2 */
#if KI_GLYPH_MUTATION == 7
    {uint16_t result=((dl_read8(m,m->cs,63393))&(0))&255;ec_logic(m,result,8);}
#else
    {uint16_t result=((dl_read8(m,m->cs,63393))&(2))&255;ec_logic(m,result,8);}
#endif
L_1F82A: ; /* IDA 0x1f82a; jnz     short locret_1F877 */
    if(!(m->flags&0x40)) goto L_1F877;
L_1F82C: ; /* IDA 0x1f82c; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1F82F: ; /* IDA 0x1f82f; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F831: ; /* IDA 0x1f831; out     dx, al; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F832: ; /* IDA 0x1f832; inc     dx */
    m->dx=(uint16_t)(ec_inc(m,m->dx,16));
L_1F833: ; /* IDA 0x1f833; mov     al, byte ptr cs:word_1F7A1+1 */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,63394)));
L_1F837: ; /* IDA 0x1f837; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_1F839: ; /* IDA 0x1f839; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_1F83B: ; /* IDA 0x1f83b; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_1F83D: ; /* IDA 0x1f83d; shr     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
L_1F83F: ; /* IDA 0x1f83f; out     dx, al; EGA port: graphics controller data register */
    dl_out(m,m->dx,(uint8_t)m->ax,8);
L_1F840: ; /* IDA 0x1f840; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_1F842: ; /* IDA 0x1f842; mov     si, bp */
    m->si=(uint16_t)(m->bp);
L_1F844: ; /* IDA 0x1f844; mov     ch, 10h */
    vg_ch(m,(uint8_t)(16));
L_1F846: ; /* IDA 0x1f846; mov     ax, es:[si] */
    m->ax=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->si)));
L_1F849: ; /* IDA 0x1f849; xchg    al, ah */
#if KI_GLYPH_MUTATION == 3
    ; /* lost byte exchange */
#else
    {uint16_t a=(uint8_t)m->ax,b=(uint8_t)(m->ax>>8);vg_al(m,(uint8_t)(b));vg_ah(m,(uint8_t)(a));}
#endif
L_1F84B: ; /* IDA 0x1f84b; not     ax */
    m->ax=(uint16_t)(~(m->ax));
L_1F84D: ; /* IDA 0x1f84d; cmp     cs:byte_1F7A3, 0 */
    ec_cmp(m,dl_read8(m,m->cs,63395),0,8);
L_1F853: ; /* IDA 0x1f853; jnz     short loc_1F858 */
    if(!(m->flags&0x40)) goto L_1F858;
L_1F855: ; /* IDA 0x1f855; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1F857: ; /* IDA 0x1f857; dec     si */
#if KI_GLYPH_MUTATION == 4
    ; /* lost narrow stride */
#else
    m->si=(uint16_t)(st_dec(m,m->si,16));
#endif
L_1F858: ; /* IDA 0x1f858; mov     dh, al */
    vg_dh(m,(uint8_t)((uint8_t)m->ax));
L_1F85A: ; /* IDA 0x1f85a; xor     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_1F85C: ; /* IDA 0x1f85c; shr     ax, cl */
    m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,5));
L_1F85E: ; /* IDA 0x1f85e; shr     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
L_1F860: ; /* IDA 0x1f860; mov     dh, [di] */
#if KI_GLYPH_MUTATION == 5
    ; /* lost latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
#endif
L_1F862: ; /* IDA 0x1f862; mov     [di], ah */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)(m->ax>>8));
L_1F864: ; /* IDA 0x1f864; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1F865: ; /* IDA 0x1f865; mov     dh, [di] */
#if KI_GLYPH_MUTATION == 5
    ; /* lost latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
#endif
L_1F867: ; /* IDA 0x1f867; mov     [di], al */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)m->ax);
L_1F869: ; /* IDA 0x1f869; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1F86A: ; /* IDA 0x1f86a; mov     dh, [di] */
#if KI_GLYPH_MUTATION == 5
    ; /* lost latch read */
#else
    vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
#endif
L_1F86C: ; /* IDA 0x1f86c; mov     [di], dl */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)m->dx);
L_1F86E: ; /* IDA 0x1f86e; add     di, 4Eh ; 'N' */
#if KI_GLYPH_MUTATION == 6
    m->di=(uint16_t)(ec_math(m,m->di,77,16,0,0));
#else
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
#endif
L_1F871: ; /* IDA 0x1f871; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_1F872: ; /* IDA 0x1f872; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_1F873: ; /* IDA 0x1f873; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_1F875: ; /* IDA 0x1f875; jnz     short loc_1F846 */
    if(!(m->flags&0x40)) goto L_1F846;
L_1F877: ; /* IDA 0x1f877; retn */
    return;
}
/* Original sub_106F5; IDA 0x106F5; raw control; source re/108 */
static void gl_body_106F5(KiMachine16 *m,const KiEconomyHooks *h) {
L_106F5: ; /* IDA 0x106f5; call    sub_1F6DC */
    gl_call(m,63196,1784,h);
L_106F8: ; /* IDA 0x106f8; retn */
    return;
}
/* Original sub_106FD; IDA 0x106FD; raw control; source re/108 */
static void gl_body_106FD(KiMachine16 *m,const KiEconomyHooks *h) {
L_106FD: ; /* IDA 0x106fd; call    loc_10701 */
    gl_call(m,1793,1792,h);
L_10700: ; /* IDA 0x10700; retn */
    return;
}
/* Original loc_1F75E; IDA 0x1F75E; raw control; source re/108 */
static void gl_body_1F75E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1F75E: ; /* IDA 0x1f75e; push    ds */
    ec_push(m,m->ds);
L_1F75F: ; /* IDA 0x1f75f; push    es */
    ec_push(m,m->es);
L_1F760: ; /* IDA 0x1f760; push    si */
    ec_push(m,m->si);
L_1F761: ; /* IDA 0x1f761; push    di */
    ec_push(m,m->di);
L_1F762: ; /* IDA 0x1f762; push    bp */
    ec_push(m,m->bp);
L_1F763: ; /* IDA 0x1f763; mov     cs:word_1F7A1, ax */
    dl_write16(m,m->cs,63393,m->ax);
L_1F767: ; /* IDA 0x1f767; mov     ax, ss */
    m->ax=(uint16_t)(m->ss);
L_1F769: ; /* IDA 0x1f769; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1F76B: ; /* IDA 0x1f76b; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1F76D: ; /* IDA 0x1f76d; sub     sp, 20h */
    m->sp=(uint16_t)(ec_math(m,m->sp,32,16,0,1));
L_1F770: ; /* IDA 0x1f770; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_1F772: ; /* IDA 0x1f772; mov     si, sp */
    m->si=(uint16_t)(m->sp);
L_1F774: ; /* IDA 0x1f774; cmp     dx, 270h */
#if KI_GLYPH_MUTATION == 1
    ec_cmp(m,m->dx,623,16);
#else
    ec_cmp(m,m->dx,624,16);
#endif
L_1F778: ; /* IDA 0x1f778; ja      short loc_1F798 */
    if(!(m->flags&(1|0x40))) goto L_1F798;
L_1F77A: ; /* IDA 0x1f77a; cmp     bx, 180h */
#if KI_GLYPH_MUTATION == 2
    ec_cmp(m,m->bx,383,16);
#else
    ec_cmp(m,m->bx,384,16);
#endif
L_1F77E: ; /* IDA 0x1f77e; ja      short loc_1F798 */
    if(!(m->flags&(1|0x40))) goto L_1F798;
L_1F780: ; /* IDA 0x1f780; mov     cs:byte_1F7A3, ch */
    dl_write8(m,m->cs,63395,(uint8_t)(m->cx>>8));
L_1F785: ; /* IDA 0x1f785; and     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))&((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_1F787: ; /* IDA 0x1f787; jz      short near ptr byte_1F790 */
    if(m->flags&0x40) goto L_1F790;
L_1F789: ; /* IDA 0x1f789; call    far ptr 0:0 */
    gl_far(m,0xf789,0xf78e,h);
L_1F78E: ; /* IDA 0x1f78e; jmp     short loc_1F795 */
    goto L_1F795;
L_1F790: ; /* IDA 0x1f790; call    far ptr 0:0 */
    gl_far(m,0xf790,0xf795,h);
L_1F795: ; /* IDA 0x1f795; call    sub_1F7A4 */
    gl_call(m,63396,63384,h);
L_1F798: ; /* IDA 0x1f798; add     sp, 20h */
    m->sp=(uint16_t)(ec_math(m,m->sp,32,16,0,0));
L_1F79B: ; /* IDA 0x1f79b; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1F79C: ; /* IDA 0x1f79c; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1F79D: ; /* IDA 0x1f79d; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1F79E: ; /* IDA 0x1f79e; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1F79F: ; /* IDA 0x1f79f; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1F7A0: ; /* IDA 0x1f7a0; retn */
    return;
}
/* Original loc_10701; IDA 0x10701; raw control; source re/108 */
static void gl_body_10701(KiMachine16 *m,const KiEconomyHooks *h) {
L_10701: ; /* IDA 0x10701; push    ax */
    ec_push(m,m->ax);
L_10702: ; /* IDA 0x10702; push    cx */
    ec_push(m,m->cx);
L_10703: ; /* IDA 0x10703; push    di */
    ec_push(m,m->di);
L_10704: ; /* IDA 0x10704; push    bp */
    ec_push(m,m->bp);
L_10705: ; /* IDA 0x10705; mov     cs:word_1071F, ax */
    dl_write16(m,m->cs,1823,m->ax);
L_10709: ; /* IDA 0x10709; mov     cs:word_10736, ax */
    dl_write16(m,m->cs,1846,m->ax);
L_1070D: ; /* IDA 0x1070d; mov     cs:word_1074D, ax */
    dl_write16(m,m->cs,1869,m->ax);
L_10711: ; /* IDA 0x10711; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_10713: ; /* IDA 0x10713; mov     bp, bx */
    m->bp=(uint16_t)(m->bx);
L_10715: ; /* IDA 0x10715; mov     ax, [si] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_10717: ; /* IDA 0x10717; xchg    al, ah */
    {uint16_t a=(uint8_t)m->ax,b=(uint8_t)(m->ax>>8);vg_al(m,(uint8_t)(b));vg_ah(m,(uint8_t)(a));}
L_10719: ; /* IDA 0x10719; call    sub_1F878 */
    gl_call(m,63608,1820,h);
L_1071C: ; /* IDA 0x1071c; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_1071E: ; /* IDA 0x1071e; mov     ax, 9001h */
#if KI_GLYPH_MUTATION == 8
    m->ax=0x9001;
#else
    m->ax=dl_read16(m,m->cs,0x71f);
#endif
L_10721: ; /* IDA 0x10721; call    loc_1F75E */
    gl_call(m,63326,1828,h);
L_10724: ; /* IDA 0x10724; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_10726: ; /* IDA 0x10726; mov     bx, bp */
    m->bx=(uint16_t)(m->bp);
L_10728: ; /* IDA 0x10728; add     dx, 10h */
#if KI_GLYPH_MUTATION == 9
    m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,0));
#else
    m->dx=(uint16_t)(ec_math(m,m->dx,16,16,0,0));
#endif
L_1072B: ; /* IDA 0x1072b; mov     ax, [si+2] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+2)));
L_1072E: ; /* IDA 0x1072e; xchg    al, ah */
    {uint16_t a=(uint8_t)m->ax,b=(uint8_t)(m->ax>>8);vg_al(m,(uint8_t)(b));vg_ah(m,(uint8_t)(a));}
L_10730: ; /* IDA 0x10730; call    sub_1F878 */
    gl_call(m,63608,1843,h);
L_10733: ; /* IDA 0x10733; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_10735: ; /* IDA 0x10735; mov     ax, 9001h */
#if KI_GLYPH_MUTATION == 8
    m->ax=0x9001;
#else
    m->ax=dl_read16(m,m->cs,0x736);
#endif
L_10738: ; /* IDA 0x10738; call    loc_1F75E */
    gl_call(m,63326,1851,h);
L_1073B: ; /* IDA 0x1073b; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_1073D: ; /* IDA 0x1073d; mov     bx, bp */
    m->bx=(uint16_t)(m->bp);
L_1073F: ; /* IDA 0x1073f; add     dx, 20h ; ' ' */
    m->dx=(uint16_t)(ec_math(m,m->dx,32,16,0,0));
L_10742: ; /* IDA 0x10742; mov     ax, [si+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+4)));
L_10745: ; /* IDA 0x10745; xchg    al, ah */
    {uint16_t a=(uint8_t)m->ax,b=(uint8_t)(m->ax>>8);vg_al(m,(uint8_t)(b));vg_ah(m,(uint8_t)(a));}
L_10747: ; /* IDA 0x10747; call    sub_1F878 */
    gl_call(m,63608,1866,h);
L_1074A: ; /* IDA 0x1074a; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_1074C: ; /* IDA 0x1074c; mov     ax, 9001h */
#if KI_GLYPH_MUTATION == 8
    m->ax=0x9001;
#else
    m->ax=dl_read16(m,m->cs,0x74d);
#endif
L_1074F: ; /* IDA 0x1074f; call    loc_1F75E */
    gl_call(m,63326,1874,h);
L_10752: ; /* IDA 0x10752; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_10754: ; /* IDA 0x10754; mov     bx, bp */
    m->bx=(uint16_t)(m->bp);
L_10756: ; /* IDA 0x10756; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_10757: ; /* IDA 0x10757; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_10758: ; /* IDA 0x10758; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_10759: ; /* IDA 0x10759; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1075A: ; /* IDA 0x1075a; retn */
    return;
}
int ki_glyph_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0xf720:gl_body_1F720(m,h);break;
case 0xf7a4:gl_body_1F7A4(m,h);break;
case 0x6f5:gl_body_106F5(m,h);break;
case 0x6fd:gl_body_106FD(m,h);break;
case 0xf75e:gl_body_1F75E(m,h);break;
case 0x701:gl_body_10701(m,h);break;
default:return 0;}return 1;}
