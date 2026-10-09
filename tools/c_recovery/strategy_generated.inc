/* Original native C strategy; re/124, spec/244. */
/* Original sub_102F5; IDA 0x102F5. */
static void sty_body_102F5(KiMachine16 *m,const KiEconomyHooks *h) {
L_102F5: ; /* IDA 0x102F5; cmp     cs:byte_1020F, al */
    ec_cmp(m,dl_read8(m,m->cs,527),(uint8_t)m->ax,8);
L_102FA: ; /* IDA 0x102FA; jnb     short loc_10306 */
    if(!(m->flags&1)) goto L_10306;
L_102FC: ; /* IDA 0x102FC; push    ax */
    ec_push(m,m->ax);
L_102FD: ; /* IDA 0x102FD; mov     ah, 0Ah */
    vg_ah(m,(uint8_t)(10));
L_102FF: ; /* IDA 0x102FF; int     61h; reserved for user interrupt */
    gl_interrupt(m,97,0x301);
L_10301: ; /* IDA 0x10301; test    al, 4 */
    {uint16_t result=(((uint8_t)m->ax)&(4))&255;ec_logic(m,result,8);}
L_10303: ; /* IDA 0x10303; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10304: ; /* IDA 0x10304; jnz     short locret_1030E */
    if(!(m->flags&0x40)) goto L_1030E;
L_10306: ; /* IDA 0x10306; mov     cs:byte_1020F, al */
    dl_write8(m,m->cs,527,(uint8_t)m->ax);
L_1030A: ; /* IDA 0x1030A; mov     ah, 5 */
    vg_ah(m,(uint8_t)(5));
L_1030C: ; /* IDA 0x1030C; int     61h; reserved for user interrupt */
    gl_interrupt(m,97,0x30e);
L_1030E: ; /* IDA 0x1030E; retn */
    return;
}
/* Original sub_1300E; IDA 0x1300E. */
static void sty_body_1300E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1300E: ; /* IDA 0x1300E; push    cx */
    ec_push(m,m->cx);
L_1300F: ; /* IDA 0x1300F; mov     cx, si */
    m->cx=(uint16_t)(m->si);
L_13011: ; /* IDA 0x13011; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_13013: ; /* IDA 0x13013; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_13015: ; /* IDA 0x13015; mov     ah, ch */
    vg_ah(m,(uint8_t)((uint8_t)(m->cx>>8)));
L_13017: ; /* IDA 0x13017; call    sub_1301C */
    sty_call(m,12316,12314,h);
L_1301A: ; /* IDA 0x1301A; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1301B: ; /* IDA 0x1301B; retn */
    return;
}
/* Original sub_1304E; IDA 0x1304E. */
static void sty_body_1304E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1304E: ; /* IDA 0x1304E; push    ds */
    ec_push(m,m->ds);
L_1304F: ; /* IDA 0x1304F; push    ax */
    ec_push(m,m->ax);
L_13050: ; /* IDA 0x13050; push    bx */
    ec_push(m,m->bx);
L_13051: ; /* IDA 0x13051; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_13053: ; /* IDA 0x13053; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_13055: ; /* IDA 0x13055; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_13057: ; /* IDA 0x13057; mov     ah, bh */
    vg_ah(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_13059: ; /* IDA 0x13059; jmp     short loc_1305E */
    goto L_1305E;
L_1305B: ; /* IDA 0x1305B; push    ds */
    ec_push(m,m->ds);
L_1305C: ; /* IDA 0x1305C; push    ax */
    ec_push(m,m->ax);
L_1305D: ; /* IDA 0x1305D; push    bx */
    ec_push(m,m->bx);
L_1305E: ; /* IDA 0x1305E; mov     ds, cs:word_10D56 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3414));
L_13063: ; /* IDA 0x13063; mov     bx, cs:word_10D20 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,3360));
L_13068: ; /* IDA 0x13068; cmp     [bx], ax */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->bx)),m->ax,16);
L_1306A: ; /* IDA 0x1306A; jnz     short loc_13080 */
    if(!(m->flags&0x40)) goto L_13080;
L_1306C: ; /* IDA 0x1306C; cmp     dl, 0FFh */
    ec_cmp(m,(uint8_t)m->dx,255,8);
L_1306F: ; /* IDA 0x1306F; jz      short loc_13076 */
    if(m->flags&0x40) goto L_13076;
L_13071: ; /* IDA 0x13071; cmp     [bx+2], dl */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+2)),(uint8_t)m->dx,8);
L_13074: ; /* IDA 0x13074; jnz     short loc_13080 */
    if(!(m->flags&0x40)) goto L_13080;
L_13076: ; /* IDA 0x13076; cmp     dh, 0FFh */
    ec_cmp(m,(uint8_t)(m->dx>>8),255,8);
L_13079: ; /* IDA 0x13079; jz      short loc_1308C */
    if(m->flags&0x40) goto L_1308C;
L_1307B: ; /* IDA 0x1307B; cmp     [bx+3], dl */
#if KI_STRATEGY_MUTATION == 3
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+3)),(uint8_t)(m->dx>>8),8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+3)),(uint8_t)m->dx,8);
#endif
L_1307E: ; /* IDA 0x1307E; jz      short loc_1308C */
    if(m->flags&0x40) goto L_1308C;
L_13080: ; /* IDA 0x13080; add     bx, 4 */
    m->bx=(uint16_t)(ec_math(m,m->bx,4,16,0,0));
L_13083: ; /* IDA 0x13083; cmp     bx, 400h */
    ec_cmp(m,m->bx,1024,16);
L_13087: ; /* IDA 0x13087; jb      short loc_13068 */
    if(m->flags&1) goto L_13068;
L_13089: ; /* IDA 0x13089; stc */
    m->flags|=1;
L_1308A: ; /* IDA 0x1308A; jmp     short loc_1308D */
    goto L_1308D;
L_1308C: ; /* IDA 0x1308C; clc */
    m->flags&=(uint16_t)~1u;
L_1308D: ; /* IDA 0x1308D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1308E: ; /* IDA 0x1308E; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1308F: ; /* IDA 0x1308F; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_13090: ; /* IDA 0x13090; retn */
    return;
}
/* Original sub_13830; IDA 0x13830. */
static void sty_body_13830(KiMachine16 *m,const KiEconomyHooks *h) {
L_13830: ; /* IDA 0x13830; push    ds */
    ec_push(m,m->ds);
L_13831: ; /* IDA 0x13831; push    es */
    ec_push(m,m->es);
L_13832: ; /* IDA 0x13832; push    bp */
    ec_push(m,m->bp);
L_13833: ; /* IDA 0x13833; call    sub_102C2 */
    sty_call(m,706,14390,h);
L_13836: ; /* IDA 0x13836; call    sub_13C1E */
    sty_call(m,15390,14393,h);
L_13839: ; /* IDA 0x13839; sub     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,1));
L_1383C: ; /* IDA 0x1383C; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_1383E: ; /* IDA 0x1383E; mov     [bp+0], cx */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->cx);
L_13841: ; /* IDA 0x13841; mov     [bp+2], ax */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->ax);
L_13844: ; /* IDA 0x13844; mov     [bp+4], dl */
    dl_write8(m,m->ss,(uint16_t)(m->bp+4),(uint8_t)m->dx);
L_13847: ; /* IDA 0x13847; mov     al, 0 */
    vg_al(m,(uint8_t)(0));
L_13849: ; /* IDA 0x13849; call    sub_13D09 */
    sty_call(m,15625,14412,h);
L_1384C: ; /* IDA 0x1384C; mov     al, 6 */
    vg_al(m,(uint8_t)(6));
L_1384E: ; /* IDA 0x1384E; call    sub_10241 */
    sty_call(m,577,14417,h);
L_13851: ; /* IDA 0x13851; mov     cx, [bp+0] */
    m->cx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_13854: ; /* IDA 0x13854; call    sub_13C99 */
    sty_call(m,15513,14423,h);
L_13857: ; /* IDA 0x13857; mov     cx, [bp+0] */
    m->cx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1385A: ; /* IDA 0x1385A; add     cx, 3 */
    m->cx=(uint16_t)(ec_math(m,m->cx,3,16,0,0));
L_1385D: ; /* IDA 0x1385D; call    sub_13CDC */
    sty_call(m,15580,14432,h);
L_13860: ; /* IDA 0x13860; mov     al, [bp+2] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+2))));
L_13863: ; /* IDA 0x13863; mov     cl, al */
    vg_cl(m,(uint8_t)((uint8_t)m->ax));
L_13865: ; /* IDA 0x13865; shl     cl, 1 */
    vg_cl(m,(uint8_t)(dl_shift(m,(uint8_t)m->cx,1,8,4)));
L_13867: ; /* IDA 0x13867; add     cl, al */
    vg_cl(m,(uint8_t)(ec_math(m,(uint8_t)m->cx,(uint8_t)m->ax,8,0,0)));
L_13869: ; /* IDA 0x13869; xor     ch, ch */
    {uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
L_1386B: ; /* IDA 0x1386B; add     cx, [bp+0] */
    m->cx=(uint16_t)(ec_math(m,m->cx,dl_read16(m,m->ss,(uint16_t)(m->bp)),16,0,0));
L_1386E: ; /* IDA 0x1386E; add     cx, 4 */
    m->cx=(uint16_t)(ec_math(m,m->cx,4,16,0,0));
L_13871: ; /* IDA 0x13871; cmp     al, 3 */
    ec_cmp(m,(uint8_t)m->ax,3,8);
L_13873: ; /* IDA 0x13873; jb      short loc_1387C */
    if(m->flags&1) goto L_1387C;
L_13875: ; /* IDA 0x13875; jz      short loc_13894 */
    if(m->flags&0x40) goto L_13894;
L_13877: ; /* IDA 0x13877; mov     cx, 53h ; 'S' */
    m->cx=(uint16_t)(83);
L_1387A: ; /* IDA 0x1387A; jmp     short loc_1389E */
    goto L_1389E;
L_1387C: ; /* IDA 0x1387C; cmp     al, 1 */
    ec_cmp(m,(uint8_t)m->ax,1,8);
L_1387E: ; /* IDA 0x1387E; jb      short loc_13894 */
    if(m->flags&1) goto L_13894;
L_13880: ; /* IDA 0x13880; jz      short loc_138A4 */
    if(m->flags&0x40) goto L_138A4;
L_13882: ; /* IDA 0x13882; call    sub_13C99 */
    sty_call(m,15513,14469,h);
L_13885: ; /* IDA 0x13885; add     word ptr [bp+0], 10h */
    dl_write16(m,m->ss,(uint16_t)(m->bp),ec_math(m,dl_read16(m,m->ss,(uint16_t)(m->bp)),16,16,0,0));
L_13889: ; /* IDA 0x13889; call    sub_13B5A */
    sty_call(m,15194,14476,h);
L_1388C: ; /* IDA 0x1388C; jb      short loc_138B2 */
    if(m->flags&1) goto L_138B2;
L_1388E: ; /* IDA 0x1388E; mov     al, 14h */
    vg_al(m,(uint8_t)(20));
L_13890: ; /* IDA 0x13890; shr     al, 1 */
#if KI_STRATEGY_MUTATION == 16
    ; /* lost accepted-reason trust halving */
#else
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,5)));
#endif
L_13892: ; /* IDA 0x13892; jmp     short loc_138A9 */
    goto L_138A9;
L_13894: ; /* IDA 0x13894; mov     al, 14h */
    vg_al(m,(uint8_t)(20));
L_13896: ; /* IDA 0x13896; push    cx */
    ec_push(m,m->cx);
L_13897: ; /* IDA 0x13897; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_1389A: ; /* IDA 0x1389A; call    sub_13DC9 */
    sty_call(m,15817,14493,h);
L_1389D: ; /* IDA 0x1389D; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1389E: ; /* IDA 0x1389E; call    sub_13C99 */
    sty_call(m,15513,14497,h);
L_138A1: ; /* IDA 0x138A1; stc */
    m->flags|=1;
L_138A2: ; /* IDA 0x138A2; jmp     short loc_138B2 */
    goto L_138B2;
L_138A4: ; /* IDA 0x138A4; call    sub_13C99 */
    sty_call(m,15513,14503,h);
L_138A7: ; /* IDA 0x138A7; mov     al, 14h */
    vg_al(m,(uint8_t)(20));
L_138A9: ; /* IDA 0x138A9; push    cx */
    ec_push(m,m->cx);
L_138AA: ; /* IDA 0x138AA; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_138AD: ; /* IDA 0x138AD; call    sub_13D91 */
    sty_call(m,15761,14512,h);
L_138B0: ; /* IDA 0x138B0; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_138B1: ; /* IDA 0x138B1; clc */
    m->flags&=(uint16_t)~1u;
L_138B2: ; /* IDA 0x138B2; lahf */
    vg_ah(m,(uint8_t)((m->flags&0xd5)|2));
L_138B3: ; /* IDA 0x138B3; call    sub_102C2 */
    sty_call(m,706,14518,h);
L_138B6: ; /* IDA 0x138B6; call    sub_13D45 */
    sty_call(m,15685,14521,h);
