/* Original map-cell native C; re/114, spec/234. No opcode fetch. */
/* Original sub_1D46A; IDA 0x1D46A. */
static void mc_body_1D46A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D46A: ; /* IDA 0x1D46A; mov     cs:word_1D84A, ax */
    dl_write16(m,m->cs,55370,m->ax);
L_1D46E: ; /* IDA 0x1D46E; add     ax, 800h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2048,16,0,0));
L_1D471: ; /* IDA 0x1D471; mov     cs:word_1D84C, ax */
    dl_write16(m,m->cs,55372,m->ax);
L_1D475: ; /* IDA 0x1D475; mov     cs:word_1D84E, bx */
    dl_write16(m,m->cs,55374,m->bx);
L_1D47A: ; /* IDA 0x1D47A; mov     cs:word_1D850, cx */
    dl_write16(m,m->cs,55376,m->cx);
L_1D47F: ; /* IDA 0x1D47F; call    sub_1D483 */
    mc_call(m,54403,54402,h);
L_1D482: ; /* IDA 0x1D482; retn */
    return;
}
/* Original sub_1D483; IDA 0x1D483. */
static void mc_body_1D483(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D483: ; /* IDA 0x1D483; push    es */
    ec_push(m,m->es);
L_1D484: ; /* IDA 0x1D484; push    ax */
    ec_push(m,m->ax);
L_1D485: ; /* IDA 0x1D485; push    cx */
    ec_push(m,m->cx);
L_1D486: ; /* IDA 0x1D486; push    di */
    ec_push(m,m->di);
L_1D487: ; /* IDA 0x1D487; mov     es, cs:word_1D84E */
    m->es=(uint16_t)(dl_read16(m,m->cs,55374));
L_1D48C: ; /* IDA 0x1D48C; xor     di, di */
    {uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_1D48E: ; /* IDA 0x1D48E; mov     cx, 398h */
    m->cx=(uint16_t)(920);
L_1D491: ; /* IDA 0x1D491; cld */
    m->flags&=(uint16_t)~0x400u;
L_1D492: ; /* IDA 0x1D492; mov     ax, 40h ; '@' */
    m->ax=(uint16_t)(64);
L_1D495: ; /* IDA 0x1D495; stosw */
    mc_string(m,2,0);
L_1D496: ; /* IDA 0x1D496; mov     ax, 0FFFFh */
    m->ax=(uint16_t)(65535);
L_1D499: ; /* IDA 0x1D499; stosw */
    mc_string(m,2,0);
L_1D49A: ; /* IDA 0x1D49A; stosw */
    mc_string(m,2,0);
L_1D49B: ; /* IDA 0x1D49B; stosw */
    mc_string(m,2,0);
L_1D49C: ; /* IDA 0x1D49C; loop    loc_1D492 */
    if(--m->cx) goto L_1D492;
L_1D49E: ; /* IDA 0x1D49E; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1D49F: ; /* IDA 0x1D49F; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1D4A0: ; /* IDA 0x1D4A0; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1D4A1: ; /* IDA 0x1D4A1; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1D4A2: ; /* IDA 0x1D4A2; retn */
    return;
}
/* Original sub_1D4C7; IDA 0x1D4C7. */
static void mc_body_1D4C7(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D4C7: ; /* IDA 0x1D4C7; sub     bx, cs:word_1D856 */
    m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->cs,55382),16,0,1));
L_1D4CC: ; /* IDA 0x1D4CC; jb      short locret_1D51E */
    if(m->flags&1) goto L_1D51E;
L_1D4CE: ; /* IDA 0x1D4CE; cmp     bx, 17h */
    ec_cmp(m,m->bx,23,16);
L_1D4D1: ; /* IDA 0x1D4D1; jnb     short locret_1D51E */
    if(!(m->flags&1)) goto L_1D51E;
L_1D4D3: ; /* IDA 0x1D4D3; sub     dx, cs:word_1D854 */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->cs,55380),16,0,1));
L_1D4D8: ; /* IDA 0x1D4D8; jb      short locret_1D51E */
    if(m->flags&1) goto L_1D51E;
L_1D4DA: ; /* IDA 0x1D4DA; cmp     dx, 28h ; '(' */
    ec_cmp(m,m->dx,40,16);
L_1D4DD: ; /* IDA 0x1D4DD; jnb     short locret_1D51E */
    if(!(m->flags&1)) goto L_1D51E;
L_1D4DF: ; /* IDA 0x1D4DF; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4E1: ; /* IDA 0x1D4E1; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4E3: ; /* IDA 0x1D4E3; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4E5: ; /* IDA 0x1D4E5; add     dx, bx */
    m->dx=(uint16_t)(ec_math(m,m->dx,m->bx,16,0,0));
