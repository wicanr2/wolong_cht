/* Original native C personnel; re/121, spec/241. */
/* Original sub_16265; IDA 0x16265. */
static void psn_body_16265(KiMachine16 *m,const KiEconomyHooks *h) {
L_16265: ; /* IDA 0x16265; mov     ax, 4 */
    m->ax=(uint16_t)(4);
L_16268: ; /* IDA 0x16268; mov     cx, 4Eh ; 'N' */
    m->cx=(uint16_t)(78);
L_1626B: ; /* IDA 0x1626B; mov     dx, 403h */
    m->dx=(uint16_t)(1027);
L_1626E: ; /* IDA 0x1626E; call    sub_193E9 */
    psn_call(m,37865,25201,h);
L_16271: ; /* IDA 0x16271; jb      short locret_1627E */
    if(m->flags&1) goto L_1627E;
L_16273: ; /* IDA 0x16273; xor     ah, ah */
    {uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
L_16275: ; /* IDA 0x16275; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_16277: ; /* IDA 0x16277; mov     bx, ax */
    m->bx=(uint16_t)(m->ax);
L_16279: ; /* IDA 0x16279; call    cs:funcs_16279[bx] */
#if KI_PERSONNEL_MUTATION == 10
    psn_call(m,dl_read16(m,m->cs,(uint16_t)(0x6280+(m->bx^2))),0x627e,h);
#else
    psn_call(m,dl_read16(m,m->cs,(uint16_t)(m->bx+25216)),25214,h);
#endif
L_1627E: ; /* IDA 0x1627E; retn */
    return;
}
/* Original sub_16A9B; IDA 0x16A9B. */
static void psn_body_16A9B(KiMachine16 *m,const KiEconomyHooks *h) {
L_16A9B: ; /* IDA 0x16A9B; push    ds */
    ec_push(m,m->ds);
L_16A9C: ; /* IDA 0x16A9C; call    sub_12078 */
    psn_call(m,8312,27295,h);
L_16A9F: ; /* IDA 0x16A9F; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16AA4: ; /* IDA 0x16AA4; mov     cx, 0Bh */
    m->cx=(uint16_t)(11);
L_16AA7: ; /* IDA 0x16AA7; call    sub_18853 */
    psn_call(m,34899,27306,h);
L_16AAA: ; /* IDA 0x16AAA; call    sub_17400 */
    psn_call(m,29696,27309,h);
L_16AAD: ; /* IDA 0x16AAD; jb      short loc_16AFA */
    if(m->flags&1) goto L_16AFA;
L_16AAF: ; /* IDA 0x16AAF; cmp     byte ptr [bx+19h], 0FFh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+25)),255,8);
L_16AB3: ; /* IDA 0x16AB3; jz      short loc_16ACB */
    if(m->flags&0x40) goto L_16ACB;
L_16AB5: ; /* IDA 0x16AB5; mov     al, [bx+19h] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+25))));
L_16AB8: ; /* IDA 0x16AB8; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_16ABA: ; /* IDA 0x16ABA; push    ax */
    ec_push(m,m->ax);
L_16ABB: ; /* IDA 0x16ABB; push    bx */
    ec_push(m,m->bx);
L_16ABC: ; /* IDA 0x16ABC; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16ABE: ; /* IDA 0x16ABE; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_16AC0: ; /* IDA 0x16AC0; mov     cx, 34h ; '4' */
    m->cx=(uint16_t)(52);
L_16AC3: ; /* IDA 0x16AC3; call    sub_18810 */
    psn_call(m,34832,27334,h);
L_16AC6: ; /* IDA 0x16AC6; add     sp, 4 */
    m->sp=(uint16_t)(ec_math(m,m->sp,4,16,0,0));
L_16AC9: ; /* IDA 0x16AC9; jmp     short loc_16AA4 */
    goto L_16AA4;
L_16ACB: ; /* IDA 0x16ACB; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16ACD: ; /* IDA 0x16ACD; mov     cx, 9 */
    m->cx=(uint16_t)(9);
L_16AD0: ; /* IDA 0x16AD0; call    sub_18853 */
    psn_call(m,34899,27347,h);
L_16AD3: ; /* IDA 0x16AD3; call    sub_17663 */
    psn_call(m,30307,27350,h);
L_16AD6: ; /* IDA 0x16AD6; jb      short loc_16AA4 */
#if KI_PERSONNEL_MUTATION == 9
    if(m->flags&1) goto L_16AFA;
#else
    if(m->flags&1) goto L_16AA4;
#endif
L_16AD8: ; /* IDA 0x16AD8; mov     byte ptr [bx+17h], 2 */
#if KI_PERSONNEL_MUTATION == 1
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),3);
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),2);
#endif
L_16ADC: ; /* IDA 0x16ADC; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_16ADE: ; /* IDA 0x16ADE; sub     ax, 4240h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16960,16,0,1));
L_16AE1: ; /* IDA 0x16AE1; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_16AE3: ; /* IDA 0x16AE3; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_16AE5: ; /* IDA 0x16AE5; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_16AE7: ; /* IDA 0x16AE7; mov     [si+19h], ah */
#if KI_PERSONNEL_MUTATION == 2
    dl_write8(m,m->ds,(uint16_t)(m->si+25),(uint8_t)m->ax);
