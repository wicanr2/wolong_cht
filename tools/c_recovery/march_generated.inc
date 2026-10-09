/* Original native C march; re/123, spec/243. */
/* Original sub_11C8D; IDA 0x11C8D. */
static void mch_body_11C8D(KiMachine16 *m,const KiEconomyHooks *h) {
L_11C8D: ; /* IDA 0x11C8D; push    ds */
    ec_push(m,m->ds);
L_11C8E: ; /* IDA 0x11C8E; push    es */
    ec_push(m,m->es);
L_11C8F: ; /* IDA 0x11C8F; push    ax */
    ec_push(m,m->ax);
L_11C90: ; /* IDA 0x11C90; push    bx */
    ec_push(m,m->bx);
L_11C91: ; /* IDA 0x11C91; push    cx */
    ec_push(m,m->cx);
L_11C92: ; /* IDA 0x11C92; push    dx */
    ec_push(m,m->dx);
L_11C93: ; /* IDA 0x11C93; push    si */
    ec_push(m,m->si);
L_11C94: ; /* IDA 0x11C94; push    di */
    ec_push(m,m->di);
L_11C95: ; /* IDA 0x11C95; mov     dx, cs:word_1988E */
    m->dx=(uint16_t)(dl_read16(m,m->cs,39054));
L_11C9A: ; /* IDA 0x11C9A; mov     bx, cs:word_19890 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,39056));
L_11C9F: ; /* IDA 0x11C9F; call    sub_1D615 */
    mch_call(m,54805,7330,h);
L_11CA2: ; /* IDA 0x11CA2; call    sub_11CC9 */
    mch_call(m,7369,7333,h);
L_11CA5: ; /* IDA 0x11CA5; call    sub_1D66A */
    mch_call(m,54890,7336,h);
L_11CA8: ; /* IDA 0x11CA8; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_11CA9: ; /* IDA 0x11CA9; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_11CAA: ; /* IDA 0x11CAA; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_11CAB: ; /* IDA 0x11CAB; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_11CAC: ; /* IDA 0x11CAC; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_11CAD: ; /* IDA 0x11CAD; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_11CAE: ; /* IDA 0x11CAE; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_11CAF: ; /* IDA 0x11CAF; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_11CB0: ; /* IDA 0x11CB0; retn */
    return;
}
/* Original sub_11F7F; IDA 0x11F7F. */
static void mch_body_11F7F(KiMachine16 *m,const KiEconomyHooks *h) {
L_11F7F: ; /* IDA 0x11F7F; push    ds */
    ec_push(m,m->ds);
L_11F80: ; /* IDA 0x11F80; push    bp */
    ec_push(m,m->bp);
L_11F81: ; /* IDA 0x11F81; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_11F83: ; /* IDA 0x11F83; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_11F85: ; /* IDA 0x11F85; xor     bp, bp */
    {uint16_t result=((m->bp)^(m->bp))&65535;ec_logic(m,result,16);m->bp=(uint16_t)(result);}
L_11F87: ; /* IDA 0x11F87; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_11F8A: ; /* IDA 0x11F8A; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_11F8C: ; /* IDA 0x11F8C; call    sub_20000 */
    mx_far(m,0x1f8c,0x1f91,h);
L_11F91: ; /* IDA 0x11F91; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_11F93: ; /* IDA 0x11F93; jz      short loc_11F98 */
    if(m->flags&0x40) goto L_11F98;
L_11F95: ; /* IDA 0x11F95; mov     bp, 1 */
    m->bp=(uint16_t)(1);
L_11F98: ; /* IDA 0x11F98; mov     ax, 5 */
    m->ax=(uint16_t)(5);
L_11F9B: ; /* IDA 0x11F9B; mov     bx, 1 */
    m->bx=(uint16_t)(1);
L_11F9E: ; /* IDA 0x11F9E; call    sub_20000 */
    mx_far(m,0x1f9e,0x1fa3,h);
L_11FA3: ; /* IDA 0x11FA3; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_11FA5: ; /* IDA 0x11FA5; jz      short loc_11FAA */
    if(m->flags&0x40) goto L_11FAA;
L_11FA7: ; /* IDA 0x11FA7; or      bp, 2 */
    {uint16_t result=((m->bp)|(2))&65535;ec_logic(m,result,16);m->bp=(uint16_t)(result);}
L_11FAA: ; /* IDA 0x11FAA; mov     ax, 3 */
    m->ax=(uint16_t)(3);
L_11FAD: ; /* IDA 0x11FAD; call    sub_20000 */
    mx_far(m,0x1fad,0x1fb2,h);
L_11FB2: ; /* IDA 0x11FB2; cmp     cx, word_1987E */
    ec_cmp(m,m->cx,dl_read16(m,m->ds,39038),16);
L_11FB6: ; /* IDA 0x11FB6; jnz     short loc_11FD0 */
    if(!(m->flags&0x40)) goto L_11FD0;
L_11FB8: ; /* IDA 0x11FB8; cmp     dx, word_19880 */
    ec_cmp(m,m->dx,dl_read16(m,m->ds,39040),16);
L_11FBC: ; /* IDA 0x11FBC; jnz     short loc_11FD0 */
    if(!(m->flags&0x40)) goto L_11FD0;
L_11FBE: ; /* IDA 0x11FBE; mov     cx, word_19886 */
    m->cx=(uint16_t)(dl_read16(m,m->ds,39046));
L_11FC2: ; /* IDA 0x11FC2; mov     dx, word_19888 */
    m->dx=(uint16_t)(dl_read16(m,m->ds,39048));
L_11FC6: ; /* IDA 0x11FC6; mov     bx, bp */
    m->bx=(uint16_t)(m->bp);
L_11FC8: ; /* IDA 0x11FC8; or      byte_198A3, 80h */
    {uint16_t result=((dl_read8(m,m->ds,39075))|(128))&255;ec_logic(m,result,8);dl_write8(m,m->ds,39075,result);}
L_11FCD: ; /* IDA 0x11FCD; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_11FCE: ; /* IDA 0x11FCE; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_11FCF: ; /* IDA 0x11FCF; retn */
    return;
L_11FD0: ; /* IDA 0x11FD0; mov     ds:987Eh, cx */
    dl_write16(m,m->ds,39038,m->cx);
L_11FD4: ; /* IDA 0x11FD4; mov     ds:9880h, dx */
    dl_write16(m,m->ds,39040,m->dx);
L_11FD8: ; /* IDA 0x11FD8; mov     byte ptr ds:98A5h, 0Ch */
    dl_write8(m,m->ds,39077,12);
L_11FDD: ; /* IDA 0x11FDD; and     byte ptr ds:98A3h, 7Fh */
    {uint16_t result=((dl_read8(m,m->ds,39075))&(127))&255;ec_logic(m,result,8);dl_write8(m,m->ds,39075,result);}
L_11FE2: ; /* IDA 0x11FE2; sub     cx, ds:9882h */
    m->cx=(uint16_t)(ec_math(m,m->cx,dl_read16(m,m->ds,39042),16,0,1));
L_11FE6: ; /* IDA 0x11FE6; jb      short loc_11FFB */
    if(m->flags&1) goto L_11FFB;
L_11FE8: ; /* IDA 0x11FE8; cmp     cx, 27Fh */
    ec_cmp(m,m->cx,639,16);
L_11FEC: ; /* IDA 0x11FEC; jb      short loc_12001 */
    if(m->flags&1) goto L_12001;
L_11FEE: ; /* IDA 0x11FEE; sub     cx, 27Fh */
    m->cx=(uint16_t)(ec_math(m,m->cx,639,16,0,1));
L_11FF2: ; /* IDA 0x11FF2; add     ds:9882h, cx */
    dl_write16(m,m->ds,39042,ec_math(m,dl_read16(m,m->ds,39042),m->cx,16,0,0));
L_11FF6: ; /* IDA 0x11FF6; mov     cx, 27Fh */
    m->cx=(uint16_t)(639);
L_11FF9: ; /* IDA 0x11FF9; jmp     short loc_12001 */
    goto L_12001;
L_11FFB: ; /* IDA 0x11FFB; add     ds:9882h, cx */
    dl_write16(m,m->ds,39042,ec_math(m,dl_read16(m,m->ds,39042),m->cx,16,0,0));
L_11FFF: ; /* IDA 0x11FFF; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_12001: ; /* IDA 0x12001; sub     dx, ds:9884h */
    m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ds,39044),16,0,1));
L_12005: ; /* IDA 0x12005; jb      short loc_1201A */
    if(m->flags&1) goto L_1201A;
L_12007: ; /* IDA 0x12007; cmp     dx, 18Fh */
    ec_cmp(m,m->dx,399,16);
L_1200B: ; /* IDA 0x1200B; jb      short loc_12020 */
    if(m->flags&1) goto L_12020;
L_1200D: ; /* IDA 0x1200D; sub     dx, 18Fh */
    m->dx=(uint16_t)(ec_math(m,m->dx,399,16,0,1));