L_1D4E7: ; /* IDA 0x1D4E7; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4E9: ; /* IDA 0x1D4E9; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4EB: ; /* IDA 0x1D4EB; add     bx, dx */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
L_1D4ED: ; /* IDA 0x1D4ED; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4EF: ; /* IDA 0x1D4EF; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4F1: ; /* IDA 0x1D4F1; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D4F3: ; /* IDA 0x1D4F3; push    ds */
    ec_push(m,m->ds);
L_1D4F4: ; /* IDA 0x1D4F4; push    si */
    ec_push(m,m->si);
L_1D4F5: ; /* IDA 0x1D4F5; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_1D4F7: ; /* IDA 0x1D4F7; mov     ds, cs:word_1D84E */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55374));
L_1D4FC: ; /* IDA 0x1D4FC; mov     bx, [si] */
    m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_1D4FE: ; /* IDA 0x1D4FE; test    bl, 10h */
#if KI_MAPCELLS_MUTATION == 12
    {uint16_t result=(((uint8_t)m->bx)&(0))&255;ec_logic(m,result,8);}
#else
    {uint16_t result=(((uint8_t)m->bx)&(16))&255;ec_logic(m,result,8);}
#endif
L_1D501: ; /* IDA 0x1D501; jnz     short loc_1D51C */
    if(!(m->flags&0x40)) goto L_1D51C;
L_1D503: ; /* IDA 0x1D503; mov     bl, bh */
    gw_bl(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_1D505: ; /* IDA 0x1D505; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1D507: ; /* IDA 0x1D507; cmp     bl, 4 */
#if KI_MAPCELLS_MUTATION == 11
    ec_cmp(m,(uint8_t)m->bx,5,8);
#else
    ec_cmp(m,(uint8_t)m->bx,4,8);
#endif
L_1D50A: ; /* IDA 0x1D50A; jnb     short loc_1D51C */
    if(!(m->flags&1)) goto L_1D51C;
L_1D50C: ; /* IDA 0x1D50C; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_1D50E: ; /* IDA 0x1D50E; xchg    al, [bx+si+3] */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->ds,(uint16_t)(m->bx+m->si+3));vg_al(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->bx+m->si+3),a);}
L_1D511: ; /* IDA 0x1D511; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1D512: ; /* IDA 0x1D512; mov     [si+1], bl */
    dl_write8(m,m->ds,(uint16_t)(m->si+1),(uint8_t)m->bx);
L_1D515: ; /* IDA 0x1D515; cmp     al, ah */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
L_1D517: ; /* IDA 0x1D517; jz      short loc_1D51C */
    if(m->flags&0x40) goto L_1D51C;
L_1D519: ; /* IDA 0x1D519; or      byte ptr [si], 20h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(32))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_1D51C: ; /* IDA 0x1D51C; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1D51D: ; /* IDA 0x1D51D; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1D51E: ; /* IDA 0x1D51E; retn */
    return;
}
/* Original sub_1D615; IDA 0x1D615. */
static void mc_body_1D615(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D615: ; /* IDA 0x1D615; push    ds */
    ec_push(m,m->ds);
L_1D616: ; /* IDA 0x1D616; push    es */
    ec_push(m,m->es);
L_1D617: ; /* IDA 0x1D617; push    si */
    ec_push(m,m->si);
L_1D618: ; /* IDA 0x1D618; push    di */
    ec_push(m,m->di);
L_1D619: ; /* IDA 0x1D619; mov     cs:word_1D854, dx */
    dl_write16(m,m->cs,55380,m->dx);
L_1D61E: ; /* IDA 0x1D61E; mov     cs:word_1D856, bx */
    dl_write16(m,m->cs,55382,m->bx);
L_1D623: ; /* IDA 0x1D623; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_1D625: ; /* IDA 0x1D625; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1D627: ; /* IDA 0x1D627; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D629: ; /* IDA 0x1D629; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D62B: ; /* IDA 0x1D62B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D62D: ; /* IDA 0x1D62D; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_1D62F: ; /* IDA 0x1D62F; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1D631: ; /* IDA 0x1D631; add     ax, bx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->bx,16,0,0));
L_1D633: ; /* IDA 0x1D633; add     ax, cs:word_1D850 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,55376),16,0,0));
L_1D638: ; /* IDA 0x1D638; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1D63A: ; /* IDA 0x1D63A; mov     ds, cs:word_1D84E */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55374));
L_1D63F: ; /* IDA 0x1D63F; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_1D641: ; /* IDA 0x1D641; mov     bx, 2 */
    m->bx=(uint16_t)(2);
L_1D644: ; /* IDA 0x1D644; mov     dx, 17h */
    m->dx=(uint16_t)(23);
L_1D647: ; /* IDA 0x1D647; mov     cx, 28h ; '(' */
    m->cx=(uint16_t)(40);