#else
    dl_write8(m,m->ds,(uint16_t)(m->si+25),(uint8_t)(m->ax>>8));
#endif
L_16AEA: ; /* IDA 0x16AEA; mov     cx, 19Ch */
    m->cx=(uint16_t)(412);
L_16AED: ; /* IDA 0x16AED; mov     di, sp */
#if KI_PERSONNEL_MUTATION == 11
    dl_write8(m,m->ds,(uint16_t)(m->bx+26),0);m->di=m->sp;
#else
    m->di=(uint16_t)(m->sp);
#endif
L_16AEF: ; /* IDA 0x16AEF; mov     ah, [bx+1Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
L_16AF2: ; /* IDA 0x16AF2; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_16AF5: ; /* IDA 0x16AF5; call    sub_18810 */
    psn_call(m,34832,27384,h);
L_16AF8: ; /* IDA 0x16AF8; jmp     short loc_16AA4 */
    goto L_16AA4;
L_16AFA: ; /* IDA 0x16AFA; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16AFD: ; /* IDA 0x16AFD; call    sub_18853 */
    psn_call(m,34899,27392,h);
L_16B00: ; /* IDA 0x16B00; call    sub_11D46 */
    psn_call(m,7494,27395,h);
L_16B03: ; /* IDA 0x16B03; call    sub_120D6 */
    psn_call(m,8406,27398,h);