L_12011: ; /* IDA 0x12011; add     ds:9884h, dx */
    dl_write16(m,m->ds,39044,ec_math(m,dl_read16(m,m->ds,39044),m->dx,16,0,0));
L_12015: ; /* IDA 0x12015; mov     dx, 18Fh */
    m->dx=(uint16_t)(399);
L_12018: ; /* IDA 0x12018; jmp     short loc_12020 */
    goto L_12020;
L_1201A: ; /* IDA 0x1201A; add     ds:9884h, dx */
    dl_write16(m,m->ds,39044,ec_math(m,dl_read16(m,m->ds,39044),m->dx,16,0,0));
L_1201E: ; /* IDA 0x1201E; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_12020: ; /* IDA 0x12020; mov     ds:9886h, cx */
    dl_write16(m,m->ds,39046,m->cx);
L_12024: ; /* IDA 0x12024; mov     ds:9888h, dx */
    dl_write16(m,m->ds,39048,m->dx);
L_12028: ; /* IDA 0x12028; mov     ax, cx */
    m->ax=(uint16_t)(m->cx);
L_1202A: ; /* IDA 0x1202A; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1202C: ; /* IDA 0x1202C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1202E: ; /* IDA 0x1202E; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12030: ; /* IDA 0x12030; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12032: ; /* IDA 0x12032; mov     ds:9896h, ax */
    dl_write16(m,m->ds,39062,m->ax);
L_12035: ; /* IDA 0x12035; mov     ax, dx */
    m->ax=(uint16_t)(m->dx);
L_12037: ; /* IDA 0x12037; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12039: ; /* IDA 0x12039; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1203B: ; /* IDA 0x1203B; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1203D: ; /* IDA 0x1203D; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1203F: ; /* IDA 0x1203F; mov     ds:9898h, ax */
    dl_write16(m,m->ds,39064,m->ax);
L_12042: ; /* IDA 0x12042; mov     ax, ds:9882h */
    m->ax=(uint16_t)(dl_read16(m,m->ds,39042));
L_12045: ; /* IDA 0x12045; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12047: ; /* IDA 0x12047; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12049: ; /* IDA 0x12049; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1204B: ; /* IDA 0x1204B; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1204D: ; /* IDA 0x1204D; mov     ds:988Eh, ax */
    dl_write16(m,m->ds,39054,m->ax);
L_12050: ; /* IDA 0x12050; add     ax, ds:9896h */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->ds,39062),16,0,0));
L_12054: ; /* IDA 0x12054; mov     ds:989Ah, ax */
    dl_write16(m,m->ds,39066,m->ax);
L_12057: ; /* IDA 0x12057; mov     ax, ds:9884h */
    m->ax=(uint16_t)(dl_read16(m,m->ds,39044));
L_1205A: ; /* IDA 0x1205A; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1205C: ; /* IDA 0x1205C; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1205E: ; /* IDA 0x1205E; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12060: ; /* IDA 0x12060; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_12062: ; /* IDA 0x12062; mov     ds:9890h, ax */
    dl_write16(m,m->ds,39056,m->ax);
L_12065: ; /* IDA 0x12065; add     ax, ds:9898h */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->ds,39064),16,0,0));
L_12069: ; /* IDA 0x12069; sub     ax, 2 */
#if KI_MARCH_MUTATION == 14
    dl_write16(m,m->ds,0x989c,m->ax);
#else
    m->ax=(uint16_t)(ec_math(m,m->ax,2,16,0,1));
#endif
L_1206C: ; /* IDA 0x1206C; jnb     short loc_12070 */
    if(!(m->flags&1)) goto L_12070;
L_1206E: ; /* IDA 0x1206E; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_12070: ; /* IDA 0x12070; mov     ds:989Ch, ax */
    dl_write16(m,m->ds,39068,m->ax);
L_12073: ; /* IDA 0x12073; mov     bx, bp */
    m->bx=(uint16_t)(m->bp);
L_12075: ; /* IDA 0x12075; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_12076: ; /* IDA 0x12076; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_12077: ; /* IDA 0x12077; retn */
    return;
}
/* Original sub_12151; IDA 0x12151. */
static void mch_body_12151(KiMachine16 *m,const KiEconomyHooks *h) {
L_12151: ; /* IDA 0x12151; mov     cs:word_19882, ax */
    dl_write16(m,m->cs,39042,m->ax);
L_12155: ; /* IDA 0x12155; mov     cs:word_19884, cx */
    dl_write16(m,m->cs,39044,m->cx);
L_1215A: ; /* IDA 0x1215A; push    dx */
    ec_push(m,m->dx);
L_1215B: ; /* IDA 0x1215B; push    bx */
    ec_push(m,m->bx);
L_1215C: ; /* IDA 0x1215C; mov     cl, 4 */
    vg_cl(m,(uint8_t)(4));
L_1215E: ; /* IDA 0x1215E; shl     bx, cl */
    m->bx=(uint16_t)(dl_shift(m,m->bx,(uint8_t)m->cx,16,4));
L_12160: ; /* IDA 0x12160; shl     dx, cl */
    m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,4));
L_12162: ; /* IDA 0x12162; add     bx, 20h ; ' ' */
    m->bx=(uint16_t)(ec_math(m,m->bx,32,16,0,0));
L_12165: ; /* IDA 0x12165; mov     cx, dx */
    m->cx=(uint16_t)(m->dx);
L_12167: ; /* IDA 0x12167; mov     dx, bx */
    m->dx=(uint16_t)(m->bx);
L_12169: ; /* IDA 0x12169; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_1216C: ; /* IDA 0x1216C; call    sub_20000 */
    mx_far(m,0x216c,0x2171,h);
L_12171: ; /* IDA 0x12171; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_12172: ; /* IDA 0x12172; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_12173: ; /* IDA 0x12173; mov     ax, dx */
    m->ax=(uint16_t)(m->dx);
L_12175: ; /* IDA 0x12175; sub     ax, cs:word_19882 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,39042),16,0,1));
L_1217A: ; /* IDA 0x1217A; jnb     short loc_1217E */
    if(!(m->flags&1)) goto L_1217E;
L_1217C: ; /* IDA 0x1217C; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_1217E: ; /* IDA 0x1217E; cmp     ax, 158h */
    ec_cmp(m,m->ax,344,16);
L_12181: ; /* IDA 0x12181; jbe     short loc_12186 */
    if(m->flags&(1|0x40)) goto L_12186;
L_12183: ; /* IDA 0x12183; mov     ax, 158h */
    m->ax=(uint16_t)(344);
L_12186: ; /* IDA 0x12186; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_12188: ; /* IDA 0x12188; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1218A: ; /* IDA 0x1218A; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1218C: ; /* IDA 0x1218C; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1218E: ; /* IDA 0x1218E; mov     cs:word_19882, ax */
    dl_write16(m,m->cs,39042,m->ax);
L_12192: ; /* IDA 0x12192; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_12194: ; /* IDA 0x12194; sub     ax, cs:word_19884 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,39044),16,0,1));
L_12199: ; /* IDA 0x12199; jnb     short loc_1219D */
    if(!(m->flags&1)) goto L_1219D;
L_1219B: ; /* IDA 0x1219B; xor     ax, ax */
    {uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
L_1219D: ; /* IDA 0x1219D; cmp     ax, 0E9h */
    ec_cmp(m,m->ax,233,16);
L_121A0: ; /* IDA 0x121A0; jbe     short loc_121A5 */
    if(m->flags&(1|0x40)) goto L_121A5;
L_121A2: ; /* IDA 0x121A2; mov     ax, 0E9h */
    m->ax=(uint16_t)(233);
L_121A5: ; /* IDA 0x121A5; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_121A7: ; /* IDA 0x121A7; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_121A9: ; /* IDA 0x121A9; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_121AB: ; /* IDA 0x121AB; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_121AD: ; /* IDA 0x121AD; mov     cs:word_19884, ax */
    dl_write16(m,m->cs,39044,m->ax);
L_121B1: ; /* IDA 0x121B1; retn */
    return;
}
/* Original sub_121B2; IDA 0x121B2. */
static void mch_body_121B2(KiMachine16 *m,const KiEconomyHooks *h) {
L_121B2: ; /* IDA 0x121B2; push    ax */
    ec_push(m,m->ax);
L_121B3: ; /* IDA 0x121B3; push    cx */
    ec_push(m,m->cx);
L_121B4: ; /* IDA 0x121B4; mov     ax, 14h */
    m->ax=(uint16_t)(20);
L_121B7: ; /* IDA 0x121B7; mov     cx, 0Ch */
    m->cx=(uint16_t)(12);
L_121BA: ; /* IDA 0x121BA; call    sub_12151 */
    mch_call(m,8529,8637,h);
L_121BD: ; /* IDA 0x121BD; push    dx */
    ec_push(m,m->dx);
L_121BE: ; /* IDA 0x121BE; mov     ax, cs:word_19882 */
    m->ax=(uint16_t)(dl_read16(m,m->cs,39042));
L_121C2: ; /* IDA 0x121C2; add     ax, cs:word_19886 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,39046),16,0,0));
L_121C7: ; /* IDA 0x121C7; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_121C9: ; /* IDA 0x121C9; mov     ax, cs:word_19884 */
    m->ax=(uint16_t)(dl_read16(m,m->cs,39044));
L_121CD: ; /* IDA 0x121CD; add     ax, cs:word_19888 */
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,39048),16,0,0));
L_121D2: ; /* IDA 0x121D2; mov     dx, ax */
    m->dx=(uint16_t)(m->ax);