L_138B9: ; /* IDA 0x138B9; call    sub_11D46 */
    sty_call(m,7494,14524,h);
L_138BC: ; /* IDA 0x138BC; add     sp, 6 */
    m->sp=(uint16_t)(ec_math(m,m->sp,6,16,0,0));
L_138BF: ; /* IDA 0x138BF; call    sub_19321 */
    sty_call(m,37665,14530,h);
L_138C2: ; /* IDA 0x138C2; sahf */
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
L_138C3: ; /* IDA 0x138C3; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_138C4: ; /* IDA 0x138C4; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_138C5: ; /* IDA 0x138C5; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_138C6: ; /* IDA 0x138C6; retn */
    return;
}
/* Original sub_13B5A; IDA 0x13B5A. */
static void sty_body_13B5A(KiMachine16 *m,const KiEconomyHooks *h) {
L_13B5A: ; /* IDA 0x13B5A; mov     byte ptr [bp+5], 0 */
    dl_write8(m,m->ss,(uint16_t)(m->bp+5),0);
L_13B5E: ; /* IDA 0x13B5E; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_13B60: ; /* IDA 0x13B60; call    sub_13B7E */
    sty_call(m,15230,15203,h);
L_13B63: ; /* IDA 0x13B63; push    ax */
    ec_push(m,m->ax);
L_13B64: ; /* IDA 0x13B64; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_13B66: ; /* IDA 0x13B66; add     ax, [bp+0] */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->ss,(uint16_t)(m->bp)),16,0,0));
L_13B69: ; /* IDA 0x13B69; inc     ax */
    m->ax=(uint16_t)(ec_inc(m,m->ax,16));
L_13B6A: ; /* IDA 0x13B6A; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_13B6C: ; /* IDA 0x13B6C; call    sub_13CDC */
    sty_call(m,15580,15215,h);
L_13B6F: ; /* IDA 0x13B6F; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_13B70: ; /* IDA 0x13B70; call    sub_13BA9 */
    sty_call(m,15273,15219,h);
L_13B73: ; /* IDA 0x13B73; jnb     short loc_13B5E */
    if(!(m->flags&1)) goto L_13B5E;
L_13B75: ; /* IDA 0x13B75; cmp     ah, 1 */
    ec_cmp(m,(uint8_t)(m->ax>>8),1,8);
L_13B78: ; /* IDA 0x13B78; jz      short loc_13B7C */
    if(m->flags&0x40) goto L_13B7C;
L_13B7A: ; /* IDA 0x13B7A; stc */
    m->flags|=1;
L_13B7B: ; /* IDA 0x13B7B; retn */
    return;
L_13B7C: ; /* IDA 0x13B7C; clc */
    m->flags&=(uint16_t)~1u;
L_13B7D: ; /* IDA 0x13B7D; retn */
    return;
}
/* Original sub_13BA9; IDA 0x13BA9. */
static void sty_body_13BA9(KiMachine16 *m,const KiEconomyHooks *h) {
L_13BA9: ; /* IDA 0x13BA9; mov     cl, al */
    vg_cl(m,(uint8_t)((uint8_t)m->ax));
L_13BAB: ; /* IDA 0x13BAB; inc     cl */
    vg_cl(m,(uint8_t)(ec_inc(m,(uint8_t)m->cx,8)));
L_13BAD: ; /* IDA 0x13BAD; cmp     al, 4 */
    ec_cmp(m,(uint8_t)m->ax,4,8);
L_13BAF: ; /* IDA 0x13BAF; jnz     short loc_13BB9 */
    if(!(m->flags&0x40)) goto L_13BB9;
L_13BB1: ; /* IDA 0x13BB1; mov     ah, 4 */
    vg_ah(m,(uint8_t)(4));
L_13BB3: ; /* IDA 0x13BB3; mov     bx, 2Ah ; '*' */
    m->bx=(uint16_t)(42);
L_13BB6: ; /* IDA 0x13BB6; stc */
    m->flags|=1;
L_13BB7: ; /* IDA 0x13BB7; jmp     short loc_13C13 */
    goto L_13C13;
L_13BB9: ; /* IDA 0x13BB9; mov     ch, [bp+5] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+5))));
L_13BBC: ; /* IDA 0x13BBC; shr     ch, cl */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),(uint8_t)m->cx,8,5)));
L_13BBE: ; /* IDA 0x13BBE; jnb     short loc_13BC5 */
    if(!(m->flags&1)) goto L_13BC5;
L_13BC0: ; /* IDA 0x13BC0; mov     ah, 3 */
    vg_ah(m,(uint8_t)(3));
L_13BC2: ; /* IDA 0x13BC2; clc */
    m->flags&=(uint16_t)~1u;
L_13BC3: ; /* IDA 0x13BC3; jmp     short loc_13BF1 */
    goto L_13BF1;
L_13BC5: ; /* IDA 0x13BC5; mov     ah, 80h */
    vg_ah(m,(uint8_t)(128));
L_13BC7: ; /* IDA 0x13BC7; rol     ah, cl */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),(uint8_t)m->cx,8,0)));
L_13BC9: ; /* IDA 0x13BC9; or      [bp+5], ah */
#if KI_STRATEGY_MUTATION == 6
    ; /* lost attempted reason bit */
#else
    {uint16_t result=((dl_read8(m,m->ss,(uint16_t)(m->bp+5)))|((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);dl_write8(m,m->ss,(uint16_t)(m->bp+5),result);}
#endif
L_13BCC: ; /* IDA 0x13BCC; mov     ch, [bp+4] */
    vg_ch(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp+4))));
L_13BCF: ; /* IDA 0x13BCF; shr     ch, cl */
    vg_ch(m,(uint8_t)(dl_shift(m,(uint8_t)(m->cx>>8),(uint8_t)m->cx,8,5)));
L_13BD1: ; /* IDA 0x13BD1; jb      short loc_13BE4 */
    if(m->flags&1) goto L_13BE4;
L_13BD3: ; /* IDA 0x13BD3; push    ax */
    ec_push(m,m->ax);
L_13BD4: ; /* IDA 0x13BD4; push    cx */
    ec_push(m,m->cx);
L_13BD5: ; /* IDA 0x13BD5; mov     al, 14h */
    vg_al(m,(uint8_t)(20));
L_13BD7: ; /* IDA 0x13BD7; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_13BDA: ; /* IDA 0x13BDA; call    sub_13DC9 */
    sty_call(m,15817,15325,h);
L_13BDD: ; /* IDA 0x13BDD; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_13BDE: ; /* IDA 0x13BDE; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_13BDF: ; /* IDA 0x13BDF; mov     ah, 0 */
    vg_ah(m,(uint8_t)(0));
L_13BE1: ; /* IDA 0x13BE1; stc */
    m->flags|=1;
L_13BE2: ; /* IDA 0x13BE2; jmp     short loc_13BF1 */
    goto L_13BF1;
L_13BE4: ; /* IDA 0x13BE4; dec     byte ptr [bp+3] */
#if KI_STRATEGY_MUTATION == 5
    ; /* lost reason budget decrement */
#else
    dl_write8(m,m->ss,(uint16_t)(m->bp+3),st_dec(m,dl_read8(m,m->ss,(uint16_t)(m->bp+3)),8));
#endif
L_13BE7: ; /* IDA 0x13BE7; jz      short loc_13BEE */
    if(m->flags&0x40) goto L_13BEE;
L_13BE9: ; /* IDA 0x13BE9; mov     ah, 2 */
    vg_ah(m,(uint8_t)(2));
L_13BEB: ; /* IDA 0x13BEB; clc */
    m->flags&=(uint16_t)~1u;
L_13BEC: ; /* IDA 0x13BEC; jmp     short loc_13BF1 */
    goto L_13BF1;
L_13BEE: ; /* IDA 0x13BEE; mov     ah, 1 */
    vg_ah(m,(uint8_t)(1));
L_13BF0: ; /* IDA 0x13BF0; stc */
#if KI_STRATEGY_MUTATION == 7
    m->flags&=(uint16_t)~1u;
#else
    m->flags|=1;
#endif
L_13BF1: ; /* IDA 0x13BF1; pushf */
    ec_push(m,m->flags);
L_13BF2: ; /* IDA 0x13BF2; mov     bx, 2Dh ; '-' */
    m->bx=(uint16_t)(45);
L_13BF5: ; /* IDA 0x13BF5; cmp     ah, 3 */
    ec_cmp(m,(uint8_t)(m->ax>>8),3,8);
L_13BF8: ; /* IDA 0x13BF8; jnb     short loc_13C14 */
    if(!(m->flags&1)) goto L_13C14;
L_13BFA: ; /* IDA 0x13BFA; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_13BFC: ; /* IDA 0x13BFC; shl     bl, 1 */
    gw_bl(m,(uint8_t)(dl_shift(m,(uint8_t)m->bx,1,8,4)));
L_13BFE: ; /* IDA 0x13BFE; shl     bl, 1 */
    gw_bl(m,(uint8_t)(dl_shift(m,(uint8_t)m->bx,1,8,4)));
L_13C00: ; /* IDA 0x13C00; shl     bl, 1 */
    gw_bl(m,(uint8_t)(dl_shift(m,(uint8_t)m->bx,1,8,4)));
L_13C02: ; /* IDA 0x13C02; add     bl, al */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,(uint8_t)m->ax,8,0,0)));
L_13C04: ; /* IDA 0x13C04; mov     bh, ah */
    gw_bh(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_13C06: ; /* IDA 0x13C06; shl     bh, 1 */
    gw_bh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->bx>>8),1,8,4)));
L_13C08: ; /* IDA 0x13C08; add     bh, ah */
    gw_bh(m,(uint8_t)(ec_math(m,(uint8_t)(m->bx>>8),(uint8_t)(m->ax>>8),8,0,0)));
L_13C0A: ; /* IDA 0x13C0A; add     bl, bh */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,(uint8_t)(m->bx>>8),8,0,0)));
L_13C0C: ; /* IDA 0x13C0C; add     bl, 6 */
    gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,6,8,0,0)));
L_13C0F: ; /* IDA 0x13C0F; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_13C11: ; /* IDA 0x13C11; jmp     short loc_13C14 */
    goto L_13C14;
L_13C13: ; /* IDA 0x13C13; pushf */
    ec_push(m,m->flags);
L_13C14: ; /* IDA 0x13C14; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_13C16: ; /* IDA 0x13C16; add     cx, [bp+0] */
    m->cx=(uint16_t)(ec_math(m,m->cx,dl_read16(m,m->ss,(uint16_t)(m->bp)),16,0,0));
L_13C19: ; /* IDA 0x13C19; call    sub_13C99 */
    sty_call(m,15513,15388,h);
L_13C1C: ; /* IDA 0x13C1C; popf */
    m->flags=ec_pop(m)|2;
L_13C1D: ; /* IDA 0x13C1D; retn */
    return;
}
/* Original sub_13C1E; IDA 0x13C1E. */
static void sty_body_13C1E(KiMachine16 *m,const KiEconomyHooks *h) {
L_13C1E: ; /* IDA 0x13C1E; push    dx */
    ec_push(m,m->dx);
L_13C1F: ; /* IDA 0x13C1F; mov     dl, cs:byte_10D00 */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,3328)));
L_13C24: ; /* IDA 0x13C24; mov     ah, 1 */
    vg_ah(m,(uint8_t)(1));
L_13C26: ; /* IDA 0x13C26; cmp     dl, 0E0h */
#if KI_STRATEGY_MUTATION == 4
    ec_cmp(m,(uint8_t)m->dx,0xe1,8);
#else
    ec_cmp(m,(uint8_t)m->dx,224,8);
#endif
L_13C29: ; /* IDA 0x13C29; jnb     short loc_13C3B */
    if(!(m->flags&1)) goto L_13C3B;
L_13C2B: ; /* IDA 0x13C2B; mov     ah, 2 */
    vg_ah(m,(uint8_t)(2));
L_13C2D: ; /* IDA 0x13C2D; cmp     dl, 90h */
    ec_cmp(m,(uint8_t)m->dx,144,8);
L_13C30: ; /* IDA 0x13C30; jnb     short loc_13C3B */
    if(!(m->flags&1)) goto L_13C3B;
L_13C32: ; /* IDA 0x13C32; mov     ah, 3 */
    vg_ah(m,(uint8_t)(3));
L_13C34: ; /* IDA 0x13C34; cmp     dl, 20h ; ' ' */
    ec_cmp(m,(uint8_t)m->dx,32,8);
L_13C37: ; /* IDA 0x13C37; jnb     short loc_13C3B */
    if(!(m->flags&1)) goto L_13C3B;
L_13C39: ; /* IDA 0x13C39; mov     ah, 4 */
    vg_ah(m,(uint8_t)(4));
