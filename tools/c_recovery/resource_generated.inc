/* Fixed IDA native C resource control; spec/233, re/113. Original locators retained. */
/* Original sub_187AF; IDA 0x187AF. */
static void rs_body_187AF(KiMachine16 *m,const KiEconomyHooks *h) {
L_187AF: ; /* IDA 0x187AF; mov     dx, 0DB8h */
    m->dx=(uint16_t)(3512);
L_187B2: ; /* IDA 0x187B2; mov     bx, ds:9876h */
    m->bx=(uint16_t)(dl_read16(m,m->ds,39030));
L_187B6: ; /* IDA 0x187B6; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_187B8: ; /* IDA 0x187B8; call    sub_1E378 */
    rs_call(m,58232,34747,h);
L_187BB: ; /* IDA 0x187BB; mov     dx, 0DC1h */
    m->dx=(uint16_t)(3521);
L_187BE: ; /* IDA 0x187BE; mov     bx, ds:9876h */
    m->bx=(uint16_t)(dl_read16(m,m->ds,39030));
L_187C2: ; /* IDA 0x187C2; add     bx, 800h */
#if KI_RESOURCE_MUTATION == 4
    m->bx=(uint16_t)(ec_math(m,m->bx,1024,16,0,0));
#else
    m->bx=(uint16_t)(ec_math(m,m->bx,2048,16,0,0));
#endif
L_187C6: ; /* IDA 0x187C6; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_187C8: ; /* IDA 0x187C8; call    sub_1E378 */
    rs_call(m,58232,34763,h);
L_187CB: ; /* IDA 0x187CB; retn */
    return;
}
/* Original sub_1E378; IDA 0x1E378. */
static void rs_body_1E378(KiMachine16 *m,const KiEconomyHooks *h) {
L_1E378: ; /* IDA 0x1E378; push    ds */
    ec_push(m,m->ds);
L_1E379: ; /* IDA 0x1E379; push    bx */
    ec_push(m,m->bx);
L_1E37A: ; /* IDA 0x1E37A; push    dx */
    ec_push(m,m->dx);
L_1E37B: ; /* IDA 0x1E37B; push    si */
    ec_push(m,m->si);
L_1E37C: ; /* IDA 0x1E37C; call    sub_1F4A2 */
    rs_call(m,62626,58239,h);
L_1E37F: ; /* IDA 0x1E37F; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1E380: ; /* IDA 0x1E380; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_1E381: ; /* IDA 0x1E381; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1E382: ; /* IDA 0x1E382; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1E383: ; /* IDA 0x1E383; jnb     short locret_1E38B */
    if(!(m->flags&1)) goto L_1E38B;
L_1E385: ; /* IDA 0x1E385; mov     ah, 0 */
    vg_ah(m,(uint8_t)(0));
L_1E387: ; /* IDA 0x1E387; int     50h; through 57 - IRQ0-IRQ7 relocated by DESQview */
    gl_interrupt(m,80,58249);
L_1E389: ; /* IDA 0x1E389; jmp     short sub_1E378 */
    goto L_1E378;
L_1E38B: ; /* IDA 0x1E38B; retn */
    return;
}
/* Original sub_1F4A2; IDA 0x1F4A2. */
static void rs_body_1F4A2(KiMachine16 *m,const KiEconomyHooks *h) {
L_1F4A2: ; /* IDA 0x1F4A2; push    cx */
    ec_push(m,m->cx);
L_1F4A3: ; /* IDA 0x1F4A3; mov     ax, 3D00h */
    m->ax=(uint16_t)(15616);
L_1F4A6: ; /* IDA 0x1F4A6; int     21h; DOS - 2+ - OPEN DISK FILE WITH HANDLE */
    gl_interrupt(m,33,62632);
L_1F4A8: ; /* IDA 0x1F4A8; jb      short loc_1F4DB */
    if(m->flags&1) goto L_1F4DB;
L_1F4AA: ; /* IDA 0x1F4AA; mov     cs:word_1F4DD, ax */
    dl_write16(m,m->cs,62685,m->ax);
L_1F4AE: ; /* IDA 0x1F4AE; push    ds */
    ec_push(m,m->ds);
L_1F4AF: ; /* IDA 0x1F4AF; mov     ah, 3Fh ; '?' */
    vg_ah(m,(uint8_t)(63));
L_1F4B1: ; /* IDA 0x1F4B1; push    bx */
    ec_push(m,m->bx);
L_1F4B2: ; /* IDA 0x1F4B2; push    si */
    ec_push(m,m->si);
L_1F4B3: ; /* IDA 0x1F4B3; mov     ds, bx */
    m->ds=(uint16_t)(m->bx);
L_1F4B5: ; /* IDA 0x1F4B5; mov     dx, si */
    m->dx=(uint16_t)(m->si);
L_1F4B7: ; /* IDA 0x1F4B7; mov     cx, 0F000h */
#if KI_RESOURCE_MUTATION == 1
    m->cx=(uint16_t)(32768);
#else
    m->cx=(uint16_t)(61440);
#endif
L_1F4BA: ; /* IDA 0x1F4BA; mov     bx, cs:word_1F4DD */
    m->bx=(uint16_t)(dl_read16(m,m->cs,62685));
L_1F4BF: ; /* IDA 0x1F4BF; int     21h; DOS - 2+ - READ FROM FILE WITH HANDLE */
    gl_interrupt(m,33,62657);
L_1F4C1: ; /* IDA 0x1F4C1; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_1F4C2: ; /* IDA 0x1F4C2; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1F4C3: ; /* IDA 0x1F4C3; cmp     ax, 0F000h */
    ec_cmp(m,m->ax,61440,16);
L_1F4C6: ; /* IDA 0x1F4C6; jb      short loc_1F4CE */
#if KI_RESOURCE_MUTATION == 3
    if(m->flags&(1|0x40)) goto L_1F4CE;
#else
    if(m->flags&1) goto L_1F4CE;
#endif
L_1F4C8: ; /* IDA 0x1F4C8; add     bx, 0F00h */
#if KI_RESOURCE_MUTATION == 2
    m->bx=(uint16_t)(ec_math(m,m->bx,4096,16,0,0));
#else
    m->bx=(uint16_t)(ec_math(m,m->bx,3840,16,0,0));
#endif
L_1F4CC: ; /* IDA 0x1F4CC; jmp     short loc_1F4AF */
    goto L_1F4AF;
L_1F4CE: ; /* IDA 0x1F4CE; mov     ah, 3Eh ; '>' */
    vg_ah(m,(uint8_t)(62));
L_1F4D0: ; /* IDA 0x1F4D0; mov     bx, cs:word_1F4DD */
    m->bx=(uint16_t)(dl_read16(m,m->cs,62685));
L_1F4D5: ; /* IDA 0x1F4D5; int     21h; DOS - 2+ - CLOSE A FILE WITH HANDLE */
    gl_interrupt(m,33,62679);
L_1F4D7: ; /* IDA 0x1F4D7; mov     ax, 0FFFFh */
#if KI_RESOURCE_MUTATION == 12
    m->ax=(uint16_t)(0);
#else
    m->ax=(uint16_t)(65535);
#endif
L_1F4DA: ; /* IDA 0x1F4DA; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1F4DB: ; /* IDA 0x1F4DB; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1F4DC: ; /* IDA 0x1F4DC; retn */
    return;
}
/* Original sub_10241; IDA 0x10241. */
static void rs_body_10241(KiMachine16 *m,const KiEconomyHooks *h) {
L_10241: ; /* IDA 0x10241; cmp     cs:byte_1020E, al */
    ec_cmp(m,dl_read8(m,m->cs,526),(uint8_t)m->ax,8);
L_10246: ; /* IDA 0x10246; jz      short locret_102C1 */
#if KI_RESOURCE_MUTATION == 5
    ; /* lost cue same-value gate */
#else
    if(m->flags&0x40) goto L_102C1;
#endif
L_10248: ; /* IDA 0x10248; mov     cs:byte_1020E, al */
    dl_write8(m,m->cs,526,(uint8_t)m->ax);
L_1024C: ; /* IDA 0x1024C; push    ds */
    ec_push(m,m->ds);
L_1024D: ; /* IDA 0x1024D; push    ax */
    ec_push(m,m->ax);
L_1024E: ; /* IDA 0x1024E; push    bx */
    ec_push(m,m->bx);
L_1024F: ; /* IDA 0x1024F; push    cx */
    ec_push(m,m->cx);
L_10250: ; /* IDA 0x10250; push    dx */
    ec_push(m,m->dx);
L_10251: ; /* IDA 0x10251; push    si */
    ec_push(m,m->si);
L_10252: ; /* IDA 0x10252; push    di */
    ec_push(m,m->di);
L_10253: ; /* IDA 0x10253; push    ax */
    ec_push(m,m->ax);
L_10254: ; /* IDA 0x10254; mov     ah, 8 */
    vg_ah(m,(uint8_t)(8));
L_10256: ; /* IDA 0x10256; int     61h; reserved for user interrupt */
    gl_interrupt(m,97,600);
L_10258: ; /* IDA 0x10258; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_1025A: ; /* IDA 0x1025A; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_1025C: ; /* IDA 0x1025C; mov     dx, 0DF0h */
    m->dx=(uint16_t)(3568);
L_1025F: ; /* IDA 0x1025F; mov     bx, cs:word_10D3C */
    m->bx=(uint16_t)(dl_read16(m,m->cs,3388));
L_10264: ; /* IDA 0x10264; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_10266: ; /* IDA 0x10266; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_10268: ; /* IDA 0x10268; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_1026A: ; /* IDA 0x1026A; mov     di, 80h */
    m->di=(uint16_t)(128);
L_1026D: ; /* IDA 0x1026D; call    sub_1E38C */
    rs_call(m,58252,624,h);
L_10270: ; /* IDA 0x10270; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10271: ; /* IDA 0x10271; push    ax */
    ec_push(m,m->ax);
L_10272: ; /* IDA 0x10272; mov     ds, cs:word_10D3C */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3388));
L_10277: ; /* IDA 0x10277; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_10279: ; /* IDA 0x10279; xor     bh, bh */
#if KI_RESOURCE_MUTATION == 6
    ; /* lost cue offset high-byte clear */
