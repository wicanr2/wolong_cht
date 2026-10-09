/* Original native C list families; re/119, spec/239. */
/* Original sub_1716D; IDA 0x1716D. */
static void cc_body_1716D(KiMachine16 *m,const KiEconomyHooks *h) {
L_1716D: ; /* IDA 0x1716D; push    ds */
    ec_push(m,m->ds);
L_1716E: ; /* IDA 0x1716E; push    es */
    ec_push(m,m->es);
L_1716F: ; /* IDA 0x1716F; push    ax */
    ec_push(m,m->ax);
L_17170: ; /* IDA 0x17170; push    cx */
    ec_push(m,m->cx);
L_17171: ; /* IDA 0x17171; push    dx */
    ec_push(m,m->dx);
L_17172: ; /* IDA 0x17172; push    si */
    ec_push(m,m->si);
L_17173: ; /* IDA 0x17173; push    di */
    ec_push(m,m->di);
L_17174: ; /* IDA 0x17174; push    bp */
    ec_push(m,m->bp);
L_17175: ; /* IDA 0x17175; mov     dx, 18h */
    m->dx=(uint16_t)(24);
L_17178: ; /* IDA 0x17178; mov     bx, 58h ; 'X' */
    m->bx=(uint16_t)(88);
L_1717B: ; /* IDA 0x1717B; mov     cx, 0B18h */
    m->cx=(uint16_t)(2840);
L_1717E: ; /* IDA 0x1717E; call    sub_181C0 */
    cc_call(m,33216,29057,h);
L_17181: ; /* IDA 0x17181; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_17186: ; /* IDA 0x17186; mov     ax, 71A8h */
    m->ax=(uint16_t)(29096);
L_17189: ; /* IDA 0x17189; mov     bx, 724Dh */
    m->bx=(uint16_t)(29261);
L_1718C: ; /* IDA 0x1718C; mov     si, 70ECh */
    m->si=(uint16_t)(28908);
L_1718F: ; /* IDA 0x1718F; mov     di, 713Dh */
    m->di=(uint16_t)(28989);
L_17192: ; /* IDA 0x17192; mov     cl, byte ptr cs:word_198A8 */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,39080)));
L_17197: ; /* IDA 0x17197; call    sub_1820E */
    cc_call(m,33294,29082,h);
L_1719A: ; /* IDA 0x1719A; mov     byte ptr cs:word_198A8, cl */
    dl_write8(m,m->cs,39080,(uint8_t)m->cx);
L_1719F: ; /* IDA 0x1719F; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_171A0: ; /* IDA 0x171A0; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_171A1: ; /* IDA 0x171A1; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_171A2: ; /* IDA 0x171A2; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_171A3: ; /* IDA 0x171A3; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_171A4: ; /* IDA 0x171A4; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_171A5: ; /* IDA 0x171A5; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_171A6: ; /* IDA 0x171A6; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_171A7: ; /* IDA 0x171A7; retn */
    return;
}
/* Original sub_171D3; IDA 0x171D3. */
static void cc_body_171D3(KiMachine16 *m,const KiEconomyHooks *h) {
L_171D3: ; /* IDA 0x171D3; push    ds */
    ec_push(m,m->ds);
L_171D4: ; /* IDA 0x171D4; push    es */
    ec_push(m,m->es);
L_171D5: ; /* IDA 0x171D5; push    ax */
    ec_push(m,m->ax);
L_171D6: ; /* IDA 0x171D6; push    cx */
    ec_push(m,m->cx);
L_171D7: ; /* IDA 0x171D7; push    dx */
    ec_push(m,m->dx);
L_171D8: ; /* IDA 0x171D8; push    si */
    ec_push(m,m->si);
L_171D9: ; /* IDA 0x171D9; push    di */
    ec_push(m,m->di);
L_171DA: ; /* IDA 0x171DA; push    bp */
    ec_push(m,m->bp);
L_171DB: ; /* IDA 0x171DB; mov     word ptr cs:loc_1722A+1, ax */
    dl_write16(m,m->cs,29227,m->ax);
L_171DF: ; /* IDA 0x171DF; mov     word ptr cs:loc_17232+1, bx */
    dl_write16(m,m->cs,29235,m->bx);
L_171E4: ; /* IDA 0x171E4; mov     dx, 18h */
    m->dx=(uint16_t)(24);
L_171E7: ; /* IDA 0x171E7; mov     bx, 58h ; 'X' */
    m->bx=(uint16_t)(88);
L_171EA: ; /* IDA 0x171EA; mov     cx, 0B18h */
    m->cx=(uint16_t)(2840);
L_171ED: ; /* IDA 0x171ED; call    sub_181C0 */
    cc_call(m,33216,29168,h);
L_171F0: ; /* IDA 0x171F0; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_171F5: ; /* IDA 0x171F5; mov     ax, 7217h */
    m->ax=(uint16_t)(29207);
L_171F8: ; /* IDA 0x171F8; mov     bx, 724Dh */
    m->bx=(uint16_t)(29261);
L_171FB: ; /* IDA 0x171FB; mov     si, 70ECh */
    m->si=(uint16_t)(28908);
L_171FE: ; /* IDA 0x171FE; mov     di, 713Dh */
    m->di=(uint16_t)(28989);
L_17201: ; /* IDA 0x17201; mov     cl, byte ptr cs:word_198A8 */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,39080)));
L_17206: ; /* IDA 0x17206; call    sub_1820E */
    cc_call(m,33294,29193,h);
L_17209: ; /* IDA 0x17209; mov     byte ptr cs:word_198A8, cl */
    dl_write8(m,m->cs,39080,(uint8_t)m->cx);
L_1720E: ; /* IDA 0x1720E; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1720F: ; /* IDA 0x1720F; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17210: ; /* IDA 0x17210; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17211: ; /* IDA 0x17211; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17212: ; /* IDA 0x17212; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17213: ; /* IDA 0x17213; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17214: ; /* IDA 0x17214; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_17215: ; /* IDA 0x17215; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17216: ; /* IDA 0x17216; retn */
    return;
}
/* Original sub_175FA; IDA 0x175FA. */
static void cc_body_175FA(KiMachine16 *m,const KiEconomyHooks *h) {
L_175FA: ; /* IDA 0x175FA; push    ds */
    ec_push(m,m->ds);
L_175FB: ; /* IDA 0x175FB; push    es */
    ec_push(m,m->es);
L_175FC: ; /* IDA 0x175FC; push    ax */
    ec_push(m,m->ax);
L_175FD: ; /* IDA 0x175FD; push    cx */
    ec_push(m,m->cx);
L_175FE: ; /* IDA 0x175FE; push    dx */
    ec_push(m,m->dx);
L_175FF: ; /* IDA 0x175FF; push    si */
    ec_push(m,m->si);
L_17600: ; /* IDA 0x17600; push    di */
    ec_push(m,m->di);
L_17601: ; /* IDA 0x17601; push    bp */
    ec_push(m,m->bp);
L_17602: ; /* IDA 0x17602; mov     al, ds:0CFFh */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,3327)));
L_17605: ; /* IDA 0x17605; mov     byte ptr cs:loc_17643+1, al */
    dl_write8(m,m->cs,30276,(uint8_t)m->ax);
L_17609: ; /* IDA 0x17609; mov     dx, 18h */
    m->dx=(uint16_t)(24);
L_1760C: ; /* IDA 0x1760C; mov     bx, 58h ; 'X' */
    m->bx=(uint16_t)(88);
L_1760F: ; /* IDA 0x1760F; mov     cx, 0B18h */
    m->cx=(uint16_t)(2840);
L_17612: ; /* IDA 0x17612; call    sub_181C0 */
    cc_call(m,33216,30229,h);
L_17615: ; /* IDA 0x17615; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_1761A: ; /* IDA 0x1761A; mov     ax, 763Ch */
    m->ax=(uint16_t)(30268);
L_1761D: ; /* IDA 0x1761D; mov     bx, 76DCh */
    m->bx=(uint16_t)(30428);
L_17620: ; /* IDA 0x17620; mov     si, 7550h */
    m->si=(uint16_t)(30032);
L_17623: ; /* IDA 0x17623; mov     di, 75C8h */
    m->di=(uint16_t)(30152);
L_17626: ; /* IDA 0x17626; mov     cl, byte ptr cs:word_198AA */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,39082)));
L_1762B: ; /* IDA 0x1762B; call    sub_1820E */
    cc_call(m,33294,30254,h);
L_1762E: ; /* IDA 0x1762E; mov     byte ptr cs:word_198AA, cl */
    dl_write8(m,m->cs,39082,(uint8_t)m->cx);
L_17633: ; /* IDA 0x17633; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17634: ; /* IDA 0x17634; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17635: ; /* IDA 0x17635; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17636: ; /* IDA 0x17636; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17637: ; /* IDA 0x17637; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17638: ; /* IDA 0x17638; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17639: ; /* IDA 0x17639; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1763A: ; /* IDA 0x1763A; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1763B: ; /* IDA 0x1763B; retn */
    return;
}
/* Original sub_17663; IDA 0x17663. */
static void cc_body_17663(KiMachine16 *m,const KiEconomyHooks *h) {
L_17663: ; /* IDA 0x17663; push    ds */
    ec_push(m,m->ds);
L_17664: ; /* IDA 0x17664; push    es */
    ec_push(m,m->es);
L_17665: ; /* IDA 0x17665; push    ax */
    ec_push(m,m->ax);
L_17666: ; /* IDA 0x17666; push    cx */
    ec_push(m,m->cx);
L_17667: ; /* IDA 0x17667; push    dx */
    ec_push(m,m->dx);
L_17668: ; /* IDA 0x17668; push    si */
    ec_push(m,m->si);
L_17669: ; /* IDA 0x17669; push    di */
    ec_push(m,m->di);
L_1766A: ; /* IDA 0x1766A; push    bp */
    ec_push(m,m->bp);
L_1766B: ; /* IDA 0x1766B; mov     dx, 18h */
    m->dx=(uint16_t)(24);
L_1766E: ; /* IDA 0x1766E; mov     bx, 58h ; 'X' */
    m->bx=(uint16_t)(88);
L_17671: ; /* IDA 0x17671; mov     cx, 0B18h */
    m->cx=(uint16_t)(2840);
L_17674: ; /* IDA 0x17674; call    sub_181C0 */
    cc_call(m,33216,30327,h);
L_17677: ; /* IDA 0x17677; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_1767C: ; /* IDA 0x1767C; mov     ax, 76A0h */
    m->ax=(uint16_t)(30368);
L_1767F: ; /* IDA 0x1767F; mov     bx, 76DCh */
    m->bx=(uint16_t)(30428);
L_17682: ; /* IDA 0x17682; mov     si, 7550h */
    m->si=(uint16_t)(30032);
L_17685: ; /* IDA 0x17685; mov     di, 75C8h */
    m->di=(uint16_t)(30152);
L_17688: ; /* IDA 0x17688; xor     cl, cl */
    {uint16_t result=(((uint8_t)m->cx)^((uint8_t)m->cx))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_1768A: ; /* IDA 0x1768A; mov     cl, byte ptr cs:word_198AA */
#if KI_CATALOG_MUTATION == 5
    vg_cl(m,0);
#else
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,39082)));
#endif
L_1768F: ; /* IDA 0x1768F; call    sub_1820E */
    cc_call(m,33294,30354,h);