L_121D4: ; /* IDA 0x121D4; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_121D7: ; /* IDA 0x121D7; call    sub_20000 */
    mx_far(m,0x21d7,0x21dc,h);
L_121DC: ; /* IDA 0x121DC; mov     cs:word_1988A, 0FFFFh */
    dl_write16(m,m->cs,39050,65535);
L_121E3: ; /* IDA 0x121E3; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_121E4: ; /* IDA 0x121E4; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_121E5: ; /* IDA 0x121E5; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_121E6: ; /* IDA 0x121E6; retn */
    return;
}
/* Original sub_14325; IDA 0x14325. */
static void mch_body_14325(KiMachine16 *m,const KiEconomyHooks *h) {
L_14325: ; /* IDA 0x14325; push    ds */
    ec_push(m,m->ds);
L_14326: ; /* IDA 0x14326; push    es */
    ec_push(m,m->es);
L_14327: ; /* IDA 0x14327; push    si */
    ec_push(m,m->si);
L_14328: ; /* IDA 0x14328; mov     ah, [si+20h] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+32))));
L_1432B: ; /* IDA 0x1432B; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_1432D: ; /* IDA 0x1432D; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_1432F: ; /* IDA 0x1432F; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14331: ; /* IDA 0x14331; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14333: ; /* IDA 0x14333; add     ax, 840h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2112,16,0,0));
L_14336: ; /* IDA 0x14336; mov     bl, [si+23h] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+35))));
L_14339: ; /* IDA 0x14339; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_1433B: ; /* IDA 0x1433B; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_1433D: ; /* IDA 0x1433D; cmp     bx, 10h */
    ec_cmp(m,m->bx,16,16);
L_14340: ; /* IDA 0x14340; jnb     short loc_1434F */
    if(!(m->flags&1)) goto L_1434F;
L_14342: ; /* IDA 0x14342; mov     cl, cs:byte_10CFF */
    vg_cl(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_14347: ; /* IDA 0x14347; cmp     cl, [si+1] */
    ec_cmp(m,(uint8_t)m->cx,dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_1434A: ; /* IDA 0x1434A; jz      short loc_1434F */
    if(m->flags&0x40) goto L_1434F;
L_1434C: ; /* IDA 0x1434C; add     bx, 8 */
#if KI_MARCH_MUTATION == 6
    m->bx=ec_math(m,m->bx,4,16,0,0);
#else
    m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,0));
#endif
L_1434F: ; /* IDA 0x1434F; call    cs:funcs_1434F[bx] */
    mch_call(m,dl_read16(m,m->cs,(uint16_t)(m->bx+17240)),17236,h);
L_14354: ; /* IDA 0x14354; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_14355: ; /* IDA 0x14355; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_14356: ; /* IDA 0x14356; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_14357: ; /* IDA 0x14357; retn */
    return;
}
/* Original sub_14370; IDA 0x14370. */
static void mch_body_14370(KiMachine16 *m,const KiEconomyHooks *h) {
L_14370: ; /* IDA 0x14370; mov     byte ptr [si+23h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),0);
L_14374: ; /* IDA 0x14374; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_14376: ; /* IDA 0x14376; call    sub_14548 */
    mch_call(m,17736,17273,h);
L_14379: ; /* IDA 0x14379; jb      short loc_1437C */
    if(m->flags&1) goto L_1437C;
L_1437B: ; /* IDA 0x1437B; retn */
    return;
L_1437C: ; /* IDA 0x1437C; cmp     word ptr [si+4], 258h */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),600,16);
L_14381: ; /* IDA 0x14381; jnb     short locret_1439C */
    if(!(m->flags&1)) goto L_1439C;
L_14383: ; /* IDA 0x14383; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_14385: ; /* IDA 0x14385; sub     ax, 840h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2112,16,0,1));
L_14388: ; /* IDA 0x14388; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1438A: ; /* IDA 0x1438A; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1438C: ; /* IDA 0x1438C; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_1438E: ; /* IDA 0x1438E; mov     bx, cs:word_10CFD */
    m->bx=(uint16_t)(dl_read16(m,m->cs,3325));
L_14393: ; /* IDA 0x14393; cmp     ah, [bx+3] */
    ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->ds,(uint16_t)(m->bx+3)),8);
L_14396: ; /* IDA 0x14396; jnz     short locret_1439C */
    if(!(m->flags&0x40)) goto L_1439C;
L_14398: ; /* IDA 0x14398; mov     byte ptr [si+23h], 9 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),9);
L_1439C: ; /* IDA 0x1439C; retn */
    return;
}
/* Original sub_1439D; IDA 0x1439D. */
static void mch_body_1439D(KiMachine16 *m,const KiEconomyHooks *h) {
L_1439D: ; /* IDA 0x1439D; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_1439F: ; /* IDA 0x1439F; call    sub_14548 */
    mch_call(m,17736,17314,h);
L_143A2: ; /* IDA 0x143A2; jb      short loc_143A5 */
    if(m->flags&1) goto L_143A5;
L_143A4: ; /* IDA 0x143A4; retn */
    return;
L_143A5: ; /* IDA 0x143A5; test    byte ptr [bx], 40h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(64))&255;ec_logic(m,result,8);}
L_143A8: ; /* IDA 0x143A8; jnz     short locret_143AE */
    if(!(m->flags&0x40)) goto L_143AE;
L_143AA: ; /* IDA 0x143AA; mov     byte ptr [si+23h], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),1);
L_143AE: ; /* IDA 0x143AE; retn */
    return;
}
/* Original sub_143AF; IDA 0x143AF. */
static void mch_body_143AF(KiMachine16 *m,const KiEconomyHooks *h) {
L_143AF: ; /* IDA 0x143AF; cmp     word ptr [si+0Eh], 800h */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+14)),2048,16);
L_143B4: ; /* IDA 0x143B4; jb      short loc_143BB */
    if(m->flags&1) goto L_143BB;
L_143B6: ; /* IDA 0x143B6; mov     byte ptr [si+23h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),0);
L_143BA: ; /* IDA 0x143BA; retn */
    return;
L_143BB: ; /* IDA 0x143BB; cmp     word ptr [si+4], 12Ch */
#if KI_MARCH_MUTATION == 7
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),299,16);
#else
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),300,16);
#endif
L_143C0: ; /* IDA 0x143C0; ja      short loc_143C7 */
    if(!(m->flags&(1|0x40))) goto L_143C7;
L_143C2: ; /* IDA 0x143C2; mov     byte ptr [si+23h], 0Ah */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),10);
L_143C6: ; /* IDA 0x143C6; retn */
    return;
L_143C7: ; /* IDA 0x143C7; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_143C9: ; /* IDA 0x143C9; test    byte ptr [bx], 40h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(64))&255;ec_logic(m,result,8);}
L_143CC: ; /* IDA 0x143CC; jnz     short loc_143B6 */
    if(!(m->flags&0x40)) goto L_143B6;
L_143CE: ; /* IDA 0x143CE; cmp     byte ptr [bx], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx)),128,8);
L_143D1: ; /* IDA 0x143D1; jb      short loc_143D9 */
    if(m->flags&1) goto L_143D9;
L_143D3: ; /* IDA 0x143D3; cmp     byte ptr [di+18h], 2 */
#if KI_MARCH_MUTATION == 15
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+24)),2,8);
#else
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+24)),2,8);
#endif
L_143D7: ; /* IDA 0x143D7; jbe     short loc_143E8 */
    if(m->flags&(1|0x40)) goto L_143E8;
L_143D9: ; /* IDA 0x143D9; call    sub_1ECE0 */
    mch_call(m,60640,17372,h);
L_143DC: ; /* IDA 0x143DC; and     al, 7 */
    {uint16_t result=(((uint8_t)m->ax)&(7))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_143DE: ; /* IDA 0x143DE; inc     al */
#if KI_MARCH_MUTATION == 8
    ; /* lost random +1 */
#else
    vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
#endif
L_143E0: ; /* IDA 0x143E0; mov     [si+0Bh], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+11),(uint8_t)m->ax);
L_143E3: ; /* IDA 0x143E3; mov     byte ptr [si+23h], 2 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),2);
L_143E7: ; /* IDA 0x143E7; retn */
    return;
L_143E8: ; /* IDA 0x143E8; cmp     word ptr [si+4], 258h */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),600,16);
L_143ED: ; /* IDA 0x143ED; jnb     short locret_1440E */
    if(!(m->flags&1)) goto L_1440E;