#else
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
#endif
L_1027B: ; /* IDA 0x1027B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1027D: ; /* IDA 0x1027D; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1027F: ; /* IDA 0x1027F; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_10281: ; /* IDA 0x10281; mov     ax, [bx] */
#if KI_RESOURCE_MUTATION == 7
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+2)));
#else
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx)));
#endif
L_10283: ; /* IDA 0x10283; mov     cx, [bx+2] */
    m->cx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+2)));
L_10286: ; /* IDA 0x10286; mov     di, [bx+4] */
    m->di=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+4)));
L_10289: ; /* IDA 0x10289; mov     dx, cs */
    m->dx=(uint16_t)(m->cs);
L_1028B: ; /* IDA 0x1028B; mov     ds, dx */
    m->ds=(uint16_t)(m->dx);
L_1028D: ; /* IDA 0x1028D; mov     dx, 0DF0h */
    m->dx=(uint16_t)(3568);
L_10290: ; /* IDA 0x10290; mov     bx, cs:word_10D3C */
    m->bx=(uint16_t)(dl_read16(m,m->cs,3388));
L_10295: ; /* IDA 0x10295; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_10297: ; /* IDA 0x10297; call    sub_1E38C */
    rs_call(m,58252,666,h);
L_1029A: ; /* IDA 0x1029A; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_1029B: ; /* IDA 0x1029B; mov     ax, 606h */
    m->ax=(uint16_t)(1542);