L_17692: ; /* IDA 0x17692; mov     byte ptr cs:word_198AA, cl */
    dl_write8(m,m->cs,39082,(uint8_t)m->cx);
L_17697: ; /* IDA 0x17697; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17698: ; /* IDA 0x17698; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17699: ; /* IDA 0x17699; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1769A: ; /* IDA 0x1769A; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1769B: ; /* IDA 0x1769B; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1769C: ; /* IDA 0x1769C; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1769D: ; /* IDA 0x1769D; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1769E: ; /* IDA 0x1769E; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1769F: ; /* IDA 0x1769F; retn */
    return;
}
/* Original sub_1770C; IDA 0x1770C. */
static void cc_body_1770C(KiMachine16 *m,const KiEconomyHooks *h) {
L_1770C: ; /* IDA 0x1770C; push    ax */
    ec_push(m,m->ax);
L_1770D: ; /* IDA 0x1770D; push    bx */
    ec_push(m,m->bx);
L_1770E: ; /* IDA 0x1770E; push    cx */
    ec_push(m,m->cx);
L_1770F: ; /* IDA 0x1770F; push    dx */
    ec_push(m,m->dx);
L_17710: ; /* IDA 0x17710; push    si */
    ec_push(m,m->si);
L_17711: ; /* IDA 0x17711; push    di */
    ec_push(m,m->di);
L_17712: ; /* IDA 0x17712; push    bp */
    ec_push(m,m->bp);
L_17713: ; /* IDA 0x17713; cmp     word ptr [bp+si], 0FFFFh */
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)),65535,16);
L_17717: ; /* IDA 0x17717; jnz     short loc_1772B */
    if(!(m->flags&0x40)) goto L_1772B;
L_17719: ; /* IDA 0x17719; push    ds */
    ec_push(m,m->ds);
L_1771A: ; /* IDA 0x1771A; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_1771C: ; /* IDA 0x1771C; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1771E: ; /* IDA 0x1771E; mov     si, 7579h */
    m->si=(uint16_t)(30073);
L_17721: ; /* IDA 0x17721; mov     ax, 9000h */
    m->ax=(uint16_t)(36864);
L_17724: ; /* IDA 0x17724; call    sub_106F5 */
    cc_call(m,1781,30503,h);
L_17727: ; /* IDA 0x17727; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17728: ; /* IDA 0x17728; jmp     loc_177E6 */
    goto L_177E6;
L_1772B: ; /* IDA 0x1772B; mov     si, [bp+si] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)));
L_1772D: ; /* IDA 0x1772D; add     si, 2 */
    m->si=(uint16_t)(ec_math(m,m->si,2,16,0,0));
L_17730: ; /* IDA 0x17730; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_17733: ; /* IDA 0x17733; call    loc_10701 */
    cc_call(m,1793,30518,h);
L_17736: ; /* IDA 0x17736; sub     si, 2 */
    m->si=(uint16_t)(ec_math(m,m->si,2,16,0,1));
L_17739: ; /* IDA 0x17739; push    si */
    ec_push(m,m->si);
L_1773A: ; /* IDA 0x1773A; push    dx */
    ec_push(m,m->dx);
L_1773B: ; /* IDA 0x1773B; mov     ah, [si+1Dh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+29))));
L_1773E: ; /* IDA 0x1773E; cmp     ah, 0FFh */
    ec_cmp(m,(uint8_t)(m->ax>>8),255,8);
L_17741: ; /* IDA 0x17741; jnz     short loc_17746 */
#if KI_CATALOG_MUTATION == 7
    ; /* lost affiliation fallback gate */
#else
    if(!(m->flags&0x40)) goto L_17746;
#endif
L_17743: ; /* IDA 0x17743; mov     ah, [si+1Ch] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+28))));
L_17746: ; /* IDA 0x17746; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17748: ; /* IDA 0x17748; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1774A: ; /* IDA 0x1774A; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1774C: ; /* IDA 0x1774C; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1774E: ; /* IDA 0x1774E; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_17751: ; /* IDA 0x17751; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17753: ; /* IDA 0x17753; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17755: ; /* IDA 0x17755; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17757: ; /* IDA 0x17757; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17759: ; /* IDA 0x17759; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_1775C: ; /* IDA 0x1775C; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1775E: ; /* IDA 0x1775E; add     dx, 0D0h */
    m->dx=(uint16_t)(ec_math(m,m->dx,208,16,0,0));
L_17762: ; /* IDA 0x17762; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_17765: ; /* IDA 0x17765; call    loc_10701 */
    cc_call(m,1793,30568,h);
L_17768: ; /* IDA 0x17768; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17769: ; /* IDA 0x17769; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1776A: ; /* IDA 0x1776A; push    ds */
    ec_push(m,m->ds);
L_1776B: ; /* IDA 0x1776B; push    dx */
    ec_push(m,m->dx);
L_1776C: ; /* IDA 0x1776C; push    si */
    ec_push(m,m->si);
L_1776D: ; /* IDA 0x1776D; add     dx, 120h */
    m->dx=(uint16_t)(ec_math(m,m->dx,288,16,0,0));
L_17771: ; /* IDA 0x17771; mov     al, [si+17h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+23))));
L_17774: ; /* IDA 0x17774; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17776: ; /* IDA 0x17776; jnz     short loc_1777F */
#if KI_CATALOG_MUTATION == 8
    ; /* lost duty priority */
#else
    if(!(m->flags&0x40)) goto L_1777F;
#endif
L_17778: ; /* IDA 0x17778; test    byte ptr [si], 40h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(64))&255;ec_logic(m,result,8);}
L_1777B: ; /* IDA 0x1777B; jz      short loc_1777F */
    if(m->flags&0x40) goto L_1777F;
L_1777D: ; /* IDA 0x1777D; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_1777F: ; /* IDA 0x1777F; mov     ah, al */
    vg_ah(m,(uint8_t)((uint8_t)m->ax));
L_17781: ; /* IDA 0x17781; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_17783: ; /* IDA 0x17783; add     al, ah */
    vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,(uint8_t)(m->ax>>8),8,0,0)));
L_17785: ; /* IDA 0x17785; shl     al, 1 */
    vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,4)));
L_17787: ; /* IDA 0x17787; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17789: ; /* IDA 0x17789; add     ax, 75A4h */
    m->ax=(uint16_t)(ec_math(m,m->ax,30116,16,0,0));
L_1778C: ; /* IDA 0x1778C; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1778E: ; /* IDA 0x1778E; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_17790: ; /* IDA 0x17790; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_17792: ; /* IDA 0x17792; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_17795: ; /* IDA 0x17795; call    loc_10701 */
    cc_call(m,1793,30616,h);
L_17798: ; /* IDA 0x17798; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17799: ; /* IDA 0x17799; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1779A: ; /* IDA 0x1779A; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1779B: ; /* IDA 0x1779B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1779D: ; /* IDA 0x1779D; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1779F: ; /* IDA 0x1779F; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_177A1: ; /* IDA 0x177A1; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_177A3: ; /* IDA 0x177A3; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_177A5: ; /* IDA 0x177A5; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_177A7: ; /* IDA 0x177A7; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_177A9: ; /* IDA 0x177A9; add     di, bx */
    m->di=(uint16_t)(ec_math(m,m->di,m->bx,16,0,0));
L_177AB: ; /* IDA 0x177AB; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_177AD: ; /* IDA 0x177AD; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_177AF: ; /* IDA 0x177AF; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_177B1: ; /* IDA 0x177B1; add     di, dx */
    m->di=(uint16_t)(ec_math(m,m->di,m->dx,16,0,0));
L_177B3: ; /* IDA 0x177B3; add     di, 8 */
    m->di=(uint16_t)(ec_math(m,m->di,8,16,0,0));
L_177B6: ; /* IDA 0x177B6; mov     al, [si+11h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+17))));
L_177B9: ; /* IDA 0x177B9; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_177BB: ; /* IDA 0x177BB; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_177BD: ; /* IDA 0x177BD; mov     bx, 9002h */
    m->bx=(uint16_t)(36866);
L_177C0: ; /* IDA 0x177C0; call    sub_1062F */
    cc_call(m,1583,30659,h);
L_177C3: ; /* IDA 0x177C3; add     di, 5 */
    m->di=(uint16_t)(ec_math(m,m->di,5,16,0,0));
L_177C6: ; /* IDA 0x177C6; mov     al, [si+12h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+18))));
L_177C9: ; /* IDA 0x177C9; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_177CB: ; /* IDA 0x177CB; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_177CD: ; /* IDA 0x177CD; mov     bx, 9002h */
    m->bx=(uint16_t)(36866);
L_177D0: ; /* IDA 0x177D0; call    sub_1062F */
    cc_call(m,1583,30675,h);
L_177D3: ; /* IDA 0x177D3; add     di, 5 */
    m->di=(uint16_t)(ec_math(m,m->di,5,16,0,0));
L_177D6: ; /* IDA 0x177D6; mov     al, [si+13h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+19))));
L_177D9: ; /* IDA 0x177D9; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_177DB: ; /* IDA 0x177DB; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_177DD: ; /* IDA 0x177DD; mov     bx, 9002h */
    m->bx=(uint16_t)(36866);
L_177E0: ; /* IDA 0x177E0; call    sub_1062F */
    cc_call(m,1583,30691,h);
L_177E3: ; /* IDA 0x177E3; add     di, 7 */
    m->di=(uint16_t)(ec_math(m,m->di,7,16,0,0));