L_13C3B: ; /* IDA 0x13C3B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_13C3C: ; /* IDA 0x13C3C; retn */
    return;
}
/* Original sub_13D91; IDA 0x13D91. */
static void sty_body_13D91(KiMachine16 *m,const KiEconomyHooks *h) {
L_13D91: ; /* IDA 0x13D91; push    ds */
    ec_push(m,m->ds);
L_13D92: ; /* IDA 0x13D92; push    ax */
    ec_push(m,m->ax);
L_13D93: ; /* IDA 0x13D93; push    bx */
    ec_push(m,m->bx);
L_13D94: ; /* IDA 0x13D94; push    cx */
    ec_push(m,m->cx);
L_13D95: ; /* IDA 0x13D95; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_13D9A: ; /* IDA 0x13D9A; add     cs:byte_10D00, al */
    dl_write8(m,m->cs,3328,ec_math(m,dl_read8(m,m->cs,3328),(uint8_t)m->ax,8,0,0));
L_13D9F: ; /* IDA 0x13D9F; jnb     short loc_13DA7 */
    if(!(m->flags&1)) goto L_13DA7;
L_13DA1: ; /* IDA 0x13DA1; mov     cs:byte_10D00, 0FFh */
#if KI_STRATEGY_MUTATION == 8
    dl_write8(m,m->cs,0xd00,0xfe);
#else
    dl_write8(m,m->cs,3328,255);
#endif
L_13DA7: ; /* IDA 0x13DA7; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_13DA9: ; /* IDA 0x13DA9; call    sub_102F5 */
    sty_call(m,757,15788,h);
L_13DAC: ; /* IDA 0x13DAC; cmp     cx, 0FFFFh */
    ec_cmp(m,m->cx,65535,16);
L_13DAF: ; /* IDA 0x13DAF; jz      short loc_13DBF */
    if(m->flags&0x40) goto L_13DBF;
L_13DB1: ; /* IDA 0x13DB1; call    sub_187FF */
    sty_call(m,34815,15796,h);
L_13DB4: ; /* IDA 0x13DB4; mov     ah, [bx+425Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16990))));
L_13DB8: ; /* IDA 0x13DB8; mov     al, [bx+4241h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16961))));
L_13DBC: ; /* IDA 0x13DBC; call    sub_18810 */
    sty_call(m,34832,15807,h);
L_13DBF: ; /* IDA 0x13DBF; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_13DC1: ; /* IDA 0x13DC1; call    sub_15E80 */
    sty_call(m,24192,15812,h);
L_13DC4: ; /* IDA 0x13DC4; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_13DC5: ; /* IDA 0x13DC5; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_13DC6: ; /* IDA 0x13DC6; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_13DC7: ; /* IDA 0x13DC7; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_13DC8: ; /* IDA 0x13DC8; retn */
    return;
}
/* Original sub_161CA; IDA 0x161CA. */
static void sty_body_161CA(KiMachine16 *m,const KiEconomyHooks *h) {
L_161CA: ; /* IDA 0x161CA; mov     ax, cx */
    m->ax=(uint16_t)(m->cx);
L_161CC: ; /* IDA 0x161CC; sub     ax, 18h */
#if KI_STRATEGY_MUTATION == 1
    m->ax=ec_math(m,m->ax,23,16,0,1);
#else
    m->ax=(uint16_t)(ec_math(m,m->ax,24,16,0,1));
#endif
L_161CF: ; /* IDA 0x161CF; jb      short locret_16213 */
    if(m->flags&1) goto L_16213;
L_161D1: ; /* IDA 0x161D1; cmp     ax, 180h */
    ec_cmp(m,m->ax,384,16);
L_161D4: ; /* IDA 0x161D4; jnb     short locret_16213 */
    if(!(m->flags&1)) goto L_16213;
L_161D6: ; /* IDA 0x161D6; mov     cl, 30h ; '0' */
    vg_cl(m,(uint8_t)(48));
L_161D8: ; /* IDA 0x161D8; div     cl */
    dl_div(m,(uint8_t)m->cx,8);
L_161DA: ; /* IDA 0x161DA; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_161DC: ; /* IDA 0x161DC; mov     dx, ax */
    m->dx=(uint16_t)(m->ax);
L_161DE: ; /* IDA 0x161DE; mul     cl */
    dl_mul(m,(uint8_t)m->cx,8);
L_161E0: ; /* IDA 0x161E0; xchg    ax, dx */
    {uint16_t a=m->ax,b=m->dx;m->ax=(uint16_t)(b);m->dx=(uint16_t)(a);}
L_161E1: ; /* IDA 0x161E1; add     dx, 18h */
    m->dx=(uint16_t)(ec_math(m,m->dx,24,16,0,0));
L_161E4: ; /* IDA 0x161E4; mov     bx, 28h ; '(' */
    m->bx=(uint16_t)(40);
L_161E7: ; /* IDA 0x161E7; mov     si, 30h ; '0' */
    m->si=(uint16_t)(48);
L_161EA: ; /* IDA 0x161EA; mov     di, 10h */
    m->di=(uint16_t)(16);
L_161ED: ; /* IDA 0x161ED; push    dx */
    ec_push(m,m->dx);
L_161EE: ; /* IDA 0x161EE; push    bx */
    ec_push(m,m->bx);
L_161EF: ; /* IDA 0x161EF; push    si */
    ec_push(m,m->si);
L_161F0: ; /* IDA 0x161F0; push    di */
    ec_push(m,m->di);
L_161F1: ; /* IDA 0x161F1; call    sub_101B4 */
    sty_call(m,436,25076,h);
L_161F4: ; /* IDA 0x161F4; call    sub_10B46 */
    sty_call(m,2886,25079,h);
L_161F7: ; /* IDA 0x161F7; call    sub_101DB */
    sty_call(m,475,25082,h);
L_161FA: ; /* IDA 0x161FA; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_161FC: ; /* IDA 0x161FC; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_161FE: ; /* IDA 0x161FE; call    cs:funcs_161FE[bx] */
    sty_call(m,dl_read16(m,m->cs,(uint16_t)(m->bx+25108)),25091,h);
L_16203: ; /* IDA 0x16203; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_16204: ; /* IDA 0x16204; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16205: ; /* IDA 0x16205; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16206: ; /* IDA 0x16206; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_16207: ; /* IDA 0x16207; call    sub_101B4 */
    sty_call(m,436,25098,h);
L_1620A: ; /* IDA 0x1620A; call    sub_11C8D */
    sty_call(m,7309,25101,h);
L_1620D: ; /* IDA 0x1620D; call    sub_10B46 */
    sty_call(m,2886,25104,h);
L_16210: ; /* IDA 0x16210; call    sub_101DB */
    sty_call(m,475,25107,h);
L_16213: ; /* IDA 0x16213; retn */
    return;
}
/* Original sub_16224; IDA 0x16224. */
static void sty_body_16224(KiMachine16 *m,const KiEconomyHooks *h) {
L_16224: ; /* IDA 0x16224; mov     ax, ds */
    m->ax=(uint16_t)(m->ds);
L_16226: ; /* IDA 0x16226; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_1622B: ; /* IDA 0x1622B; call    sub_187FF */
    sty_call(m,34815,25134,h);
L_1622E: ; /* IDA 0x1622E; cmp     byte ptr [bx+4257h], 0 */
#if KI_STRATEGY_MUTATION == 2
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+0x4257)),1,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+16983)),0,8);
#endif
L_16233: ; /* IDA 0x16233; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_16235: ; /* IDA 0x16235; jz      short loc_16241 */
    if(m->flags&0x40) goto L_16241;
L_16237: ; /* IDA 0x16237; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_16239: ; /* IDA 0x16239; mov     cx, 40h ; '@' */
    m->cx=(uint16_t)(64);
L_1623C: ; /* IDA 0x1623C; call    sub_18810 */
    sty_call(m,34832,25151,h);
L_1623F: ; /* IDA 0x1623F; jmp     short locret_1625A */
    goto L_1625A;
L_16241: ; /* IDA 0x16241; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_16244: ; /* IDA 0x16244; mov     cx, 4Dh ; 'M' */
    m->cx=(uint16_t)(77);
L_16247: ; /* IDA 0x16247; mov     dx, 400h */
    m->dx=(uint16_t)(1024);
L_1624A: ; /* IDA 0x1624A; call    sub_193E9 */
    sty_call(m,37865,25165,h);
L_1624D: ; /* IDA 0x1624D; jb      short locret_1625A */
    if(m->flags&1) goto L_1625A;
L_1624F: ; /* IDA 0x1624F; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_16251: ; /* IDA 0x16251; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_16253: ; /* IDA 0x16253; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_16255: ; /* IDA 0x16255; call    cs:funcs_16255[bx] */
    sty_call(m,dl_read16(m,m->cs,(uint16_t)(m->bx+25179)),25178,h);
L_1625A: ; /* IDA 0x1625A; retn */
    return;
}
/* Original sub_16288; IDA 0x16288. */
static void sty_body_16288(KiMachine16 *m,const KiEconomyHooks *h) {
L_16288: ; /* IDA 0x16288; mov     cx, 1 */
    m->cx=(uint16_t)(1);
L_1628B: ; /* IDA 0x1628B; call    sub_16C5E */
    sty_call(m,27742,25230,h);
L_1628E: ; /* IDA 0x1628E; retn */
    return;
}
/* Original sub_1628F; IDA 0x1628F. */
static void sty_body_1628F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1628F: ; /* IDA 0x1628F; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_16292: ; /* IDA 0x16292; mov     cx, 4Fh ; 'O' */
    m->cx=(uint16_t)(79);
L_16295: ; /* IDA 0x16295; mov     dx, 40Ch */
    m->dx=(uint16_t)(1036);
L_16298: ; /* IDA 0x16298; call    sub_193E9 */
    sty_call(m,37865,25243,h);
L_1629B: ; /* IDA 0x1629B; jnb     short loc_1629E */
    if(!(m->flags&1)) goto L_1629E;
L_1629D: ; /* IDA 0x1629D; retn */
    return;
L_1629E: ; /* IDA 0x1629E; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_162A0: ; /* IDA 0x162A0; jz      short loc_162A4 */
    if(m->flags&0x40) goto L_162A4;
L_162A2: ; /* IDA 0x162A2; jmp     short loc_162D2 */
    goto L_162D2;
L_162A4: ; /* IDA 0x162A4; call    sub_12078 */
    sty_call(m,8312,25255,h);
L_162A7: ; /* IDA 0x162A7; mov     cx, 16h */
    m->cx=(uint16_t)(22);
L_162AA: ; /* IDA 0x162AA; call    sub_18853 */
    sty_call(m,34899,25261,h);
L_162AD: ; /* IDA 0x162AD; call    sub_1716D */
    sty_call(m,29037,25264,h);
L_162B0: ; /* IDA 0x162B0; pushf */
    ec_push(m,m->flags);
L_162B1: ; /* IDA 0x162B1; call    sub_120D6 */
    sty_call(m,8406,25268,h);
L_162B4: ; /* IDA 0x162B4; popf */
    m->flags=ec_pop(m)|2;
L_162B5: ; /* IDA 0x162B5; jb      short loc_162CB */
    if(m->flags&1) goto L_162CB;
L_162B7: ; /* IDA 0x162B7; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_162BC: ; /* IDA 0x162BC; mov     dx, [bx+10h] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+16)));
L_162BF: ; /* IDA 0x162BF; mov     bx, [bx+12h] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+18)));
L_162C2: ; /* IDA 0x162C2; mov     ax, 14h */
    m->ax=(uint16_t)(20);
L_162C5: ; /* IDA 0x162C5; mov     cx, 0Ch */
    m->cx=(uint16_t)(12);
L_162C8: ; /* IDA 0x162C8; call    sub_12151 */
    sty_call(m,8529,25291,h);
L_162CB: ; /* IDA 0x162CB; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_162CE: ; /* IDA 0x162CE; call    sub_18853 */
    sty_call(m,34899,25297,h);
L_162D1: ; /* IDA 0x162D1; retn */
    return;
L_162D2: ; /* IDA 0x162D2; call    sub_12078 */
    sty_call(m,8312,25301,h);
L_162D5: ; /* IDA 0x162D5; mov     cx, 2 */
    m->cx=(uint16_t)(2);
L_162D8: ; /* IDA 0x162D8; call    sub_18853 */
    sty_call(m,34899,25307,h);
L_162DB: ; /* IDA 0x162DB; call    sub_1716D */
    sty_call(m,29037,25310,h);
L_162DE: ; /* IDA 0x162DE; jb      short loc_162F1 */
    if(m->flags&1) goto L_162F1;
L_162E0: ; /* IDA 0x162E0; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_162E2: ; /* IDA 0x162E2; sub     si, 2240h */
    m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,1));
L_162E6: ; /* IDA 0x162E6; call    sub_120D6 */
    sty_call(m,8406,25321,h);
L_162E9: ; /* IDA 0x162E9; call    sub_17F90 */
    sty_call(m,32656,25324,h);
L_162EC: ; /* IDA 0x162EC; call    sub_12078 */
    sty_call(m,8312,25327,h);
L_162EF: ; /* IDA 0x162EF; jmp     short loc_162D5 */
    goto L_162D5;
L_162F1: ; /* IDA 0x162F1; call    sub_120D6 */
    sty_call(m,8406,25332,h);
L_162F4: ; /* IDA 0x162F4; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_162F7: ; /* IDA 0x162F7; call    sub_18853 */
    sty_call(m,34899,25338,h);