L_143EF: ; /* IDA 0x143EF; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_143F2: ; /* IDA 0x143F2; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_143F4: ; /* IDA 0x143F4; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_143F6: ; /* IDA 0x143F6; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_143F8: ; /* IDA 0x143F8; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_143FA: ; /* IDA 0x143FA; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_143FC: ; /* IDA 0x143FC; sub     ax, 840h */
    m->ax=(uint16_t)(ec_math(m,m->ax,2112,16,0,1));
L_143FF: ; /* IDA 0x143FF; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14401: ; /* IDA 0x14401; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14403: ; /* IDA 0x14403; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_14405: ; /* IDA 0x14405; cmp     ah, [di+3] */
    ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->ds,(uint16_t)(m->di+3)),8);
L_14408: ; /* IDA 0x14408; jnz     short locret_1440E */
    if(!(m->flags&0x40)) goto L_1440E;
L_1440A: ; /* IDA 0x1440A; mov     byte ptr [si+23h], 9 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),9);
L_1440E: ; /* IDA 0x1440E; retn */
    return;
}
/* Original sub_1440F; IDA 0x1440F. */
static void mch_body_1440F(KiMachine16 *m,const KiEconomyHooks *h) {
L_1440F: ; /* IDA 0x1440F; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_14411: ; /* IDA 0x14411; test    byte ptr [di], 40h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->di)))&(64))&255;ec_logic(m,result,8);}
L_14414: ; /* IDA 0x14414; jnz     short loc_14421 */
    if(!(m->flags&0x40)) goto L_14421;
L_14416: ; /* IDA 0x14416; cmp     byte ptr [di], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di)),128,8);
L_14419: ; /* IDA 0x14419; jb      short loc_14426 */
    if(m->flags&1) goto L_14426;
L_1441B: ; /* IDA 0x1441B; cmp     byte ptr [di+18h], 1 */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+24)),1,8);
L_1441F: ; /* IDA 0x1441F; ja      short loc_14426 */
    if(!(m->flags&(1|0x40))) goto L_14426;
L_14421: ; /* IDA 0x14421; mov     byte ptr [si+23h], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),1);
L_14425: ; /* IDA 0x14425; retn */
    return;
L_14426: ; /* IDA 0x14426; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14429: ; /* IDA 0x14429; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_1442B: ; /* IDA 0x1442B; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1442D: ; /* IDA 0x1442D; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1442F: ; /* IDA 0x1442F; test    byte ptr [bx], 40h */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->bx)))&(64))&255;ec_logic(m,result,8);}
L_14432: ; /* IDA 0x14432; jnz     short loc_1443D */
    if(!(m->flags&0x40)) goto L_1443D;
L_14434: ; /* IDA 0x14434; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_14436: ; /* IDA 0x14436; xchg    al, [bx+17h] */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->ds,(uint16_t)(m->bx+23));vg_al(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->bx+23),a);}
L_14439: ; /* IDA 0x14439; cmp     al, 0FFh */
#if KI_MARCH_MUTATION == 9
    ec_cmp(m,(uint8_t)m->ax,0,8);
#else
    ec_cmp(m,(uint8_t)m->ax,255,8);
#endif
L_1443B: ; /* IDA 0x1443B; jnz     short loc_14446 */
    if(!(m->flags&0x40)) goto L_14446;
L_1443D: ; /* IDA 0x1443D; mov     al, 0FFh */
    vg_al(m,(uint8_t)(255));
L_1443F: ; /* IDA 0x1443F; xchg    al, [bx+16h] */
    {uint16_t a=(uint8_t)m->ax,b=dl_read8(m,m->ds,(uint16_t)(m->bx+22));vg_al(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->bx+22),a);}
L_14442: ; /* IDA 0x14442; cmp     al, 0FFh */
    ec_cmp(m,(uint8_t)m->ax,255,8);
L_14444: ; /* IDA 0x14444; jz      short loc_14459 */
    if(m->flags&0x40) goto L_14459;
L_14446: ; /* IDA 0x14446; cmp     [si+20h], al */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+32)),(uint8_t)m->ax,8);
L_14449: ; /* IDA 0x14449; jz      short loc_14451 */
    if(m->flags&0x40) goto L_14451;
L_1444B: ; /* IDA 0x1444B; mov     [si+20h], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)m->ax);
L_1444E: ; /* IDA 0x1444E; or      byte ptr [si], 2 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(2))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_14451: ; /* IDA 0x14451; mov     byte ptr [si+23h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),0);
L_14455: ; /* IDA 0x14455; dec     byte ptr [di+18h] */
    dl_write8(m,m->ds,(uint16_t)(m->di+24),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->di+24)),8));
L_14458: ; /* IDA 0x14458; retn */
    return;
L_14459: ; /* IDA 0x14459; cmp     byte ptr [di], 80h */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di)),128,8);
L_1445C: ; /* IDA 0x1445C; jnb     short locret_14465 */
    if(!(m->flags&1)) goto L_14465;
L_1445E: ; /* IDA 0x1445E; mov     byte ptr [si+23h], 0Bh */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),11);
L_14462: ; /* IDA 0x14462; dec     byte ptr [di+18h] */
    dl_write8(m,m->ds,(uint16_t)(m->di+24),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->di+24)),8));
L_14465: ; /* IDA 0x14465; retn */
    return;
}
/* Original sub_14466; IDA 0x14466. */
static void mch_body_14466(KiMachine16 *m,const KiEconomyHooks *h) {
L_14466: ; /* IDA 0x14466; mov     al, 1Eh */
    vg_al(m,(uint8_t)(30));
L_14468: ; /* IDA 0x14468; mov     cx, 6 */
    m->cx=(uint16_t)(6);
L_1446B: ; /* IDA 0x1446B; mov     bx, si */
    m->bx=(uint16_t)(m->si);
L_1446D: ; /* IDA 0x1446D; add     bx, 29h ; ')' */
    m->bx=(uint16_t)(ec_math(m,m->bx,41,16,0,0));
L_14470: ; /* IDA 0x14470; cmp     al, [bx] */
#if KI_MARCH_MUTATION == 10
    ec_cmp(m,29,dl_read8(m,m->ds,m->bx),8);
#else
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx)),8);
#endif
L_14472: ; /* IDA 0x14472; ja      short loc_1447E */
    if(!(m->flags&(1|0x40))) goto L_1447E;
L_14474: ; /* IDA 0x14474; add     bx, 4 */
    m->bx=(uint16_t)(ec_math(m,m->bx,4,16,0,0));
L_14477: ; /* IDA 0x14477; loop    loc_14470 */
    if(--m->cx) goto L_14470;
L_14479: ; /* IDA 0x14479; mov     byte ptr [si+23h], 8 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),8);
L_1447D: ; /* IDA 0x1447D; retn */
    return;
L_1447E: ; /* IDA 0x1447E; mov     byte ptr [si+23h], 0Bh */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),11);
L_14482: ; /* IDA 0x14482; retn */
    return;
}
/* Original sub_14483; IDA 0x14483. */
static void mch_body_14483(KiMachine16 *m,const KiEconomyHooks *h) {
L_14483: ; /* IDA 0x14483; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14486: ; /* IDA 0x14486; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_14488: ; /* IDA 0x14488; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1448A: ; /* IDA 0x1448A; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1448C: ; /* IDA 0x1448C; mov     al, [bx+1Dh] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+29))));
L_1448F: ; /* IDA 0x1448F; cmp     [si+6], al */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+6)),(uint8_t)m->ax,8);
L_14492: ; /* IDA 0x14492; jb      short locret_14498 */
    if(m->flags&1) goto L_14498;
L_14494: ; /* IDA 0x14494; mov     byte ptr [si+23h], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),1);
L_14498: ; /* IDA 0x14498; retn */
    return;
}
/* Original sub_14499; IDA 0x14499. */
static void mch_body_14499(KiMachine16 *m,const KiEconomyHooks *h) {
L_14499: ; /* IDA 0x14499; call    sub_1461D */
    mch_call(m,17949,17564,h);
L_1449C: ; /* IDA 0x1449C; call    sub_16FD2 */
    mch_call(m,28626,17567,h);
L_1449F: ; /* IDA 0x1449F; mov     byte ptr [si+23h], 3 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),3);
L_144A3: ; /* IDA 0x144A3; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_144A5: ; /* IDA 0x144A5; call    sub_15E80 */
    mch_call(m,24192,17576,h);
L_144A8: ; /* IDA 0x144A8; retn */
    return;
}
/* Original sub_144A9; IDA 0x144A9. */
static void mch_body_144A9(KiMachine16 *m,const KiEconomyHooks *h) {
L_144A9: ; /* IDA 0x144A9; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_144AC: ; /* IDA 0x144AC; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_144AE: ; /* IDA 0x144AE; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144B0: ; /* IDA 0x144B0; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144B2: ; /* IDA 0x144B2; mov     bh, [bx+3] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_144B5: ; /* IDA 0x144B5; cmp     [si+20h], bh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+32)),(uint8_t)(m->bx>>8),8);
L_144B8: ; /* IDA 0x144B8; jz      short loc_144C0 */
    if(m->flags&0x40) goto L_144C0;
L_144BA: ; /* IDA 0x144BA; mov     [si+20h], bh */
    dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)(m->bx>>8));