L_177E6: ; /* IDA 0x177E6; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_177E7: ; /* IDA 0x177E7; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_177E8: ; /* IDA 0x177E8; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_177E9: ; /* IDA 0x177E9; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_177EA: ; /* IDA 0x177EA; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_177EB: ; /* IDA 0x177EB; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_177EC: ; /* IDA 0x177EC; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_177ED: ; /* IDA 0x177ED; retn */
    return;
}
/* Original sub_178A7; IDA 0x178A7. */
static void cc_body_178A7(KiMachine16 *m,const KiEconomyHooks *h) {
L_178A7: ; /* IDA 0x178A7; push    ds */
    ec_push(m,m->ds);
L_178A8: ; /* IDA 0x178A8; push    es */
    ec_push(m,m->es);
L_178A9: ; /* IDA 0x178A9; push    ax */
    ec_push(m,m->ax);
L_178AA: ; /* IDA 0x178AA; push    cx */
    ec_push(m,m->cx);
L_178AB: ; /* IDA 0x178AB; push    dx */
    ec_push(m,m->dx);
L_178AC: ; /* IDA 0x178AC; push    si */
    ec_push(m,m->si);
L_178AD: ; /* IDA 0x178AD; push    di */
    ec_push(m,m->di);
L_178AE: ; /* IDA 0x178AE; push    bp */
    ec_push(m,m->bp);
L_178AF: ; /* IDA 0x178AF; call    sub_17A7A */
    cc_call(m,31354,30898,h);
L_178B2: ; /* IDA 0x178B2; mov     dx, 18h */
    m->dx=(uint16_t)(24);
L_178B5: ; /* IDA 0x178B5; mov     bx, 58h ; 'X' */
    m->bx=(uint16_t)(88);
L_178B8: ; /* IDA 0x178B8; mov     cx, 0B18h */
    m->cx=(uint16_t)(2840);
L_178BB: ; /* IDA 0x178BB; call    sub_181C0 */
    cc_call(m,33216,30910,h);
L_178BE: ; /* IDA 0x178BE; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_178C3: ; /* IDA 0x178C3; mov     ax, 78E5h */
    m->ax=(uint16_t)(30949);
L_178C6: ; /* IDA 0x178C6; mov     bx, 796Ch */
    m->bx=(uint16_t)(31084);
L_178C9: ; /* IDA 0x178C9; mov     si, 77EEh */
    m->si=(uint16_t)(30702);
L_178CC: ; /* IDA 0x178CC; mov     di, 7875h */
    m->di=(uint16_t)(30837);
L_178CF: ; /* IDA 0x178CF; mov     cl, byte ptr cs:word_198AA+1 */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,39083)));
L_178D4: ; /* IDA 0x178D4; call    sub_1820E */
    cc_call(m,33294,30935,h);
L_178D7: ; /* IDA 0x178D7; mov     byte ptr cs:word_198AA+1, cl */
    dl_write8(m,m->cs,39083,(uint8_t)m->cx);
L_178DC: ; /* IDA 0x178DC; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_178DD: ; /* IDA 0x178DD; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_178DE: ; /* IDA 0x178DE; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_178DF: ; /* IDA 0x178DF; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_178E0: ; /* IDA 0x178E0; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_178E1: ; /* IDA 0x178E1; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_178E2: ; /* IDA 0x178E2; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_178E3: ; /* IDA 0x178E3; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_178E4: ; /* IDA 0x178E4; retn */
    return;
}
/* Original sub_17906; IDA 0x17906. */
static void cc_body_17906(KiMachine16 *m,const KiEconomyHooks *h) {
L_17906: ; /* IDA 0x17906; push    ds */
    ec_push(m,m->ds);
L_17907: ; /* IDA 0x17907; push    es */
    ec_push(m,m->es);
L_17908: ; /* IDA 0x17908; push    ax */
    ec_push(m,m->ax);
L_17909: ; /* IDA 0x17909; push    cx */
    ec_push(m,m->cx);
L_1790A: ; /* IDA 0x1790A; push    dx */
    ec_push(m,m->dx);
L_1790B: ; /* IDA 0x1790B; push    si */
    ec_push(m,m->si);
L_1790C: ; /* IDA 0x1790C; push    di */
    ec_push(m,m->di);
L_1790D: ; /* IDA 0x1790D; push    bp */
    ec_push(m,m->bp);
L_1790E: ; /* IDA 0x1790E; call    sub_17A7A */
    cc_call(m,31354,30993,h);
L_17911: ; /* IDA 0x17911; mov     dx, 18h */
    m->dx=(uint16_t)(24);
L_17914: ; /* IDA 0x17914; mov     bx, 58h ; 'X' */
    m->bx=(uint16_t)(88);
L_17917: ; /* IDA 0x17917; mov     cx, 0B18h */
    m->cx=(uint16_t)(2840);
L_1791A: ; /* IDA 0x1791A; call    sub_181C0 */
    cc_call(m,33216,31005,h);
L_1791D: ; /* IDA 0x1791D; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_17922: ; /* IDA 0x17922; mov     ax, 7944h */
    m->ax=(uint16_t)(31044);
L_17925: ; /* IDA 0x17925; mov     bx, 796Ch */
    m->bx=(uint16_t)(31084);
L_17928: ; /* IDA 0x17928; mov     si, 77EEh */
    m->si=(uint16_t)(30702);
L_1792B: ; /* IDA 0x1792B; mov     di, 7875h */
    m->di=(uint16_t)(30837);
L_1792E: ; /* IDA 0x1792E; mov     cl, byte ptr cs:word_198AA+1 */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,39083)));
L_17933: ; /* IDA 0x17933; call    sub_1820E */
    cc_call(m,33294,31030,h);
L_17936: ; /* IDA 0x17936; mov     byte ptr cs:word_198AA+1, cl */
    dl_write8(m,m->cs,39083,(uint8_t)m->cx);
L_1793B: ; /* IDA 0x1793B; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1793C: ; /* IDA 0x1793C; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1793D: ; /* IDA 0x1793D; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1793E: ; /* IDA 0x1793E; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1793F: ; /* IDA 0x1793F; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17940: ; /* IDA 0x17940; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17941: ; /* IDA 0x17941; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_17942: ; /* IDA 0x17942; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17943: ; /* IDA 0x17943; retn */
    return;
}
/* Original sub_1799C; IDA 0x1799C. */
static void cc_body_1799C(KiMachine16 *m,const KiEconomyHooks *h) {
L_1799C: ; /* IDA 0x1799C; push    ax */
    ec_push(m,m->ax);
L_1799D: ; /* IDA 0x1799D; push    bx */
    ec_push(m,m->bx);
L_1799E: ; /* IDA 0x1799E; push    cx */
    ec_push(m,m->cx);
L_1799F: ; /* IDA 0x1799F; push    dx */
    ec_push(m,m->dx);
L_179A0: ; /* IDA 0x179A0; push    si */
    ec_push(m,m->si);
L_179A1: ; /* IDA 0x179A1; push    di */
    ec_push(m,m->di);
L_179A2: ; /* IDA 0x179A2; push    bp */
    ec_push(m,m->bp);
L_179A3: ; /* IDA 0x179A3; cmp     word ptr [bp+si], 0FFFFh */
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)),65535,16);
L_179A7: ; /* IDA 0x179A7; jnz     short loc_179BB */
    if(!(m->flags&0x40)) goto L_179BB;
L_179A9: ; /* IDA 0x179A9; push    ds */
    ec_push(m,m->ds);
L_179AA: ; /* IDA 0x179AA; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_179AC: ; /* IDA 0x179AC; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_179AE: ; /* IDA 0x179AE; mov     si, 7819h */
    m->si=(uint16_t)(30745);
L_179B1: ; /* IDA 0x179B1; mov     ax, 9000h */
    m->ax=(uint16_t)(36864);
L_179B4: ; /* IDA 0x179B4; call    sub_1F6DC */
    cc_call(m,63196,31159,h);
L_179B7: ; /* IDA 0x179B7; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_179B8: ; /* IDA 0x179B8; jmp     loc_17A72 */
    goto L_17A72;
L_179BB: ; /* IDA 0x179BB; push    dx */
    ec_push(m,m->dx);
L_179BC: ; /* IDA 0x179BC; mov     si, [bp+si] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)));
L_179BE: ; /* IDA 0x179BE; mov     di, si */
    m->di=(uint16_t)(m->si);
L_179C0: ; /* IDA 0x179C0; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_179C3: ; /* IDA 0x179C3; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_179C5: ; /* IDA 0x179C5; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_179C7: ; /* IDA 0x179C7; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_179C9: ; /* IDA 0x179C9; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_179CB: ; /* IDA 0x179CB; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_179CE: ; /* IDA 0x179CE; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_179D0: ; /* IDA 0x179D0; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_179D3: ; /* IDA 0x179D3; call    loc_10701 */
    cc_call(m,1793,31190,h);
L_179D6: ; /* IDA 0x179D6; add     dx, 0A0h */
    m->dx=(uint16_t)(ec_math(m,m->dx,160,16,0,0));
L_179DA: ; /* IDA 0x179DA; mov     ah, [di+3] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+3))));
L_179DD: ; /* IDA 0x179DD; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_179DF: ; /* IDA 0x179DF; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_179E1: ; /* IDA 0x179E1; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_179E3: ; /* IDA 0x179E3; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_179E5: ; /* IDA 0x179E5; add     ax, 842h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2114,16,0,0));
L_179E8: ; /* IDA 0x179E8; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_179EA: ; /* IDA 0x179EA; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_179ED: ; /* IDA 0x179ED; call    loc_10701 */
    cc_call(m,1793,31216,h);
L_179F0: ; /* IDA 0x179F0; add     dx, 40h ; '@' */
    m->dx=(uint16_t)(ec_math(m,m->dx,64,16,0,0));
L_179F3: ; /* IDA 0x179F3; mov     al, [di+3Fh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+63))));
L_179F6: ; /* IDA 0x179F6; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_179F8: ; /* IDA 0x179F8; add     ax, 7844h */
    m->ax=(uint16_t)(ec_math(m,m->ax,30788,16,0,0));
L_179FB: ; /* IDA 0x179FB; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_179FD: ; /* IDA 0x179FD; push    ds */
    ec_push(m,m->ds);