L_1D64A: ; /* IDA 0x1D64A; mov     al, es:[di] */
    vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->di))));
L_1D64D: ; /* IDA 0x1D64D; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_1D64F: ; /* IDA 0x1D64F; xchg    al, [bx+si] */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->ds,(uint16_t)(m->bx+m->si));vg_al(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->bx+m->si),a);}
L_1D651: ; /* IDA 0x1D651; cmp     al, ah */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8);
L_1D653: ; /* IDA 0x1D653; jz      short loc_1D658 */
    if(m->flags&0x40) goto L_1D658;
L_1D655: ; /* IDA 0x1D655; or      byte ptr [si], 60h */
#if KI_MAPCELLS_MUTATION == 2
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(64))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(96))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#endif
L_1D658: ; /* IDA 0x1D658; add     si, 8 */
    m->si=(uint16_t)(ec_math(m,m->si,8,16,0,0));
L_1D65B: ; /* IDA 0x1D65B; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1D65C: ; /* IDA 0x1D65C; loop    loc_1D64A */
    if(--m->cx) goto L_1D64A;
L_1D65E: ; /* IDA 0x1D65E; add     di, 158h */
#if KI_MAPCELLS_MUTATION == 1
    m->di=(uint16_t)(ec_math(m,m->di,345,16,0,0));
#else
    m->di=(uint16_t)(ec_math(m,m->di,344,16,0,0));
#endif
L_1D662: ; /* IDA 0x1D662; dec     dx */
    m->dx=(uint16_t)(st_dec(m,m->dx,16));
L_1D663: ; /* IDA 0x1D663; jnz     short loc_1D647 */
    if(!(m->flags&0x40)) goto L_1D647;