L_144BD: ; /* IDA 0x144BD; or      byte ptr [si], 2 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(2))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_144C0: ; /* IDA 0x144C0; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_144C2: ; /* IDA 0x144C2; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144C4: ; /* IDA 0x144C4; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144C6: ; /* IDA 0x144C6; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144C8: ; /* IDA 0x144C8; add     bx, 840h */
    m->bx=(uint16_t)(ec_math(m,m->bx,2112,16,0,0));
L_144CC: ; /* IDA 0x144CC; call    sub_14548 */
    mch_call(m,17736,17615,h);
L_144CF: ; /* IDA 0x144CF; jnb     short locret_144D5 */
    if(!(m->flags&1)) goto L_144D5;
L_144D1: ; /* IDA 0x144D1; mov     byte ptr [si+23h], 9 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),9);
L_144D5: ; /* IDA 0x144D5; retn */
    return;
}
/* Original sub_144D6; IDA 0x144D6. */
static void mch_body_144D6(KiMachine16 *m,const KiEconomyHooks *h) {
L_144D6: ; /* IDA 0x144D6; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_144D9: ; /* IDA 0x144D9; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_144DB: ; /* IDA 0x144DB; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144DD: ; /* IDA 0x144DD; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144DF: ; /* IDA 0x144DF; mov     bh, [bx+3] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_144E2: ; /* IDA 0x144E2; cmp     [si+20h], bh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+32)),(uint8_t)(m->bx>>8),8);
L_144E5: ; /* IDA 0x144E5; jz      short loc_144ED */
    if(m->flags&0x40) goto L_144ED;
L_144E7: ; /* IDA 0x144E7; mov     [si+20h], bh */
    dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)(m->bx>>8));
L_144EA: ; /* IDA 0x144EA; or      byte ptr [si], 2 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(2))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_144ED: ; /* IDA 0x144ED; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_144EF: ; /* IDA 0x144EF; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144F1: ; /* IDA 0x144F1; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144F3: ; /* IDA 0x144F3; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_144F5: ; /* IDA 0x144F5; add     bx, 840h */
    m->bx=(uint16_t)(ec_math(m,m->bx,2112,16,0,0));
L_144F9: ; /* IDA 0x144F9; call    sub_14548 */
    mch_call(m,17736,17660,h);
L_144FC: ; /* IDA 0x144FC; jnb     short locret_14501 */
    if(!(m->flags&1)) goto L_14501;
L_144FE: ; /* IDA 0x144FE; call    sub_1463E */
    mch_call(m,17982,17665,h);
L_14501: ; /* IDA 0x14501; retn */
    return;
}
/* Original sub_14548; IDA 0x14548. */
static void mch_body_14548(KiMachine16 *m,const KiEconomyHooks *h) {
L_14548: ; /* IDA 0x14548; push    bx */
    ec_push(m,m->bx);
L_14549: ; /* IDA 0x14549; mov     ax, [bx+8] */
    m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+8)));
L_1454C: ; /* IDA 0x1454C; mov     dx, [bx+0Ah] */
    m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+10)));
L_1454F: ; /* IDA 0x1454F; sub     bx, 840h */
    m->bx=(uint16_t)(ec_math(m,m->bx,2112,16,0,1));
L_14553: ; /* IDA 0x14553; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14555: ; /* IDA 0x14555; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14557: ; /* IDA 0x14557; cmp     [si+10h], ax */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+16)),m->ax,16);
L_1455A: ; /* IDA 0x1455A; jnz     short loc_14569 */
    if(!(m->flags&0x40)) goto L_14569;
L_1455C: ; /* IDA 0x1455C; cmp     [si+12h], dx */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+18)),m->dx,16);
L_1455F: ; /* IDA 0x1455F; jnz     short loc_14569 */
    if(!(m->flags&0x40)) goto L_14569;
L_14561: ; /* IDA 0x14561; cmp     [si+0Eh], bx */
#if KI_MARCH_MUTATION == 16
    ; /* lost node comparison */
#else
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+14)),m->bx,16);
#endif
L_14564: ; /* IDA 0x14564; jnz     short loc_14569 */
    if(!(m->flags&0x40)) goto L_14569;
L_14566: ; /* IDA 0x14566; stc */
    m->flags|=1;
L_14567: ; /* IDA 0x14567; jmp     short loc_1456A */
    goto L_1456A;
L_14569: ; /* IDA 0x14569; clc */
    m->flags&=(uint16_t)~1u;
L_1456A: ; /* IDA 0x1456A; mov     [si+16h], ax */
    dl_write16(m,m->ds,(uint16_t)(m->si+22),m->ax);
L_1456D: ; /* IDA 0x1456D; mov     [si+18h], dx */
    dl_write16(m,m->ds,(uint16_t)(m->si+24),m->dx);
L_14570: ; /* IDA 0x14570; mov     [si+14h], bx */
    dl_write16(m,m->ds,(uint16_t)(m->si+20),m->bx);
L_14573: ; /* IDA 0x14573; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14574: ; /* IDA 0x14574; retn */
    return;
}
/* Original sub_1463E; IDA 0x1463E. */
static void mch_body_1463E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1463E: ; /* IDA 0x1463E; call    sub_14651 */
    mch_call(m,18001,17985,h);
L_14641: ; /* IDA 0x14641; mov     al, [si+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_14644: ; /* IDA 0x14644; cmp     al, cs:byte_10CFF */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
L_14649: ; /* IDA 0x14649; jnz     short locret_14650 */
    if(!(m->flags&0x40)) goto L_14650;
L_1464B: ; /* IDA 0x1464B; mov     al, 8 */
    vg_al(m,(uint8_t)(8));
L_1464D: ; /* IDA 0x1464D; call    sub_15E80 */
    mch_call(m,24192,18000,h);
L_14650: ; /* IDA 0x14650; retn */
    return;
}
/* Original sub_14651; IDA 0x14651. */
static void mch_body_14651(KiMachine16 *m,const KiEconomyHooks *h) {
L_14651: ; /* IDA 0x14651; push    es */
    ec_push(m,m->es);
L_14652: ; /* IDA 0x14652; push    ax */
    ec_push(m,m->ax);
L_14653: ; /* IDA 0x14653; push    di */
    ec_push(m,m->di);
L_14654: ; /* IDA 0x14654; mov     ax, ds */
    m->ax=(uint16_t)(m->ds);
L_14656: ; /* IDA 0x14656; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_14658: ; /* IDA 0x14658; call    sub_14689 */
    mch_call(m,18057,18011,h);
L_1465B: ; /* IDA 0x1465B; mov     ah, [si+1] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_1465E: ; /* IDA 0x1465E; xor     al, al */
    {uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_14660: ; /* IDA 0x14660; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14662: ; /* IDA 0x14662; shr     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
L_14664: ; /* IDA 0x14664; mov     di, ax */
    m->di=(uint16_t)(m->ax);
L_14666: ; /* IDA 0x14666; add     ax, 4 */
    m->ax=(uint16_t)(ec_math(m,m->ax,4,16,0,0));
L_14669: ; /* IDA 0x14669; add     di, 4 */
    m->di=(uint16_t)(ec_math(m,m->di,4,16,0,0));
L_1466C: ; /* IDA 0x1466C; call    sub_14717 */
    mch_call(m,18199,18031,h);
L_1466F: ; /* IDA 0x1466F; mov     byte ptr [si], 0 */
#if KI_MARCH_MUTATION == 11
    dl_write16(m,m->ds,m->si,0);
#else
    dl_write8(m,m->ds,(uint16_t)(m->si),0);
#endif
L_14672: ; /* IDA 0x14672; mov     di, si */
    m->di=(uint16_t)(m->si);
L_14674: ; /* IDA 0x14674; sub     di, 2240h */
    m->di=(uint16_t)(ec_math(m,m->di,8768,16,0,1));
L_14678: ; /* IDA 0x14678; shr     di, 1 */
    m->di=(uint16_t)(dl_shift(m,m->di,1,16,5));
L_1467A: ; /* IDA 0x1467A; mov     byte ptr [di+4257h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->di+16983),0);
L_1467F: ; /* IDA 0x1467F; les     di, [si+1Ah] */
    {uint16_t a=(uint16_t)(m->si+26),s=m->ds;uint16_t v=dl_read16(m,s,a);m->di=(uint16_t)(v);m->es=dl_read16(m,s,(uint16_t)(a+2));}
L_14682: ; /* IDA 0x14682; dec     byte ptr es:[di] */
    dl_write8(m,m->es,(uint16_t)(m->di),st_dec(m,dl_read8(m,m->es,(uint16_t)(m->di)),8));
L_14685: ; /* IDA 0x14685; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_14686: ; /* IDA 0x14686; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_14687: ; /* IDA 0x14687; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_14688: ; /* IDA 0x14688; retn */
    return;
}
/* Original sub_14689; IDA 0x14689. */
static void mch_body_14689(KiMachine16 *m,const KiEconomyHooks *h) {
L_14689: ; /* IDA 0x14689; push    bx */
    ec_push(m,m->bx);
L_1468A: ; /* IDA 0x1468A; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_1468D: ; /* IDA 0x1468D; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_1468F: ; /* IDA 0x1468F; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14691: ; /* IDA 0x14691; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_14693: ; /* IDA 0x14693; dec     byte ptr [bx+14h] */
    dl_write8(m,m->ds,(uint16_t)(m->bx+20),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->bx+20)),8));