L_179FE: ; /* IDA 0x179FE; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_17A00: ; /* IDA 0x17A00; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_17A02: ; /* IDA 0x17A02; mov     ah, [si+6] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+6))));
L_17A05: ; /* IDA 0x17A05; or      ah, 90h */
    {uint16_t result=(((uint8_t)(m->ax>>8))|(144))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17A08: ; /* IDA 0x17A08; mov     al, 1 */
    vg_al(m,(uint8_t)(1));
L_17A0A: ; /* IDA 0x17A0A; call    loc_10701 */
    cc_call(m,1793,31245,h);
L_17A0D: ; /* IDA 0x17A0D; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17A0E: ; /* IDA 0x17A0E; push    ds */
    ec_push(m,m->ds);
L_17A0F: ; /* IDA 0x17A0F; mov     ah, [di+2Ah] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+42))));
L_17A12: ; /* IDA 0x17A12; cmp     ah, 0FFh */
    ec_cmp(m,(uint8_t)(m->ax>>8),255,8);
L_17A15: ; /* IDA 0x17A15; jz      short loc_17A26 */
    if(m->flags&0x40) goto L_17A26;
L_17A17: ; /* IDA 0x17A17; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17A19: ; /* IDA 0x17A19; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17A1B: ; /* IDA 0x17A1B; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17A1D: ; /* IDA 0x17A1D; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17A1F: ; /* IDA 0x17A1F; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_17A22: ; /* IDA 0x17A22; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17A24: ; /* IDA 0x17A24; jmp     short loc_17A2D */
    goto L_17A2D;
L_17A26: ; /* IDA 0x17A26; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_17A28: ; /* IDA 0x17A28; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_17A2A: ; /* IDA 0x17A2A; mov     si, 7819h */
    m->si=(uint16_t)(30745);
L_17A2D: ; /* IDA 0x17A2D; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_17A30: ; /* IDA 0x17A30; add     dx, 40h ; '@' */
    m->dx=(uint16_t)(ec_math(m,m->dx,64,16,0,0));
L_17A33: ; /* IDA 0x17A33; call    loc_10701 */
    cc_call(m,1793,31286,h);
L_17A36: ; /* IDA 0x17A36; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17A37: ; /* IDA 0x17A37; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17A38: ; /* IDA 0x17A38; mov     si, di */
    m->si=(uint16_t)(m->di);
L_17A3A: ; /* IDA 0x17A3A; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17A3C: ; /* IDA 0x17A3C; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17A3E: ; /* IDA 0x17A3E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17A40: ; /* IDA 0x17A40; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17A42: ; /* IDA 0x17A42; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_17A44: ; /* IDA 0x17A44; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17A46: ; /* IDA 0x17A46; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17A48: ; /* IDA 0x17A48; add     di, bx */
    m->di=(uint16_t)(ec_math(m,m->di,m->bx,16,0,0));
L_17A4A: ; /* IDA 0x17A4A; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_17A4C: ; /* IDA 0x17A4C; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_17A4E: ; /* IDA 0x17A4E; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_17A50: ; /* IDA 0x17A50; add     di, dx */
    m->di=(uint16_t)(ec_math(m,m->di,m->dx,16,0,0));
L_17A52: ; /* IDA 0x17A52; add     di, 8 */
    m->di=(uint16_t)(ec_math(m,m->di,8,16,0,0));
L_17A55: ; /* IDA 0x17A55; mov     al, [si+18h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+24))));
L_17A58: ; /* IDA 0x17A58; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17A5A: ; /* IDA 0x17A5A; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_17A5C: ; /* IDA 0x17A5C; mov     bx, 9003h */
#if KI_CATALOG_MUTATION == 13
    m->bx=(uint16_t)(36868);
#else
    m->bx=(uint16_t)(36867);
#endif
L_17A5F: ; /* IDA 0x17A5F; call    sub_1062F */
    cc_call(m,1583,31330,h);
L_17A62: ; /* IDA 0x17A62; add     di, 6 */
    m->di=(uint16_t)(ec_math(m,m->di,6,16,0,0));
L_17A65: ; /* IDA 0x17A65; mov     al, [si+23h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+35))));
L_17A68: ; /* IDA 0x17A68; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17A6A: ; /* IDA 0x17A6A; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_17A6C: ; /* IDA 0x17A6C; mov     bx, 9003h */
    m->bx=(uint16_t)(36867);
L_17A6F: ; /* IDA 0x17A6F; call    sub_1062F */
    cc_call(m,1583,31346,h);
L_17A72: ; /* IDA 0x17A72; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17A73: ; /* IDA 0x17A73; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17A74: ; /* IDA 0x17A74; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17A75: ; /* IDA 0x17A75; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17A76: ; /* IDA 0x17A76; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17A77: ; /* IDA 0x17A77; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17A78: ; /* IDA 0x17A78; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17A79: ; /* IDA 0x17A79; retn */
    return;
}
/* Original sub_17A7A; IDA 0x17A7A. */
static void cc_body_17A7A(KiMachine16 *m,const KiEconomyHooks *h) {
L_17A7A: ; /* IDA 0x17A7A; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_17A7F: ; /* IDA 0x17A7F; mov     si, cs:word_10CFD */
    m->si=(uint16_t)(dl_read16(m,m->cs,3325));
L_17A84: ; /* IDA 0x17A84; mov     ax, si */
    m->ax=(uint16_t)(m->si);
L_17A86: ; /* IDA 0x17A86; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17A88: ; /* IDA 0x17A88; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17A8A: ; /* IDA 0x17A8A; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17A8C: ; /* IDA 0x17A8C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17A8E: ; /* IDA 0x17A8E; add     si, ax */
    m->si=(uint16_t)(ec_math(m,m->si,m->ax,16,0,0));
L_17A90: ; /* IDA 0x17A90; add     si, 600h */
    m->si=(uint16_t)(ec_math(m,m->si,1536,16,0,0));
L_17A94: ; /* IDA 0x17A94; mov     di, 3Fh ; '?' */
    m->di=(uint16_t)(63);
L_17A97: ; /* IDA 0x17A97; mov     cx, 16h */
    m->cx=(uint16_t)(22);
L_17A9A: ; /* IDA 0x17A9A; xor     dl, dl */
#if KI_CATALOG_MUTATION == 11
    if(m->cx==22) {vg_dl(m,0);ec_logic(m,0,8);} /* erroneous cross-row state */
#else
    {uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
#endif
L_17A9C: ; /* IDA 0x17A9C; mov     al, [si] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si))));
L_17A9E: ; /* IDA 0x17A9E; cmp     al, 80h */
    ec_cmp(m,(uint8_t)m->ax,128,8);
L_17AA0: ; /* IDA 0x17AA0; jb      short loc_17ABC */
    if(m->flags&1) goto L_17ABC;
L_17AA2: ; /* IDA 0x17AA2; and     al, 7Fh */
    {uint16_t result=(((uint8_t)m->ax)&(127))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17AA4: ; /* IDA 0x17AA4; mov     dl, 2Ah ; '*' */
    vg_dl(m,(uint8_t)(42));
L_17AA6: ; /* IDA 0x17AA6; cmp     al, 64h ; 'd' */
    ec_cmp(m,(uint8_t)m->ax,100,8);
L_17AA8: ; /* IDA 0x17AA8; ja      short loc_17ABC */
    if(!(m->flags&(1|0x40))) goto L_17ABC;
L_17AAA: ; /* IDA 0x17AAA; jb      short loc_17AAE */
#if KI_CATALOG_MUTATION == 10
    ; /* lost strict relation boundary */
#else
    if(m->flags&1) goto L_17AAE;
#endif
L_17AAC: ; /* IDA 0x17AAC; dec     al */
    vg_al(m,(uint8_t)(st_dec(m,(uint8_t)m->ax,8)));
L_17AAE: ; /* IDA 0x17AAE; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17AB0: ; /* IDA 0x17AB0; mov     dh, 14h */
    vg_dh(m,(uint8_t)(20));
L_17AB2: ; /* IDA 0x17AB2; div     dh */
    dl_div(m,(uint8_t)(m->dx>>8),8);
L_17AB4: ; /* IDA 0x17AB4; inc     al */
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
L_17AB6: ; /* IDA 0x17AB6; mov     ah, 7 */
    vg_ah(m,(uint8_t)(7));
L_17AB8: ; /* IDA 0x17AB8; mul     ah */
    dl_mul(m,(uint8_t)(m->ax>>8),8);
L_17ABA: ; /* IDA 0x17ABA; mov     dl, al */
    vg_dl(m,(uint8_t)((uint8_t)m->ax));
L_17ABC: ; /* IDA 0x17ABC; mov     [di], dl */
    dl_write8(m,m->ds,(uint16_t)(m->di),(uint8_t)m->dx);
L_17ABE: ; /* IDA 0x17ABE; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17ABF: ; /* IDA 0x17ABF; add     di, 40h ; '@' */
    m->di=(uint16_t)(ec_math(m,m->di,64,16,0,0));
L_17AC2: ; /* IDA 0x17AC2; loop    loc_17A9A */
    if(--m->cx) goto L_17A9A;
L_17AC4: ; /* IDA 0x17AC4; retn */
    return;
}
/* Original sub_17B3C; IDA 0x17B3C. */
static void cc_body_17B3C(KiMachine16 *m,const KiEconomyHooks *h) {
L_17B3C: ; /* IDA 0x17B3C; push    ds */
    ec_push(m,m->ds);
L_17B3D: ; /* IDA 0x17B3D; push    es */
    ec_push(m,m->es);
L_17B3E: ; /* IDA 0x17B3E; push    ax */
    ec_push(m,m->ax);
L_17B3F: ; /* IDA 0x17B3F; push    cx */
    ec_push(m,m->cx);
L_17B40: ; /* IDA 0x17B40; push    dx */
    ec_push(m,m->dx);
L_17B41: ; /* IDA 0x17B41; push    si */
    ec_push(m,m->si);
L_17B42: ; /* IDA 0x17B42; push    di */
    ec_push(m,m->di);
L_17B43: ; /* IDA 0x17B43; push    bp */
    ec_push(m,m->bp);
L_17B44: ; /* IDA 0x17B44; mov     dx, 88h */
#if KI_CATALOG_MUTATION == 12
    m->dx=(uint16_t)(24);
#else
    m->dx=(uint16_t)(136);
#endif
L_17B47: ; /* IDA 0x17B47; mov     bx, 68h ; 'h' */
    m->bx=(uint16_t)(104);
L_17B4A: ; /* IDA 0x17B4A; mov     cx, 0B18h */
    m->cx=(uint16_t)(2840);
L_17B4D: ; /* IDA 0x17B4D; call    sub_181C0 */
    cc_call(m,33216,31568,h);
L_17B50: ; /* IDA 0x17B50; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_17B55: ; /* IDA 0x17B55; mov     ax, 7B6Fh */
    m->ax=(uint16_t)(31599);
L_17B58: ; /* IDA 0x17B58; mov     bx, 7B90h */
    m->bx=(uint16_t)(31632);
L_17B5B: ; /* IDA 0x17B5B; mov     si, 7AC6h */
    m->si=(uint16_t)(31430);
L_17B5E: ; /* IDA 0x17B5E; mov     di, 7B12h */
    m->di=(uint16_t)(31506);
L_17B61: ; /* IDA 0x17B61; xor     cl, cl */
    {uint16_t result=(((uint8_t)m->cx)^((uint8_t)m->cx))&255;ec_logic(m,result,8);vg_cl(m,(uint8_t)(result));}
L_17B63: ; /* IDA 0x17B63; call    sub_1820E */
    cc_call(m,33294,31590,h);
L_17B66: ; /* IDA 0x17B66; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17B67: ; /* IDA 0x17B67; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17B68: ; /* IDA 0x17B68; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17B69: ; /* IDA 0x17B69; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17B6A: ; /* IDA 0x17B6A; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17B6B: ; /* IDA 0x17B6B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17B6C: ; /* IDA 0x17B6C; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_17B6D: ; /* IDA 0x17B6D; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17B6E: ; /* IDA 0x17B6E; retn */
    return;
}
/* Original sub_17BC0; IDA 0x17BC0. */
static void cc_body_17BC0(KiMachine16 *m,const KiEconomyHooks *h) {
L_17BC0: ; /* IDA 0x17BC0; push    ax */
    ec_push(m,m->ax);
L_17BC1: ; /* IDA 0x17BC1; push    bx */
    ec_push(m,m->bx);
L_17BC2: ; /* IDA 0x17BC2; push    cx */
    ec_push(m,m->cx);
L_17BC3: ; /* IDA 0x17BC3; push    dx */
    ec_push(m,m->dx);
L_17BC4: ; /* IDA 0x17BC4; push    si */
    ec_push(m,m->si);
L_17BC5: ; /* IDA 0x17BC5; push    di */
    ec_push(m,m->di);
L_17BC6: ; /* IDA 0x17BC6; push    bp */
    ec_push(m,m->bp);
L_17BC7: ; /* IDA 0x17BC7; cmp     word ptr [bp+si], 0FFFFh */
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)),65535,16);
L_17BCB: ; /* IDA 0x17BCB; jnz     short loc_17BDF */
    if(!(m->flags&0x40)) goto L_17BDF;