L_1029E: ; /* IDA 0x1029E; cmp     bl, 2 */
    ec_cmp(m,(uint8_t)m->bx,2,8);
L_102A1: ; /* IDA 0x102A1; jb      short loc_102A5 */
    if(m->flags&1) goto L_102A5;
L_102A3: ; /* IDA 0x102A3; mov     al, 5 */
    vg_al(m,(uint8_t)(5));
L_102A5: ; /* IDA 0x102A5; mov     ds, cs:word_10D3C */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3388));
L_102AA: ; /* IDA 0x102AA; xor     si, si */
    {uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
L_102AC: ; /* IDA 0x102AC; int     61h; reserved for user interrupt */
    gl_interrupt(m,97,686);
L_102AE: ; /* IDA 0x102AE; cmp     cs:byte_10CF9, 0 */
    ec_cmp(m,dl_read8(m,m->cs,3321),0,8);
L_102B4: ; /* IDA 0x102B4; jz      short loc_102BA */
#if KI_RESOURCE_MUTATION == 8
    if(!(m->flags&0x40)) goto L_102BA;
#else
    if(m->flags&0x40) goto L_102BA;
#endif
L_102B6: ; /* IDA 0x102B6; mov     ah, 7 */
    vg_ah(m,(uint8_t)(7));
L_102B8: ; /* IDA 0x102B8; int     61h; reserved for user interrupt */
    gl_interrupt(m,97,698);
L_102BA: ; /* IDA 0x102BA; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_102BB: ; /* IDA 0x102BB; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_102BC: ; /* IDA 0x102BC; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_102BD: ; /* IDA 0x102BD; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_102BE: ; /* IDA 0x102BE; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_102BF: ; /* IDA 0x102BF; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_102C0: ; /* IDA 0x102C0; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_102C1: ; /* IDA 0x102C1; retn */
    return;
}
/* Original sub_102C2; IDA 0x102C2. */
static void rs_body_102C2(KiMachine16 *m,const KiEconomyHooks *h) {
L_102C2: ; /* IDA 0x102C2; push    ax */
    ec_push(m,m->ax);
L_102C3: ; /* IDA 0x102C3; mov     cs:byte_1020E, 0FFh */
#if KI_RESOURCE_MUTATION == 11
    dl_write8(m,m->cs,526,0);
#else
    dl_write8(m,m->cs,526,255);
#endif
L_102C9: ; /* IDA 0x102C9; mov     ax, 9F2h */
    m->ax=(uint16_t)(2546);
L_102CC: ; /* IDA 0x102CC; int     61h; reserved for user interrupt */
    gl_interrupt(m,97,718);
L_102CE: ; /* IDA 0x102CE; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_102CF: ; /* IDA 0x102CF; retn */
    return;
}
/* Original sub_100DF; IDA 0x100DF. */
static void rs_body_100DF(KiMachine16 *m,const KiEconomyHooks *h) {
L_100DF: ; /* IDA 0x100DF; push    ax */
    ec_push(m,m->ax);
L_100E0: ; /* IDA 0x100E0; push    bx */
    ec_push(m,m->bx);
L_100E1: ; /* IDA 0x100E1; mov     ah, 48h ; 'H' */
    vg_ah(m,(uint8_t)(72));
L_100E3: ; /* IDA 0x100E3; mov     bx, 1D5Eh */
    m->bx=(uint16_t)(7518);
L_100E6: ; /* IDA 0x100E6; int     21h; DOS - 2+ - ALLOCATE MEMORY */
    gl_interrupt(m,33,232);
L_100E8: ; /* IDA 0x100E8; jb      short loc_10146 */
    if(m->flags&1) goto L_10146;
L_100EA: ; /* IDA 0x100EA; mov     word_10D38, ax */
    dl_write16(m,m->ds,3384,m->ax);
L_100ED: ; /* IDA 0x100ED; add     ax, 0D00h */
#if KI_RESOURCE_MUTATION == 9
    m->ax=(uint16_t)(ec_math(m,m->ax,3327,16,0,0));
#else
    m->ax=(uint16_t)(ec_math(m,m->ax,3328,16,0,0));
#endif
L_100F0: ; /* IDA 0x100F0; mov     word_10D3A, ax */
    dl_write16(m,m->ds,3386,m->ax);
L_100F3: ; /* IDA 0x100F3; add     ax, 300h */
    m->ax=(uint16_t)(ec_math(m,m->ax,768,16,0,0));
L_100F6: ; /* IDA 0x100F6; mov     word_10D3C, ax */
    dl_write16(m,m->ds,3388,m->ax);
L_100F9: ; /* IDA 0x100F9; add     ax, 1C0h */
    m->ax=(uint16_t)(ec_math(m,m->ax,448,16,0,0));
L_100FC: ; /* IDA 0x100FC; mov     word_10D3E, ax */
    dl_write16(m,m->ds,3390,m->ax);
L_100FF: ; /* IDA 0x100FF; add     ax, 200h */
    m->ax=(uint16_t)(ec_math(m,m->ax,512,16,0,0));
L_10102: ; /* IDA 0x10102; mov     word_10D56, ax */
    dl_write16(m,m->ds,3414,m->ax);
L_10105: ; /* IDA 0x10105; add     ax, 40h ; '@' */
    m->ax=(uint16_t)(ec_math(m,m->ax,64,16,0,0));
L_10108: ; /* IDA 0x10108; mov     word_10D40, ax */
    dl_write16(m,m->ds,3392,m->ax);
L_1010B: ; /* IDA 0x1010B; add     ax, 200h */
    m->ax=(uint16_t)(ec_math(m,m->ax,512,16,0,0));
L_1010E: ; /* IDA 0x1010E; mov     word_10D52, ax */
    dl_write16(m,m->ds,3410,m->ax);
L_10111: ; /* IDA 0x10111; add     ax, 524h */
    m->ax=(uint16_t)(ec_math(m,m->ax,1316,16,0,0));
L_10114: ; /* IDA 0x10114; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_10116: ; /* IDA 0x10116; mov     word_10D48, ax */
    dl_write16(m,m->ds,3400,m->ax);
L_10119: ; /* IDA 0x10119; add     ax, 6Ch ; 'l' */
    m->ax=(uint16_t)(ec_math(m,m->ax,108,16,0,0));
L_1011C: ; /* IDA 0x1011C; mov     word_10D4A, ax */
    dl_write16(m,m->ds,3402,m->ax);
L_1011F: ; /* IDA 0x1011F; add     ax, 18h */
#if KI_RESOURCE_MUTATION == 10
    m->ax=(uint16_t)(ec_math(m,m->ax,23,16,0,0));
#else
    m->ax=(uint16_t)(ec_math(m,m->ax,24,16,0,0));
#endif
L_10122: ; /* IDA 0x10122; mov     word_10D54, ax */
    dl_write16(m,m->ds,3412,m->ax);
L_10125: ; /* IDA 0x10125; add     ax, 0Bh */
    m->ax=(uint16_t)(ec_math(m,m->ax,11,16,0,0));
L_10128: ; /* IDA 0x10128; mov     word_10D4C, ax */
    dl_write16(m,m->ds,3404,m->ax);
L_1012B: ; /* IDA 0x1012B; add     ax, 0Bh */
    m->ax=(uint16_t)(ec_math(m,m->ax,11,16,0,0));
L_1012E: ; /* IDA 0x1012E; mov     word_10D50, ax */
    dl_write16(m,m->ds,3408,m->ax);
L_10131: ; /* IDA 0x10131; add     ax, 1A0h */
    m->ax=(uint16_t)(ec_math(m,m->ax,416,16,0,0));
L_10134: ; /* IDA 0x10134; mov     ah, 48h ; 'H' */
    vg_ah(m,(uint8_t)(72));
L_10136: ; /* IDA 0x10136; mov     bx, 5356h */
    m->bx=(uint16_t)(21334);
L_10139: ; /* IDA 0x10139; int     21h; DOS - 2+ - ALLOCATE MEMORY */
    gl_interrupt(m,33,315);
L_1013B: ; /* IDA 0x1013B; jb      short loc_10146 */
    if(m->flags&1) goto L_10146;
L_1013D: ; /* IDA 0x1013D; mov     word_10D46, ax */
    dl_write16(m,m->ds,3398,m->ax);
L_10140: ; /* IDA 0x10140; mov     word_10D44, ax */
    dl_write16(m,m->ds,3396,m->ax);
L_10143: ; /* IDA 0x10143; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_10144: ; /* IDA 0x10144; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10145: ; /* IDA 0x10145; retn */
    return;
L_10146: ; /* IDA 0x10146; push    ax */
    ec_push(m,m->ax);
L_10147: ; /* IDA 0x10147; push    dx */
    ec_push(m,m->dx);
L_10148: ; /* IDA 0x10148; mov     dx, 1974h */
    m->dx=(uint16_t)(6516);
L_1014B: ; /* IDA 0x1014B; mov     ah, 9 */
    vg_ah(m,(uint8_t)(9));
L_1014D: ; /* IDA 0x1014D; int     21h; DOS - PRINT STRING */
    gl_interrupt(m,33,335);
L_1014F: ; /* IDA 0x1014F; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_10150: ; /* IDA 0x10150; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_10151: ; /* IDA 0x10151; mov     ah, 4Ch */
    vg_ah(m,(uint8_t)(76));
L_10153: ; /* IDA 0x10153; int     21h; DOS - 2+ - QUIT WITH EXIT CODE (EXIT) */
    gl_interrupt(m,33,341);
}
int ki_resource_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x87af:rs_body_187AF(m,h);break;
case 0xe378:rs_body_1E378(m,h);break;
case 0xf4a2:rs_body_1F4A2(m,h);break;
case 0x241:rs_body_10241(m,h);break;
case 0x2c2:rs_body_102C2(m,h);break;
case 0xdf:rs_body_100DF(m,h);break;
default:return 0;}return 1;}