L_14696: ; /* IDA 0x14696; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_14697: ; /* IDA 0x14697; retn */
    return;
}
/* Original sub_159A6; IDA 0x159A6. */
static void mch_body_159A6(KiMachine16 *m,const KiEconomyHooks *h) {
L_159A6: ; /* IDA 0x159A6; mov     bl, al */
    gw_bl(m,(uint8_t)((uint8_t)m->ax));
L_159A8: ; /* IDA 0x159A8; xor     bh, bh */
    {uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
L_159AA: ; /* IDA 0x159AA; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_159AC: ; /* IDA 0x159AC; call    word ptr [bx+59D2h] */
    mch_call(m,dl_read16(m,m->ds,(uint16_t)(m->bx+22994)),22960,h);
L_159B0: ; /* IDA 0x159B0; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_159B2: ; /* IDA 0x159B2; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_159B4: ; /* IDA 0x159B4; mov     es, ax */
    m->es=(uint16_t)(m->ax);
L_159B6: ; /* IDA 0x159B6; retn */
    return;
}
/* Original sub_15AB6; IDA 0x15AB6. */
static void mch_body_15AB6(KiMachine16 *m,const KiEconomyHooks *h) {
L_15AB6: ; /* IDA 0x15AB6; sub     cx, 1B8h */
    m->cx=(uint16_t)(ec_math(m,m->cx,440,16,0,1));
L_15ABA: ; /* IDA 0x15ABA; jnb     short loc_15ABE */
    if(!(m->flags&1)) goto L_15ABE;
L_15ABC: ; /* IDA 0x15ABC; xor     cx, cx */
    {uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
L_15ABE: ; /* IDA 0x15ABE; sub     dx, 28h ; '(' */
    m->dx=(uint16_t)(ec_math(m,m->dx,40,16,0,1));
L_15AC1: ; /* IDA 0x15AC1; jnb     short loc_15AC5 */
    if(!(m->flags&1)) goto L_15AC5;
L_15AC3: ; /* IDA 0x15AC3; xor     dx, dx */
    {uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
L_15AC5: ; /* IDA 0x15AC5; shl     cx, 1 */
    m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
L_15AC7: ; /* IDA 0x15AC7; shl     dx, 1 */
    m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
L_15AC9: ; /* IDA 0x15AC9; mov     bx, dx */
    m->bx=(uint16_t)(m->dx);
L_15ACB: ; /* IDA 0x15ACB; mov     dx, cx */
    m->dx=(uint16_t)(m->cx);
L_15ACD: ; /* IDA 0x15ACD; call    sub_121B2 */
    mch_call(m,8626,23248,h);
L_15AD0: ; /* IDA 0x15AD0; retn */
    return;
}
/* Original sub_1703C; IDA 0x1703C. */
static void mch_body_1703C(KiMachine16 *m,const KiEconomyHooks *h) {
L_1703C: ; /* IDA 0x1703C; push    ds */
    ec_push(m,m->ds);
L_1703D: ; /* IDA 0x1703D; push    es */
    ec_push(m,m->es);
L_1703E: ; /* IDA 0x1703E; push    ax */
    ec_push(m,m->ax);
L_1703F: ; /* IDA 0x1703F; push    cx */
    ec_push(m,m->cx);
L_17040: ; /* IDA 0x17040; push    dx */
    ec_push(m,m->dx);
L_17041: ; /* IDA 0x17041; push    si */
    ec_push(m,m->si);
L_17042: ; /* IDA 0x17042; push    di */
    ec_push(m,m->di);
L_17043: ; /* IDA 0x17043; push    bp */
    ec_push(m,m->bp);
L_17044: ; /* IDA 0x17044; mov     ax, cs */
    m->ax=(uint16_t)(m->cs);
L_17046: ; /* IDA 0x17046; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_17048: ; /* IDA 0x17048; call    sub_11F7F */
    mch_call(m,8063,28747,h);
L_1704B: ; /* IDA 0x1704B; push    cx */
    ec_push(m,m->cx);
L_1704C: ; /* IDA 0x1704C; push    dx */
    ec_push(m,m->dx);
L_1704D: ; /* IDA 0x1704D; call    sub_1E453 */
    mch_call(m,58451,28752,h);
L_17050: ; /* IDA 0x17050; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17051: ; /* IDA 0x17051; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_17052: ; /* IDA 0x17052; cmp     bx, 2 */
#if KI_MARCH_MUTATION == 1
    ec_cmp(m,m->bx,3,16);
#else
    ec_cmp(m,m->bx,2,16);
#endif
L_17055: ; /* IDA 0x17055; jz      short loc_17094 */
    if(m->flags&0x40) goto L_17094;
L_17057: ; /* IDA 0x17057; and     al, al */
    {uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
L_17059: ; /* IDA 0x17059; jz      short loc_1706F */
    if(m->flags&0x40) goto L_1706F;
L_1705B: ; /* IDA 0x1705B; cmp     al, 16h */
    ec_cmp(m,(uint8_t)m->ax,22,8);
L_1705D: ; /* IDA 0x1705D; jnz     short loc_17067 */
    if(!(m->flags&0x40)) goto L_17067;
L_1705F: ; /* IDA 0x1705F; cmp     bx, 1 */
    ec_cmp(m,m->bx,1,16);
L_17062: ; /* IDA 0x17062; jnz     short loc_17067 */
    if(!(m->flags&0x40)) goto L_17067;
L_17064: ; /* IDA 0x17064; call    sub_159A6 */
    mch_call(m,22950,28775,h);
L_17067: ; /* IDA 0x17067; mov     cs:byte_198A2, 1 */
    dl_write8(m,m->cs,39074,1);
L_1706D: ; /* IDA 0x1706D; jmp     short loc_1707B */
    goto L_1707B;
L_1706F: ; /* IDA 0x1706F; and     bx, bx */
    {uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_17071: ; /* IDA 0x17071; jz      short loc_17078 */
    if(m->flags&0x40) goto L_17078;
L_17073: ; /* IDA 0x17073; call    sub_1709E */
    mch_call(m,28830,28790,h);
L_17076: ; /* IDA 0x17076; jnb     short loc_17095 */
    if(!(m->flags&1)) goto L_17095;
L_17078: ; /* IDA 0x17078; call    sub_11F30 */
    mch_call(m,7984,28795,h);
L_1707B: ; /* IDA 0x1707B; test    byte_198A6, 4 */
    {uint16_t result=((dl_read8(m,m->ds,39078))&(4))&255;ec_logic(m,result,8);}
L_17080: ; /* IDA 0x17080; jz      short loc_17085 */
    if(m->flags&0x40) goto L_17085;
L_17082: ; /* IDA 0x17082; call    sub_15C58 */
    mch_call(m,23640,28805,h);
L_17085: ; /* IDA 0x17085; call    sub_11C8D */
    mch_call(m,7309,28808,h);
L_17088: ; /* IDA 0x17088; cmp     byte_198A2, 0 */
    ec_cmp(m,dl_read8(m,m->ds,39074),0,8);
L_1708D: ; /* IDA 0x1708D; jz      short loc_17048 */
    if(m->flags&0x40) goto L_17048;
L_1708F: ; /* IDA 0x1708F; call    sub_11F5A */
    mch_call(m,8026,28818,h);
L_17092: ; /* IDA 0x17092; jmp     short loc_17048 */
    goto L_17048;
L_17094: ; /* IDA 0x17094; stc */
    m->flags|=1;
L_17095: ; /* IDA 0x17095; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_17096: ; /* IDA 0x17096; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_17097: ; /* IDA 0x17097; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_17098: ; /* IDA 0x17098; pop     dx */
    m->dx=(uint16_t)(ec_pop(m));
L_17099: ; /* IDA 0x17099; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_1709A: ; /* IDA 0x1709A; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_1709B: ; /* IDA 0x1709B; pop     es */
    m->es=(uint16_t)(ec_pop(m));
L_1709C: ; /* IDA 0x1709C; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_1709D: ; /* IDA 0x1709D; retn */
    return;
}
/* Original sub_1709E; IDA 0x1709E. */
static void mch_body_1709E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1709E: ; /* IDA 0x1709E; push    ds */
    ec_push(m,m->ds);
L_1709F: ; /* IDA 0x1709F; push    ax */
    ec_push(m,m->ax);
L_170A0: ; /* IDA 0x170A0; push    cx */
    ec_push(m,m->cx);
L_170A1: ; /* IDA 0x170A1; mov     ax, cs:word_1989C */
    m->ax=(uint16_t)(dl_read16(m,m->cs,39068));
L_170A5: ; /* IDA 0x170A5; mov     cx, ax */
    m->cx=(uint16_t)(m->ax);
L_170A7: ; /* IDA 0x170A7; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_170A9: ; /* IDA 0x170A9; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_170AB: ; /* IDA 0x170AB; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_170AD: ; /* IDA 0x170AD; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_170AF: ; /* IDA 0x170AF; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_170B1: ; /* IDA 0x170B1; add     ax, bx */
    m->ax=(uint16_t)(ec_math(m,m->ax,m->bx,16,0,0));
L_170B3: ; /* IDA 0x170B3; add     ax, cs:word_10D44 */
#if KI_MARCH_MUTATION == 13
    m->ax=ec_math(m,m->ax,0,16,0,0);
#else
    m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,3396),16,0,0));
#endif
L_170B8: ; /* IDA 0x170B8; mov     ds, ax */
    m->ds=(uint16_t)(m->ax);
L_170BA: ; /* IDA 0x170BA; mov     bx, cs:word_1989A */
    m->bx=(uint16_t)(dl_read16(m,m->cs,39066));
L_170BF: ; /* IDA 0x170BF; mov     al, [bx] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx))));
L_170C1: ; /* IDA 0x170C1; cmp     al, 0CBh */
    ec_cmp(m,(uint8_t)m->ax,203,8);
L_170C3: ; /* IDA 0x170C3; jb      short loc_170E6 */
    if(m->flags&1) goto L_170E6;
L_170C5: ; /* IDA 0x170C5; cmp     al, 0D4h */
#if KI_MARCH_MUTATION == 2
    ec_cmp(m,(uint8_t)m->ax,0xd3,8);
#else
    ec_cmp(m,(uint8_t)m->ax,212,8);
#endif
L_170C7: ; /* IDA 0x170C7; jnb     short loc_170E6 */
    if(!(m->flags&1)) goto L_170E6;
L_170C9: ; /* IDA 0x170C9; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_170CE: ; /* IDA 0x170CE; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_170D0: ; /* IDA 0x170D0; xor     bx, bx */
    {uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
L_170D2: ; /* IDA 0x170D2; cmp     [bx+848h], ax */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->bx+2120)),m->ax,16);
L_170D6: ; /* IDA 0x170D6; jnz     short loc_170DE */
    if(!(m->flags&0x40)) goto L_170DE;