L_17BCD: ; /* IDA 0x17BCD; push    ds */
    ec_push(m,m->ds);
L_17BCE: ; /* IDA 0x17BCE; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_17BD0: ; /* IDA 0x17BD0; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_17BD2: ; /* IDA 0x17BD2; mov     si, 7AEBh */
    m->si=(uint16_t)(31467);
L_17BD5: ; /* IDA 0x17BD5; mov     ax, 9000h */
    m->ax=(uint16_t)(36864);
L_17BD8: ; /* IDA 0x17BD8; call    sub_1F6DC */
    cc_call(m,63196,31707,h);
L_17BDB: ; /* IDA 0x17BDB; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17BDC: ; /* IDA 0x17BDC; jmp     loc_17C66 */
    goto L_17C66;
L_17BDF: ; /* IDA 0x17BDF; push    dx */
    ec_push(m,m->dx);
L_17BE0: ; /* IDA 0x17BE0; mov     di, [bp+si] */
    m->di=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)));
L_17BE2: ; /* IDA 0x17BE2; mov     ah, [di+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+1))));
L_17BE5: ; /* IDA 0x17BE5; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17BE7: ; /* IDA 0x17BE7; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17BE9: ; /* IDA 0x17BE9; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17BEB: ; /* IDA 0x17BEB; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17BED: ; /* IDA 0x17BED; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_17BF0: ; /* IDA 0x17BF0; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17BF2: ; /* IDA 0x17BF2; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_17BF5: ; /* IDA 0x17BF5; call    loc_10701 */
    cc_call(m,1793,31736,h);
L_17BF8: ; /* IDA 0x17BF8; add     dx, 40h ; '@' */
    m->dx=(uint16_t)(ec_math(m,m->dx,64,16,0,0));
L_17BFB: ; /* IDA 0x17BFB; mov     ah, [di+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+2))));
L_17BFE: ; /* IDA 0x17BFE; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17C00: ; /* IDA 0x17C00; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17C02: ; /* IDA 0x17C02; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17C04: ; /* IDA 0x17C04; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17C06: ; /* IDA 0x17C06; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_17C09: ; /* IDA 0x17C09; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17C0B: ; /* IDA 0x17C0B; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_17C0E: ; /* IDA 0x17C0E; call    loc_10701 */
    cc_call(m,1793,31761,h);
L_17C11: ; /* IDA 0x17C11; add     dx, 0C0h */
    m->dx=(uint16_t)(ec_math(m,m->dx,192,16,0,0));
L_17C15: ; /* IDA 0x17C15; mov     ah, [di+3] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+3))));
L_17C18: ; /* IDA 0x17C18; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17C1A: ; /* IDA 0x17C1A; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17C1C: ; /* IDA 0x17C1C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17C1E: ; /* IDA 0x17C1E; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_17C20: ; /* IDA 0x17C20; add     ax, 842h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2114,16,0,0));
L_17C23: ; /* IDA 0x17C23; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17C25: ; /* IDA 0x17C25; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_17C28: ; /* IDA 0x17C28; call    loc_10701 */
    cc_call(m,1793,31787,h);
L_17C2B: ; /* IDA 0x17C2B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17C2C: ; /* IDA 0x17C2C; mov     si, di */
    m->si=(uint16_t)(m->di);
L_17C2E: ; /* IDA 0x17C2E; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17C30: ; /* IDA 0x17C30; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17C32: ; /* IDA 0x17C32; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17C34: ; /* IDA 0x17C34; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17C36: ; /* IDA 0x17C36; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_17C38: ; /* IDA 0x17C38; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17C3A: ; /* IDA 0x17C3A; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17C3C: ; /* IDA 0x17C3C; add     di, bx */
    m->di=(uint16_t)(ec_math(m,m->di,m->bx,16,0,0));
L_17C3E: ; /* IDA 0x17C3E; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_17C40: ; /* IDA 0x17C40; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_17C42: ; /* IDA 0x17C42; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_17C44: ; /* IDA 0x17C44; add     di, dx */
    m->di=(uint16_t)(ec_math(m,m->di,m->dx,16,0,0));
L_17C46: ; /* IDA 0x17C46; add     di, 12h */
    m->di=(uint16_t)(ec_math(m,m->di,18,16,0,0));
L_17C49: ; /* IDA 0x17C49; mov     al, [si+18h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+24))));
L_17C4C: ; /* IDA 0x17C4C; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17C4E: ; /* IDA 0x17C4E; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_17C50: ; /* IDA 0x17C50; mov     bx, 9003h */
    m->bx=(uint16_t)(36867);
L_17C53: ; /* IDA 0x17C53; call    sub_1062F */
    cc_call(m,1583,31830,h);
L_17C56: ; /* IDA 0x17C56; add     di, 6 */
    m->di=(uint16_t)(ec_math(m,m->di,6,16,0,0));
L_17C59: ; /* IDA 0x17C59; mov     al, [si+23h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+35))));
L_17C5C: ; /* IDA 0x17C5C; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17C5E: ; /* IDA 0x17C5E; mov     ah, dl */
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
L_17C60: ; /* IDA 0x17C60; mov     bx, 9003h */
    m->bx=(uint16_t)(36867);
L_17C63: ; /* IDA 0x17C63; call    sub_1062F */
    cc_call(m,1583,31846,h);
L_17C66: ; /* IDA 0x17C66; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17C67: ; /* IDA 0x17C67; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17C68: ; /* IDA 0x17C68; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17C69: ; /* IDA 0x17C69; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17C6A: ; /* IDA 0x17C6A; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17C6B: ; /* IDA 0x17C6B; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17C6C: ; /* IDA 0x17C6C; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17C6D: ; /* IDA 0x17C6D; retn */
    return;
}
/* Original unnamed code; IDA 0x171A8. */
static void cc_body_171A8(KiMachine16 *m,const KiEconomyHooks *h) {
L_171A8: ; /* IDA 0x171A8; push    dx */
    ec_push(m,m->dx);
L_171A9: ; /* IDA 0x171A9; push    si */
    ec_push(m,m->si);
L_171AA: ; /* IDA 0x171AA; push    di */
    ec_push(m,m->di);
L_171AB: ; /* IDA 0x171AB; mov     si, 2240h */
    m->si=(uint16_t)(8768);
L_171AE: ; /* IDA 0x171AE; xor     di, di */
    {uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_171B0: ; /* IDA 0x171B0; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_171B2: ; /* IDA 0x171B2; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_171B6: ; /* IDA 0x171B6; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_171B9: ; /* IDA 0x171B9; jb      short loc_171C6 */
    if(m->flags&1) goto L_171C6;
L_171BB: ; /* IDA 0x171BB; cmp     al, [si+1] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_171BE: ; /* IDA 0x171BE; jnz     short loc_171C6 */
    if(!(m->flags&0x40)) goto L_171C6;
L_171C0: ; /* IDA 0x171C0; mov     [bp+di], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->si);
L_171C2: ; /* IDA 0x171C2; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_171C4: ; /* IDA 0x171C4; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_171C5: ; /* IDA 0x171C5; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_171C6: ; /* IDA 0x171C6; add     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
L_171C9: ; /* IDA 0x171C9; cmp     si, 41C0h */
#if KI_CATALOG_MUTATION == 1
    ec_cmp(m,m->si,16896,16);
#else
    ec_cmp(m,m->si,16832,16);
#endif
L_171CD: ; /* IDA 0x171CD; jnz     short loc_171B6 */
    if(!(m->flags&0x40)) goto L_171B6;
L_171CF: ; /* IDA 0x171CF; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_171D0: ; /* IDA 0x171D0; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_171D1: ; /* IDA 0x171D1; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_171D2: ; /* IDA 0x171D2; retn */
    return;
}
/* Original unnamed code; IDA 0x17217. */
static void cc_body_17217(KiMachine16 *m,const KiEconomyHooks *h) {
L_17217: ; /* IDA 0x17217; push    dx */
    ec_push(m,m->dx);
L_17218: ; /* IDA 0x17218; push    si */
    ec_push(m,m->si);
L_17219: ; /* IDA 0x17219; push    di */
    ec_push(m,m->di);
L_1721A: ; /* IDA 0x1721A; mov     si, 2240h */
    m->si=(uint16_t)(8768);
L_1721D: ; /* IDA 0x1721D; xor     di, di */
    {uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_1721F: ; /* IDA 0x1721F; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17221: ; /* IDA 0x17221; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_17225: ; /* IDA 0x17225; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_17228: ; /* IDA 0x17228; jb      short near ptr byte_17235+0Bh */
    if(m->flags&1) goto L_17240;
L_1722A: ; /* IDA 0x1722A; mov     dx, 1234h */
#if KI_CATALOG_MUTATION == 2
    m->dx=0x1234;
#else
    m->dx=dl_read16(m,m->cs,0x722b);
#endif
L_1722D: ; /* IDA 0x1722D; cmp     dx, [si+10h] */
    ec_cmp(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+16)),16);
L_17230: ; /* IDA 0x17230; jnz     short near ptr byte_17235+0Bh */
    if(!(m->flags&0x40)) goto L_17240;
L_17232: ; /* IDA 0x17232; mov     dx, 1234h */
    m->dx=dl_read16(m,m->cs,0x7233);
L_17235: ; /* IDA 0x17235; cmp     dx, [si+12h] */
    ec_cmp(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->si+18)),16);
L_17238: ; /* IDA 0x17238; jnz     short near ptr unk_17240 */
    if(!(m->flags&0x40)) goto L_17240;
L_1723A: ; /* IDA 0x1723A; mov     [bp+di], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->si);
L_1723C: ; /* IDA 0x1723C; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_1723E: ; /* IDA 0x1723E; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1723F: ; /* IDA 0x1723F; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_17240: ; /* IDA 0x17240; add     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
L_17243: ; /* IDA 0x17243; cmp     si, 41C0h */
    ec_cmp(m,m->si,16832,16);
L_17247: ; /* IDA 0x17247; jnz     short loc_17225 */
    if(!(m->flags&0x40)) goto L_17225;