L_162FA: ; /* IDA 0x162FA; retn */
    return;
}
/* Original sub_162FB; IDA 0x162FB. */
static void sty_body_162FB(KiMachine16 *m,const KiEconomyHooks *h) {
L_162FB: ; /* IDA 0x162FB; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_162FE: ; /* IDA 0x162FE; mov     cx, 52h ; 'R' */
    m->cx=(uint16_t)(82);
L_16301: ; /* IDA 0x16301; mov     dx, 40Fh */
    m->dx=(uint16_t)(1039);
L_16304: ; /* IDA 0x16304; call    sub_193E9 */
    sty_call(m,37865,25351,h);
L_16307: ; /* IDA 0x16307; jnb     short loc_1630A */
    if(!(m->flags&1)) goto L_1630A;
L_16309: ; /* IDA 0x16309; retn */
    return;
L_1630A: ; /* IDA 0x1630A; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1630C: ; /* IDA 0x1630C; jnz     short loc_16329 */
    if(!(m->flags&0x40)) goto L_16329;
L_1630E: ; /* IDA 0x1630E; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16313: ; /* IDA 0x16313; mov     bx, cs:word_10CFD */
    m->bx=(uint16_t)(dl_read16(m,m->cs,3325));
L_16318: ; /* IDA 0x16318; mov     bh, [bx+3] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_1631B: ; /* IDA 0x1631B; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_1631D: ; /* IDA 0x1631D; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1631F: ; /* IDA 0x1631F; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16321: ; /* IDA 0x16321; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16323: ; /* IDA 0x16323; add     bx, 840h */
    m->bx=(uint16_t)(ec_math(m,m->bx,2112,16,0,0));
L_16327: ; /* IDA 0x16327; jmp     short loc_1633C */
    goto L_1633C;
L_16329: ; /* IDA 0x16329; call    sub_12078 */
    sty_call(m,8312,25388,h);
L_1632C: ; /* IDA 0x1632C; mov     cx, 17h */
    m->cx=(uint16_t)(23);
L_1632F: ; /* IDA 0x1632F; call    sub_18853 */
    sty_call(m,34899,25394,h);
L_16332: ; /* IDA 0x16332; call    sub_17400 */
    sty_call(m,29696,25397,h);
L_16335: ; /* IDA 0x16335; pushf */
    ec_push(m,m->flags);
L_16336: ; /* IDA 0x16336; call    sub_120D6 */
    sty_call(m,8406,25401,h);
L_16339: ; /* IDA 0x16339; popf */
    m->flags=ec_pop(m)|2;
L_1633A: ; /* IDA 0x1633A; jb      short loc_1635F */
    if(m->flags&1) goto L_1635F;
L_1633C: ; /* IDA 0x1633C; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_1633E: ; /* IDA 0x1633E; sub     si, 840h */
    m->si=(uint16_t)(ec_math(m,m->si,2112,16,0,1));
L_16342: ; /* IDA 0x16342; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16347: ; /* IDA 0x16347; mov     dx, [bx+8] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+8)));
L_1634A: ; /* IDA 0x1634A; mov     bx, [bx+0Ah] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+10)));
L_1634D: ; /* IDA 0x1634D; mov     ax, 14h */
    m->ax=(uint16_t)(20);
L_16350: ; /* IDA 0x16350; mov     cx, 0Ch */
    m->cx=(uint16_t)(12);
L_16353: ; /* IDA 0x16353; call    sub_12151 */
    sty_call(m,8529,25430,h);
L_16356: ; /* IDA 0x16356; call    sub_11F7F */
    sty_call(m,8063,25433,h);
L_16359: ; /* IDA 0x16359; call    sub_11D46 */
    sty_call(m,7494,25436,h);
L_1635C: ; /* IDA 0x1635C; call    sub_17E1F */
    sty_call(m,32287,25439,h);
L_1635F: ; /* IDA 0x1635F; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16362: ; /* IDA 0x16362; call    sub_18853 */
    sty_call(m,34899,25445,h);
L_16365: ; /* IDA 0x16365; retn */
    return;
}
/* Original sub_16366; IDA 0x16366. */
static void sty_body_16366(KiMachine16 *m,const KiEconomyHooks *h) {
L_16366: ; /* IDA 0x16366; call    sub_12078 */
    sty_call(m,8312,25449,h);
L_16369: ; /* IDA 0x16369; mov     cx, 18h */
    m->cx=(uint16_t)(24);
L_1636C: ; /* IDA 0x1636C; call    sub_18853 */
    sty_call(m,34899,25455,h);
L_1636F: ; /* IDA 0x1636F; call    sub_175FA */
    sty_call(m,30202,25458,h);
L_16372: ; /* IDA 0x16372; pushf */
    ec_push(m,m->flags);
L_16373: ; /* IDA 0x16373; call    sub_120D6 */
    sty_call(m,8406,25462,h);
L_16376: ; /* IDA 0x16376; popf */
    m->flags=ec_pop(m)|2;
L_16377: ; /* IDA 0x16377; jnb     short loc_16380 */
    if(!(m->flags&1)) goto L_16380;
L_16379: ; /* IDA 0x16379; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_1637C: ; /* IDA 0x1637C; call    sub_18853 */
    sty_call(m,34899,25471,h);
L_1637F: ; /* IDA 0x1637F; retn */
    return;
L_16380: ; /* IDA 0x16380; push    ds */
    ec_push(m,m->ds);
L_16381: ; /* IDA 0x16381; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16386: ; /* IDA 0x16386; mov     cx, 1A8h */
    m->cx=(uint16_t)(424);
L_16389: ; /* IDA 0x16389; cmp     byte ptr [bx+1Dh], 0FFh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+29)),255,8);
L_1638D: ; /* IDA 0x1638D; jnz     short loc_163B3 */
    if(!(m->flags&0x40)) goto L_163B3;
L_1638F: ; /* IDA 0x1638F; mov     al, [bx+0Eh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+14))));
L_16392: ; /* IDA 0x16392; mov     ah, [bx+0Fh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+15))));
L_16395: ; /* IDA 0x16395; mov     cl, [bx+10h] */
    vg_cl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16))));
L_16398: ; /* IDA 0x16398; cmp     al, ah */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
L_1639A: ; /* IDA 0x1639A; jnb     short loc_1639E */
    if(!(m->flags&1)) goto L_1639E;
L_1639C: ; /* IDA 0x1639C; mov     al, ah */
    vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1639E: ; /* IDA 0x1639E; cmp     al, cl */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)m->cx,8);
L_163A0: ; /* IDA 0x163A0; jnb     short loc_163A4 */
#if KI_STRATEGY_MUTATION == 14
    if(m->flags&1) goto L_163A4;
#else
    if(!(m->flags&1)) goto L_163A4;
#endif
L_163A2: ; /* IDA 0x163A2; mov     al, cl */
    vg_al(m,(uint8_t)((uint8_t)m->cx));
L_163A4: ; /* IDA 0x163A4; mov     cx, 1A9h */
    m->cx=(uint16_t)(425);
L_163A7: ; /* IDA 0x163A7; cmp     al, [bx+0Eh] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+14)),8);
L_163AA: ; /* IDA 0x163AA; jz      short loc_163B3 */
    if(m->flags&0x40) goto L_163B3;
L_163AC: ; /* IDA 0x163AC; inc     cx */
    m->cx=(uint16_t)(ec_inc(m,m->cx,16));
L_163AD: ; /* IDA 0x163AD; cmp     al, [bx+0Fh] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+15)),8);
L_163B0: ; /* IDA 0x163B0; jz      short loc_163B3 */
    if(m->flags&0x40) goto L_163B3;
L_163B2: ; /* IDA 0x163B2; inc     cx */
    m->cx=(uint16_t)(ec_inc(m,m->cx,16));
L_163B3: ; /* IDA 0x163B3; mov     ah, [bx+1Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
L_163B6: ; /* IDA 0x163B6; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_163B9: ; /* IDA 0x163B9; call    sub_18810 */
    sty_call(m,34832,25532,h);
L_163BC: ; /* IDA 0x163BC; pop     ds */
#if KI_STRATEGY_MUTATION == 15
    m->es=ec_pop(m);
#else
    m->ds=(uint16_t)(ec_pop(m));
#endif
L_163BD: ; /* IDA 0x163BD; jmp     short sub_16366 */
    m->ip=0x6366;if(h&&h->enter)h->enter(m,m->ip,h->user);goto L_16366;
}
/* Original sub_163BF; IDA 0x163BF. */
static void sty_body_163BF(KiMachine16 *m,const KiEconomyHooks *h) {
L_163BF: ; /* IDA 0x163BF; call    sub_12078 */
    sty_call(m,8312,25538,h);
L_163C2: ; /* IDA 0x163C2; mov     cx, 19h */
    m->cx=(uint16_t)(25);
L_163C5: ; /* IDA 0x163C5; call    sub_18853 */
    sty_call(m,34899,25544,h);
L_163C8: ; /* IDA 0x163C8; call    sub_178A7 */
    sty_call(m,30887,25547,h);
L_163CB: ; /* IDA 0x163CB; pushf */
    ec_push(m,m->flags);
L_163CC: ; /* IDA 0x163CC; call    sub_120D6 */
    sty_call(m,8406,25551,h);
L_163CF: ; /* IDA 0x163CF; popf */
    m->flags=ec_pop(m)|2;
L_163D0: ; /* IDA 0x163D0; jb      short loc_163FE */
    if(m->flags&1) goto L_163FE;
L_163D2: ; /* IDA 0x163D2; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_163D7: ; /* IDA 0x163D7; mov     bh, [bx+3] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_163DA: ; /* IDA 0x163DA; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_163DC: ; /* IDA 0x163DC; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_163DE: ; /* IDA 0x163DE; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_163E0: ; /* IDA 0x163E0; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_163E2: ; /* IDA 0x163E2; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_163E4: ; /* IDA 0x163E4; mov     dx, [bx+848h] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+2120)));
L_163E8: ; /* IDA 0x163E8; mov     bx, [bx+84Ah] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+2122)));
L_163EC: ; /* IDA 0x163EC; mov     ax, 14h */
    m->ax=(uint16_t)(20);
L_163EF: ; /* IDA 0x163EF; mov     cx, 0Ch */
    m->cx=(uint16_t)(12);
L_163F2: ; /* IDA 0x163F2; call    sub_12151 */
    sty_call(m,8529,25589,h);
L_163F5: ; /* IDA 0x163F5; call    sub_11F7F */
    sty_call(m,8063,25592,h);
L_163F8: ; /* IDA 0x163F8; call    sub_11D46 */
    sty_call(m,7494,25595,h);
L_163FB: ; /* IDA 0x163FB; call    sub_17E1F */
    sty_call(m,32287,25598,h);
L_163FE: ; /* IDA 0x163FE; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16401: ; /* IDA 0x16401; call    sub_18853 */
    sty_call(m,34899,25604,h);
L_16404: ; /* IDA 0x16404; retn */
    return;
}
/* Original sub_16405; IDA 0x16405. */
static void sty_body_16405(KiMachine16 *m,const KiEconomyHooks *h) {
L_16405: ; /* IDA 0x16405; push    ds */
    ec_push(m,m->ds);
L_16406: ; /* IDA 0x16406; call    sub_12078 */
    sty_call(m,8312,25609,h);
L_16409: ; /* IDA 0x16409; call    sub_101B4 */
    sty_call(m,436,25612,h);
L_1640C: ; /* IDA 0x1640C; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16411: ; /* IDA 0x16411; mov     cx, 5 */
    m->cx=(uint16_t)(5);
L_16414: ; /* IDA 0x16414; call    sub_18853 */
    sty_call(m,34899,25623,h);
L_16417: ; /* IDA 0x16417; call    sub_101DB */
    sty_call(m,475,25626,h);
L_1641A: ; /* IDA 0x1641A; call    sub_17906 */
    sty_call(m,30982,25629,h);
L_1641D: ; /* IDA 0x1641D; jb      short loc_16460 */
    if(m->flags&1) goto L_16460;
L_1641F: ; /* IDA 0x1641F; push    ax */
    ec_push(m,m->ax);
L_16420: ; /* IDA 0x16420; push    bx */
    ec_push(m,m->bx);
L_16421: ; /* IDA 0x16421; push    dx */
    ec_push(m,m->dx);
L_16422: ; /* IDA 0x16422; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_16425: ; /* IDA 0x16425; call    sub_20000 */
    mx_far(m,0x6425,0x642a,h);
L_1642A: ; /* IDA 0x1642A; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1642B: ; /* IDA 0x1642B; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1642C: ; /* IDA 0x1642C; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1642D: ; /* IDA 0x1642D; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16430: ; /* IDA 0x16430; call    sub_18853 */
    sty_call(m,34899,25651,h);
L_16433: ; /* IDA 0x16433; call    sub_11D46 */
    sty_call(m,7494,25654,h);
L_16436: ; /* IDA 0x16436; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_16438: ; /* IDA 0x16438; call    sub_187FF */
    sty_call(m,34815,25659,h);
L_1643B: ; /* IDA 0x1643B; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_1643D: ; /* IDA 0x1643D; add     si, 4240h */
    m->si=(uint16_t)(ec_math(m,m->si,16960,16,0,0));
L_16441: ; /* IDA 0x16441; mov     bx, cx */
    m->bx=(uint16_t)(m->cx);
L_16443: ; /* IDA 0x16443; call    sub_16475 */
    sty_call(m,25717,25670,h);
L_16446: ; /* IDA 0x16446; push    bx */
    ec_push(m,m->bx);
L_16447: ; /* IDA 0x16447; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16449: ; /* IDA 0x16449; mov     cx, 56h ; 'V' */
    m->cx=(uint16_t)(86);
L_1644C: ; /* IDA 0x1644C; call    sub_13830 */
    sty_call(m,14384,25679,h);
L_1644F: ; /* IDA 0x1644F; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16450: ; /* IDA 0x16450; jb      short loc_16460 */
    if(m->flags&1) goto L_16460;
L_16452: ; /* IDA 0x16452; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_16457: ; /* IDA 0x16457; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_16459: ; /* IDA 0x16459; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1645B: ; /* IDA 0x1645B; mov     dl, bh */
    vg_dl(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_1645D: ; /* IDA 0x1645D; call    sub_13526 */
    sty_call(m,13606,25696,h);
L_16460: ; /* IDA 0x16460; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16463: ; /* IDA 0x16463; call    sub_18853 */
    sty_call(m,34899,25702,h);
L_16466: ; /* IDA 0x16466; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_16469: ; /* IDA 0x16469; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_1646B: ; /* IDA 0x1646B; call    sub_20000 */
    mx_far(m,0x646b,0x6470,h);
L_16470: ; /* IDA 0x16470; call    sub_120D6 */
    sty_call(m,8406,25715,h);
L_16473: ; /* IDA 0x16473; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16474: ; /* IDA 0x16474; retn */
    return;
}
/* Original sub_16475; IDA 0x16475. */
static void sty_body_16475(KiMachine16 *m,const KiEconomyHooks *h) {
L_16475: ; /* IDA 0x16475; push    bx */
    ec_push(m,m->bx);
L_16476: ; /* IDA 0x16476; push    cx */
    ec_push(m,m->cx);
L_16477: ; /* IDA 0x16477; push    si */
    ec_push(m,m->si);
L_16478: ; /* IDA 0x16478; push    di */
    ec_push(m,m->di);
L_16479: ; /* IDA 0x16479; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_1647E: ; /* IDA 0x1647E; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_16480: ; /* IDA 0x16480; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_16482: ; /* IDA 0x16482; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_16484: ; /* IDA 0x16484; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_16486: ; /* IDA 0x16486; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_16488: ; /* IDA 0x16488; mov     dl, dh */
    vg_dl(m,(uint8_t)((uint8_t)(m->dx>>8)));
L_1648A: ; /* IDA 0x1648A; mov     dh, 0FFh */
    vg_dh(m,(uint8_t)(255));
L_1648C: ; /* IDA 0x1648C; call    sub_1304E */
    sty_call(m,12366,25743,h);
L_1648F: ; /* IDA 0x1648F; jb      short loc_16495 */
    if(m->flags&1) goto L_16495;
L_16491: ; /* IDA 0x16491; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_16493: ; /* IDA 0x16493; jmp     short loc_164EC */
    goto L_164EC;
L_16495: ; /* IDA 0x16495; call    sub_130CB */
    sty_call(m,12491,25752,h);
L_16498: ; /* IDA 0x16498; sub     al, 80h */
#if KI_STRATEGY_MUTATION == 11
    vg_al(m,ec_math(m,(uint8_t)m->ax,0x7f,8,0,1));
#else
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,128,8,0,1)));
#endif
L_1649A: ; /* IDA 0x1649A; jnb     short loc_164A0 */
    if(!(m->flags&1)) goto L_164A0;