L_170D8: ; /* IDA 0x170D8; cmp     [bx+84Ah], cx */
    ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->bx+2122)),m->cx,16);
L_170DC: ; /* IDA 0x170DC; jz      short loc_170E3 */
    if(m->flags&0x40) goto L_170E3;
L_170DE: ; /* IDA 0x170DE; add     bx, 20h ; ' ' */
    m->bx=(uint16_t)(ec_math(m,m->bx,32,16,0,0));
L_170E1: ; /* IDA 0x170E1; jmp     short loc_170D2 */
    goto L_170D2;
L_170E3: ; /* IDA 0x170E3; clc */
    m->flags&=(uint16_t)~1u;
L_170E4: ; /* IDA 0x170E4; jmp     short loc_170E7 */
    goto L_170E7;
L_170E6: ; /* IDA 0x170E6; stc */
    m->flags|=1;
L_170E7: ; /* IDA 0x170E7; pop     cx */
    m->cx=(uint16_t)(ec_pop(m));
L_170E8: ; /* IDA 0x170E8; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_170E9: ; /* IDA 0x170E9; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_170EA: ; /* IDA 0x170EA; retn */
    return;
}
/* Original sub_17F90; IDA 0x17F90. */
static void mch_body_17F90(KiMachine16 *m,const KiEconomyHooks *h) {
L_17F90: ; /* IDA 0x17F90; call    sub_12078 */
    mch_call(m,8312,32659,h);
L_17F93: ; /* IDA 0x17F93; call    sub_101B4 */
    mch_call(m,436,32662,h);
L_17F96: ; /* IDA 0x17F96; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_17F9B: ; /* IDA 0x17F9B; add     si, 2240h */
    m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,0));
L_17F9F: ; /* IDA 0x17F9F; call    sub_1807B */
    mch_call(m,32891,32674,h);
L_17FA2: ; /* IDA 0x17FA2; call    sub_1812A */
    mch_call(m,33066,32677,h);
L_17FA5: ; /* IDA 0x17FA5; mov     al, cs:byte_10CFF */
    vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
L_17FA9: ; /* IDA 0x17FA9; cmp     al, [si+1] */
    ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+1)),8);
L_17FAC: ; /* IDA 0x17FAC; jnz     short loc_17FBC */
    if(!(m->flags&0x40)) goto L_17FBC;
L_17FAE: ; /* IDA 0x17FAE; call    sub_101DB */
    mch_call(m,475,32689,h);
L_17FB1: ; /* IDA 0x17FB1; call    sub_17FDB */
    mch_call(m,32731,32692,h);
L_17FB4: ; /* IDA 0x17FB4; call    sub_14325 */
    mch_call(m,17189,32695,h);
L_17FB7: ; /* IDA 0x17FB7; or      byte ptr [si], 2 */
#if KI_MARCH_MUTATION == 12
    ; /* lost ordered bit */
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(2))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#endif
L_17FBA: ; /* IDA 0x17FBA; jmp     short loc_17FC8 */
    goto L_17FC8;
L_17FBC: ; /* IDA 0x17FBC; mov     cx, 4 */
    m->cx=(uint16_t)(4);
L_17FBF: ; /* IDA 0x17FBF; call    sub_18853 */
    mch_call(m,34899,32706,h);
L_17FC2: ; /* IDA 0x17FC2; call    sub_101DB */
    mch_call(m,475,32709,h);
L_17FC5: ; /* IDA 0x17FC5; call    sub_121E7 */
    mch_call(m,8679,32712,h);
L_17FC8: ; /* IDA 0x17FC8; call    sub_101B4 */
    mch_call(m,436,32715,h);
L_17FCB: ; /* IDA 0x17FCB; call    sub_1817D */
    mch_call(m,33149,32718,h);
L_17FCE: ; /* IDA 0x17FCE; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_17FD1: ; /* IDA 0x17FD1; call    sub_18853 */
    mch_call(m,34899,32724,h);
L_17FD4: ; /* IDA 0x17FD4; call    sub_101DB */
    mch_call(m,475,32727,h);
L_17FD7: ; /* IDA 0x17FD7; call    sub_120D6 */
    mch_call(m,8406,32730,h);