L_17249: ; /* IDA 0x17249; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_1724A: ; /* IDA 0x1724A; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1724B: ; /* IDA 0x1724B; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1724C: ; /* IDA 0x1724C; retn */
    return;
}
/* Original unnamed code; IDA 0x1724D. */
static void cc_body_1724D(KiMachine16 *m,const KiEconomyHooks *h) {
L_1724D: ; /* IDA 0x1724D; push    ax */
    ec_push(m,m->ax);
L_1724E: ; /* IDA 0x1724E; push    bx */
    ec_push(m,m->bx);
L_1724F: ; /* IDA 0x1724F; push    dx */
    ec_push(m,m->dx);
L_17250: ; /* IDA 0x17250; push    si */
    ec_push(m,m->si);
L_17251: ; /* IDA 0x17251; call    sub_18607 */
    cc_call(m,34311,29268,h);
L_17254: ; /* IDA 0x17254; mov     dx, cs:word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33198));
L_17259: ; /* IDA 0x17259; add     dx, 18h */
    m->dx=(uint16_t)(ec_math(m,m->dx,24,16,0,0));
L_1725C: ; /* IDA 0x1725C; mov     bx, cs:word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33200));
L_17261: ; /* IDA 0x17261; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17264: ; /* IDA 0x17264; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17266: ; /* IDA 0x17266; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_17268: ; /* IDA 0x17268; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_1726A: ; /* IDA 0x1726A; mov     ax, 0Ah */
    m->ax=(uint16_t)(10);
L_1726D: ; /* IDA 0x1726D; call    near ptr unk_1727D */
    cc_call(m,29309,29296,h);
L_17270: ; /* IDA 0x17270; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17271: ; /* IDA 0x17271; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17272: ; /* IDA 0x17272; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17275: ; /* IDA 0x17275; dec     ax */
    m->ax=(uint16_t)(st_dec(m,m->ax,16));
L_17276: ; /* IDA 0x17276; jnz     short loc_1726D */
    if(!(m->flags&0x40)) goto L_1726D;
L_17278: ; /* IDA 0x17278; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17279: ; /* IDA 0x17279; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1727A: ; /* IDA 0x1727A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1727B: ; /* IDA 0x1727B; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1727C: ; /* IDA 0x1727C; retn */
    return;
}
/* Original sub_1727D; IDA 0x1727D. */
static void cc_body_1727D(KiMachine16 *m,const KiEconomyHooks *h) {
L_1727D: ; /* IDA 0x1727D; push    ax */
    ec_push(m,m->ax);
L_1727E: ; /* IDA 0x1727E; push    bx */
    ec_push(m,m->bx);
L_1727F: ; /* IDA 0x1727F; push    cx */
    ec_push(m,m->cx);
L_17280: ; /* IDA 0x17280; push    dx */
    ec_push(m,m->dx);
L_17281: ; /* IDA 0x17281; push    si */
    ec_push(m,m->si);
L_17282: ; /* IDA 0x17282; push    di */
    ec_push(m,m->di);
L_17283: ; /* IDA 0x17283; push    bp */
    ec_push(m,m->bp);
L_17284: ; /* IDA 0x17284; cmp     word ptr [bp+si], 0FFFFh */
    ec_cmp(m,dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)),65535,16);
L_17288: ; /* IDA 0x17288; jnz     short near ptr unk_1729C */
    if(!(m->flags&0x40)) goto L_1729C;
L_1728A: ; /* IDA 0x1728A; push    ds */
    ec_push(m,m->ds);
L_1728B: ; /* IDA 0x1728B; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_1728D: ; /* IDA 0x1728D; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1728F: ; /* IDA 0x1728F; mov     si, 7115h */
    m->si=(uint16_t)(28949);
L_17292: ; /* IDA 0x17292; mov     ax, 9000h */
    m->ax=(uint16_t)(36864);
L_17295: ; /* IDA 0x17295; call    sub_1F6DC */
    cc_call(m,63196,29336,h);
L_17298: ; /* IDA 0x17298; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_17299: ; /* IDA 0x17299; jmp     near ptr unk_1736F */
    goto L_1736F;
L_1729C: ; /* IDA 0x1729C; push    dx */
    ec_push(m,m->dx);
L_1729D: ; /* IDA 0x1729D; mov     si, [bp+si] */
    m->si=(uint16_t)(dl_read16(m,m->ss,(uint16_t)(m->bp+m->si)));
L_1729F: ; /* IDA 0x1729F; mov     di, si */
    m->di=(uint16_t)(m->si);
L_172A1: ; /* IDA 0x172A1; mov     ah, [si+2] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+2))));
L_172A4: ; /* IDA 0x172A4; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_172A6: ; /* IDA 0x172A6; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_172A8: ; /* IDA 0x172A8; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_172AA: ; /* IDA 0x172AA; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_172AC: ; /* IDA 0x172AC; add     ax, 4242h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16962,16,0,0));
L_172AF: ; /* IDA 0x172AF; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_172B1: ; /* IDA 0x172B1; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_172B4: ; /* IDA 0x172B4; call    loc_10701 */
    cc_call(m,1793,29367,h);
L_172B7: ; /* IDA 0x172B7; add     dx, 0C0h */
    m->dx=(uint16_t)(ec_math(m,m->dx,192,16,0,0));
L_172BB: ; /* IDA 0x172BB; mov     si, [di+0Eh] */
    m->si=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->di+14)));
L_172BE: ; /* IDA 0x172BE; cmp     si, 800h */
    ec_cmp(m,m->si,2048,16);
L_172C2: ; /* IDA 0x172C2; jnb     short near ptr unk_172D4 */
    if(!(m->flags&1)) goto L_172D4;
L_172C4: ; /* IDA 0x172C4; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_172C6: ; /* IDA 0x172C6; shl     si, 1 */
    m->si=(uint16_t)(dl_shift(m,m->si,1,16,4));
L_172C8: ; /* IDA 0x172C8; add     si, 842h */
    m->si=(uint16_t)(ec_math(m,m->si,2114,16,0,0));
L_172CC: ; /* IDA 0x172CC; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_172CF: ; /* IDA 0x172CF; call    loc_10701 */
    cc_call(m,1793,29394,h);
L_172D2: ; /* IDA 0x172D2; jmp     short near ptr unk_172E3 */
    goto L_172E3;
L_172D4: ; /* IDA 0x172D4; push    ds */
    ec_push(m,m->ds);
L_172D5: ; /* IDA 0x172D5; mov     si, 7167h */
    m->si=(uint16_t)(29031);
L_172D8: ; /* IDA 0x172D8; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_172DA: ; /* IDA 0x172DA; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_172DC: ; /* IDA 0x172DC; mov     ax, 9A01h */
    m->ax=(uint16_t)(39425);
L_172DF: ; /* IDA 0x172DF; call    loc_10701 */
    cc_call(m,1793,29410,h);
L_172E2: ; /* IDA 0x172E2; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_172E3: ; /* IDA 0x172E3; add     dx, 48h ; 'H' */
    m->dx=(uint16_t)(ec_math(m,m->dx,72,16,0,0));
L_172E6: ; /* IDA 0x172E6; mov     ah, [di+20h] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+32))));
L_172E9: ; /* IDA 0x172E9; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_172EB: ; /* IDA 0x172EB; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_172ED: ; /* IDA 0x172ED; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_172EF: ; /* IDA 0x172EF; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_172F1: ; /* IDA 0x172F1; add     ax, 842h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2114,16,0,0));
L_172F4: ; /* IDA 0x172F4; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_172F6: ; /* IDA 0x172F6; mov     ax, 9001h */
    m->ax=(uint16_t)(36865);
L_172F9: ; /* IDA 0x172F9; call    loc_10701 */
    cc_call(m,1793,29436,h);
L_172FC: ; /* IDA 0x172FC; add     dx, 38h ; '8' */
    m->dx=(uint16_t)(ec_math(m,m->dx,56,16,0,0));
L_172FF: ; /* IDA 0x172FF; mov     cx, 0A965h */
    m->cx=(uint16_t)(43365);
L_17302: ; /* IDA 0x17302; mov     si, 0A5F4h */
    m->si=(uint16_t)(42484);
L_17305: ; /* IDA 0x17305; test    byte ptr [di], 4 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->di)))&(4))&255;ec_logic(m,result,8);}
L_17308: ; /* IDA 0x17308; jnz     short near ptr unk_1730F */
    if(!(m->flags&0x40)) goto L_1730F;
L_1730A: ; /* IDA 0x1730A; mov     cx, 0A140h */
    m->cx=(uint16_t)(41280);
L_1730D: ; /* IDA 0x1730D; mov     si, cx */
    m->si=(uint16_t)(m->cx);
L_1730F: ; /* IDA 0x1730F; push    bx */
    ec_push(m,m->bx);
L_17310: ; /* IDA 0x17310; push    dx */
    ec_push(m,m->dx);
L_17311: ; /* IDA 0x17311; mov     ax, 9A01h */
    m->ax=(uint16_t)(39425);
L_17314: ; /* IDA 0x17314; call    loc_1F75E */
    cc_call(m,63326,29463,h);
L_17317: ; /* IDA 0x17317; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17318: ; /* IDA 0x17318; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17319: ; /* IDA 0x17319; add     dx, 10h */
    m->dx=(uint16_t)(ec_math(m,m->dx,16,16,0,0));
L_1731C: ; /* IDA 0x1731C; mov     cx, si */
    m->cx=(uint16_t)(m->si);
L_1731E: ; /* IDA 0x1731E; mov     ax, 9A01h */
    m->ax=(uint16_t)(39425);
L_17321: ; /* IDA 0x17321; push    bx */
    ec_push(m,m->bx);
L_17322: ; /* IDA 0x17322; call    loc_1F75E */
    cc_call(m,63326,29477,h);
L_17325: ; /* IDA 0x17325; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17326: ; /* IDA 0x17326; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17327: ; /* IDA 0x17327; mov     si, di */
    m->si=(uint16_t)(m->di);
L_17329: ; /* IDA 0x17329; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1732B: ; /* IDA 0x1732B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1732D: ; /* IDA 0x1732D; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1732F: ; /* IDA 0x1732F; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17331: ; /* IDA 0x17331; mov     di, bx */
    m->di=(uint16_t)(m->bx);
L_17333: ; /* IDA 0x17333; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17335: ; /* IDA 0x17335; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17337: ; /* IDA 0x17337; add     di, bx */
    m->di=(uint16_t)(ec_math(m,m->di,m->bx,16,0,0));
L_17339: ; /* IDA 0x17339; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1733B: ; /* IDA 0x1733B; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1733D: ; /* IDA 0x1733D; shr     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
L_1733F: ; /* IDA 0x1733F; add     di, dx */
    m->di=(uint16_t)(ec_math(m,m->di,m->dx,16,0,0));
L_17341: ; /* IDA 0x17341; add     di, 9 */
    m->di=(uint16_t)(ec_math(m,m->di,9,16,0,0));