L_1D665: ; /* IDA 0x1D665; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1D666: ; /* IDA 0x1D666; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1D667: ; /* IDA 0x1D667; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1D668: ; /* IDA 0x1D668; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1D669: ; /* IDA 0x1D669; retn */
    return;
}
/* Original sub_1D66A; IDA 0x1D66A. */
static void mc_body_1D66A(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D66A: ; /* IDA 0x1D66A; push    ds */
    ec_push(m,m->ds);
L_1D66B: ; /* IDA 0x1D66B; mov     dx, 3CEh */
    m->dx=(uint16_t)(974);
L_1D66E: ; /* IDA 0x1D66E; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_1D671: ; /* IDA 0x1D671; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_1D672: ; /* IDA 0x1D672; mov     ax, 1 */
    m->ax=(uint16_t)(1);
L_1D675: ; /* IDA 0x1D675; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_1D676: ; /* IDA 0x1D676; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_1D679: ; /* IDA 0x1D679; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_1D67A: ; /* IDA 0x1D67A; mov     ax, 0FF08h */
    m->ax=(uint16_t)(65288);
L_1D67D: ; /* IDA 0x1D67D; out     dx, ax; EGA: graph 1 and 2 addr reg: */
    dl_out(m,m->dx,m->ax,16);
L_1D67E: ; /* IDA 0x1D67E; cld */
    m->flags&=(uint16_t)~0x400u;
L_1D67F: ; /* IDA 0x1D67F; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_1D682: ; /* IDA 0x1D682; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1D684: ; /* IDA 0x1D684; mov     bp, cs:word_1D84E */
    m->bp=(uint16_t)(dl_read16(m,m->cs,55374));
L_1D689: ; /* IDA 0x1D689; mov     ds, bp */
    m->ds=(uint16_t)(m->bp);
L_1D68B: ; /* IDA 0x1D68B; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_1D68D: ; /* IDA 0x1D68D; mov     di, 0A00h */
    m->di=(uint16_t)(2560);
L_1D690: ; /* IDA 0x1D690; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1D692: ; /* IDA 0x1D692; mov     ch, 17h */
    vg_ch(m,(uint8_t)(23));
L_1D694: ; /* IDA 0x1D694; mov     cl, 28h ; '(' */
    vg_cl(m,(uint8_t)(40));
L_1D696: ; /* IDA 0x1D696; mov     dx, [si] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si)));
L_1D698: ; /* IDA 0x1D698; test    dl, 10h */
#if KI_MAPCELLS_MUTATION == 3
    {uint16_t result=(((uint8_t)m->dx)&(0))&255;ec_logic(m,result,8);}
#else
    {uint16_t result=(((uint8_t)m->dx)&(16))&255;ec_logic(m,result,8);}
#endif
L_1D69B: ; /* IDA 0x1D69B; jz      short loc_1D6A0 */
    if(m->flags&0x40) goto L_1D6A0;
L_1D69D: ; /* IDA 0x1D69D; jmp     loc_1D747 */
    goto L_1D747;
L_1D6A0: ; /* IDA 0x1D6A0; and     dh, dh */
    {uint16_t result=(((uint8_t)(m->dx>>8))&((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);vg_dh(m,(uint8_t)(result));}
L_1D6A2: ; /* IDA 0x1D6A2; jz      short loc_1D6B4 */
    if(m->flags&0x40) goto L_1D6B4;
L_1D6A4: ; /* IDA 0x1D6A4; test    dl, 20h */
    {uint16_t result=(((uint8_t)m->dx)&(32))&255;ec_logic(m,result,8);}
L_1D6A7: ; /* IDA 0x1D6A7; jnz     short loc_1D6DF */
    if(!(m->flags&0x40)) goto L_1D6DF;
L_1D6A9: ; /* IDA 0x1D6A9; mov     al, dl */
    vg_al(m,(uint8_t)((uint8_t)m->dx));
L_1D6AB: ; /* IDA 0x1D6AB; and     al, 0Fh */
    {uint16_t result=(((uint8_t)m->ax)&(15))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1D6AD: ; /* IDA 0x1D6AD; cmp     al, dh */
    ec_cmp(m,(uint8_t)m->ax,(uint8_t)(m->dx>>8),8);
L_1D6AF: ; /* IDA 0x1D6AF; jnz     short loc_1D6DF */
    if(!(m->flags&0x40)) goto L_1D6DF;
L_1D6B1: ; /* IDA 0x1D6B1; jmp     loc_1D747 */
    goto L_1D747;
L_1D6B4: ; /* IDA 0x1D6B4; test    dl, 40h */
    {uint16_t result=(((uint8_t)m->dx)&(64))&255;ec_logic(m,result,8);}
L_1D6B7: ; /* IDA 0x1D6B7; jnz     short loc_1D6BC */
    if(!(m->flags&0x40)) goto L_1D6BC;
L_1D6B9: ; /* IDA 0x1D6B9; jmp     loc_1D747 */
    goto L_1D747;
L_1D6BC: ; /* IDA 0x1D6BC; cmp     dl, 80h */
    ec_cmp(m,(uint8_t)m->dx,128,8);
L_1D6BF: ; /* IDA 0x1D6BF; jnb     short loc_1D71B */
    if(!(m->flags&1)) goto L_1D71B;
L_1D6C1: ; /* IDA 0x1D6C1; mov     byte ptr [si], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si),0);
L_1D6C4: ; /* IDA 0x1D6C4; mov     ah, [si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_1D6C7: ; /* IDA 0x1D6C7; mov     ds, cs:word_1D84A */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55370));
L_1D6CC: ; /* IDA 0x1D6CC; call    sub_1D782 */
    mc_call(m,55170,54991,h);
L_1D6CF: ; /* IDA 0x1D6CF; mov     ds, bp */
    m->ds=(uint16_t)(m->bp);
L_1D6D1: ; /* IDA 0x1D6D1; mov     ax, 0FFFFh */
    m->ax=(uint16_t)(65535);
L_1D6D4: ; /* IDA 0x1D6D4; mov     [si+3], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+3),(uint8_t)m->ax);
L_1D6D7: ; /* IDA 0x1D6D7; mov     [si+4], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+4),m->ax);
L_1D6DA: ; /* IDA 0x1D6DA; mov     [si+6], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+6),m->ax);
L_1D6DD: ; /* IDA 0x1D6DD; jmp     short loc_1D747 */
    goto L_1D747;
L_1D6DF: ; /* IDA 0x1D6DF; mov     al, dh */
    vg_al(m,(uint8_t)((uint8_t)(m->dx>>8)));
L_1D6E1: ; /* IDA 0x1D6E1; or      al, 40h */
    {uint16_t result=(((uint8_t)m->ax)|(64))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1D6E3: ; /* IDA 0x1D6E3; mov     [si], al */
    dl_write8(m,m->ds,(uint16_t)(m->si),(uint8_t)m->ax);
L_1D6E5: ; /* IDA 0x1D6E5; mov     ah, [si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_1D6E8: ; /* IDA 0x1D6E8; mov     ds, cs:word_1D84A */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55370));
L_1D6ED: ; /* IDA 0x1D6ED; call    sub_1D7E7 */
    mc_call(m,55271,55024,h);
L_1D6F0: ; /* IDA 0x1D6F0; mov     ds, bp */
    m->ds=(uint16_t)(m->bp);
L_1D6F2: ; /* IDA 0x1D6F2; mov     bl, 3 */
    gw_bl(m,(uint8_t)(3));
L_1D6F4: ; /* IDA 0x1D6F4; mov     ah, [bx+si] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+m->si))));
L_1D6F6: ; /* IDA 0x1D6F6; mov     ds, cs:word_1D84C */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55372));
L_1D6FB: ; /* IDA 0x1D6FB; call    sub_1D804 */
    mc_call(m,55300,55038,h);