L_1649C: ; /* IDA 0x1649C; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1649E: ; /* IDA 0x1649E; jmp     short loc_164EC */
    goto L_164EC;
L_164A0: ; /* IDA 0x164A0; mov     ah, [si+28h] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+40))));
L_164A3: ; /* IDA 0x164A3; shl     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,4)));
L_164A5: ; /* IDA 0x164A5; add     ah, 14h */
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),20,8,0,0)));
L_164A8: ; /* IDA 0x164A8; cmp     al, ah */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
L_164AA: ; /* IDA 0x164AA; jb      short loc_164B0 */
    if(m->flags&1) goto L_164B0;
L_164AC: ; /* IDA 0x164AC; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_164AE: ; /* IDA 0x164AE; jmp     short loc_164EC */
    goto L_164EC;
L_164B0: ; /* IDA 0x164B0; xor     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_164B2: ; /* IDA 0x164B2; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_164B4: ; /* IDA 0x164B4; add     ah, 5 */
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),5,8,0,0)));
L_164B7: ; /* IDA 0x164B7; cmp     al, ah */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
L_164B9: ; /* IDA 0x164B9; jnb     short loc_164BE */
    if(!(m->flags&1)) goto L_164BE;
L_164BB: ; /* IDA 0x164BB; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_164BE: ; /* IDA 0x164BE; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_164C0: ; /* IDA 0x164C0; call    sub_16A28 */
    sty_call(m,27176,25795,h);
L_164C3: ; /* IDA 0x164C3; cmp     ax, cx */
    ec_cmp(m,m->ax,m->cx,16);
L_164C5: ; /* IDA 0x164C5; jbe     short loc_164CA */
    if(m->flags&(1|0x40)) goto L_164CA;
L_164C7: ; /* IDA 0x164C7; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_164CA: ; /* IDA 0x164CA; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_164CC: ; /* IDA 0x164CC; mov     al, [di+19h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+25))));
L_164CF: ; /* IDA 0x164CF; cmp     al, 0FFh */
    ec_cmp(m,(uint8_t)m->ax,255,8);
L_164D1: ; /* IDA 0x164D1; jz      short loc_164DD */
    if(m->flags&0x40) goto L_164DD;
L_164D3: ; /* IDA 0x164D3; cmp     al, cs:byte_10CFF */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
L_164D8: ; /* IDA 0x164D8; jz      short loc_164DD */
    if(m->flags&0x40) goto L_164DD;
L_164DA: ; /* IDA 0x164DA; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_164DD: ; /* IDA 0x164DD; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_164DF: ; /* IDA 0x164DF; cmp     word ptr [di+21h], 0 */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->di+33)),0,16);
L_164E3: ; /* IDA 0x164E3; jge     short loc_164E8 */
    if(!ec_less(m)) goto L_164E8;
L_164E5: ; /* IDA 0x164E5; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_164E8: ; /* IDA 0x164E8; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_164EA: ; /* IDA 0x164EA; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_164EC: ; /* IDA 0x164EC; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_164ED: ; /* IDA 0x164ED; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_164EE: ; /* IDA 0x164EE; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_164EF: ; /* IDA 0x164EF; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_164F0: ; /* IDA 0x164F0; retn */
    return;
}
/* Original sub_164F1; IDA 0x164F1. */
static void sty_body_164F1(KiMachine16 *m,const KiEconomyHooks *h) {
L_164F1: ; /* IDA 0x164F1; push    ds */
    ec_push(m,m->ds);
L_164F2: ; /* IDA 0x164F2; call    sub_12078 */
    sty_call(m,8312,25845,h);
L_164F5: ; /* IDA 0x164F5; call    sub_101B4 */
    sty_call(m,436,25848,h);
L_164F8: ; /* IDA 0x164F8; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_164FD: ; /* IDA 0x164FD; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_16500: ; /* IDA 0x16500; call    sub_18853 */
    sty_call(m,34899,25859,h);
L_16503: ; /* IDA 0x16503; call    sub_101DB */
    sty_call(m,475,25862,h);
L_16506: ; /* IDA 0x16506; call    sub_17906 */
    sty_call(m,30982,25865,h);
L_16509: ; /* IDA 0x16509; jb      short loc_16562 */
    if(m->flags&1) goto L_16562;
L_1650B: ; /* IDA 0x1650B; call    sub_165EF */
    sty_call(m,26095,25870,h);
L_1650E: ; /* IDA 0x1650E; jb      short loc_16506 */
    if(m->flags&1) goto L_16506;
L_16510: ; /* IDA 0x16510; call    sub_16605 */
    sty_call(m,26117,25875,h);
L_16513: ; /* IDA 0x16513; jb      short loc_16506 */
    if(m->flags&1) goto L_16506;
L_16515: ; /* IDA 0x16515; push    ax */
    ec_push(m,m->ax);
L_16516: ; /* IDA 0x16516; push    bx */
    ec_push(m,m->bx);
L_16517: ; /* IDA 0x16517; push    dx */
    ec_push(m,m->dx);
L_16518: ; /* IDA 0x16518; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_1651B: ; /* IDA 0x1651B; call    sub_20000 */
    mx_far(m,0x651b,0x6520,h);
L_16520: ; /* IDA 0x16520; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_16521: ; /* IDA 0x16521; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16522: ; /* IDA 0x16522; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16523: ; /* IDA 0x16523; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16526: ; /* IDA 0x16526; call    sub_18853 */
    sty_call(m,34899,25897,h);
L_16529: ; /* IDA 0x16529; call    sub_11D46 */
    sty_call(m,7494,25900,h);
L_1652C: ; /* IDA 0x1652C; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_1652E: ; /* IDA 0x1652E; call    sub_187FF */
    sty_call(m,34815,25905,h);
L_16531: ; /* IDA 0x16531; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16533: ; /* IDA 0x16533; add     si, 4240h */
    m->si=(uint16_t)(ec_math(m,m->si,16960,16,0,0));
L_16537: ; /* IDA 0x16537; mov     bx, cx */
    m->bx=(uint16_t)(m->cx);
L_16539: ; /* IDA 0x16539; call    sub_16577 */
    sty_call(m,25975,25916,h);
L_1653C: ; /* IDA 0x1653C; push    bx */
    ec_push(m,m->bx);
L_1653D: ; /* IDA 0x1653D; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_1653F: ; /* IDA 0x1653F; mov     cx, 96h */
    m->cx=(uint16_t)(150);
L_16542: ; /* IDA 0x16542; call    sub_13830 */
    sty_call(m,14384,25925,h);
L_16545: ; /* IDA 0x16545; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16546: ; /* IDA 0x16546; jb      short loc_16562 */
    if(m->flags&1) goto L_16562;
L_16548: ; /* IDA 0x16548; mov     al, 6 */
    vg_al(m,(uint8_t)(6));
L_1654A: ; /* IDA 0x1654A; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_1654C: ; /* IDA 0x1654C; mov     bl, 14h */
    gw_bl(m,(uint8_t)(20));
L_1654E: ; /* IDA 0x1654E; call    sub_1300E */
    sty_call(m,12302,25937,h);
L_16551: ; /* IDA 0x16551; mov     al, [si+2Ah] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+42))));
L_16554: ; /* IDA 0x16554; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_16556: ; /* IDA 0x16556; push    ax */
    ec_push(m,m->ax);
L_16557: ; /* IDA 0x16557; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16559: ; /* IDA 0x16559; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_1655B: ; /* IDA 0x1655B; mov     cx, 3Bh ; ';' */
    m->cx=(uint16_t)(59);
L_1655E: ; /* IDA 0x1655E; call    sub_18810 */
    sty_call(m,34832,25953,h);
L_16561: ; /* IDA 0x16561; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16562: ; /* IDA 0x16562; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16565: ; /* IDA 0x16565; call    sub_18853 */
    sty_call(m,34899,25960,h);
L_16568: ; /* IDA 0x16568; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_1656B: ; /* IDA 0x1656B; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_1656D: ; /* IDA 0x1656D; call    sub_20000 */
    mx_far(m,0x656d,0x6572,h);
L_16572: ; /* IDA 0x16572; call    sub_120D6 */
    sty_call(m,8406,25973,h);
L_16575: ; /* IDA 0x16575; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16576: ; /* IDA 0x16576; retn */
    return;
}
/* Original sub_16577; IDA 0x16577. */
static void sty_body_16577(KiMachine16 *m,const KiEconomyHooks *h) {
L_16577: ; /* IDA 0x16577; push    bx */
    ec_push(m,m->bx);
L_16578: ; /* IDA 0x16578; push    cx */
    ec_push(m,m->cx);
L_16579: ; /* IDA 0x16579; push    si */
    ec_push(m,m->si);
L_1657A: ; /* IDA 0x1657A; push    di */
    ec_push(m,m->di);
L_1657B: ; /* IDA 0x1657B; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_16580: ; /* IDA 0x16580; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_16582: ; /* IDA 0x16582; call    sub_130CB */
    sty_call(m,12491,25989,h);
L_16585: ; /* IDA 0x16585; cmp     al, 80h */
    ec_cmp(m,(uint8_t)m->ax,128,8);
L_16587: ; /* IDA 0x16587; jb      short loc_1658D */
    if(m->flags&1) goto L_1658D;
L_16589: ; /* IDA 0x16589; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1658B: ; /* IDA 0x1658B; jmp     short loc_165EA */
    goto L_165EA;
L_1658D: ; /* IDA 0x1658D; mov     ah, [si+28h] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+40))));
L_16590: ; /* IDA 0x16590; shr     ah, 1 */
    vg_ah(m,(uint8_t)(dl_shift(m,(uint8_t)(m->ax>>8),1,8,5)));
L_16592: ; /* IDA 0x16592; cmp     al, ah */
#if KI_STRATEGY_MUTATION == 12
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+0x28)),8);
#else
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
#endif
L_16594: ; /* IDA 0x16594; jnb     short loc_1659A */
    if(!(m->flags&1)) goto L_1659A;