L_17344: ; /* IDA 0x17344; mov     ax, [si+4] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+4)));
L_17347: ; /* IDA 0x17347; mov     bx, 9004h */
    m->bx=(uint16_t)(36868);
L_1734A: ; /* IDA 0x1734A; cmp     ax, 12Ch */
    ec_cmp(m,m->ax,300,16);
L_1734D: ; /* IDA 0x1734D; jnb     short near ptr unk_17351 */
    if(!(m->flags&1)) goto L_17351;
L_1734F: ; /* IDA 0x1734F; mov     bh, 9Ah */
    gw_bh(m,(uint8_t)(154));
L_17351: ; /* IDA 0x17351; mov     dx, 0Ah */
    m->dx=(uint16_t)(10);
L_17354: ; /* IDA 0x17354; mul     dx */
    dl_mul(m,m->dx,16);
L_17356: ; /* IDA 0x17356; call    sub_1062F */
    cc_call(m,1583,29529,h);
L_17359: ; /* IDA 0x17359; add     di, 8 */
    m->di=(uint16_t)(ec_math(m,m->di,8,16,0,0));
L_1735C: ; /* IDA 0x1735C; mov     al, [si+6] */
#if KI_CATALOG_MUTATION == 6
    m->ax=dl_read16(m,m->ds,(uint16_t)(m->si+6));
#else
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+6))));
#endif
L_1735F: ; /* IDA 0x1735F; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_17361: ; /* IDA 0x17361; mov     ah, dl */
#if KI_CATALOG_MUTATION == 6
    ; /* erroneous word morale retained */
#else
    vg_ah(m,(uint8_t)((uint8_t)m->dx));
#endif
L_17363: ; /* IDA 0x17363; mov     bx, 9003h */
    m->bx=(uint16_t)(36867);
L_17366: ; /* IDA 0x17366; cmp     al, 64h ; 'd' */
    ec_cmp(m,(uint8_t)m->ax,100,8);
L_17368: ; /* IDA 0x17368; jnb     short near ptr unk_1736C */
    if(!(m->flags&1)) goto L_1736C;
L_1736A: ; /* IDA 0x1736A; mov     bh, 9Ah */
    gw_bh(m,(uint8_t)(154));
L_1736C: ; /* IDA 0x1736C; call    sub_1062F */
    cc_call(m,1583,29551,h);