L_1D6FE: ; /* IDA 0x1D6FE; mov     ds, bp */
    m->ds=(uint16_t)(m->bp);
L_1D700: ; /* IDA 0x1D700; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1D701: ; /* IDA 0x1D701; dec     dh */
    vg_dh(m,(uint8_t)(st_dec(m,(uint8_t)(m->dx>>8),8)));
L_1D703: ; /* IDA 0x1D703; jnz     short loc_1D6F4 */
    if(!(m->flags&0x40)) goto L_1D6F4;
L_1D705: ; /* IDA 0x1D705; cmp     bl, 8 */
#if KI_MAPCELLS_MUTATION == 9
    ec_cmp(m,(uint8_t)m->bx,7,8);
#else
    ec_cmp(m,(uint8_t)m->bx,8,8);
#endif
L_1D708: ; /* IDA 0x1D708; jnb     short loc_1D714 */
    if(!(m->flags&1)) goto L_1D714;
L_1D70A: ; /* IDA 0x1D70A; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1D70C: ; /* IDA 0x1D70C; mov     [bx+si], al */
    dl_write8(m,m->ds,(uint16_t)(m->bx+m->si),(uint8_t)m->ax);
L_1D70E: ; /* IDA 0x1D70E; inc     bx */
    m->bx=(uint16_t)(ec_inc(m,m->bx,16));
L_1D70F: ; /* IDA 0x1D70F; cmp     bl, 8 */
    ec_cmp(m,(uint8_t)m->bx,8,8);
L_1D712: ; /* IDA 0x1D712; jb      short loc_1D70C */
    if(m->flags&1) goto L_1D70C;
L_1D714: ; /* IDA 0x1D714; cmp     dl, 80h */
    ec_cmp(m,(uint8_t)m->dx,128,8);
L_1D717: ; /* IDA 0x1D717; jnb     short loc_1D735 */
    if(!(m->flags&1)) goto L_1D735;
L_1D719: ; /* IDA 0x1D719; jmp     short loc_1D744 */
    goto L_1D744;
L_1D71B: ; /* IDA 0x1D71B; mov     byte ptr [si], 40h ; '@' */
    dl_write8(m,m->ds,(uint16_t)(m->si),64);
L_1D71E: ; /* IDA 0x1D71E; mov     ah, [si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_1D721: ; /* IDA 0x1D721; mov     dx, 0FFFFh */
    m->dx=(uint16_t)(65535);
L_1D724: ; /* IDA 0x1D724; mov     [si+3], dl */
    dl_write8(m,m->ds,(uint16_t)(m->si+3),(uint8_t)m->dx);
L_1D727: ; /* IDA 0x1D727; mov     [si+4], dx */
    dl_write16(m,m->ds,(uint16_t)(m->si+4),m->dx);
L_1D72A: ; /* IDA 0x1D72A; mov     [si+6], dx */
    dl_write16(m,m->ds,(uint16_t)(m->si+6),m->dx);
L_1D72D: ; /* IDA 0x1D72D; mov     ds, cs:word_1D84A */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55370));
L_1D732: ; /* IDA 0x1D732; call    sub_1D7E7 */
    mc_call(m,55271,55093,h);
L_1D735: ; /* IDA 0x1D735; mov     ah, 0FFh */
#if KI_MAPCELLS_MUTATION == 10
    vg_ah(m,(uint8_t)(254));
#else
    vg_ah(m,(uint8_t)(255));
#endif
L_1D737: ; /* IDA 0x1D737; mov     ds, cs:word_1D84C */
    m->ds=(uint16_t)(dl_read16(m,m->cs,55372));
L_1D73C: ; /* IDA 0x1D73C; call    sub_1D804 */
    mc_call(m,55300,55103,h);
L_1D73F: ; /* IDA 0x1D73F; mov     ds, bp */
    m->ds=(uint16_t)(m->bp);
L_1D741: ; /* IDA 0x1D741; or      byte ptr [si], 60h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(96))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_1D744: ; /* IDA 0x1D744; call    sub_1D76B */
    mc_call(m,55147,55111,h);
L_1D747: ; /* IDA 0x1D747; mov     byte ptr [si+1], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+1),0);
L_1D74B: ; /* IDA 0x1D74B; add     si, 8 */
    m->si=(uint16_t)(ec_math(m,m->si,8,16,0,0));
L_1D74E: ; /* IDA 0x1D74E; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1D74F: ; /* IDA 0x1D74F; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1D750: ; /* IDA 0x1D750; dec     cl */
    vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
L_1D752: ; /* IDA 0x1D752; jz      short loc_1D757 */
    if(m->flags&0x40) goto L_1D757;
L_1D754: ; /* IDA 0x1D754; jmp     loc_1D696 */
    goto L_1D696;