L_16596: ; /* IDA 0x16596; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_16598: ; /* IDA 0x16598; jmp     short loc_165EA */
    goto L_165EA;
L_1659A: ; /* IDA 0x1659A; xor     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_1659C: ; /* IDA 0x1659C; call    sub_16A28 */
    sty_call(m,27176,26015,h);
L_1659F: ; /* IDA 0x1659F; cmp     ax, cx */
    ec_cmp(m,m->ax,m->cx,16);
L_165A1: ; /* IDA 0x165A1; jnb     short loc_165A6 */
    if(!(m->flags&1)) goto L_165A6;
L_165A3: ; /* IDA 0x165A3; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_165A6: ; /* IDA 0x165A6; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_165A8: ; /* IDA 0x165A8; mov     bx, 0 */
    m->bx=(uint16_t)(0);
L_165AB: ; /* IDA 0x165AB; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_165AF: ; /* IDA 0x165AF; mov     cx, 16h */
    m->cx=(uint16_t)(22);
L_165B2: ; /* IDA 0x165B2; cmp     byte ptr [bx], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx)),128,8);
L_165B5: ; /* IDA 0x165B5; jb      short loc_165C3 */
    if(m->flags&1) goto L_165C3;
L_165B7: ; /* IDA 0x165B7; cmp     bx, di */
    ec_cmp(m,m->bx,m->di,16);
L_165B9: ; /* IDA 0x165B9; jz      short loc_165C3 */
    if(m->flags&0x40) goto L_165C3;
L_165BB: ; /* IDA 0x165BB; cmp     al, [bx+19h] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+25)),8);
L_165BE: ; /* IDA 0x165BE; jnz     short loc_165C3 */
    if(!(m->flags&0x40)) goto L_165C3;
L_165C0: ; /* IDA 0x165C0; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_165C3: ; /* IDA 0x165C3; add     bx, 40h ; '@' */
    m->bx=(uint16_t)(ec_math(m,m->bx,64,16,0,0));
L_165C6: ; /* IDA 0x165C6; loop    loc_165B2 */
    if(--m->cx) goto L_165B2;
L_165C8: ; /* IDA 0x165C8; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_165CA: ; /* IDA 0x165CA; mov     al, [di+19h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+25))));
L_165CD: ; /* IDA 0x165CD; cmp     al, 0FFh */
    ec_cmp(m,(uint8_t)m->ax,255,8);
L_165CF: ; /* IDA 0x165CF; jz      short loc_165DB */
    if(m->flags&0x40) goto L_165DB;
L_165D1: ; /* IDA 0x165D1; cmp     al, cs:byte_10CFF */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
L_165D6: ; /* IDA 0x165D6; jz      short loc_165DB */
    if(m->flags&0x40) goto L_165DB;
L_165D8: ; /* IDA 0x165D8; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_165DB: ; /* IDA 0x165DB; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_165DD: ; /* IDA 0x165DD; cmp     word ptr [si+21h], 0 */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+33)),0,16);
L_165E1: ; /* IDA 0x165E1; jge     short loc_165E6 */
    if(!ec_less(m)) goto L_165E6;
L_165E3: ; /* IDA 0x165E3; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_165E6: ; /* IDA 0x165E6; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_165E8: ; /* IDA 0x165E8; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_165EA: ; /* IDA 0x165EA; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_165EB: ; /* IDA 0x165EB; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_165EC: ; /* IDA 0x165EC; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_165ED: ; /* IDA 0x165ED; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_165EE: ; /* IDA 0x165EE; retn */
    return;
}
/* Original sub_165EF; IDA 0x165EF. */
static void sty_body_165EF(KiMachine16 *m,const KiEconomyHooks *h) {
L_165EF: ; /* IDA 0x165EF; cmp     byte ptr [bx+2Ah], 0FFh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+42)),255,8);
L_165F3: ; /* IDA 0x165F3; jz      short loc_165F7 */
    if(m->flags&0x40) goto L_165F7;
L_165F5: ; /* IDA 0x165F5; clc */
    m->flags&=(uint16_t)~1u;
L_165F6: ; /* IDA 0x165F6; retn */
    return;
L_165F7: ; /* IDA 0x165F7; push    bx */
    ec_push(m,m->bx);
L_165F8: ; /* IDA 0x165F8; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_165FA: ; /* IDA 0x165FA; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_165FC: ; /* IDA 0x165FC; mov     cx, 37h ; '7' */
    m->cx=(uint16_t)(55);
L_165FF: ; /* IDA 0x165FF; call    sub_18810 */
    sty_call(m,34832,26114,h);
L_16602: ; /* IDA 0x16602; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16603: ; /* IDA 0x16603; stc */
    m->flags|=1;
L_16604: ; /* IDA 0x16604; retn */
    return;
}
/* Original sub_16605; IDA 0x16605. */
static void sty_body_16605(KiMachine16 *m,const KiEconomyHooks *h) {
L_16605: ; /* IDA 0x16605; push    si */
    ec_push(m,m->si);
L_16606: ; /* IDA 0x16606; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16608: ; /* IDA 0x16608; mov     al, 6 */
    vg_al(m,(uint8_t)(6));
L_1660A: ; /* IDA 0x1660A; mov     dx, 0FFFFh */
    m->dx=(uint16_t)(65535);
L_1660D: ; /* IDA 0x1660D; call    sub_1304E */
    sty_call(m,12366,26128,h);
L_16610: ; /* IDA 0x16610; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16611: ; /* IDA 0x16611; jnb     short loc_16615 */
    if(!(m->flags&1)) goto L_16615;
L_16613: ; /* IDA 0x16613; clc */
    m->flags&=(uint16_t)~1u;
L_16614: ; /* IDA 0x16614; retn */
    return;
L_16615: ; /* IDA 0x16615; push    bx */
    ec_push(m,m->bx);
L_16616: ; /* IDA 0x16616; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16618: ; /* IDA 0x16618; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_1661A: ; /* IDA 0x1661A; mov     cx, 49h ; 'I' */
    m->cx=(uint16_t)(73);
L_1661D: ; /* IDA 0x1661D; call    sub_18810 */
    sty_call(m,34832,26144,h);
L_16620: ; /* IDA 0x16620; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16621: ; /* IDA 0x16621; stc */
    m->flags|=1;
L_16622: ; /* IDA 0x16622; retn */
    return;
}
/* Original sub_16623; IDA 0x16623. */
static void sty_body_16623(KiMachine16 *m,const KiEconomyHooks *h) {
L_16623: ; /* IDA 0x16623; push    ds */
    ec_push(m,m->ds);
L_16624: ; /* IDA 0x16624; push    bp */
    ec_push(m,m->bp);
L_16625: ; /* IDA 0x16625; sub     sp, 4 */
    m->sp=(uint16_t)(ec_math(m,m->sp,4,16,0,1));
L_16628: ; /* IDA 0x16628; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_1662A: ; /* IDA 0x1662A; call    sub_12078 */
    sty_call(m,8312,26157,h);
L_1662D: ; /* IDA 0x1662D; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16632: ; /* IDA 0x16632; call    sub_101B4 */
    sty_call(m,436,26165,h);
L_16635: ; /* IDA 0x16635; mov     cx, 8 */
    m->cx=(uint16_t)(8);
L_16638: ; /* IDA 0x16638; call    sub_18853 */
    sty_call(m,34899,26171,h);
L_1663B: ; /* IDA 0x1663B; call    sub_101DB */
    sty_call(m,475,26174,h);
L_1663E: ; /* IDA 0x1663E; call    sub_17906 */
    sty_call(m,30982,26177,h);
L_16641: ; /* IDA 0x16641; jb      short loc_166C0 */
    if(m->flags&1) goto L_166C0;
L_16643: ; /* IDA 0x16643; call    sub_165EF */
    sty_call(m,26095,26182,h);
L_16646: ; /* IDA 0x16646; jb      short loc_1663E */
    if(m->flags&1) goto L_1663E;
L_16648: ; /* IDA 0x16648; call    sub_1676F */
    sty_call(m,26479,26187,h);
L_1664B: ; /* IDA 0x1664B; jb      short loc_1663E */
    if(m->flags&1) goto L_1663E;
L_1664D: ; /* IDA 0x1664D; mov     [bp+0], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp),m->bx);
L_16650: ; /* IDA 0x16650; call    sub_101B4 */
    sty_call(m,436,26195,h);
L_16653: ; /* IDA 0x16653; mov     cx, 7 */
    m->cx=(uint16_t)(7);
L_16656: ; /* IDA 0x16656; call    sub_18853 */
    sty_call(m,34899,26201,h);
L_16659: ; /* IDA 0x16659; call    sub_101DB */
    sty_call(m,475,26204,h);
L_1665C: ; /* IDA 0x1665C; call    sub_17906 */
    sty_call(m,30982,26207,h);
L_1665F: ; /* IDA 0x1665F; jb      short loc_166C0 */
    if(m->flags&1) goto L_166C0;
L_16661: ; /* IDA 0x16661; mov     [bp+2], bx */
    dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->bx);
L_16664: ; /* IDA 0x16664; push    ax */
    ec_push(m,m->ax);
L_16665: ; /* IDA 0x16665; push    bx */
    ec_push(m,m->bx);
L_16666: ; /* IDA 0x16666; push    dx */
    ec_push(m,m->dx);
L_16667: ; /* IDA 0x16667; mov     ax, 2 */
    m->ax=(uint16_t)(2);
L_1666A: ; /* IDA 0x1666A; call    sub_20000 */
    mx_far(m,0x666a,0x666f,h);
L_1666F: ; /* IDA 0x1666F; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_16670: ; /* IDA 0x16670; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_16671: ; /* IDA 0x16671; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_16672: ; /* IDA 0x16672; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16675: ; /* IDA 0x16675; call    sub_18853 */
    sty_call(m,34899,26232,h);
L_16678: ; /* IDA 0x16678; call    sub_11D46 */
    sty_call(m,7494,26235,h);
L_1667B: ; /* IDA 0x1667B; mov     cx, bx */
    m->cx=(uint16_t)(m->bx);
L_1667D: ; /* IDA 0x1667D; call    sub_187FF */
    sty_call(m,34815,26240,h);
L_16680: ; /* IDA 0x16680; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16682: ; /* IDA 0x16682; add     si, 4240h */
    m->si=(uint16_t)(ec_math(m,m->si,16960,16,0,0));
L_16686: ; /* IDA 0x16686; mov     bx, cx */
    m->bx=(uint16_t)(m->cx);
L_16688: ; /* IDA 0x16688; call    sub_166D9 */
    sty_call(m,26329,26251,h);
L_1668B: ; /* IDA 0x1668B; mov     di, bp */
    m->di=(uint16_t)(m->bp);
L_1668D: ; /* IDA 0x1668D; mov     cx, 0D6h */
    m->cx=(uint16_t)(214);
L_16690: ; /* IDA 0x16690; call    sub_13830 */
    sty_call(m,14384,26259,h);
L_16693: ; /* IDA 0x16693; jb      short loc_166C0 */
    if(m->flags&1) goto L_166C0;
L_16695: ; /* IDA 0x16695; mov     ax, [bp+0] */
    m->ax=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_16698: ; /* IDA 0x16698; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1669A: ; /* IDA 0x1669A; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1669C: ; /* IDA 0x1669C; mov     al, 7 */
    vg_al(m,(uint8_t)(7));
L_1669E: ; /* IDA 0x1669E; mov     dx, [bp+2] */
    m->dx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_166A1: ; /* IDA 0x166A1; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_166A3: ; /* IDA 0x166A3; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_166A5: ; /* IDA 0x166A5; mov     dl, dh */
    vg_dl(m,(uint8_t)((uint8_t)(m->dx>>8)));
L_166A7: ; /* IDA 0x166A7; mov     bl, 14h */
    gw_bl(m,(uint8_t)(20));
L_166A9: ; /* IDA 0x166A9; call    sub_1301C */
    sty_call(m,12316,26284,h);
L_166AC: ; /* IDA 0x166AC; mov     bx, [bp+0] */
    m->bx=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_166AF: ; /* IDA 0x166AF; mov     al, [bx+2Ah] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+42))));
L_166B2: ; /* IDA 0x166B2; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_166B4: ; /* IDA 0x166B4; push    ax */
    ec_push(m,m->ax);
L_166B5: ; /* IDA 0x166B5; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_166B7: ; /* IDA 0x166B7; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_166B9: ; /* IDA 0x166B9; mov     cx, 3Ch ; '<' */
    m->cx=(uint16_t)(60);
L_166BC: ; /* IDA 0x166BC; call    sub_18810 */
    sty_call(m,34832,26303,h);