L_1736F: ; /* IDA 0x1736F; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17370: ; /* IDA 0x17370; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17371: ; /* IDA 0x17371; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17372: ; /* IDA 0x17372; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17373: ; /* IDA 0x17373; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17374: ; /* IDA 0x17374; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17375: ; /* IDA 0x17375; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17376: ; /* IDA 0x17376; retn */
    return;
}
/* Original unnamed code; IDA 0x1763C. */
static void cc_body_1763C(KiMachine16 *m,const KiEconomyHooks *h) {
L_1763C: ; /* IDA 0x1763C; push    si */
    ec_push(m,m->si);
L_1763D: ; /* IDA 0x1763D; push    di */
    ec_push(m,m->di);
L_1763E: ; /* IDA 0x1763E; mov     si, 4240h */
    m->si=(uint16_t)(16960);
L_17641: ; /* IDA 0x17641; xor     di, di */
    {uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_17643: ; /* IDA 0x17643; mov     al, 0 */
#if KI_CATALOG_MUTATION == 3
    vg_al(m,0);
#else
    vg_al(m,dl_read8(m,m->cs,0x7644));
#endif
L_17645: ; /* IDA 0x17645; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17647: ; /* IDA 0x17647; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_1764A: ; /* IDA 0x1764A; jb      short loc_17657 */
    if(m->flags&1) goto L_17657;
L_1764C: ; /* IDA 0x1764C; cmp     al, [si+1Ch] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+28)),8);
L_1764F: ; /* IDA 0x1764F; jnz     short loc_17657 */
    if(!(m->flags&0x40)) goto L_17657;
L_17651: ; /* IDA 0x17651; mov     [bp+di], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->si);
L_17653: ; /* IDA 0x17653; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_17655: ; /* IDA 0x17655; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_17656: ; /* IDA 0x17656; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_17657: ; /* IDA 0x17657; add     si, 20h ; ' ' */
    m->si=(uint16_t)(ec_math(m,m->si,32,16,0,0));
L_1765A: ; /* IDA 0x1765A; cmp     si, 5240h */
    ec_cmp(m,m->si,21056,16);
L_1765E: ; /* IDA 0x1765E; jnz     short loc_17647 */
    if(!(m->flags&0x40)) goto L_17647;
L_17660: ; /* IDA 0x17660; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17661: ; /* IDA 0x17661; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17662: ; /* IDA 0x17662; retn */
    return;
}
/* Original unnamed code; IDA 0x176A0. */
static void cc_body_176A0(KiMachine16 *m,const KiEconomyHooks *h) {
L_176A0: ; /* IDA 0x176A0; push    si */
    ec_push(m,m->si);
L_176A1: ; /* IDA 0x176A1; push    di */
    ec_push(m,m->di);
L_176A2: ; /* IDA 0x176A2; push    dx */
    ec_push(m,m->dx);
L_176A3: ; /* IDA 0x176A3; call    sub_187FF */
    cc_call(m,34815,30374,h);
L_176A6: ; /* IDA 0x176A6; mov     si, 4240h */
    m->si=(uint16_t)(16960);
L_176A9: ; /* IDA 0x176A9; add     bx, si */
    m->bx=(uint16_t)(ec_math(m,m->bx,m->si,16,0,0));
L_176AB: ; /* IDA 0x176AB; mov     dx, bx */
    m->dx=(uint16_t)(m->bx);
L_176AD: ; /* IDA 0x176AD; xor     di, di */
    {uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
L_176AF: ; /* IDA 0x176AF; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_176B3: ; /* IDA 0x176B3; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_176B5: ; /* IDA 0x176B5; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_176B8: ; /* IDA 0x176B8; jb      short loc_176CF */
    if(m->flags&1) goto L_176CF;
L_176BA: ; /* IDA 0x176BA; cmp     al, [si+1Ch] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+28)),8);
L_176BD: ; /* IDA 0x176BD; jnz     short loc_176CF */
    if(!(m->flags&0x40)) goto L_176CF;
L_176BF: ; /* IDA 0x176BF; cmp     byte ptr [si+17h], 0 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+23)),0,8);
L_176C3: ; /* IDA 0x176C3; jnz     short loc_176CF */
    if(!(m->flags&0x40)) goto L_176CF;
L_176C5: ; /* IDA 0x176C5; cmp     dx, si */
    ec_cmp(m,m->dx,m->si,16);
L_176C7: ; /* IDA 0x176C7; jz      short loc_176CF */
#if KI_CATALOG_MUTATION == 4
    ; /* lost exclusion of ruler */
#else
    if(m->flags&0x40) goto L_176CF;
#endif
L_176C9: ; /* IDA 0x176C9; mov     [bp+di], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->si);
L_176CB: ; /* IDA 0x176CB; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_176CD: ; /* IDA 0x176CD; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_176CE: ; /* IDA 0x176CE; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_176CF: ; /* IDA 0x176CF; add     si, 20h ; ' ' */
    m->si=(uint16_t)(ec_math(m,m->si,32,16,0,0));
L_176D2: ; /* IDA 0x176D2; cmp     si, 5240h */
    ec_cmp(m,m->si,21056,16);
L_176D6: ; /* IDA 0x176D6; jnz     short loc_176B5 */
    if(!(m->flags&0x40)) goto L_176B5;
L_176D8: ; /* IDA 0x176D8; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_176D9: ; /* IDA 0x176D9; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_176DA: ; /* IDA 0x176DA; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_176DB: ; /* IDA 0x176DB; retn */
    return;
}
/* Original unnamed code; IDA 0x176DC. */
static void cc_body_176DC(KiMachine16 *m,const KiEconomyHooks *h) {
L_176DC: ; /* IDA 0x176DC; push    ax */
    ec_push(m,m->ax);
L_176DD: ; /* IDA 0x176DD; push    bx */
    ec_push(m,m->bx);
L_176DE: ; /* IDA 0x176DE; push    dx */
    ec_push(m,m->dx);
L_176DF: ; /* IDA 0x176DF; push    si */
    ec_push(m,m->si);
L_176E0: ; /* IDA 0x176E0; call    sub_18607 */
    cc_call(m,34311,30435,h);
L_176E3: ; /* IDA 0x176E3; mov     dx, cs:word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33198));
L_176E8: ; /* IDA 0x176E8; add     dx, 18h */
    m->dx=(uint16_t)(ec_math(m,m->dx,24,16,0,0));
L_176EB: ; /* IDA 0x176EB; mov     bx, cs:word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33200));
L_176F0: ; /* IDA 0x176F0; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_176F3: ; /* IDA 0x176F3; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_176F5: ; /* IDA 0x176F5; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_176F7: ; /* IDA 0x176F7; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_176F9: ; /* IDA 0x176F9; mov     ax, 0Ah */
    m->ax=(uint16_t)(10);
L_176FC: ; /* IDA 0x176FC; call    sub_1770C */
    cc_call(m,30476,30463,h);
L_176FF: ; /* IDA 0x176FF; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17700: ; /* IDA 0x17700; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17701: ; /* IDA 0x17701; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17704: ; /* IDA 0x17704; dec     ax */
    m->ax=(uint16_t)(st_dec(m,m->ax,16));
L_17705: ; /* IDA 0x17705; jnz     short loc_176FC */
    if(!(m->flags&0x40)) goto L_176FC;
L_17707: ; /* IDA 0x17707; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17708: ; /* IDA 0x17708; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17709: ; /* IDA 0x17709; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1770A: ; /* IDA 0x1770A; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1770B: ; /* IDA 0x1770B; retn */
    return;
}
/* Original unnamed code; IDA 0x178E5. */
static void cc_body_178E5(KiMachine16 *m,const KiEconomyHooks *h) {
L_178E5: ; /* IDA 0x178E5; push    dx */
    ec_push(m,m->dx);
L_178E6: ; /* IDA 0x178E6; push    si */
    ec_push(m,m->si);
L_178E7: ; /* IDA 0x178E7; push    di */
    ec_push(m,m->di);
L_178E8: ; /* IDA 0x178E8; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_178EA: ; /* IDA 0x178EA; mov     di, si */
    m->di=(uint16_t)(m->si);
L_178EC: ; /* IDA 0x178EC; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_178EE: ; /* IDA 0x178EE; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_178F1: ; /* IDA 0x178F1; jb      short loc_178F9 */
    if(m->flags&1) goto L_178F9;
L_178F3: ; /* IDA 0x178F3; mov     [bp+di], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->si);
L_178F5: ; /* IDA 0x178F5; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_178F7: ; /* IDA 0x178F7; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_178F8: ; /* IDA 0x178F8; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_178F9: ; /* IDA 0x178F9; add     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
L_178FC: ; /* IDA 0x178FC; cmp     si, 580h */
    ec_cmp(m,m->si,1408,16);
L_17900: ; /* IDA 0x17900; jnz     short loc_178EE */
    if(!(m->flags&0x40)) goto L_178EE;
L_17902: ; /* IDA 0x17902; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17903: ; /* IDA 0x17903; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17904: ; /* IDA 0x17904; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17905: ; /* IDA 0x17905; retn */
    return;
}
/* Original unnamed code; IDA 0x17944. */
static void cc_body_17944(KiMachine16 *m,const KiEconomyHooks *h) {
L_17944: ; /* IDA 0x17944; push    dx */
    ec_push(m,m->dx);
L_17945: ; /* IDA 0x17945; push    si */
    ec_push(m,m->si);
L_17946: ; /* IDA 0x17946; push    di */
    ec_push(m,m->di);
L_17947: ; /* IDA 0x17947; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_17949: ; /* IDA 0x17949; mov     di, si */
    m->di=(uint16_t)(m->si);
L_1794B: ; /* IDA 0x1794B; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_1794D: ; /* IDA 0x1794D; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_17950: ; /* IDA 0x17950; jb      short loc_1795F */
    if(m->flags&1) goto L_1795F;
L_17952: ; /* IDA 0x17952; cmp     si, cs:word_10CFD */
    ec_cmp(m,m->si,dl_read16(m,m->cs,3325),16);
L_17957: ; /* IDA 0x17957; jz      short loc_1795F */
#if KI_CATALOG_MUTATION == 9
    ; /* lost exclusion of player */
#else
    if(m->flags&0x40) goto L_1795F;
#endif
L_17959: ; /* IDA 0x17959; mov     [bp+di], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->si);
L_1795B: ; /* IDA 0x1795B; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_1795D: ; /* IDA 0x1795D; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1795E: ; /* IDA 0x1795E; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_1795F: ; /* IDA 0x1795F; add     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
L_17962: ; /* IDA 0x17962; cmp     si, 580h */
    ec_cmp(m,m->si,1408,16);
L_17966: ; /* IDA 0x17966; jnz     short loc_1794D */
    if(!(m->flags&0x40)) goto L_1794D;
L_17968: ; /* IDA 0x17968; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17969: ; /* IDA 0x17969; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1796A: ; /* IDA 0x1796A; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1796B: ; /* IDA 0x1796B; retn */
    return;
}
/* Original unnamed code; IDA 0x1796C. */
static void cc_body_1796C(KiMachine16 *m,const KiEconomyHooks *h) {
L_1796C: ; /* IDA 0x1796C; push    ax */
    ec_push(m,m->ax);
L_1796D: ; /* IDA 0x1796D; push    bx */
    ec_push(m,m->bx);
L_1796E: ; /* IDA 0x1796E; push    dx */
    ec_push(m,m->dx);
L_1796F: ; /* IDA 0x1796F; push    si */
    ec_push(m,m->si);
L_17970: ; /* IDA 0x17970; call    sub_18607 */
    cc_call(m,34311,31091,h);
L_17973: ; /* IDA 0x17973; mov     dx, cs:word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33198));
L_17978: ; /* IDA 0x17978; add     dx, 18h */
    m->dx=(uint16_t)(ec_math(m,m->dx,24,16,0,0));
L_1797B: ; /* IDA 0x1797B; mov     bx, cs:word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33200));
L_17980: ; /* IDA 0x17980; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17983: ; /* IDA 0x17983; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17985: ; /* IDA 0x17985; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_17987: ; /* IDA 0x17987; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17989: ; /* IDA 0x17989; mov     ax, 0Ah */
    m->ax=(uint16_t)(10);
L_1798C: ; /* IDA 0x1798C; call    sub_1799C */
    cc_call(m,31132,31119,h);
L_1798F: ; /* IDA 0x1798F; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17990: ; /* IDA 0x17990; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17991: ; /* IDA 0x17991; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17994: ; /* IDA 0x17994; dec     ax */
    m->ax=(uint16_t)(st_dec(m,m->ax,16));
L_17995: ; /* IDA 0x17995; jnz     short loc_1798C */
    if(!(m->flags&0x40)) goto L_1798C;
L_17997: ; /* IDA 0x17997; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17998: ; /* IDA 0x17998; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17999: ; /* IDA 0x17999; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1799A: ; /* IDA 0x1799A; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1799B: ; /* IDA 0x1799B; retn */
    return;
}
/* Original unnamed code; IDA 0x17B6F. */
static void cc_body_17B6F(KiMachine16 *m,const KiEconomyHooks *h) {
L_17B6F: ; /* IDA 0x17B6F; push    dx */
    ec_push(m,m->dx);
L_17B70: ; /* IDA 0x17B70; push    si */
    ec_push(m,m->si);
L_17B71: ; /* IDA 0x17B71; push    di */
    ec_push(m,m->di);
L_17B72: ; /* IDA 0x17B72; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_17B74: ; /* IDA 0x17B74; mov     di, si */
    m->di=(uint16_t)(m->si);
L_17B76: ; /* IDA 0x17B76; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17B78: ; /* IDA 0x17B78; cmp     byte ptr [si], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
L_17B7B: ; /* IDA 0x17B7B; jb      short loc_17B83 */
    if(m->flags&1) goto L_17B83;
L_17B7D: ; /* IDA 0x17B7D; mov     [bp+di], si */
    dl_write16(m,m->ss,(uint16_t)(m->bp+m->di),m->si);
L_17B7F: ; /* IDA 0x17B7F; inc     ah */
    vg_ah(m,(uint8_t)(ec_inc(m,(uint8_t)(m->ax>>8),8)));
L_17B81: ; /* IDA 0x17B81; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_17B82: ; /* IDA 0x17B82; inc     di */
    m->di=(uint16_t)(ec_inc(m,m->di,16));
L_17B83: ; /* IDA 0x17B83; add     si, 40h ; '@' */
    m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
L_17B86: ; /* IDA 0x17B86; cmp     si, 580h */
    ec_cmp(m,m->si,1408,16);
L_17B8A: ; /* IDA 0x17B8A; jnz     short loc_17B78 */
    if(!(m->flags&0x40)) goto L_17B78;
L_17B8C: ; /* IDA 0x17B8C; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17B8D: ; /* IDA 0x17B8D; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17B8E: ; /* IDA 0x17B8E; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17B8F: ; /* IDA 0x17B8F; retn */
    return;
}
/* Original unnamed code; IDA 0x17B90. */
static void cc_body_17B90(KiMachine16 *m,const KiEconomyHooks *h) {
L_17B90: ; /* IDA 0x17B90; push    ax */
    ec_push(m,m->ax);
L_17B91: ; /* IDA 0x17B91; push    bx */
    ec_push(m,m->bx);
L_17B92: ; /* IDA 0x17B92; push    dx */
    ec_push(m,m->dx);
L_17B93: ; /* IDA 0x17B93; push    si */
    ec_push(m,m->si);
L_17B94: ; /* IDA 0x17B94; call    sub_18607 */
    cc_call(m,34311,31639,h);
L_17B97: ; /* IDA 0x17B97; mov     dx, cs:word_181AE */
    m->dx=(uint16_t)(dl_read16(m,m->cs,33198));
L_17B9C: ; /* IDA 0x17B9C; add     dx, 18h */
    m->dx=(uint16_t)(ec_math(m,m->dx,24,16,0,0));
L_17B9F: ; /* IDA 0x17B9F; mov     bx, cs:word_181B0 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,33200));
L_17BA4: ; /* IDA 0x17BA4; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17BA7: ; /* IDA 0x17BA7; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_17BA9: ; /* IDA 0x17BA9; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_17BAB: ; /* IDA 0x17BAB; mov     si, ax */
    m->si=(uint16_t)(m->ax);
L_17BAD: ; /* IDA 0x17BAD; mov     ax, 0Ah */
    m->ax=(uint16_t)(10);
L_17BB0: ; /* IDA 0x17BB0; call    sub_17BC0 */
    cc_call(m,31680,31667,h);
L_17BB3: ; /* IDA 0x17BB3; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17BB4: ; /* IDA 0x17BB4; inc     si */
    m->si=(uint16_t)(ec_inc(m,m->si,16));
L_17BB5: ; /* IDA 0x17BB5; add     bx, 10h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16,16,0,0));
L_17BB8: ; /* IDA 0x17BB8; dec     ax */
    m->ax=(uint16_t)(st_dec(m,m->ax,16));
L_17BB9: ; /* IDA 0x17BB9; jnz     short loc_17BB0 */
    if(!(m->flags&0x40)) goto L_17BB0;
L_17BBB: ; /* IDA 0x17BBB; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17BBC: ; /* IDA 0x17BBC; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17BBD: ; /* IDA 0x17BBD; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_17BBE: ; /* IDA 0x17BBE; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_17BBF: ; /* IDA 0x17BBF; retn */
    return;
}
int ki_catalog_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x716d:cc_body_1716D(m,h);break;
case 0x71d3:cc_body_171D3(m,h);break;
case 0x75fa:cc_body_175FA(m,h);break;
case 0x7663:cc_body_17663(m,h);break;
case 0x770c:cc_body_1770C(m,h);break;
case 0x78a7:cc_body_178A7(m,h);break;
case 0x7906:cc_body_17906(m,h);break;
case 0x799c:cc_body_1799C(m,h);break;
case 0x7a7a:cc_body_17A7A(m,h);break;
case 0x7b3c:cc_body_17B3C(m,h);break;
case 0x7bc0:cc_body_17BC0(m,h);break;
case 0x71a8:cc_body_171A8(m,h);break;
case 0x7217:cc_body_17217(m,h);break;
case 0x724d:cc_body_1724D(m,h);break;
case 0x727d:cc_body_1727D(m,h);break;
case 0x763c:cc_body_1763C(m,h);break;
case 0x76a0:cc_body_176A0(m,h);break;
case 0x76dc:cc_body_176DC(m,h);break;
case 0x78e5:cc_body_178E5(m,h);break;
case 0x7944:cc_body_17944(m,h);break;
case 0x796c:cc_body_1796C(m,h);break;
case 0x7b6f:cc_body_17B6F(m,h);break;
case 0x7b90:cc_body_17B90(m,h);break;
default:return 0;}return 1;}
static int cc_known(uint16_t t) {switch(t) {case 0x716d: case 0x71d3: case 0x75fa: case 0x7663: case 0x770c: case 0x78a7: case 0x7906: case 0x799c: case 0x7a7a: case 0x7b3c: case 0x7bc0: case 0x71a8: case 0x7217: case 0x724d: case 0x727d: case 0x763c: case 0x76a0: case 0x76dc: case 0x78e5: case 0x7944: case 0x796c: case 0x7b6f: case 0x7b90: return 1;default:return 0;}}