L_1D757: ; /* IDA 0x1D757; add     di, 4B0h */
    m->di=(uint16_t)(ec_math(m,m->di,1200,16,0,0));
L_1D75B: ; /* IDA 0x1D75B; dec     ch */
    vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
L_1D75D: ; /* IDA 0x1D75D; jz      short loc_1D762 */
    if(m->flags&0x40) goto L_1D762;
L_1D75F: ; /* IDA 0x1D75F; jmp     loc_1D694 */
    goto L_1D694;
L_1D762: ; /* IDA 0x1D762; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_1D765: ; /* IDA 0x1D765; mov     ax, 0F02h */
    m->ax=(uint16_t)(3842);
L_1D768: ; /* IDA 0x1D768; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_1D769: ; /* IDA 0x1D769; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1D76A: ; /* IDA 0x1D76A; retn */
    return;
}
/* Original sub_1D76B; IDA 0x1D76B. */
static void mc_body_1D76B(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D76B: ; /* IDA 0x1D76B; push    ds */
    ec_push(m,m->ds);
L_1D76C: ; /* IDA 0x1D76C; push    bx */
    ec_push(m,m->bx);
L_1D76D: ; /* IDA 0x1D76D; push    cx */
    ec_push(m,m->cx);
L_1D76E: ; /* IDA 0x1D76E; push    dx */
    ec_push(m,m->dx);
L_1D76F: ; /* IDA 0x1D76F; push    si */
    ec_push(m,m->si);
L_1D770: ; /* IDA 0x1D770; push    di */
    ec_push(m,m->di);
L_1D771: ; /* IDA 0x1D771; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_1D773: ; /* IDA 0x1D773; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1D775: ; /* IDA 0x1D775; mov     si, 0D858h */
#if KI_MAPCELLS_MUTATION == 4
    m->si=(uint16_t)(55386);
#else
    m->si=(uint16_t)(55384);
#endif
L_1D778: ; /* IDA 0x1D778; call    sub_1D796 */
    mc_call(m,55190,55163,h);
L_1D77B: ; /* IDA 0x1D77B; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1D77C: ; /* IDA 0x1D77C; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1D77D: ; /* IDA 0x1D77D; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1D77E: ; /* IDA 0x1D77E; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1D77F: ; /* IDA 0x1D77F; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1D780: ; /* IDA 0x1D780; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1D781: ; /* IDA 0x1D781; retn */
    return;
}
/* Original sub_1D782; IDA 0x1D782. */
static void mc_body_1D782(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D782: ; /* IDA 0x1D782; push    bx */
    ec_push(m,m->bx);
L_1D783: ; /* IDA 0x1D783; push    cx */
    ec_push(m,m->cx);
L_1D784: ; /* IDA 0x1D784; push    dx */
    ec_push(m,m->dx);
L_1D785: ; /* IDA 0x1D785; push    si */
    ec_push(m,m->si);
L_1D786: ; /* IDA 0x1D786; push    di */
    ec_push(m,m->di);
L_1D787: ; /* IDA 0x1D787; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1D789: ; /* IDA 0x1D789; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1D78B: ; /* IDA 0x1D78B; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1D78D: ; /* IDA 0x1D78D; call    sub_1D796 */
    mc_call(m,55190,55184,h);
L_1D790: ; /* IDA 0x1D790; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1D791: ; /* IDA 0x1D791; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1D792: ; /* IDA 0x1D792; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1D793: ; /* IDA 0x1D793; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1D794: ; /* IDA 0x1D794; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1D795: ; /* IDA 0x1D795; retn */
    return;
}
/* Original sub_1D796; IDA 0x1D796. */
static void mc_body_1D796(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D796: ; /* IDA 0x1D796; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_1D799: ; /* IDA 0x1D799; mov     ax, 102h */
    m->ax=(uint16_t)(258);
L_1D79C: ; /* IDA 0x1D79C; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_1D79D: ; /* IDA 0x1D79D; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_1D79F: ; /* IDA 0x1D79F; mov     cx, 10h */
    m->cx=(uint16_t)(16);
L_1D7A2: ; /* IDA 0x1D7A2; movsw */
    mc_string(m,0,0);
L_1D7A3: ; /* IDA 0x1D7A3; add     di, 4Eh ; 'N' */
#if KI_MAPCELLS_MUTATION == 8
    m->di=(uint16_t)(ec_math(m,m->di,79,16,0,0));
#else
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
#endif
L_1D7A6: ; /* IDA 0x1D7A6; loop    loc_1D7A2 */
    if(--m->cx) goto L_1D7A2;
L_1D7A8: ; /* IDA 0x1D7A8; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_1D7AA: ; /* IDA 0x1D7AA; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_1D7AD: ; /* IDA 0x1D7AD; mov     ax, 202h */
    m->ax=(uint16_t)(514);
L_1D7B0: ; /* IDA 0x1D7B0; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_1D7B1: ; /* IDA 0x1D7B1; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_1D7B3: ; /* IDA 0x1D7B3; mov     cx, 10h */
    m->cx=(uint16_t)(16);
L_1D7B6: ; /* IDA 0x1D7B6; movsw */
    mc_string(m,0,0);
L_1D7B7: ; /* IDA 0x1D7B7; add     di, 4Eh ; 'N' */
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
L_1D7BA: ; /* IDA 0x1D7BA; loop    loc_1D7B6 */
    if(--m->cx) goto L_1D7B6;
L_1D7BC: ; /* IDA 0x1D7BC; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_1D7BE: ; /* IDA 0x1D7BE; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_1D7C1: ; /* IDA 0x1D7C1; mov     ax, 402h */
    m->ax=(uint16_t)(1026);
L_1D7C4: ; /* IDA 0x1D7C4; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_1D7C5: ; /* IDA 0x1D7C5; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_1D7C7: ; /* IDA 0x1D7C7; mov     cx, 10h */
    m->cx=(uint16_t)(16);
L_1D7CA: ; /* IDA 0x1D7CA; movsw */
    mc_string(m,0,0);
L_1D7CB: ; /* IDA 0x1D7CB; add     di, 4Eh ; 'N' */
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
L_1D7CE: ; /* IDA 0x1D7CE; loop    loc_1D7CA */
    if(--m->cx) goto L_1D7CA;
L_1D7D0: ; /* IDA 0x1D7D0; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_1D7D2: ; /* IDA 0x1D7D2; mov     dx, 3C4h */
    m->dx=(uint16_t)(964);
L_1D7D5: ; /* IDA 0x1D7D5; mov     ax, 802h */
    m->ax=(uint16_t)(2050);
L_1D7D8: ; /* IDA 0x1D7D8; out     dx, ax; EGA: sequencer address reg */
    dl_out(m,m->dx,m->ax,16);
L_1D7D9: ; /* IDA 0x1D7D9; mov     dx, di */
    m->dx=(uint16_t)(m->di);
L_1D7DB: ; /* IDA 0x1D7DB; mov     cx, 10h */
    m->cx=(uint16_t)(16);
L_1D7DE: ; /* IDA 0x1D7DE; movsw */
    mc_string(m,0,0);
L_1D7DF: ; /* IDA 0x1D7DF; add     di, 4Eh ; 'N' */
    m->di=(uint16_t)(ec_math(m,m->di,78,16,0,0));
L_1D7E2: ; /* IDA 0x1D7E2; loop    loc_1D7DE */
    if(--m->cx) goto L_1D7DE;
L_1D7E4: ; /* IDA 0x1D7E4; mov     di, dx */
    m->di=(uint16_t)(m->dx);
L_1D7E6: ; /* IDA 0x1D7E6; retn */
    return;
}
/* Original sub_1D7E7; IDA 0x1D7E7. */
static void mc_body_1D7E7(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D7E7: ; /* IDA 0x1D7E7; push    cx */
    ec_push(m,m->cx);
L_1D7E8: ; /* IDA 0x1D7E8; push    si */
    ec_push(m,m->si);
L_1D7E9: ; /* IDA 0x1D7E9; push    di */
    ec_push(m,m->di);
L_1D7EA: ; /* IDA 0x1D7EA; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1D7EC: ; /* IDA 0x1D7EC; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1D7EE: ; /* IDA 0x1D7EE; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1D7F0: ; /* IDA 0x1D7F0; mov     ax, es */
    m->ax=(uint16_t)(m->es);
L_1D7F2: ; /* IDA 0x1D7F2; mov     di, cs */
    m->di=(uint16_t)(m->cs);
L_1D7F4: ; /* IDA 0x1D7F4; mov     es, di */
    m->es=(uint16_t)(m->di);
L_1D7F6: ; /* IDA 0x1D7F6; mov     di, 0D858h */
    m->di=(uint16_t)(55384);
L_1D7F9: ; /* IDA 0x1D7F9; mov     cx, 40h ; '@' */
    m->cx=(uint16_t)(64);
L_1D7FC: ; /* IDA 0x1D7FC; rep movsw */
#if KI_MAPCELLS_MUTATION == 5
    mc_string(m,0,0); /* lost REP */
#else
    mc_string(m,0,1);
#endif
L_1D7FE: ; /* IDA 0x1D7FE; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1D800: ; /* IDA 0x1D800; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1D801: ; /* IDA 0x1D801; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1D802: ; /* IDA 0x1D802; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1D803: ; /* IDA 0x1D803; retn */
    return;
}
/* Original sub_1D804; IDA 0x1D804. */
static void mc_body_1D804(KiMachine16 *m,const KiEconomyHooks *h) {
L_1D804: ; /* IDA 0x1D804; push    cx */
    ec_push(m,m->cx);
L_1D805: ; /* IDA 0x1D805; push    si */
    ec_push(m,m->si);
L_1D806: ; /* IDA 0x1D806; push    di */
    ec_push(m,m->di);
L_1D807: ; /* IDA 0x1D807; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1D809: ; /* IDA 0x1D809; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1D80B: ; /* IDA 0x1D80B; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1D80D: ; /* IDA 0x1D80D; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1D80F: ; /* IDA 0x1D80F; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1D811: ; /* IDA 0x1D811; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_1D813: ; /* IDA 0x1D813; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_1D815: ; /* IDA 0x1D815; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1D817: ; /* IDA 0x1D817; mov     di, 0D858h */
    m->di=(uint16_t)(55384);
L_1D81A: ; /* IDA 0x1D81A; mov     cx, 10h */
#if KI_MAPCELLS_MUTATION == 7
    m->cx=(uint16_t)(15);
#else
    m->cx=(uint16_t)(16);
#endif
L_1D81D: ; /* IDA 0x1D81D; lodsw */
    mc_string(m,1,0);
L_1D81E: ; /* IDA 0x1D81E; not     ax */
#if KI_MAPCELLS_MUTATION == 6
    ; /* lost mask inversion */
#else
    m->ax=(uint16_t)(~(m->ax));
#endif
L_1D820: ; /* IDA 0x1D820; and     es:[di], ax */
    {uint16_t result=((dl_read16(m,m->es,(uint16_t)(m->di)))&(m->ax))&65535;ec_logic(m,result,16);dl_write16(m,m->es,(uint16_t)(m->di),result);}
L_1D823: ; /* IDA 0x1D823; and     es:[di+20h], ax */
    {uint16_t result=((dl_read16(m,m->es,(uint16_t)(m->di+32)))&(m->ax))&65535;ec_logic(m,result,16);dl_write16(m,m->es,(uint16_t)(m->di+32),result);}
L_1D827: ; /* IDA 0x1D827; and     es:[di+40h], ax */
    {uint16_t result=((dl_read16(m,m->es,(uint16_t)(m->di+64)))&(m->ax))&65535;ec_logic(m,result,16);dl_write16(m,m->es,(uint16_t)(m->di+64),result);}
L_1D82B: ; /* IDA 0x1D82B; and     es:[di+60h], ax */
    {uint16_t result=((dl_read16(m,m->es,(uint16_t)(m->di+96)))&(m->ax))&65535;ec_logic(m,result,16);dl_write16(m,m->es,(uint16_t)(m->di+96),result);}
L_1D82F: ; /* IDA 0x1D82F; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1D830: ; /* IDA 0x1D830; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1D831: ; /* IDA 0x1D831; loop    loc_1D81D */
    if(--m->cx) goto L_1D81D;
L_1D833: ; /* IDA 0x1D833; mov     di, 0D858h */
    m->di=(uint16_t)(55384);
L_1D836: ; /* IDA 0x1D836; mov     cx, 40h ; '@' */
    m->cx=(uint16_t)(64);
L_1D839: ; /* IDA 0x1D839; lodsw */
    mc_string(m,1,0);
L_1D83A: ; /* IDA 0x1D83A; or      es:[di], ax */
    {uint16_t result=((dl_read16(m,m->es,(uint16_t)(m->di)))|(m->ax))&65535;ec_logic(m,result,16);dl_write16(m,m->es,(uint16_t)(m->di),result);}
L_1D83D: ; /* IDA 0x1D83D; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1D83E: ; /* IDA 0x1D83E; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1D83F: ; /* IDA 0x1D83F; loop    loc_1D839 */
    if(--m->cx) goto L_1D839;
L_1D841: ; /* IDA 0x1D841; mov     ax, 0A0C8h */
    m->ax=(uint16_t)(41160);
L_1D844: ; /* IDA 0x1D844; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_1D846: ; /* IDA 0x1D846; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1D847: ; /* IDA 0x1D847; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1D848: ; /* IDA 0x1D848; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1D849: ; /* IDA 0x1D849; retn */
    return;
}
int ki_mapcells_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0xd46a:mc_body_1D46A(m,h);break;
case 0xd483:mc_body_1D483(m,h);break;
case 0xd4c7:mc_body_1D4C7(m,h);break;
case 0xd615:mc_body_1D615(m,h);break;
case 0xd66a:mc_body_1D66A(m,h);break;
case 0xd76b:mc_body_1D76B(m,h);break;
case 0xd782:mc_body_1D782(m,h);break;
case 0xd796:mc_body_1D796(m,h);break;
case 0xd7e7:mc_body_1D7E7(m,h);break;
case 0xd804:mc_body_1D804(m,h);break;
default:return 0;}return 1;}