L_16B06: ; /* IDA 0x16B06; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16B07: ; /* IDA 0x16B07; retn */
    return;
}
/* Original sub_16B08; IDA 0x16B08. */
static void psn_body_16B08(KiMachine16 *m,const KiEconomyHooks *h) {
L_16B08: ; /* IDA 0x16B08; push    ds */
    ec_push(m,m->ds);
L_16B09: ; /* IDA 0x16B09; call    sub_12078 */
    psn_call(m,8312,27404,h);
L_16B0C: ; /* IDA 0x16B0C; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16B11: ; /* IDA 0x16B11; mov     cx, 0Dh */
    m->cx=(uint16_t)(13);
L_16B14: ; /* IDA 0x16B14; call    sub_18853 */
    psn_call(m,34899,27415,h);
L_16B17: ; /* IDA 0x16B17; call    sub_17400 */
    psn_call(m,29696,27418,h);
L_16B1A: ; /* IDA 0x16B1A; jb      short loc_16B41 */
    if(m->flags&1) goto L_16B41;
L_16B1C: ; /* IDA 0x16B1C; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16B1E: ; /* IDA 0x16B1E; call    sub_16B4F */
    psn_call(m,27471,27425,h);
L_16B21: ; /* IDA 0x16B21; jnb     short loc_16B31 */
    if(!(m->flags&1)) goto L_16B31;
L_16B23: ; /* IDA 0x16B23; push    si */
    ec_push(m,m->si);
L_16B24: ; /* IDA 0x16B24; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16B26: ; /* IDA 0x16B26; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_16B28: ; /* IDA 0x16B28; mov     cx, 36h ; '6' */
    m->cx=(uint16_t)(54);
L_16B2B: ; /* IDA 0x16B2B; call    sub_18810 */
    psn_call(m,34832,27438,h);
L_16B2E: ; /* IDA 0x16B2E; pop     bx */
#if KI_PERSONNEL_MUTATION == 12
    m->si=ec_pop(m);
#else
    m->bx=(uint16_t)(ec_pop(m));
#endif
L_16B2F: ; /* IDA 0x16B2F; jmp     short loc_16B11 */
    goto L_16B11;
L_16B31: ; /* IDA 0x16B31; mov     cx, 1A2h */
    m->cx=(uint16_t)(418);
L_16B34: ; /* IDA 0x16B34; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16B36: ; /* IDA 0x16B36; mov     ah, [bx+1Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
L_16B39: ; /* IDA 0x16B39; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_16B3C: ; /* IDA 0x16B3C; call    sub_18810 */
    psn_call(m,34832,27455,h);
L_16B3F: ; /* IDA 0x16B3F; jmp     short loc_16B11 */
    goto L_16B11;
L_16B41: ; /* IDA 0x16B41; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16B44: ; /* IDA 0x16B44; call    sub_18853 */
    psn_call(m,34899,27463,h);
L_16B47: ; /* IDA 0x16B47; call    sub_11D46 */
    psn_call(m,7494,27466,h);
L_16B4A: ; /* IDA 0x16B4A; call    sub_120D6 */
    psn_call(m,8406,27469,h);
L_16B4D: ; /* IDA 0x16B4D; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16B4E: ; /* IDA 0x16B4E; retn */
    return;
}
/* Original sub_16B4F; IDA 0x16B4F. */
static void psn_body_16B4F(KiMachine16 *m,const KiEconomyHooks *h) {
L_16B4F: ; /* IDA 0x16B4F; mov     bh, 0FFh */
#if KI_PERSONNEL_MUTATION == 7
    m->bx&=0x00ff; /* wrong sentinel */
#else
    gw_bh(m,(uint8_t)(255));
#endif
L_16B51: ; /* IDA 0x16B51; xchg    bh, [si+19h] */
    {uint16_t a=(uint8_t)(m->bx>>8),b=dl_read8(m,m->ds,(uint16_t)(m->si+25));gw_bh(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->si+25),a);}
L_16B54: ; /* IDA 0x16B54; cmp     bh, 0FFh */
    ec_cmp(m,(uint8_t)(m->bx>>8),255,8);
L_16B57: ; /* IDA 0x16B57; jnz     short loc_16B5B */
    if(!(m->flags&0x40)) goto L_16B5B;
L_16B59: ; /* IDA 0x16B59; stc */
    m->flags|=1;
L_16B5A: ; /* IDA 0x16B5A; retn */
    return;
L_16B5B: ; /* IDA 0x16B5B; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_16B5D: ; /* IDA 0x16B5D; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16B5F: ; /* IDA 0x16B5F; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16B61: ; /* IDA 0x16B61; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16B63: ; /* IDA 0x16B63; add     bx, 4240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16960,16,0,0));
L_16B67: ; /* IDA 0x16B67; mov     byte ptr [bx+1Ah], 0 */
#if KI_PERSONNEL_MUTATION == 5
    ; /* lost administrator budget clear */
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+26),0);
#endif
L_16B6B: ; /* IDA 0x16B6B; mov     byte ptr [bx+17h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),0);
L_16B6F: ; /* IDA 0x16B6F; clc */
    m->flags&=(uint16_t)~1u;