L_166BF: ; /* IDA 0x166BF; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_166C0: ; /* IDA 0x166C0; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_166C3: ; /* IDA 0x166C3; call    sub_18853 */
    sty_call(m,34899,26310,h);
L_166C6: ; /* IDA 0x166C6; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_166C9: ; /* IDA 0x166C9; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_166CB: ; /* IDA 0x166CB; call    sub_20000 */
    mx_far(m,0x66cb,0x66d0,h);
L_166D0: ; /* IDA 0x166D0; call    sub_120D6 */
    sty_call(m,8406,26323,h);
L_166D3: ; /* IDA 0x166D3; add     sp, 4 */
    m->sp=(uint16_t)(ec_math(m,m->sp,4,16,0,0));
L_166D6: ; /* IDA 0x166D6; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_166D7: ; /* IDA 0x166D7; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_166D8: ; /* IDA 0x166D8; retn */
    return;
}
/* Original sub_166D9; IDA 0x166D9. */
static void sty_body_166D9(KiMachine16 *m,const KiEconomyHooks *h) {
L_166D9: ; /* IDA 0x166D9; push    bx */
    ec_push(m,m->bx);
L_166DA: ; /* IDA 0x166DA; push    cx */
    ec_push(m,m->cx);
L_166DB: ; /* IDA 0x166DB; push    si */
    ec_push(m,m->si);
L_166DC: ; /* IDA 0x166DC; push    di */
    ec_push(m,m->di);
L_166DD: ; /* IDA 0x166DD; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_166E2: ; /* IDA 0x166E2; mov     di, [bp+0] */
    m->di=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_166E5: ; /* IDA 0x166E5; cmp     di, [bp+2] */
#if KI_STRATEGY_MUTATION == 13
    ec_cmp(m,m->di,m->di,16);
#else
    ec_cmp(m,m->di,dl_read16(m,m->ss,(uint16_t)(m->bp+2)),16);
#endif
L_166E8: ; /* IDA 0x166E8; jnz     short loc_166EE */
    if(!(m->flags&0x40)) goto L_166EE;
L_166EA: ; /* IDA 0x166EA; mov     al, 4 */
    vg_al(m,(uint8_t)(4));
L_166EC: ; /* IDA 0x166EC; jmp     short loc_1676A */
    goto L_1676A;
L_166EE: ; /* IDA 0x166EE; call    sub_130CB */
    sty_call(m,12491,26353,h);
L_166F1: ; /* IDA 0x166F1; mov     bh, [si+28h] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+40))));
L_166F4: ; /* IDA 0x166F4; shl     bh, 1 */
    gw_bh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->bx>>8),1,8,4)));
L_166F6: ; /* IDA 0x166F6; shl     bh, 1 */
    gw_bh(m,(uint8_t)(dl_shift(m,(uint8_t)(m->bx>>8),1,8,4)));
L_166F8: ; /* IDA 0x166F8; add     bh, 1Eh */
    gw_bh(m,(uint8_t)(ec_math(m,(uint8_t)(m->bx>>8),30,8,0,0)));
L_166FB: ; /* IDA 0x166FB; or      bh, 80h */
    {uint16_t result=(((uint8_t)(m->bx>>8))|(128))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_166FE: ; /* IDA 0x166FE; cmp     al, bh */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->bx>>8),8);
L_16700: ; /* IDA 0x16700; jnb     short loc_16706 */
    if(!(m->flags&1)) goto L_16706;
L_16702: ; /* IDA 0x16702; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_16704: ; /* IDA 0x16704; jmp     short loc_1676A */
    goto L_1676A;
L_16706: ; /* IDA 0x16706; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_16708: ; /* IDA 0x16708; mov     di, [bp+2] */
    m->di=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_1670B: ; /* IDA 0x1670B; call    sub_130CB */
    sty_call(m,12491,26382,h);
L_1670E: ; /* IDA 0x1670E; cmp     al, 80h */
    ec_cmp(m,(uint8_t)m->ax,128,8);
L_16710: ; /* IDA 0x16710; jb      short loc_16716 */
    if(m->flags&1) goto L_16716;
L_16712: ; /* IDA 0x16712; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_16714: ; /* IDA 0x16714; jmp     short loc_1676A */
    goto L_1676A;
L_16716: ; /* IDA 0x16716; mov     al, [di+19h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+25))));
L_16719: ; /* IDA 0x16719; cmp     al, cs:byte_10CFF */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
L_1671E: ; /* IDA 0x1671E; jnz     short loc_1672D */
    if(!(m->flags&0x40)) goto L_1672D;
L_16720: ; /* IDA 0x16720; call    sub_16A28 */
    sty_call(m,27176,26403,h);
L_16723: ; /* IDA 0x16723; shr     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
L_16725: ; /* IDA 0x16725; cmp     ax, cx */
    ec_cmp(m,m->ax,m->cx,16);
L_16727: ; /* IDA 0x16727; jnb     short loc_1672D */
    if(!(m->flags&1)) goto L_1672D;
L_16729: ; /* IDA 0x16729; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_1672B: ; /* IDA 0x1672B; jmp     short loc_1676A */
    goto L_1676A;
L_1672D: ; /* IDA 0x1672D; xor     dl, dl */
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_1672F: ; /* IDA 0x1672F; add     bh, 1Eh */
    gw_bh(m,(uint8_t)(ec_math(m,(uint8_t)(m->bx>>8),30,8,0,0)));
L_16732: ; /* IDA 0x16732; cmp     bl, bh */
    ec_cmp(m,(uint8_t)m->bx,(uint8_t)(m->bx>>8),8);
L_16734: ; /* IDA 0x16734; jb      short loc_16739 */
    if(m->flags&1) goto L_16739;
L_16736: ; /* IDA 0x16736; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_16739: ; /* IDA 0x16739; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_1673B: ; /* IDA 0x1673B; mov     di, [bp+0] */
    m->di=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp)));
L_1673E: ; /* IDA 0x1673E; call    sub_16A28 */
    sty_call(m,27176,26433,h);
L_16741: ; /* IDA 0x16741; cmp     ax, cx */
    ec_cmp(m,m->ax,m->cx,16);
L_16743: ; /* IDA 0x16743; jnb     short loc_16748 */
    if(!(m->flags&1)) goto L_16748;
L_16745: ; /* IDA 0x16745; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_16748: ; /* IDA 0x16748; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_1674A: ; /* IDA 0x1674A; mov     di, [bp+2] */
    m->di=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+2)));
L_1674D: ; /* IDA 0x1674D; call    sub_16A28 */
    sty_call(m,27176,26448,h);
L_16750: ; /* IDA 0x16750; cmp     ax, cx */
    ec_cmp(m,m->ax,m->cx,16);
L_16752: ; /* IDA 0x16752; jnb     short loc_16757 */
    if(!(m->flags&1)) goto L_16757;
L_16754: ; /* IDA 0x16754; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_16757: ; /* IDA 0x16757; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_16759: ; /* IDA 0x16759; mov     al, [di+19h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+25))));
L_1675C: ; /* IDA 0x1675C; cmp     al, cs:byte_10CFF */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
L_16761: ; /* IDA 0x16761; jnz     short loc_16766 */
    if(!(m->flags&0x40)) goto L_16766;
L_16763: ; /* IDA 0x16763; or      dl, 10h */
    {uint16_t result=(((uint8_t)m->dx)|(16))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
L_16766: ; /* IDA 0x16766; shr     dl, 1 */
    vg_dl(m,(uint8_t)(dl_shift(m,(uint8_t)m->dx,1,8,5)));
L_16768: ; /* IDA 0x16768; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_1676A: ; /* IDA 0x1676A; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1676B: ; /* IDA 0x1676B; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1676C: ; /* IDA 0x1676C; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1676D: ; /* IDA 0x1676D; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1676E: ; /* IDA 0x1676E; retn */
    return;
}
/* Original sub_1676F; IDA 0x1676F. */
static void sty_body_1676F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1676F: ; /* IDA 0x1676F; push    si */
    ec_push(m,m->si);
L_16770: ; /* IDA 0x16770; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16772: ; /* IDA 0x16772; mov     al, 7 */
    vg_al(m,(uint8_t)(7));
L_16774: ; /* IDA 0x16774; mov     dx, 0FFFFh */
    m->dx=(uint16_t)(65535);
L_16777: ; /* IDA 0x16777; call    sub_1304E */
    sty_call(m,12366,26490,h);
L_1677A: ; /* IDA 0x1677A; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1677B: ; /* IDA 0x1677B; jnb     short loc_1677F */
    if(!(m->flags&1)) goto L_1677F;
L_1677D: ; /* IDA 0x1677D; clc */
    m->flags&=(uint16_t)~1u;
L_1677E: ; /* IDA 0x1677E; retn */
    return;
L_1677F: ; /* IDA 0x1677F; push    bx */
    ec_push(m,m->bx);
L_16780: ; /* IDA 0x16780; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16782: ; /* IDA 0x16782; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_16784: ; /* IDA 0x16784; mov     cx, 4Ah ; 'J' */
    m->cx=(uint16_t)(74);
L_16787: ; /* IDA 0x16787; call    sub_18810 */
    sty_call(m,34832,26506,h);
L_1678A: ; /* IDA 0x1678A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1678B: ; /* IDA 0x1678B; stc */
    m->flags|=1;
L_1678C: ; /* IDA 0x1678C; retn */
    return;
}
/* Original sub_1678D; IDA 0x1678D. */
static void sty_body_1678D(KiMachine16 *m,const KiEconomyHooks *h) {
L_1678D: ; /* IDA 0x1678D; call    sub_168C3 */
    sty_call(m,26819,26512,h);
L_16790: ; /* IDA 0x16790; call    sub_16846 */
    sty_call(m,26694,26515,h);
L_16793: ; /* IDA 0x16793; call    sub_12078 */
    sty_call(m,8312,26518,h);
L_16796: ; /* IDA 0x16796; mov     cx, 10h */
    m->cx=(uint16_t)(16);
L_16799: ; /* IDA 0x16799; call    sub_18853 */
    sty_call(m,34899,26524,h);
L_1679C: ; /* IDA 0x1679C; call    sub_121E7 */
    sty_call(m,8679,26527,h);
L_1679F: ; /* IDA 0x1679F; jb      short loc_167B8 */
    if(m->flags&1) goto L_167B8;
L_167A1: ; /* IDA 0x167A1; call    sub_1E453 */
    sty_call(m,58451,26532,h);
L_167A4: ; /* IDA 0x167A4; sub     al, 20h ; ' ' */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,32,8,0,1)));
L_167A6: ; /* IDA 0x167A6; jb      short loc_16796 */
    if(m->flags&1) goto L_16796;
L_167A8: ; /* IDA 0x167A8; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_167AA: ; /* IDA 0x167AA; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_167AC: ; /* IDA 0x167AC; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_167AE: ; /* IDA 0x167AE; call    cs:funcs_167AE[bx] */
    sty_call(m,dl_read16(m,m->cs,(uint16_t)(m->bx+26565)),26547,h);
L_167B3: ; /* IDA 0x167B3; call    sub_16846 */
    sty_call(m,26694,26550,h);
L_167B6: ; /* IDA 0x167B6; jmp     short loc_16796 */
    goto L_16796;
L_167B8: ; /* IDA 0x167B8; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_167BB: ; /* IDA 0x167BB; call    sub_18853 */
    sty_call(m,34899,26558,h);
L_167BE: ; /* IDA 0x167BE; call    sub_120D6 */
    sty_call(m,8406,26561,h);
L_167C1: ; /* IDA 0x167C1; call    sub_168FA */
    sty_call(m,26874,26564,h);