L_17FDA: ; /* IDA 0x17FDA; retn */
    return;
}
/* Original sub_17FDB; IDA 0x17FDB. */
static void mch_body_17FDB(KiMachine16 *m,const KiEconomyHooks *h) {
L_17FDB: ; /* IDA 0x17FDB; push    bp */
    ec_push(m,m->bp);
L_17FDC: ; /* IDA 0x17FDC; sub     sp, 2 */
    m->sp=(uint16_t)(ec_math(m,m->sp,2,16,0,1));
L_17FDF: ; /* IDA 0x17FDF; mov     bp, sp */
    m->bp=(uint16_t)(m->sp);
L_17FE1: ; /* IDA 0x17FE1; call    sub_120D6 */
    mch_call(m,8406,32740,h);
L_17FE4: ; /* IDA 0x17FE4; mov     cx, 3 */
    m->cx=(uint16_t)(3);
L_17FE7: ; /* IDA 0x17FE7; call    sub_18853 */
    mch_call(m,34899,32746,h);
L_17FEA: ; /* IDA 0x17FEA; call    sub_1703C */
    mch_call(m,28732,32749,h);
L_17FED: ; /* IDA 0x17FED; jb      short loc_18046 */
    if(m->flags&1) goto L_18046;
L_17FEF: ; /* IDA 0x17FEF; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17FF1: ; /* IDA 0x17FF1; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17FF3: ; /* IDA 0x17FF3; shl     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
L_17FF5: ; /* IDA 0x17FF5; mov     [bp+0], bh */
    dl_write8(m,m->ss,(uint16_t)(m->bp),(uint8_t)(m->bx>>8));
L_17FF8: ; /* IDA 0x17FF8; mov     al, bh */
    vg_al(m,(uint8_t)((uint8_t)(m->bx>>8)));
L_17FFA: ; /* IDA 0x17FFA; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_17FFC: ; /* IDA 0x17FFC; push    di */
    ec_push(m,m->di);
L_17FFD: ; /* IDA 0x17FFD; push    ax */
    ec_push(m,m->ax);
L_17FFE: ; /* IDA 0x17FFE; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_18000: ; /* IDA 0x18000; mov     cx, 15h */
    m->cx=(uint16_t)(21);
L_18003: ; /* IDA 0x18003; call    sub_18853 */
    mch_call(m,34899,32774,h);
L_18006: ; /* IDA 0x18006; pop     ax */
    m->ax=(uint16_t)(ec_pop(m));
L_18007: ; /* IDA 0x18007; pop     di */
    m->di=(uint16_t)(ec_pop(m));
L_18008: ; /* IDA 0x18008; mov     bh, [si+1] */
    gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
L_1800B: ; /* IDA 0x1800B; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_1800D: ; /* IDA 0x1800D; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_1800F: ; /* IDA 0x1800F; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_18011: ; /* IDA 0x18011; mov     bl, [bx+3] */
    gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
L_18014: ; /* IDA 0x18014; mov     al, 2 */
    vg_al(m,(uint8_t)(2));
L_18016: ; /* IDA 0x18016; cmp     bl, [bp+0] */
    ec_cmp(m,(uint8_t)m->bx,dl_read8(m,m->ss,(uint16_t)(m->bp)),8);
L_18019: ; /* IDA 0x18019; jnz     short loc_1801D */
    if(!(m->flags&0x40)) goto L_1801D;
L_1801B: ; /* IDA 0x1801B; mov     al, 3 */
    vg_al(m,(uint8_t)(3));
L_1801D: ; /* IDA 0x1801D; call    sub_1804E */
    mch_call(m,32846,32800,h);
L_18020: ; /* IDA 0x18020; jb      short loc_17FE4 */
    if(m->flags&1) goto L_17FE4;
L_18022: ; /* IDA 0x18022; cmp     al, 1 */
    ec_cmp(m,(uint8_t)m->ax,1,8);
L_18024: ; /* IDA 0x18024; jz      short loc_1802D */
    if(m->flags&0x40) goto L_1802D;
L_18026: ; /* IDA 0x18026; ja      short loc_18032 */
    if(!(m->flags&(1|0x40))) goto L_18032;
L_18028: ; /* IDA 0x18028; and     byte ptr [si], 0FBh */
#if KI_MARCH_MUTATION == 3
    ; /* lost delegation clear */
#else
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(251))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
#endif
L_1802B: ; /* IDA 0x1802B; jmp     short loc_18038 */
    goto L_18038;
L_1802D: ; /* IDA 0x1802D; or      byte ptr [si], 4 */
    {uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(4))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
L_18030: ; /* IDA 0x18030; jmp     short loc_18038 */
    goto L_18038;
L_18032: ; /* IDA 0x18032; mov     byte ptr [si+23h], 0Bh */
#if KI_MARCH_MUTATION == 4
    dl_write8(m,m->ds,(uint16_t)(m->si+35),0);
#else
    dl_write8(m,m->ds,(uint16_t)(m->si+35),11);
#endif
L_18036: ; /* IDA 0x18036; jmp     short loc_1803C */
    goto L_1803C;
L_18038: ; /* IDA 0x18038; mov     byte ptr [si+23h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->si+35),0);
L_1803C: ; /* IDA 0x1803C; mov     byte ptr [si+0Bh], 1 */
    dl_write8(m,m->ds,(uint16_t)(m->si+11),1);
L_18040: ; /* IDA 0x18040; mov     al, [bp+0] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ss,(uint16_t)(m->bp))));
L_18043: ; /* IDA 0x18043; mov     [si+20h], al */
    dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)m->ax);
L_18046: ; /* IDA 0x18046; call    sub_12078 */
    mch_call(m,8312,32841,h);
L_18049: ; /* IDA 0x18049; add     sp, 2 */
    m->sp=(uint16_t)(ec_math(m,m->sp,2,16,0,0));
L_1804C: ; /* IDA 0x1804C; pop     bp */
    m->bp=(uint16_t)(ec_pop(m));
L_1804D: ; /* IDA 0x1804D; retn */
    return;
}
/* Original sub_1804E; IDA 0x1804E. */
static void mch_body_1804E(KiMachine16 *m,const KiEconomyHooks *h) {
L_1804E: ; /* IDA 0x1804E; push    ax */
    ec_push(m,m->ax);
L_1804F: ; /* IDA 0x1804F; mov     ah, 18h */
    vg_ah(m,(uint8_t)(24));
L_18051: ; /* IDA 0x18051; sub     ah, al */
    vg_ah(m,(uint8_t)(ec_math(m,(uint8_t)(m->ax>>8),(uint8_t)m->ax,8,0,1)));
L_18053: ; /* IDA 0x18053; mov     bx, cs:word_19898 */
    m->bx=(uint16_t)(dl_read16(m,m->cs,39064));
L_18058: ; /* IDA 0x18058; cmp     bl, ah */
    ec_cmp(m,(uint8_t)m->bx,(uint8_t)(m->ax>>8),8);
L_1805A: ; /* IDA 0x1805A; jbe     short loc_1805E */
    if(m->flags&(1|0x40)) goto L_1805E;
L_1805C: ; /* IDA 0x1805C; mov     bl, ah */
    gw_bl(m,(uint8_t)((uint8_t)(m->ax>>8)));
L_1805E: ; /* IDA 0x1805E; mov     dx, cs:word_19896 */
    m->dx=(uint16_t)(dl_read16(m,m->cs,39062));
L_18063: ; /* IDA 0x18063; cmp     dl, 21h ; '!' */
    ec_cmp(m,(uint8_t)m->dx,33,8);
L_18066: ; /* IDA 0x18066; jbe     short loc_1806A */
    if(m->flags&(1|0x40)) goto L_1806A;
L_18068: ; /* IDA 0x18068; mov     dl, 21h ; '!' */
    vg_dl(m,(uint8_t)(33));
L_1806A: ; /* IDA 0x1806A; mov     dh, bl */
    vg_dh(m,(uint8_t)((uint8_t)m->bx));
L_1806C: ; /* IDA 0x1806C; mov     cx, 4Ch ; 'L' */
    m->cx=(uint16_t)(76);
L_1806F: ; /* IDA 0x1806F; mov     ah, 1 */
    vg_ah(m,(uint8_t)(1));
L_18071: ; /* IDA 0x18071; call    sub_193E9 */
    mch_call(m,37865,32884,h);
L_18074: ; /* IDA 0x18074; pop     bx */
    m->bx=(uint16_t)(ec_pop(m));
L_18075: ; /* IDA 0x18075; sahf */
#if KI_MARCH_MUTATION == 5
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|3);
#else
    m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);
#endif
L_18076: ; /* IDA 0x18076; jnb     short locret_1807A */
    if(!(m->flags&1)) goto L_1807A;
L_18078: ; /* IDA 0x18078; mov     al, bl */
    vg_al(m,(uint8_t)((uint8_t)m->bx));
L_1807A: ; /* IDA 0x1807A; retn */
    return;
}
int ki_march_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x1c8d:mch_body_11C8D(m,h);break;
case 0x1f7f:mch_body_11F7F(m,h);break;
case 0x2151:mch_body_12151(m,h);break;
case 0x21b2:mch_body_121B2(m,h);break;
case 0x4325:mch_body_14325(m,h);break;
case 0x4370:mch_body_14370(m,h);break;
case 0x439d:mch_body_1439D(m,h);break;
case 0x43af:mch_body_143AF(m,h);break;
case 0x440f:mch_body_1440F(m,h);break;
case 0x4466:mch_body_14466(m,h);break;
case 0x4483:mch_body_14483(m,h);break;
case 0x4499:mch_body_14499(m,h);break;
case 0x44a9:mch_body_144A9(m,h);break;
case 0x44d6:mch_body_144D6(m,h);break;
case 0x4548:mch_body_14548(m,h);break;
case 0x463e:mch_body_1463E(m,h);break;
case 0x4651:mch_body_14651(m,h);break;
case 0x4689:mch_body_14689(m,h);break;
case 0x59a6:mch_body_159A6(m,h);break;
case 0x5ab6:mch_body_15AB6(m,h);break;
case 0x703c:mch_body_1703C(m,h);break;
case 0x709e:mch_body_1709E(m,h);break;
case 0x7f90:mch_body_17F90(m,h);break;
case 0x7fdb:mch_body_17FDB(m,h);break;
case 0x804e:mch_body_1804E(m,h);break;
default:return 0;}return 1;}
static int mch_known(uint16_t t) {switch(t) {case 0x1c8d: case 0x1f7f: case 0x2151: case 0x21b2: case 0x4325: case 0x4370: case 0x439d: case 0x43af: case 0x440f: case 0x4466: case 0x4483: case 0x4499: case 0x44a9: case 0x44d6: case 0x4548: case 0x463e: case 0x4651: case 0x4689: case 0x59a6: case 0x5ab6: case 0x703c: case 0x709e: case 0x7f90: case 0x7fdb: case 0x804e: return 1;default:return 0;}}