L_16B70: ; /* IDA 0x16B70; retn */
    return;
}
/* Original sub_16B71; IDA 0x16B71. */
static void psn_body_16B71(KiMachine16 *m,const KiEconomyHooks *h) {
L_16B71: ; /* IDA 0x16B71; push    ds */
    ec_push(m,m->ds);
L_16B72: ; /* IDA 0x16B72; call    sub_12078 */
    psn_call(m,8312,27509,h);
L_16B75: ; /* IDA 0x16B75; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16B7A: ; /* IDA 0x16B7A; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16B7F: ; /* IDA 0x16B7F; mov     cx, 0Ch */
    m->cx=(uint16_t)(12);
L_16B82: ; /* IDA 0x16B82; call    sub_18853 */
    psn_call(m,34899,27525,h);
L_16B85: ; /* IDA 0x16B85; call    sub_17906 */
    psn_call(m,30982,27528,h);
L_16B88: ; /* IDA 0x16B88; jb      short loc_16BD5 */
    if(m->flags&1) goto L_16BD5;
L_16B8A: ; /* IDA 0x16B8A; cmp     byte ptr [bx+2Ah], 0FFh */
    ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+42)),255,8);
L_16B8E: ; /* IDA 0x16B8E; jz      short loc_16BA6 */
    if(m->flags&0x40) goto L_16BA6;
L_16B90: ; /* IDA 0x16B90; mov     al, [bx+2Ah] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+42))));
L_16B93: ; /* IDA 0x16B93; mov     ah, 0FFh */
    vg_ah(m,(uint8_t)(255));
L_16B95: ; /* IDA 0x16B95; push    ax */
    ec_push(m,m->ax);
L_16B96: ; /* IDA 0x16B96; push    bx */
    ec_push(m,m->bx);
L_16B97: ; /* IDA 0x16B97; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16B99: ; /* IDA 0x16B99; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_16B9B: ; /* IDA 0x16B9B; mov     cx, 35h ; '5' */
    m->cx=(uint16_t)(53);
L_16B9E: ; /* IDA 0x16B9E; call    sub_18810 */
    psn_call(m,34832,27553,h);
L_16BA1: ; /* IDA 0x16BA1; add     sp, 4 */
    m->sp=(uint16_t)(ec_math(m,m->sp,4,16,0,0));
L_16BA4: ; /* IDA 0x16BA4; jmp     short loc_16B7A */
    goto L_16B7A;
L_16BA6: ; /* IDA 0x16BA6; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16BA8: ; /* IDA 0x16BA8; mov     cx, 9 */
    m->cx=(uint16_t)(9);
L_16BAB: ; /* IDA 0x16BAB; call    sub_18853 */
    psn_call(m,34899,27566,h);
L_16BAE: ; /* IDA 0x16BAE; call    sub_17663 */
    psn_call(m,30307,27569,h);
L_16BB1: ; /* IDA 0x16BB1; jb      short loc_16B7A */
    if(m->flags&1) goto L_16B7A;
L_16BB3: ; /* IDA 0x16BB3; mov     byte ptr [bx+17h], 3 */
#if KI_PERSONNEL_MUTATION == 3
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),2);
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),3);
#endif
L_16BB7: ; /* IDA 0x16BB7; mov     ax, bx */
    m->ax=(uint16_t)(m->bx);
L_16BB9: ; /* IDA 0x16BB9; sub     ax, 4240h */
    m->ax=(uint16_t)(ec_math(m,m->ax,16960,16,0,1));
L_16BBC: ; /* IDA 0x16BBC; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_16BBE: ; /* IDA 0x16BBE; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_16BC0: ; /* IDA 0x16BC0; shl     ax, 1 */
    m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
L_16BC2: ; /* IDA 0x16BC2; mov     [si+2Ah], ah */
#if KI_PERSONNEL_MUTATION == 4
    dl_write8(m,m->ds,(uint16_t)(m->si+41),(uint8_t)(m->ax>>8));