L_167C4: ; /* IDA 0x167C4; retn */
    return;
}
/* Original sub_16846; IDA 0x16846. */
static void sty_body_16846(KiMachine16 *m,const KiEconomyHooks *h) {
L_16846: ; /* IDA 0x16846; push    ds */
    ec_push(m,m->ds);
L_16847: ; /* IDA 0x16847; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_1684C: ; /* IDA 0x1684C; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_16851: ; /* IDA 0x16851; mov     al, [si+22h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+34))));
L_16854: ; /* IDA 0x16854; cbw */
#if KI_STRATEGY_MUTATION == 9
    m->ax=(uint8_t)m->ax;
#else
    m->ax=(uint16_t)(int16_t)(int8_t)m->ax;
#endif
L_16855: ; /* IDA 0x16855; mov     dx, ax */
    m->dx=(uint16_t)(m->ax);
L_16857: ; /* IDA 0x16857; mov     ax, [si+20h] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+32)));
L_1685A: ; /* IDA 0x1685A; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1685B: ; /* IDA 0x1685B; mov     bx, 907h */
    m->bx=(uint16_t)(2311);
L_1685E: ; /* IDA 0x1685E; cmp     dh, 50h ; 'P' */
    ec_cmp(m,(uint8_t)(m->dx>>8),80,8);
L_16861: ; /* IDA 0x16861; jb      short loc_16866 */
    if(m->flags&1) goto L_16866;
L_16863: ; /* IDA 0x16863; mov     bx, 0A07h */
    m->bx=(uint16_t)(2567);
L_16866: ; /* IDA 0x16866; mov     di, 230Dh */
    m->di=(uint16_t)(8973);
L_16869: ; /* IDA 0x16869; call    sub_1062F */
    sty_call(m,1583,26732,h);
L_1686C: ; /* IDA 0x1686C; mov     ax, cs:word_10D02 */
    m->ax=(uint16_t)(dl_read16(m,m->cs,3330));
L_16870: ; /* IDA 0x16870; mov     dl, cs:byte_10D04 */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,3332)));
L_16875: ; /* IDA 0x16875; xor     dh, dh */
    {uint16_t result=(((uint8_t)(m->dx>>8))^((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);vg_dh(m,(uint8_t)(result));}
L_16877: ; /* IDA 0x16877; mov     bx, 906h */
    m->bx=(uint16_t)(2310);
L_1687A: ; /* IDA 0x1687A; mov     di, 1E24h */
    m->di=(uint16_t)(7716);
L_1687D: ; /* IDA 0x1687D; call    sub_1062F */
    sty_call(m,1583,26752,h);
L_16880: ; /* IDA 0x16880; mov     ax, cs:word_10D05 */
    m->ax=(uint16_t)(dl_read16(m,m->cs,3333));
L_16884: ; /* IDA 0x16884; mov     dl, cs:byte_10D07 */
    vg_dl(m,(uint8_t)(dl_read8(m,m->cs,3335)));
L_16889: ; /* IDA 0x16889; xor     dh, dh */
    {uint16_t result=(((uint8_t)(m->dx>>8))^((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);vg_dh(m,(uint8_t)(result));}
L_1688B: ; /* IDA 0x1688B; mov     bx, 906h */
    m->bx=(uint16_t)(2310);
L_1688E: ; /* IDA 0x1688E; mov     di, 2324h */
    m->di=(uint16_t)(8996);
L_16891: ; /* IDA 0x16891; call    sub_1062F */
    sty_call(m,1583,26772,h);
L_16894: ; /* IDA 0x16894; mov     si, 0D08h */
    m->si=(uint16_t)(3336);
L_16897: ; /* IDA 0x16897; mov     di, 320Fh */
    m->di=(uint16_t)(12815);
L_1689A: ; /* IDA 0x1689A; mov     cx, 8 */
    m->cx=(uint16_t)(8);
L_1689D: ; /* IDA 0x1689D; mov     ax, cs:[si] */
    m->ax=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->si)));
L_168A0: ; /* IDA 0x168A0; test    cl, 3 */
    {uint16_t result=(((uint8_t)m->cx)&(3))&255;ec_logic(m,result,8);}
L_168A3: ; /* IDA 0x168A3; jz      short loc_168AA */
    if(m->flags&0x40) goto L_168AA;
L_168A5: ; /* IDA 0x168A5; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_168A8: ; /* IDA 0x168A8; mul     dx */
    dl_mul(m,m->dx,16);
L_168AA: ; /* IDA 0x168AA; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_168AC: ; /* IDA 0x168AC; mov     bx, 905h */
    m->bx=(uint16_t)(2309);
L_168AF: ; /* IDA 0x168AF; call    sub_1062F */
    sty_call(m,1583,26802,h);
L_168B2: ; /* IDA 0x168B2; add     di, 500h */
    m->di=(uint16_t)(ec_math(m,m->di,1280,16,0,0));
L_168B6: ; /* IDA 0x168B6; cmp     cx, 5 */
    ec_cmp(m,m->cx,5,16);
L_168B9: ; /* IDA 0x168B9; jnz     short loc_168BE */
    if(!(m->flags&0x40)) goto L_168BE;
L_168BB: ; /* IDA 0x168BB; mov     di, 3223h */
    m->di=(uint16_t)(12835);
L_168BE: ; /* IDA 0x168BE; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_168BF: ; /* IDA 0x168BF; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_168C0: ; /* IDA 0x168C0; loop    loc_1689D */
    if(--m->cx) goto L_1689D;
L_168C2: ; /* IDA 0x168C2; retn */
    return;
}
/* Original sub_168C3; IDA 0x168C3. */
static void sty_body_168C3(KiMachine16 *m,const KiEconomyHooks *h) {
L_168C3: ; /* IDA 0x168C3; mov     al, 1Fh */
    vg_al(m,(uint8_t)(31));
L_168C5: ; /* IDA 0x168C5; mov     dx, 1 */
    m->dx=(uint16_t)(1);
L_168C8: ; /* IDA 0x168C8; mov     bx, 3 */
    m->bx=(uint16_t)(3);
L_168CB: ; /* IDA 0x168CB; mov     cx, 0A15h */
    m->cx=(uint16_t)(2581);
L_168CE: ; /* IDA 0x168CE; call    sub_1895D */
    sty_call(m,35165,26833,h);
L_168D1: ; /* IDA 0x168D1; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_168D3: ; /* IDA 0x168D3; mov     dx, 18h */
    m->dx=(uint16_t)(24);
L_168D6: ; /* IDA 0x168D6; mov     bx, 58h ; 'X' */
    m->bx=(uint16_t)(88);
L_168D9: ; /* IDA 0x168D9; call    sub_10337 */
    sty_call(m,823,26844,h);
L_168DC: ; /* IDA 0x168DC; mov     al, 20h ; ' ' */
    vg_al(m,(uint8_t)(32));
L_168DE: ; /* IDA 0x168DE; mov     si, 0A0h */
    m->si=(uint16_t)(160);
L_168E1: ; /* IDA 0x168E1; mov     di, 4 */
    m->di=(uint16_t)(4);
L_168E4: ; /* IDA 0x168E4; push    ax */
    ec_push(m,m->ax);
L_168E5: ; /* IDA 0x168E5; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_168E7: ; /* IDA 0x168E7; mov     dx, 100h */
    m->dx=(uint16_t)(256);
L_168EA: ; /* IDA 0x168EA; mov     cx, 203h */
    m->cx=(uint16_t)(515);
L_168ED: ; /* IDA 0x168ED; call    sub_1E3D7 */
    sty_call(m,58327,26864,h);
L_168F0: ; /* IDA 0x168F0; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_168F1: ; /* IDA 0x168F1; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_168F3: ; /* IDA 0x168F3; add     si, 10h */
    m->si=(uint16_t)(ec_math(m,m->si,16,16,0,0));
L_168F6: ; /* IDA 0x168F6; dec     di */
    m->di=(uint16_t)(st_dec(m,m->di,16));
L_168F7: ; /* IDA 0x168F7; jnz     short loc_168E4 */
    if(!(m->flags&0x40)) goto L_168E4;
L_168F9: ; /* IDA 0x168F9; retn */
    return;
}
/* Original sub_168FA; IDA 0x168FA. */
static void sty_body_168FA(KiMachine16 *m,const KiEconomyHooks *h) {
L_168FA: ; /* IDA 0x168FA; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_168FC: ; /* IDA 0x168FC; mov     dx, 1 */
    m->dx=(uint16_t)(1);
L_168FF: ; /* IDA 0x168FF; mov     bx, 3 */
    m->bx=(uint16_t)(3);
L_16902: ; /* IDA 0x16902; mov     cx, 0A15h */
    m->cx=(uint16_t)(2581);
L_16905: ; /* IDA 0x16905; call    sub_1895D */
    sty_call(m,35165,26888,h);
L_16908: ; /* IDA 0x16908; retn */
    return;
}
/* Original sub_16A28; IDA 0x16A28. */
static void sty_body_16A28(KiMachine16 *m,const KiEconomyHooks *h) {
L_16A28: ; /* IDA 0x16A28; mov     al, [di+23h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+35))));
L_16A2B: ; /* IDA 0x16A2B; mov     ah, 19h */
    vg_ah(m,(uint8_t)(25));
L_16A2D: ; /* IDA 0x16A2D; mul     ah */
    dl_mul(m,(uint8_t)(m->ax>>8),8);
L_16A2F: ; /* IDA 0x16A2F; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_16A31: ; /* IDA 0x16A31; mov     al, [si+23h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+35))));
L_16A34: ; /* IDA 0x16A34; mov     ah, [si+28h] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+40))));
L_16A37: ; /* IDA 0x16A37; add     ah, 14h */
#if KI_STRATEGY_MUTATION == 10
    vg_ah(m,ec_math(m,(uint8_t)(m->ax>>8),19,8,0,0));
#else
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),20,8,0,0)));
#endif
L_16A3A: ; /* IDA 0x16A3A; mul     ah */
    dl_mul(m,(uint8_t)(m->ax>>8),8);
L_16A3C: ; /* IDA 0x16A3C; retn */
    return;
}
int ki_strategy_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x2f5:sty_body_102F5(m,h);break;
case 0x300e:sty_body_1300E(m,h);break;
case 0x304e:sty_body_1304E(m,h);break;
case 0x3830:sty_body_13830(m,h);break;
case 0x3b5a:sty_body_13B5A(m,h);break;
case 0x3ba9:sty_body_13BA9(m,h);break;
case 0x3c1e:sty_body_13C1E(m,h);break;
case 0x3d91:sty_body_13D91(m,h);break;
case 0x61ca:sty_body_161CA(m,h);break;
case 0x6224:sty_body_16224(m,h);break;
case 0x6288:sty_body_16288(m,h);break;
case 0x628f:sty_body_1628F(m,h);break;
case 0x62fb:sty_body_162FB(m,h);break;
case 0x6366:sty_body_16366(m,h);break;
case 0x63bf:sty_body_163BF(m,h);break;
case 0x6405:sty_body_16405(m,h);break;
case 0x6475:sty_body_16475(m,h);break;
case 0x64f1:sty_body_164F1(m,h);break;
case 0x6577:sty_body_16577(m,h);break;
case 0x65ef:sty_body_165EF(m,h);break;
case 0x6605:sty_body_16605(m,h);break;
case 0x6623:sty_body_16623(m,h);break;
case 0x66d9:sty_body_166D9(m,h);break;
case 0x676f:sty_body_1676F(m,h);break;
case 0x678d:sty_body_1678D(m,h);break;
case 0x6846:sty_body_16846(m,h);break;
case 0x68c3:sty_body_168C3(m,h);break;
case 0x68fa:sty_body_168FA(m,h);break;
case 0x6a28:sty_body_16A28(m,h);break;
default:return 0;}return 1;}
static int sty_known(uint16_t t) {switch(t) {case 0x2f5: case 0x300e: case 0x304e: case 0x3830: case 0x3b5a: case 0x3ba9: case 0x3c1e: case 0x3d91: case 0x61ca: case 0x6224: case 0x6288: case 0x628f: case 0x62fb: case 0x6366: case 0x63bf: case 0x6405: case 0x6475: case 0x64f1: case 0x6577: case 0x65ef: case 0x6605: case 0x6623: case 0x66d9: case 0x676f: case 0x678d: case 0x6846: case 0x68c3: case 0x68fa: case 0x6a28: return 1;default:return 0;}}
void sub_102F5(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x2f5,h);m->ip=ec_pop(m);}
void sub_1300E(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x300e,h);m->ip=ec_pop(m);}
void sub_1304E(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x304e,h);m->ip=ec_pop(m);}
void sub_13830(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x3830,h);m->ip=ec_pop(m);}
void sub_13B5A(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x3b5a,h);m->ip=ec_pop(m);}
void sub_13BA9(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x3ba9,h);m->ip=ec_pop(m);}
void sub_13C1E(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x3c1e,h);m->ip=ec_pop(m);}
void sub_13D91(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x3d91,h);m->ip=ec_pop(m);}
void sub_161CA(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x61ca,h);m->ip=ec_pop(m);}
void sub_16224(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6224,h);m->ip=ec_pop(m);}
void sub_16288(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6288,h);m->ip=ec_pop(m);}
void sub_1628F(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x628f,h);m->ip=ec_pop(m);}
void sub_162FB(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x62fb,h);m->ip=ec_pop(m);}
void sub_16366(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6366,h);m->ip=ec_pop(m);}
void sub_163BF(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x63bf,h);m->ip=ec_pop(m);}
void sub_16405(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6405,h);m->ip=ec_pop(m);}
void sub_16475(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6475,h);m->ip=ec_pop(m);}
void sub_164F1(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x64f1,h);m->ip=ec_pop(m);}
void sub_16577(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6577,h);m->ip=ec_pop(m);}
void sub_165EF(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x65ef,h);m->ip=ec_pop(m);}
void sub_16605(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6605,h);m->ip=ec_pop(m);}
void sub_16623(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6623,h);m->ip=ec_pop(m);}
void sub_166D9(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x66d9,h);m->ip=ec_pop(m);}
void sub_1676F(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x676f,h);m->ip=ec_pop(m);}
void sub_1678D(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x678d,h);m->ip=ec_pop(m);}
void sub_16846(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6846,h);m->ip=ec_pop(m);}
void sub_168C3(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x68c3,h);m->ip=ec_pop(m);}
void sub_168FA(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x68fa,h);m->ip=ec_pop(m);}
void sub_16A28(KiMachine16 *m,const KiEconomyHooks *h) {ki_strategy_body(m,0x6a28,h);m->ip=ec_pop(m);}