#else
    dl_write8(m,m->ds,(uint16_t)(m->si+42),(uint8_t)(m->ax>>8));
#endif
L_16BC5: ; /* IDA 0x16BC5; mov     cx, 19Dh */
    m->cx=(uint16_t)(413);
L_16BC8: ; /* IDA 0x16BC8; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16BCA: ; /* IDA 0x16BCA; mov     ah, [bx+1Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
L_16BCD: ; /* IDA 0x16BCD; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_16BD0: ; /* IDA 0x16BD0; call    sub_18810 */
    psn_call(m,34832,27603,h);
L_16BD3: ; /* IDA 0x16BD3; jmp     short loc_16B7A */
    goto L_16B7A;
L_16BD5: ; /* IDA 0x16BD5; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16BD8: ; /* IDA 0x16BD8; call    sub_18853 */
    psn_call(m,34899,27611,h);
L_16BDB: ; /* IDA 0x16BDB; call    sub_11D46 */
    psn_call(m,7494,27614,h);
L_16BDE: ; /* IDA 0x16BDE; call    sub_120D6 */
    psn_call(m,8406,27617,h);
L_16BE1: ; /* IDA 0x16BE1; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16BE2: ; /* IDA 0x16BE2; retn */
    return;
}
/* Original sub_16BE3; IDA 0x16BE3. */
static void psn_body_16BE3(KiMachine16 *m,const KiEconomyHooks *h) {
L_16BE3: ; /* IDA 0x16BE3; push    ds */
    ec_push(m,m->ds);
L_16BE4: ; /* IDA 0x16BE4; call    sub_12078 */
    psn_call(m,8312,27623,h);
L_16BE7: ; /* IDA 0x16BE7; mov     ds, cs:word_10D52 */
    m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
L_16BEC: ; /* IDA 0x16BEC; mov     cx, 0Eh */
    m->cx=(uint16_t)(14);
L_16BEF: ; /* IDA 0x16BEF; call    sub_18853 */
    psn_call(m,34899,27634,h);
L_16BF2: ; /* IDA 0x16BF2; call    sub_17906 */
    psn_call(m,30982,27637,h);
L_16BF5: ; /* IDA 0x16BF5; jb      short loc_16C1C */
    if(m->flags&1) goto L_16C1C;
L_16BF7: ; /* IDA 0x16BF7; mov     si, bx */
    m->si=(uint16_t)(m->bx);
L_16BF9: ; /* IDA 0x16BF9; call    sub_16C2A */
    psn_call(m,27690,27644,h);
L_16BFC: ; /* IDA 0x16BFC; jnb     short loc_16C0C */
    if(!(m->flags&1)) goto L_16C0C;
L_16BFE: ; /* IDA 0x16BFE; push    si */
    ec_push(m,m->si);
L_16BFF: ; /* IDA 0x16BFF; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16C01: ; /* IDA 0x16C01; mov     al, 93h */
    vg_al(m,(uint8_t)(147));
L_16C03: ; /* IDA 0x16C03; mov     cx, 37h ; '7' */
    m->cx=(uint16_t)(55);
L_16C06: ; /* IDA 0x16C06; call    sub_18810 */
    psn_call(m,34832,27657,h);
L_16C09: ; /* IDA 0x16C09; pop     si */
    m->si=(uint16_t)(ec_pop(m));
L_16C0A: ; /* IDA 0x16C0A; jmp     short loc_16BEC */
    goto L_16BEC;
L_16C0C: ; /* IDA 0x16C0C; mov     cx, 1A3h */
    m->cx=(uint16_t)(419);
L_16C0F: ; /* IDA 0x16C0F; mov     di, sp */
    m->di=(uint16_t)(m->sp);
L_16C11: ; /* IDA 0x16C11; mov     ah, [bx+1Eh] */
    vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
L_16C14: ; /* IDA 0x16C14; mov     al, [bx+1] */
    vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
L_16C17: ; /* IDA 0x16C17; call    sub_18810 */
    psn_call(m,34832,27674,h);
L_16C1A: ; /* IDA 0x16C1A; jmp     short loc_16BEC */
    goto L_16BEC;
L_16C1C: ; /* IDA 0x16C1C; mov     cx, 0FFFFh */
    m->cx=(uint16_t)(65535);
L_16C1F: ; /* IDA 0x16C1F; call    sub_18853 */
    psn_call(m,34899,27682,h);
L_16C22: ; /* IDA 0x16C22; call    sub_11D46 */
    psn_call(m,7494,27685,h);
L_16C25: ; /* IDA 0x16C25; call    sub_120D6 */
    psn_call(m,8406,27688,h);
L_16C28: ; /* IDA 0x16C28; pop     ds */
    m->ds=(uint16_t)(ec_pop(m));
L_16C29: ; /* IDA 0x16C29; retn */
    return;
}
/* Original sub_16C2A; IDA 0x16C2A. */
static void psn_body_16C2A(KiMachine16 *m,const KiEconomyHooks *h) {
L_16C2A: ; /* IDA 0x16C2A; mov     bh, 0FFh */
    gw_bh(m,(uint8_t)(255));
L_16C2C: ; /* IDA 0x16C2C; xchg    bh, [si+2Ah] */
    {uint16_t a=(uint8_t)(m->bx>>8),b=dl_read8(m,m->ds,(uint16_t)(m->si+42));gw_bh(m,(uint8_t)(b));dl_write8(m,m->ds,(uint16_t)(m->si+42),a);}
L_16C2F: ; /* IDA 0x16C2F; cmp     bh, 0FFh */
    ec_cmp(m,(uint8_t)(m->bx>>8),255,8);
L_16C32: ; /* IDA 0x16C32; jnz     short loc_16C36 */
    if(!(m->flags&0x40)) goto L_16C36;
L_16C34: ; /* IDA 0x16C34; stc */
    m->flags|=1;
L_16C35: ; /* IDA 0x16C35; retn */
    return;
L_16C36: ; /* IDA 0x16C36; xor     bl, bl */
    {uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
L_16C38: ; /* IDA 0x16C38; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16C3A: ; /* IDA 0x16C3A; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16C3C: ; /* IDA 0x16C3C; shr     bx, 1 */
    m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
L_16C3E: ; /* IDA 0x16C3E; add     bx, 4240h */
    m->bx=(uint16_t)(ec_math(m,m->bx,16960,16,0,0));
L_16C42: ; /* IDA 0x16C42; mov     byte ptr [bx+1Ah], 0 */
#if KI_PERSONNEL_MUTATION == 6
    ; /* lost diplomat budget clear */
#else
    dl_write8(m,m->ds,(uint16_t)(m->bx+26),0);
#endif
L_16C46: ; /* IDA 0x16C46; mov     byte ptr [bx+17h], 0 */
    dl_write8(m,m->ds,(uint16_t)(m->bx+23),0);
L_16C4A: ; /* IDA 0x16C4A; clc */
#if KI_PERSONNEL_MUTATION == 8
    m->flags|=1;
#else
    m->flags&=(uint16_t)~1u;
#endif
L_16C4B: ; /* IDA 0x16C4B; retn */
    return;
}
int ki_personnel_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {
case 0x6265:psn_body_16265(m,h);break;
case 0x6a9b:psn_body_16A9B(m,h);break;
case 0x6b08:psn_body_16B08(m,h);break;
case 0x6b4f:psn_body_16B4F(m,h);break;
case 0x6b71:psn_body_16B71(m,h);break;
case 0x6be3:psn_body_16BE3(m,h);break;
case 0x6c2a:psn_body_16C2A(m,h);break;
default:return 0;}return 1;}
static int psn_known(uint16_t t) {switch(t) {case 0x6265: case 0x6a9b: case 0x6b08: case 0x6b4f: case 0x6b71: case 0x6be3: case 0x6c2a: return 1;default:return 0;}}
