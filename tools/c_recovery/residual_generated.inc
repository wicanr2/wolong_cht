/* Original native C residual graph; re/134, spec/254. No CPU.Step or guest decoder. */
static int rs_known(uint16_t t) {switch(t) {
case 0x45f8:
case 0x45f9:
case 0x45fb:
case 0x45fd:
case 0x45ff:
case 0x4601:
case 0x4604:
case 0x4606:
case 0x460a:
case 0x460c:
case 0x460e:
case 0x4610:
case 0x4612:
case 0x4615:
case 0x4617:
case 0x461b:
case 0x461c:
case 0x5ad1:
case 0x5ad4:
case 0x5ad7:
case 0x5ada:
case 0x5add:
case 0x5ae0:
case 0x5ae3:
case 0x5ae6:
case 0x5aec:
case 0x5aee:
case 0x5af1:
case 0x5af3:
case 0x5af5:
case 0x5af8:
case 0x5afb:
case 0x5afc:
case 0x5afd:
case 0x5b00:
case 0x5b02:
case 0x5b05:
case 0x5b08:
case 0x5b0b:
case 0x5b0e:
case 0x5b10:
case 0x5b13:
case 0x5b16:
case 0x5b19:
case 0x5b1c:
case 0x5b1f:
case 0x5b22:
case 0x5b25:
case 0x5b28:
case 0x5b2a:
case 0x5b2d:
case 0x5b30:
case 0x5b33:
case 0x5b36:
case 0x5b39:
case 0x5b3b:
case 0x5b3e:
case 0x5b40:
case 0x5b42:
case 0x5b45:
case 0x5b48:
case 0x5b4b:
case 0x5b4e:
case 0x5b51:
case 0x5b56:
case 0x5b58:
case 0x5b5a:
case 0x5b5f:
case 0x5b61:
case 0x5b63:
case 0x5b68:
case 0x5b6a:
case 0x5b6c:
case 0x5b6f:
case 0x5b71:
case 0x5b73:
case 0x5b75:
case 0x5b78:
case 0x5b7b:
case 0x5b7e:
case 0x5b81:
case 0x5b86:
case 0x5b88:
case 0x5b8a:
case 0x5b8d:
case 0x5b8f:
case 0x5b94:
case 0x5b96:
case 0x5b98:
case 0x5b99:
case 0x5b9a:
case 0x5b9d:
case 0x5b9e:
case 0x5b9f:
case 0x5ba1:
case 0x5ba3:
case 0x5ba7:
case 0x5bab:
case 0x5bad:
case 0x5bb1:
case 0x5bb3:
case 0x5bb5:
case 0x5bb7:
case 0x5bb9:
case 0x5bbb:
case 0x5bbd:
case 0x5bc0:
case 0x5bc2:
case 0x5bc7:
case 0x5bc9:
case 0x5bcc:
case 0x5bcd:
case 0x5bd1:
case 0x5bd3:
case 0x5bd6:
case 0x5bd7:
case 0x5bdb:
case 0x5bdd:
case 0x5be0:
case 0x5be3:
case 0x5be6:
case 0x5be9:
case 0x5beb:
case 0x5bee:
case 0x5bf0:
case 0x5bf3:
case 0x5bf6:
case 0x5bf9:
case 0x5bfc:
case 0x5bff:
case 0x5c02:
case 0x5c03:
case 0x5c14:
case 0x5c15:
case 0x5c17:
case 0x5c19:
case 0x5c1b:
case 0x5c1d:
case 0x5c1f:
case 0x5c22:
case 0x5c24:
case 0x5c26:
case 0x5c29:
case 0x5c2b:
case 0x5c2e:
case 0x5c31:
case 0x5c33:
case 0x5c35:
case 0x5c37:
case 0x5c39:
case 0x5c3d:
case 0x5c40:
case 0x5c42:
case 0x5c45:
case 0x5c47:
case 0x5c48:
case 0x5c4a:
case 0x5c4c:
case 0x5c4f:
case 0x5c52:
case 0x5c55:
case 0x5c56:
case 0x5c57:
case 0x9377:
case 0x937c:
case 0x937e:
case 0x9382:
case 0x9385:
case 0x9387:
case 0x938a:
case 0x938c:
case 0x938e:
case 0x9390:
case 0x9392:
case 0x9394:
case 0x9396:
case 0x9398:
case 0x939a:
case 0x939b:
case 0x939e:
case 0x93a1:
case 0x93a4:
case 0x93a6:
case 0x93a9:
case 0x93aa:
case 0x93ad:
case 0x93b0:
case 0x93b3:
case 0x93b5:
case 0x93b8:
case 0x93b9:
case 0x93bc:
case 0x93bf:
case 0x93c2:
case 0x93c4:
case 0x93c7:
case 0x93c8:
case 0x93cb:
case 0x93ce:
case 0x93d1:
case 0x93d3:
case 0x93d6:
case 0x93d7:
case 0x93da:
case 0x93dc:
case 0x93df:
case 0x93e1:
case 0x93e4:
case 0x93e5:
case 0x93e8:
case 0xab48:
case 0xab7c:
case 0xab80:
case 0xab82:
case 0xab86:
case 0xab89:
case 0xab8c:
case 0xab8e:
case 0xab90:
case 0xab93:
case 0xab94:
case 0xab97:
case 0xab98:
case 0xab9b:
case 0xc300:
case 0xc302:
case 0xc305:
case 0xc307:
case 0xc309:
case 0xc30c:
case 0xe567:
case 0xe569:
case 0xe56c:
case 0xe56f:
case 0xe572:
case 0xe575:
case 0xe577:
case 0xe57b:
case 0xe57e:
case 0xe6ff:
case 0xe701:
case 0xe704:
case 0xe707:
case 0xe70a:
case 0xe70d:
case 0xe70f:
case 0xe713:
case 0xe716:
return 1;default:return 0;}}
static void rs_graph(KiMachine16 *m,uint16_t entry,const KiEconomyHooks *h,int rs_skip_first_enter) {
const uint16_t rs_core_cs=m->cs;int rs_invoke_old=0;m->ip=entry;
RS_DISPATCH:
if(m->cs!=rs_core_cs) {assert(!rs_invoke_old);return;}
switch(m->ip) {
case 0x45f8:goto L_145F8;
case 0x45f9:goto L_145F9;
case 0x45fb:goto L_145FB;
case 0x45fd:goto L_145FD;
case 0x45ff:goto L_145FF;
case 0x4601:goto L_14601;
case 0x4604:goto L_14604;
case 0x4606:goto L_14606;
case 0x460a:goto L_1460A;
case 0x460c:goto L_1460C;
case 0x460e:goto L_1460E;
case 0x4610:goto L_14610;
case 0x4612:goto L_14612;
case 0x4615:goto L_14615;
case 0x4617:goto L_14617;
case 0x461b:goto L_1461B;
case 0x461c:goto L_1461C;
case 0x5ad1:goto L_15AD1;
case 0x5ad4:goto L_15AD4;
case 0x5ad7:goto L_15AD7;
case 0x5ada:goto L_15ADA;
case 0x5add:goto L_15ADD;
case 0x5ae0:goto L_15AE0;
case 0x5ae3:goto L_15AE3;
case 0x5ae6:goto L_15AE6;
case 0x5aec:goto L_15AEC;
case 0x5aee:goto L_15AEE;
case 0x5af1:goto L_15AF1;
case 0x5af3:goto L_15AF3;
case 0x5af5:goto L_15AF5;
case 0x5af8:goto L_15AF8;
case 0x5afb:goto L_15AFB;
case 0x5afc:goto L_15AFC;
case 0x5afd:goto L_15AFD;
case 0x5b00:goto L_15B00;
case 0x5b02:goto L_15B02;
case 0x5b05:goto L_15B05;
case 0x5b08:goto L_15B08;
case 0x5b0b:goto L_15B0B;
case 0x5b0e:goto L_15B0E;
case 0x5b10:goto L_15B10;
case 0x5b13:goto L_15B13;
case 0x5b16:goto L_15B16;
case 0x5b19:goto L_15B19;
case 0x5b1c:goto L_15B1C;
case 0x5b1f:goto L_15B1F;
case 0x5b22:goto L_15B22;
case 0x5b25:goto L_15B25;
case 0x5b28:goto L_15B28;
case 0x5b2a:goto L_15B2A;
case 0x5b2d:goto L_15B2D;
case 0x5b30:goto L_15B30;
case 0x5b33:goto L_15B33;
case 0x5b36:goto L_15B36;
case 0x5b39:goto L_15B39;
case 0x5b3b:goto L_15B3B;
case 0x5b3e:goto L_15B3E;
case 0x5b40:goto L_15B40;
case 0x5b42:goto L_15B42;
case 0x5b45:goto L_15B45;
case 0x5b48:goto L_15B48;
case 0x5b4b:goto L_15B4B;
case 0x5b4e:goto L_15B4E;
case 0x5b51:goto L_15B51;
case 0x5b56:goto L_15B56;
case 0x5b58:goto L_15B58;
case 0x5b5a:goto L_15B5A;
case 0x5b5f:goto L_15B5F;
case 0x5b61:goto L_15B61;
case 0x5b63:goto L_15B63;
case 0x5b68:goto L_15B68;
case 0x5b6a:goto L_15B6A;
case 0x5b6c:goto L_15B6C;
case 0x5b6f:goto L_15B6F;
case 0x5b71:goto L_15B71;
case 0x5b73:goto L_15B73;
case 0x5b75:goto L_15B75;
case 0x5b78:goto L_15B78;
case 0x5b7b:goto L_15B7B;
case 0x5b7e:goto L_15B7E;
case 0x5b81:goto L_15B81;
case 0x5b86:goto L_15B86;
case 0x5b88:goto L_15B88;
case 0x5b8a:goto L_15B8A;
case 0x5b8d:goto L_15B8D;
case 0x5b8f:goto L_15B8F;
case 0x5b94:goto L_15B94;
case 0x5b96:goto L_15B96;
case 0x5b98:goto L_15B98;
case 0x5b99:goto L_15B99;
case 0x5b9a:goto L_15B9A;
case 0x5b9d:goto L_15B9D;
case 0x5b9e:goto L_15B9E;
case 0x5b9f:goto L_15B9F;
case 0x5ba1:goto L_15BA1;
case 0x5ba3:goto L_15BA3;
case 0x5ba7:goto L_15BA7;
case 0x5bab:goto L_15BAB;
case 0x5bad:goto L_15BAD;
case 0x5bb1:goto L_15BB1;
case 0x5bb3:goto L_15BB3;
case 0x5bb5:goto L_15BB5;
case 0x5bb7:goto L_15BB7;
case 0x5bb9:goto L_15BB9;
case 0x5bbb:goto L_15BBB;
case 0x5bbd:goto L_15BBD;
case 0x5bc0:goto L_15BC0;
case 0x5bc2:goto L_15BC2;
case 0x5bc7:goto L_15BC7;
case 0x5bc9:goto L_15BC9;
case 0x5bcc:goto L_15BCC;
case 0x5bcd:goto L_15BCD;
case 0x5bd1:goto L_15BD1;
case 0x5bd3:goto L_15BD3;
case 0x5bd6:goto L_15BD6;
case 0x5bd7:goto L_15BD7;
case 0x5bdb:goto L_15BDB;
case 0x5bdd:goto L_15BDD;
case 0x5be0:goto L_15BE0;
case 0x5be3:goto L_15BE3;
case 0x5be6:goto L_15BE6;
case 0x5be9:goto L_15BE9;
case 0x5beb:goto L_15BEB;
case 0x5bee:goto L_15BEE;
case 0x5bf0:goto L_15BF0;
case 0x5bf3:goto L_15BF3;
case 0x5bf6:goto L_15BF6;
case 0x5bf9:goto L_15BF9;
case 0x5bfc:goto L_15BFC;
case 0x5bff:goto L_15BFF;
case 0x5c02:goto L_15C02;
case 0x5c03:goto L_15C03;
case 0x5c14:goto L_15C14;
case 0x5c15:goto L_15C15;
case 0x5c17:goto L_15C17;
case 0x5c19:goto L_15C19;
case 0x5c1b:goto L_15C1B;
case 0x5c1d:goto L_15C1D;
case 0x5c1f:goto L_15C1F;
case 0x5c22:goto L_15C22;
case 0x5c24:goto L_15C24;
case 0x5c26:goto L_15C26;
case 0x5c29:goto L_15C29;
case 0x5c2b:goto L_15C2B;
case 0x5c2e:goto L_15C2E;
case 0x5c31:goto L_15C31;
case 0x5c33:goto L_15C33;
case 0x5c35:goto L_15C35;
case 0x5c37:goto L_15C37;
case 0x5c39:goto L_15C39;
case 0x5c3d:goto L_15C3D;
case 0x5c40:goto L_15C40;
case 0x5c42:goto L_15C42;
case 0x5c45:goto L_15C45;
case 0x5c47:goto L_15C47;
case 0x5c48:goto L_15C48;
case 0x5c4a:goto L_15C4A;
case 0x5c4c:goto L_15C4C;
case 0x5c4f:goto L_15C4F;
case 0x5c52:goto L_15C52;
case 0x5c55:goto L_15C55;
case 0x5c56:goto L_15C56;
case 0x5c57:goto L_15C57;
case 0x9377:goto L_19377;
case 0x937c:goto L_1937C;
case 0x937e:goto L_1937E;
case 0x9382:goto L_19382;
case 0x9385:goto L_19385;
case 0x9387:goto L_19387;
case 0x938a:goto L_1938A;
case 0x938c:goto L_1938C;
case 0x938e:goto L_1938E;
case 0x9390:goto L_19390;
case 0x9392:goto L_19392;
case 0x9394:goto L_19394;
case 0x9396:goto L_19396;
case 0x9398:goto L_19398;
case 0x939a:goto L_1939A;
case 0x939b:goto L_1939B;
case 0x939e:goto L_1939E;
case 0x93a1:goto L_193A1;
case 0x93a4:goto L_193A4;
case 0x93a6:goto L_193A6;
case 0x93a9:goto L_193A9;
case 0x93aa:goto L_193AA;
case 0x93ad:goto L_193AD;
case 0x93b0:goto L_193B0;
case 0x93b3:goto L_193B3;
case 0x93b5:goto L_193B5;
case 0x93b8:goto L_193B8;
case 0x93b9:goto L_193B9;
case 0x93bc:goto L_193BC;
case 0x93bf:goto L_193BF;
case 0x93c2:goto L_193C2;
case 0x93c4:goto L_193C4;
case 0x93c7:goto L_193C7;
case 0x93c8:goto L_193C8;
case 0x93cb:goto L_193CB;
case 0x93ce:goto L_193CE;
case 0x93d1:goto L_193D1;
case 0x93d3:goto L_193D3;
case 0x93d6:goto L_193D6;
case 0x93d7:goto L_193D7;
case 0x93da:goto L_193DA;
case 0x93dc:goto L_193DC;
case 0x93df:goto L_193DF;
case 0x93e1:goto L_193E1;
case 0x93e4:goto L_193E4;
case 0x93e5:goto L_193E5;
case 0x93e8:goto L_193E8;
case 0xab48:goto L_1AB48;
case 0xab7c:goto L_1AB7C;
case 0xab80:goto L_1AB80;
case 0xab82:goto L_1AB82;
case 0xab86:goto L_1AB86;
case 0xab89:goto L_1AB89;
case 0xab8c:goto L_1AB8C;
case 0xab8e:goto L_1AB8E;
case 0xab90:goto L_1AB90;
case 0xab93:goto L_1AB93;
case 0xab94:goto L_1AB94;
case 0xab97:goto L_1AB97;
case 0xab98:goto L_1AB98;
case 0xab9b:goto L_1AB9B;
case 0xc300:goto L_1C300;
case 0xc302:goto L_1C302;
case 0xc305:goto L_1C305;
case 0xc307:goto L_1C307;
case 0xc309:goto L_1C309;
case 0xc30c:goto L_1C30C;
case 0xe567:goto L_1E567;
case 0xe569:goto L_1E569;
case 0xe56c:goto L_1E56C;
case 0xe56f:goto L_1E56F;
case 0xe572:goto L_1E572;
case 0xe575:goto L_1E575;
case 0xe577:goto L_1E577;
case 0xe57b:goto L_1E57B;
case 0xe57e:goto L_1E57E;
case 0xe6ff:goto L_1E6FF;
case 0xe701:goto L_1E701;
case 0xe704:goto L_1E704;
case 0xe707:goto L_1E707;
case 0xe70a:goto L_1E70A;
case 0xe70d:goto L_1E70D;
case 0xe70f:goto L_1E70F;
case 0xe713:goto L_1E713;
case 0xe716:goto L_1E716;
default:break;}
if(!rs_invoke_old)return;
ki_startup_invoke(m,m->ip,h);rs_invoke_old=0;goto RS_DISPATCH;
L_145F8: ; /* IDA 0x145F8; push    bx */
m->ip=0x45f8;
/* Original entry: sub_145F8 */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x45f8,h->user);
ec_push(m,m->bx);
goto L_145F9;
L_145F9: ; /* IDA 0x145F9; mov     bh, al */
m->ip=0x45f9;
rs_skip_first_enter=0;
gw_bh(m,(uint8_t)((uint8_t)m->ax));
goto L_145FB;
L_145FB: ; /* IDA 0x145FB; xor     bl, bl */
m->ip=0x45fb;
rs_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
goto L_145FD;
L_145FD: ; /* IDA 0x145FD; shr     bx, 1 */
m->ip=0x45fd;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_145FF;
L_145FF: ; /* IDA 0x145FF; shr     bx, 1 */
m->ip=0x45ff;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_14601;
L_14601: ; /* IDA 0x14601; cmp     ah, [bx+19h] */
m->ip=0x4601;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)(m->ax>>8),dl_read8(m,m->ds,(uint16_t)(m->bx+25)),8);
goto L_14604;
L_14604: ; /* IDA 0x14604; jnz     short loc_1460A */
m->ip=0x4604;
rs_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_1460A;}
goto L_14606;
L_14606: ; /* IDA 0x14606; mov     byte ptr [bx+19h], 0FFh */
m->ip=0x4606;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 1
dl_write8(m,m->ds,(uint16_t)(m->bx+25),254);
#else
dl_write8(m,m->ds,(uint16_t)(m->bx+25),255);
#endif
goto L_1460A;
L_1460A: ; /* IDA 0x1460A; mov     bh, ah */
m->ip=0x460a;
rs_skip_first_enter=0;
gw_bh(m,(uint8_t)((uint8_t)(m->ax>>8)));
goto L_1460C;
L_1460C: ; /* IDA 0x1460C; xor     bl, bl */
m->ip=0x460c;
rs_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
goto L_1460E;
L_1460E: ; /* IDA 0x1460E; shr     bx, 1 */
m->ip=0x460e;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_14610;
L_14610: ; /* IDA 0x14610; shr     bx, 1 */
m->ip=0x4610;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_14612;
L_14612: ; /* IDA 0x14612; cmp     al, [bx+19h] */
m->ip=0x4612;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+25)),8);
goto L_14615;
L_14615: ; /* IDA 0x14615; jnz     short loc_1461B */
m->ip=0x4615;
rs_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_1461B;}
goto L_14617;
L_14617: ; /* IDA 0x14617; mov     byte ptr [bx+19h], 0FFh */
m->ip=0x4617;
rs_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->bx+25),255);
goto L_1461B;
L_1461B: ; /* IDA 0x1461B; pop     bx */
m->ip=0x461b;
rs_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_1461C;
L_1461C: ; /* IDA 0x1461C; retn */
m->ip=0x461c;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_15AD1: ; /* IDA 0x15AD1; call    sub_15CA4 */
m->ip=0x5ad1;
/* Original entry: sub_15AD1 */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x5ad1,h->user);
{uint16_t target=23716;ec_push(m,0x5ad4);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15AD4: ; /* IDA 0x15AD4; call    sub_19528 */
m->ip=0x5ad4;
rs_skip_first_enter=0;
{uint16_t target=38184;ec_push(m,0x5ad7);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15AD7: ; /* IDA 0x15AD7; mov     cx, 4 */
m->ip=0x5ad7;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 2
m->cx=(uint16_t)(5);
#else
m->cx=(uint16_t)(4);
#endif
goto L_15ADA;
L_15ADA: ; /* IDA 0x15ADA; call    sub_18853 */
m->ip=0x5ada;
rs_skip_first_enter=0;
{uint16_t target=34899;ec_push(m,0x5add);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15ADD: ; /* IDA 0x15ADD; call    sub_15AFC */
m->ip=0x5add;
rs_skip_first_enter=0;
{uint16_t target=23292;ec_push(m,0x5ae0);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15AE0: ; /* IDA 0x15AE0; mov     cx, 0FFFFh */
m->ip=0x5ae0;
rs_skip_first_enter=0;
m->cx=(uint16_t)(65535);
goto L_15AE3;
L_15AE3: ; /* IDA 0x15AE3; call    sub_18853 */
m->ip=0x5ae3;
rs_skip_first_enter=0;
{uint16_t target=34899;ec_push(m,0x5ae6);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15AE6: ; /* IDA 0x15AE6; test    cs:byte_198A6, 2 */
m->ip=0x5ae6;
rs_skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->cs,39078))&(2))&255;ec_logic(m,result,8);}
goto L_15AEC;
L_15AEC: ; /* IDA 0x15AEC; jz      short loc_15AF5 */
m->ip=0x5aec;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_15AF5;}
goto L_15AEE;
L_15AEE: ; /* IDA 0x15AEE; call    sub_15E1E */
m->ip=0x5aee;
rs_skip_first_enter=0;
{uint16_t target=24094;ec_push(m,0x5af1);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15AF1: ; /* IDA 0x15AF1; mov     ax, cs */
m->ip=0x5af1;
rs_skip_first_enter=0;
m->ax=(uint16_t)(m->cs);
goto L_15AF3;
L_15AF3: ; /* IDA 0x15AF3; mov     ds, ax */
m->ip=0x5af3;
rs_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_15AF5;
L_15AF5: ; /* IDA 0x15AF5; call    sub_1950F */
m->ip=0x5af5;
rs_skip_first_enter=0;
{uint16_t target=38159;ec_push(m,0x5af8);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15AF8: ; /* IDA 0x15AF8; call    sub_15C58 */
m->ip=0x5af8;
rs_skip_first_enter=0;
{uint16_t target=23640;ec_push(m,0x5afb);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15AFB: ; /* IDA 0x15AFB; retn */
m->ip=0x5afb;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_15AFC: ; /* IDA 0x15AFC; push    ds */
m->ip=0x5afc;
/* Original entry: sub_15AFC */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x5afc,h->user);
ec_push(m,m->ds);
goto L_15AFD;
L_15AFD: ; /* IDA 0x15AFD; call    sub_101B4 */
m->ip=0x5afd;
rs_skip_first_enter=0;
{uint16_t target=436;ec_push(m,0x5b00);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B00: ; /* IDA 0x15B00; mov     al, 1Fh */
m->ip=0x5b00;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(31));
goto L_15B02;
L_15B02: ; /* IDA 0x15B02; mov     dx, 20h ; ' ' */
m->ip=0x5b02;
rs_skip_first_enter=0;
m->dx=(uint16_t)(32);
goto L_15B05;
L_15B05: ; /* IDA 0x15B05; mov     bx, 0Ah */
m->ip=0x5b05;
rs_skip_first_enter=0;
m->bx=(uint16_t)(10);
goto L_15B08;
L_15B08: ; /* IDA 0x15B08; mov     cx, 0D08h */
m->ip=0x5b08;
rs_skip_first_enter=0;
m->cx=(uint16_t)(3336);
goto L_15B0B;
L_15B0B: ; /* IDA 0x15B0B; call    sub_1895D */
m->ip=0x5b0b;
rs_skip_first_enter=0;
{uint16_t target=35165;ec_push(m,0x5b0e);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B0E: ; /* IDA 0x15B0E; mov     al, 20h ; ' ' */
m->ip=0x5b0e;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(32));
goto L_15B10;
L_15B10: ; /* IDA 0x15B10; mov     dx, 208h */
m->ip=0x5b10;
rs_skip_first_enter=0;
m->dx=(uint16_t)(520);
goto L_15B13;
L_15B13: ; /* IDA 0x15B13; mov     bx, 0D8h */
m->ip=0x5b13;
rs_skip_first_enter=0;
m->bx=(uint16_t)(216);
goto L_15B16;
L_15B16: ; /* IDA 0x15B16; mov     cx, 160Eh */
m->ip=0x5b16;
rs_skip_first_enter=0;
m->cx=(uint16_t)(5646);
goto L_15B19;
L_15B19: ; /* IDA 0x15B19; call    sub_1E3D7 */
m->ip=0x5b19;
rs_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x5b1c);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B1C: ; /* IDA 0x15B1C; mov     dx, 208h */
m->ip=0x5b1c;
rs_skip_first_enter=0;
m->dx=(uint16_t)(520);
goto L_15B1F;
L_15B1F: ; /* IDA 0x15B1F; mov     bx, 0C8h */
m->ip=0x5b1f;
rs_skip_first_enter=0;
m->bx=(uint16_t)(200);
goto L_15B22;
L_15B22: ; /* IDA 0x15B22; mov     si, 277h */
m->ip=0x5b22;
rs_skip_first_enter=0;
m->si=(uint16_t)(631);
goto L_15B25;
L_15B25: ; /* IDA 0x15B25; mov     di, 0D7h */
m->ip=0x5b25;
rs_skip_first_enter=0;
m->di=(uint16_t)(215);
goto L_15B28;
L_15B28: ; /* IDA 0x15B28; mov     ah, 7 */
m->ip=0x5b28;
rs_skip_first_enter=0;
vg_ah(m,(uint8_t)(7));
goto L_15B2A;
L_15B2A: ; /* IDA 0x15B2A; call    sub_1F1A3 */
m->ip=0x5b2a;
rs_skip_first_enter=0;
{uint16_t target=61859;ec_push(m,0x5b2d);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B2D: ; /* IDA 0x15B2D; mov     dx, 208h */
m->ip=0x5b2d;
rs_skip_first_enter=0;
m->dx=(uint16_t)(520);
goto L_15B30;
L_15B30: ; /* IDA 0x15B30; mov     bx, 0D8h */
m->ip=0x5b30;
rs_skip_first_enter=0;
m->bx=(uint16_t)(216);
goto L_15B33;
L_15B33: ; /* IDA 0x15B33; mov     si, 277h */
m->ip=0x5b33;
rs_skip_first_enter=0;
m->si=(uint16_t)(631);
goto L_15B36;
L_15B36: ; /* IDA 0x15B36; mov     di, 187h */
m->ip=0x5b36;
rs_skip_first_enter=0;
m->di=(uint16_t)(391);
goto L_15B39;
L_15B39: ; /* IDA 0x15B39; mov     ah, 9 */
m->ip=0x5b39;
rs_skip_first_enter=0;
vg_ah(m,(uint8_t)(9));
goto L_15B3B;
L_15B3B: ; /* IDA 0x15B3B; call    sub_1F1A3 */
m->ip=0x5b3b;
rs_skip_first_enter=0;
{uint16_t target=61859;ec_push(m,0x5b3e);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B3E: ; /* IDA 0x15B3E; mov     ax, cs */
m->ip=0x5b3e;
rs_skip_first_enter=0;
m->ax=(uint16_t)(m->cs);
goto L_15B40;
L_15B40: ; /* IDA 0x15B40; mov     ds, ax */
m->ip=0x5b40;
rs_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_15B42;
L_15B42: ; /* IDA 0x15B42; mov     si, 5C04h */
m->ip=0x5b42;
rs_skip_first_enter=0;
m->si=(uint16_t)(23556);
goto L_15B45;
L_15B45: ; /* IDA 0x15B45; mov     dx, 220h */
m->ip=0x5b45;
rs_skip_first_enter=0;
m->dx=(uint16_t)(544);
goto L_15B48;
L_15B48: ; /* IDA 0x15B48; mov     bx, 0C8h */
m->ip=0x5b48;
rs_skip_first_enter=0;
m->bx=(uint16_t)(200);
goto L_15B4B;
L_15B4B: ; /* IDA 0x15B4B; mov     ax, 7001h */
m->ip=0x5b4b;
rs_skip_first_enter=0;
m->ax=(uint16_t)(28673);
goto L_15B4E;
L_15B4E: ; /* IDA 0x15B4E; call    sub_106F5 */
m->ip=0x5b4e;
rs_skip_first_enter=0;
{uint16_t target=1781;ec_push(m,0x5b51);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B51: ; /* IDA 0x15B51; mov     ds, cs:word_10D52 */
m->ip=0x5b51;
rs_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_15B56;
L_15B56: ; /* IDA 0x15B56; xor     al, al */
m->ip=0x5b56;
rs_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_15B58;
L_15B58: ; /* IDA 0x15B58; mov     ah, 90h */
m->ip=0x5b58;
rs_skip_first_enter=0;
vg_ah(m,(uint8_t)(144));
goto L_15B5A;
L_15B5A: ; /* IDA 0x15B5A; cmp     al, cs:byte_10CFF */
m->ip=0x5b5a;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
goto L_15B5F;
L_15B5F: ; /* IDA 0x15B5F; jnz     short loc_15B63 */
m->ip=0x5b5f;
rs_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_15B63;}
goto L_15B61;
L_15B61: ; /* IDA 0x15B61; mov     ah, 9Ah */
m->ip=0x5b61;
rs_skip_first_enter=0;
vg_ah(m,(uint8_t)(154));
goto L_15B63;
L_15B63: ; /* IDA 0x15B63; cmp     al, cs:byte_198A7 */
m->ip=0x5b63;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,39079),8);
goto L_15B68;
L_15B68: ; /* IDA 0x15B68; jnz     short loc_15B6C */
m->ip=0x5b68;
rs_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_15B6C;}
goto L_15B6A;
L_15B6A: ; /* IDA 0x15B6A; mov     ah, 93h */
m->ip=0x5b6a;
rs_skip_first_enter=0;
vg_ah(m,(uint8_t)(147));
goto L_15B6C;
L_15B6C: ; /* IDA 0x15B6C; call    sub_15C14 */
m->ip=0x5b6c;
rs_skip_first_enter=0;
{uint16_t target=23572;ec_push(m,0x5b6f);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B6F: ; /* IDA 0x15B6F; inc     al */
m->ip=0x5b6f;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
goto L_15B71;
L_15B71: ; /* IDA 0x15B71; cmp     al, 16h */
m->ip=0x5b71;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 3
ec_cmp(m,(uint8_t)m->ax,23,8);
#else
ec_cmp(m,(uint8_t)m->ax,22,8);
#endif
goto L_15B73;
L_15B73: ; /* IDA 0x15B73; jb      short loc_15B58 */
m->ip=0x5b73;
rs_skip_first_enter=0;
if(m->flags&1) {goto L_15B58;}
goto L_15B75;
L_15B75: ; /* IDA 0x15B75; call    sub_101DB */
m->ip=0x5b75;
rs_skip_first_enter=0;
{uint16_t target=475;ec_push(m,0x5b78);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B78: ; /* IDA 0x15B78; call    sub_12078 */
m->ip=0x5b78;
rs_skip_first_enter=0;
{uint16_t target=8312;ec_push(m,0x5b7b);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B7B: ; /* IDA 0x15B7B; mov     ax, 5 */
m->ip=0x5b7b;
/* Original entry: inner call target */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x5b7b,h->user);
m->ax=(uint16_t)(5);
goto L_15B7E;
L_15B7E: ; /* IDA 0x15B7E; mov     bx, 1 */
m->ip=0x5b7e;
rs_skip_first_enter=0;
m->bx=(uint16_t)(1);
goto L_15B81;
L_15B81: ; /* IDA 0x15B81; call    sub_20000 */
m->ip=0x5b81;
rs_skip_first_enter=0;
mx_far(m,0x5b81,0x5b86,h);rs_invoke_old=0;goto RS_DISPATCH;
L_15B86: ; /* IDA 0x15B86; and     bx, bx */
m->ip=0x5b86;
rs_skip_first_enter=0;
{uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
goto L_15B88;
L_15B88: ; /* IDA 0x15B88; jnz     short loc_15BEB */
m->ip=0x5b88;
rs_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_15BEB;}
goto L_15B8A;
L_15B8A: ; /* IDA 0x15B8A; mov     ax, 5 */
m->ip=0x5b8a;
rs_skip_first_enter=0;
m->ax=(uint16_t)(5);
goto L_15B8D;
L_15B8D: ; /* IDA 0x15B8D; xor     bx, bx */
m->ip=0x5b8d;
rs_skip_first_enter=0;
{uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
goto L_15B8F;
L_15B8F: ; /* IDA 0x15B8F; call    sub_20000 */
m->ip=0x5b8f;
rs_skip_first_enter=0;
mx_far(m,0x5b8f,0x5b94,h);rs_invoke_old=0;goto RS_DISPATCH;
L_15B94: ; /* IDA 0x15B94; and     bx, bx */
m->ip=0x5b94;
rs_skip_first_enter=0;
{uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
goto L_15B96;
L_15B96: ; /* IDA 0x15B96; jz      short loc_15B7B */
m->ip=0x5b96;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_15B7B;}
goto L_15B98;
L_15B98: ; /* IDA 0x15B98; push    cx */
m->ip=0x5b98;
rs_skip_first_enter=0;
ec_push(m,m->cx);
goto L_15B99;
L_15B99: ; /* IDA 0x15B99; push    dx */
m->ip=0x5b99;
rs_skip_first_enter=0;
ec_push(m,m->dx);
goto L_15B9A;
L_15B9A: ; /* IDA 0x15B9A; call    sub_1E453 */
m->ip=0x5b9a;
rs_skip_first_enter=0;
{uint16_t target=58451;ec_push(m,0x5b9d);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15B9D: ; /* IDA 0x15B9D; pop     dx */
m->ip=0x5b9d;
rs_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_15B9E;
L_15B9E: ; /* IDA 0x15B9E; pop     cx */
m->ip=0x5b9e;
rs_skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_15B9F;
L_15B9F: ; /* IDA 0x15B9F; cmp     al, 20h ; ' ' */
m->ip=0x5b9f;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,32,8);
goto L_15BA1;
L_15BA1: ; /* IDA 0x15BA1; jnz     short loc_15B7B */
m->ip=0x5ba1;
rs_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_15B7B;}
goto L_15BA3;
L_15BA3: ; /* IDA 0x15BA3; sub     dx, 0D8h */
m->ip=0x5ba3;
rs_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,216,16,0,1));
goto L_15BA7;
L_15BA7: ; /* IDA 0x15BA7; cmp     cx, 240h */
m->ip=0x5ba7;
rs_skip_first_enter=0;
ec_cmp(m,m->cx,576,16);
goto L_15BAB;
L_15BAB: ; /* IDA 0x15BAB; jb      short loc_15BB1 */
m->ip=0x5bab;
rs_skip_first_enter=0;
if(m->flags&1) {goto L_15BB1;}
goto L_15BAD;
L_15BAD: ; /* IDA 0x15BAD; add     dx, 0B0h */
m->ip=0x5bad;
rs_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,176,16,0,0));
goto L_15BB1;
L_15BB1: ; /* IDA 0x15BB1; mov     cl, 4 */
m->ip=0x5bb1;
rs_skip_first_enter=0;
vg_cl(m,(uint8_t)(4));
goto L_15BB3;
L_15BB3: ; /* IDA 0x15BB3; shr     dx, cl */
m->ip=0x5bb3;
rs_skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
goto L_15BB5;
L_15BB5: ; /* IDA 0x15BB5; mov     ax, dx */
m->ip=0x5bb5;
rs_skip_first_enter=0;
m->ax=(uint16_t)(m->dx);
goto L_15BB7;
L_15BB7: ; /* IDA 0x15BB7; mov     cl, 6 */
m->ip=0x5bb7;
rs_skip_first_enter=0;
vg_cl(m,(uint8_t)(6));
goto L_15BB9;
L_15BB9: ; /* IDA 0x15BB9; shl     dx, cl */
m->ip=0x5bb9;
rs_skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,4));
goto L_15BBB;
L_15BBB: ; /* IDA 0x15BBB; mov     si, dx */
m->ip=0x5bbb;
rs_skip_first_enter=0;
m->si=(uint16_t)(m->dx);
goto L_15BBD;
L_15BBD: ; /* IDA 0x15BBD; cmp     byte ptr [si], 80h */
m->ip=0x5bbd;
rs_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
goto L_15BC0;
L_15BC0: ; /* IDA 0x15BC0; jb      short loc_15B7B */
m->ip=0x5bc0;
rs_skip_first_enter=0;
if(m->flags&1) {goto L_15B7B;}
goto L_15BC2;
L_15BC2: ; /* IDA 0x15BC2; cmp     al, cs:byte_10CFF */
m->ip=0x5bc2;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
goto L_15BC7;
L_15BC7: ; /* IDA 0x15BC7; jz      short loc_15B7B */
m->ip=0x5bc7;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_15B7B;}
goto L_15BC9;
L_15BC9: ; /* IDA 0x15BC9; call    sub_101B4 */
m->ip=0x5bc9;
rs_skip_first_enter=0;
{uint16_t target=436;ec_push(m,0x5bcc);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BCC: ; /* IDA 0x15BCC; push    ax */
m->ip=0x5bcc;
rs_skip_first_enter=0;
ec_push(m,m->ax);
goto L_15BCD;
L_15BCD: ; /* IDA 0x15BCD; mov     al, cs:byte_198A7 */
m->ip=0x5bcd;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->cs,39079)));
goto L_15BD1;
L_15BD1: ; /* IDA 0x15BD1; mov     ah, 90h */
m->ip=0x5bd1;
rs_skip_first_enter=0;
vg_ah(m,(uint8_t)(144));
goto L_15BD3;
L_15BD3: ; /* IDA 0x15BD3; call    sub_15C14 */
m->ip=0x5bd3;
rs_skip_first_enter=0;
{uint16_t target=23572;ec_push(m,0x5bd6);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BD6: ; /* IDA 0x15BD6; pop     ax */
m->ip=0x5bd6;
rs_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_15BD7;
L_15BD7: ; /* IDA 0x15BD7; mov     cs:byte_198A7, al */
m->ip=0x5bd7;
rs_skip_first_enter=0;
dl_write8(m,m->cs,39079,(uint8_t)m->ax);
goto L_15BDB;
L_15BDB: ; /* IDA 0x15BDB; mov     ah, 93h */
m->ip=0x5bdb;
rs_skip_first_enter=0;
vg_ah(m,(uint8_t)(147));
goto L_15BDD;
L_15BDD: ; /* IDA 0x15BDD; call    sub_15C14 */
m->ip=0x5bdd;
rs_skip_first_enter=0;
{uint16_t target=23572;ec_push(m,0x5be0);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BE0: ; /* IDA 0x15BE0; call    sub_15DBB */
m->ip=0x5be0;
rs_skip_first_enter=0;
{uint16_t target=23995;ec_push(m,0x5be3);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BE3: ; /* IDA 0x15BE3; call    sub_15CC6 */
m->ip=0x5be3;
rs_skip_first_enter=0;
{uint16_t target=23750;ec_push(m,0x5be6);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BE6: ; /* IDA 0x15BE6; call    sub_101DB */
m->ip=0x5be6;
rs_skip_first_enter=0;
{uint16_t target=475;ec_push(m,0x5be9);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BE9: ; /* IDA 0x15BE9; jmp     short loc_15B7B */
m->ip=0x5be9;
rs_skip_first_enter=0;
goto L_15B7B;
L_15BEB: ; /* IDA 0x15BEB; call    sub_101B4 */
m->ip=0x5beb;
rs_skip_first_enter=0;
{uint16_t target=436;ec_push(m,0x5bee);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BEE: ; /* IDA 0x15BEE; xor     al, al */
m->ip=0x5bee;
rs_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_15BF0;
L_15BF0: ; /* IDA 0x15BF0; mov     dx, 20h ; ' ' */
m->ip=0x5bf0;
rs_skip_first_enter=0;
m->dx=(uint16_t)(32);
goto L_15BF3;
L_15BF3: ; /* IDA 0x15BF3; mov     bx, 0Ah */
m->ip=0x5bf3;
rs_skip_first_enter=0;
m->bx=(uint16_t)(10);
goto L_15BF6;
L_15BF6: ; /* IDA 0x15BF6; mov     cx, 0D08h */
m->ip=0x5bf6;
rs_skip_first_enter=0;
m->cx=(uint16_t)(3336);
goto L_15BF9;
L_15BF9: ; /* IDA 0x15BF9; call    sub_1895D */
m->ip=0x5bf9;
rs_skip_first_enter=0;
{uint16_t target=35165;ec_push(m,0x5bfc);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BFC: ; /* IDA 0x15BFC; call    sub_101DB */
m->ip=0x5bfc;
rs_skip_first_enter=0;
{uint16_t target=475;ec_push(m,0x5bff);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15BFF: ; /* IDA 0x15BFF; call    sub_120D6 */
m->ip=0x5bff;
rs_skip_first_enter=0;
{uint16_t target=8406;ec_push(m,0x5c02);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15C02: ; /* IDA 0x15C02; pop     ds */
m->ip=0x5c02;
rs_skip_first_enter=0;
m->ds=(uint16_t)(ec_pop(m));
goto L_15C03;
L_15C03: ; /* IDA 0x15C03; retn */
m->ip=0x5c03;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_15C14: ; /* IDA 0x15C14; push    ax */
m->ip=0x5c14;
/* Original entry: sub_15C14 */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x5c14,h->user);
ec_push(m,m->ax);
goto L_15C15;
L_15C15: ; /* IDA 0x15C15; mov     bh, al */
m->ip=0x5c15;
rs_skip_first_enter=0;
gw_bh(m,(uint8_t)((uint8_t)m->ax));
goto L_15C17;
L_15C17: ; /* IDA 0x15C17; xor     bl, bl */
m->ip=0x5c17;
rs_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
goto L_15C19;
L_15C19: ; /* IDA 0x15C19; shr     bx, 1 */
m->ip=0x5c19;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_15C1B;
L_15C1B: ; /* IDA 0x15C1B; shr     bx, 1 */
m->ip=0x5c1b;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_15C1D;
L_15C1D: ; /* IDA 0x15C1D; mov     si, bx */
m->ip=0x5c1d;
rs_skip_first_enter=0;
m->si=(uint16_t)(m->bx);
goto L_15C1F;
L_15C1F: ; /* IDA 0x15C1F; mov     dx, 20Ch */
m->ip=0x5c1f;
rs_skip_first_enter=0;
m->dx=(uint16_t)(524);
goto L_15C22;
L_15C22: ; /* IDA 0x15C22; mov     bl, al */
m->ip=0x5c22;
rs_skip_first_enter=0;
gw_bl(m,(uint8_t)((uint8_t)m->ax));
goto L_15C24;
L_15C24: ; /* IDA 0x15C24; xor     bh, bh */
m->ip=0x5c24;
rs_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
goto L_15C26;
L_15C26: ; /* IDA 0x15C26; cmp     bl, 0Bh */
m->ip=0x5c26;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 4
ec_cmp(m,(uint8_t)m->bx,12,8);
#else
ec_cmp(m,(uint8_t)m->bx,11,8);
#endif
goto L_15C29;
L_15C29: ; /* IDA 0x15C29; jb      short loc_15C31 */
m->ip=0x5c29;
rs_skip_first_enter=0;
if(m->flags&1) {goto L_15C31;}
goto L_15C2B;
L_15C2B: ; /* IDA 0x15C2B; mov     dx, 244h */
m->ip=0x5c2b;
rs_skip_first_enter=0;
m->dx=(uint16_t)(580);
goto L_15C2E;
L_15C2E: ; /* IDA 0x15C2E; sub     bx, 0Bh */
m->ip=0x5c2e;
rs_skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,11,16,0,1));
goto L_15C31;
L_15C31: ; /* IDA 0x15C31; shl     bx, 1 */
m->ip=0x5c31;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_15C33;
L_15C33: ; /* IDA 0x15C33; shl     bx, 1 */
m->ip=0x5c33;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_15C35;
L_15C35: ; /* IDA 0x15C35; shl     bx, 1 */
m->ip=0x5c35;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_15C37;
L_15C37: ; /* IDA 0x15C37; shl     bx, 1 */
m->ip=0x5c37;
rs_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_15C39;
L_15C39: ; /* IDA 0x15C39; add     bx, 0D8h */
m->ip=0x5c39;
rs_skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,216,16,0,0));
goto L_15C3D;
L_15C3D: ; /* IDA 0x15C3D; cmp     byte ptr [si], 80h */
m->ip=0x5c3d;
rs_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si)),128,8);
goto L_15C40;
L_15C40: ; /* IDA 0x15C40; jb      short loc_15C47 */
m->ip=0x5c40;
rs_skip_first_enter=0;
if(m->flags&1) {goto L_15C47;}
goto L_15C42;
L_15C42: ; /* IDA 0x15C42; call    sub_188B0 */
m->ip=0x5c42;
rs_skip_first_enter=0;
{uint16_t target=34992;ec_push(m,0x5c45);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15C45: ; /* IDA 0x15C45; jmp     short loc_15C56 */
m->ip=0x5c45;
rs_skip_first_enter=0;
goto L_15C56;
L_15C47: ; /* IDA 0x15C47; push    ds */
m->ip=0x5c47;
rs_skip_first_enter=0;
ec_push(m,m->ds);
goto L_15C48;
L_15C48: ; /* IDA 0x15C48; mov     ax, cs */
m->ip=0x5c48;
rs_skip_first_enter=0;
m->ax=(uint16_t)(m->cs);
goto L_15C4A;
L_15C4A: ; /* IDA 0x15C4A; mov     ds, ax */
m->ip=0x5c4a;
rs_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_15C4C;
L_15C4C: ; /* IDA 0x15C4C; mov     ax, 9001h */
m->ip=0x5c4c;
rs_skip_first_enter=0;
m->ax=(uint16_t)(36865);
goto L_15C4F;
L_15C4F: ; /* IDA 0x15C4F; mov     si, 5C0Dh */
m->ip=0x5c4f;
rs_skip_first_enter=0;
m->si=(uint16_t)(23565);
goto L_15C52;
L_15C52: ; /* IDA 0x15C52; call    sub_106F5 */
m->ip=0x5c52;
rs_skip_first_enter=0;
{uint16_t target=1781;ec_push(m,0x5c55);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_15C55: ; /* IDA 0x15C55; pop     ds */
m->ip=0x5c55;
rs_skip_first_enter=0;
m->ds=(uint16_t)(ec_pop(m));
goto L_15C56;
L_15C56: ; /* IDA 0x15C56; pop     ax */
m->ip=0x5c56;
rs_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_15C57;
L_15C57: ; /* IDA 0x15C57; retn */
m->ip=0x5c57;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_19377: ; /* IDA 0x19377; cmp     byte ptr ds:0CF3h, 1 */
m->ip=0x9377;
/* Original entry: sub_19377 */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x9377,h->user);
ec_cmp(m,dl_read8(m,m->ds,3315),1,8);
goto L_1937C;
L_1937C: ; /* IDA 0x1937C; jnz     short locret_1939A */
m->ip=0x937c;
rs_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_1939A;}
goto L_1937E;
L_1937E: ; /* IDA 0x1937E; mov     ah, ds:0CF0h */
m->ip=0x937e;
rs_skip_first_enter=0;
vg_ah(m,(uint8_t)(dl_read8(m,m->ds,3312)));
goto L_19382;
L_19382: ; /* IDA 0x19382; cmp     ah, 10h */
m->ip=0x9382;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)(m->ax>>8),16,8);
goto L_19385;
L_19385: ; /* IDA 0x19385; ja      short locret_1939A */
m->ip=0x9385;
rs_skip_first_enter=0;
if(!(m->flags&(1|0x40))) {goto L_1939A;}
goto L_19387;
L_19387: ; /* IDA 0x19387; mov     al, ds:0CF4h */
m->ip=0x9387;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,3316)));
goto L_1938A;
L_1938A: ; /* IDA 0x1938A; cmp     al, 3 */
m->ip=0x938a;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 5
ec_cmp(m,(uint8_t)m->ax,2,8);
#else
ec_cmp(m,(uint8_t)m->ax,3,8);
#endif
goto L_1938C;
L_1938C: ; /* IDA 0x1938C; jz      short loc_1939B */
m->ip=0x938c;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_1939B;}
goto L_1938E;
L_1938E: ; /* IDA 0x1938E; cmp     al, 6 */
m->ip=0x938e;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,6,8);
goto L_19390;
L_19390: ; /* IDA 0x19390; jz      short loc_193AA */
m->ip=0x9390;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_193AA;}
goto L_19392;
L_19392: ; /* IDA 0x19392; cmp     al, 9 */
m->ip=0x9392;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,9,8);
goto L_19394;
L_19394: ; /* IDA 0x19394; jz      short loc_193B9 */
m->ip=0x9394;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_193B9;}
goto L_19396;
L_19396: ; /* IDA 0x19396; cmp     al, 0Ch */
m->ip=0x9396;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,12,8);
goto L_19398;
L_19398: ; /* IDA 0x19398; jz      short loc_193C8 */
m->ip=0x9398;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_193C8;}
goto L_1939A;
L_1939A: ; /* IDA 0x1939A; retn */
m->ip=0x939a;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_1939B: ; /* IDA 0x1939B; mov     si, 1934h */
m->ip=0x939b;
rs_skip_first_enter=0;
m->si=(uint16_t)(6452);
goto L_1939E;
L_1939E: ; /* IDA 0x1939E; mov     di, 18A4h */
m->ip=0x939e;
rs_skip_first_enter=0;
m->di=(uint16_t)(6308);
goto L_193A1;
L_193A1: ; /* IDA 0x193A1; call    sub_10A65 */
m->ip=0x93a1;
rs_skip_first_enter=0;
{uint16_t target=2661;ec_push(m,0x93a4);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193A4: ; /* IDA 0x193A4; mov     al, 2 */
m->ip=0x93a4;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(2));
goto L_193A6;
L_193A6: ; /* IDA 0x193A6; call    sub_193D7 */
m->ip=0x93a6;
rs_skip_first_enter=0;
{uint16_t target=37847;ec_push(m,0x93a9);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193A9: ; /* IDA 0x193A9; retn */
m->ip=0x93a9;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_193AA: ; /* IDA 0x193AA; mov     si, 18A4h */
m->ip=0x93aa;
rs_skip_first_enter=0;
m->si=(uint16_t)(6308);
goto L_193AD;
L_193AD: ; /* IDA 0x193AD; mov     di, 18D4h */
m->ip=0x93ad;
rs_skip_first_enter=0;
m->di=(uint16_t)(6356);
goto L_193B0;
L_193B0: ; /* IDA 0x193B0; call    sub_10A65 */
m->ip=0x93b0;
rs_skip_first_enter=0;
{uint16_t target=2661;ec_push(m,0x93b3);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193B3: ; /* IDA 0x193B3; mov     al, 3 */
m->ip=0x93b3;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(3));
goto L_193B5;
L_193B5: ; /* IDA 0x193B5; call    sub_193D7 */
m->ip=0x93b5;
rs_skip_first_enter=0;
{uint16_t target=37847;ec_push(m,0x93b8);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193B8: ; /* IDA 0x193B8; retn */
m->ip=0x93b8;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_193B9: ; /* IDA 0x193B9; mov     si, 18D4h */
m->ip=0x93b9;
rs_skip_first_enter=0;
m->si=(uint16_t)(6356);
goto L_193BC;
L_193BC: ; /* IDA 0x193BC; mov     di, 1904h */
m->ip=0x93bc;
rs_skip_first_enter=0;
m->di=(uint16_t)(6404);
goto L_193BF;
L_193BF: ; /* IDA 0x193BF; call    sub_10A65 */
m->ip=0x93bf;
rs_skip_first_enter=0;
{uint16_t target=2661;ec_push(m,0x93c2);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193C2: ; /* IDA 0x193C2; mov     al, 4 */
m->ip=0x93c2;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(4));
goto L_193C4;
L_193C4: ; /* IDA 0x193C4; call    sub_193D7 */
m->ip=0x93c4;
rs_skip_first_enter=0;
{uint16_t target=37847;ec_push(m,0x93c7);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193C7: ; /* IDA 0x193C7; retn */
m->ip=0x93c7;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_193C8: ; /* IDA 0x193C8; mov     si, 1904h */
m->ip=0x93c8;
rs_skip_first_enter=0;
m->si=(uint16_t)(6404);
goto L_193CB;
L_193CB: ; /* IDA 0x193CB; mov     di, 1934h */
m->ip=0x93cb;
rs_skip_first_enter=0;
m->di=(uint16_t)(6452);
goto L_193CE;
L_193CE: ; /* IDA 0x193CE; call    sub_10A65 */
m->ip=0x93ce;
rs_skip_first_enter=0;
{uint16_t target=2661;ec_push(m,0x93d1);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193D1: ; /* IDA 0x193D1; mov     al, 5 */
m->ip=0x93d1;
rs_skip_first_enter=0;
vg_al(m,(uint8_t)(5));
goto L_193D3;
L_193D3: ; /* IDA 0x193D3; call    sub_193D7 */
m->ip=0x93d3;
rs_skip_first_enter=0;
{uint16_t target=37847;ec_push(m,0x93d6);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193D6: ; /* IDA 0x193D6; retn */
m->ip=0x93d6;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_193D7: ; /* IDA 0x193D7; cmp     ah, 1 */
m->ip=0x93d7;
/* Original entry: sub_193D7 */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x93d7,h->user);
ec_cmp(m,(uint8_t)(m->ax>>8),1,8);
goto L_193DA;
L_193DA: ; /* IDA 0x193DA; jz      short loc_193E5 */
m->ip=0x93da;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_193E5;}
goto L_193DC;
L_193DC: ; /* IDA 0x193DC; cmp     ah, 2 */
m->ip=0x93dc;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 6
ec_cmp(m,(uint8_t)(m->ax>>8),3,8);
#else
ec_cmp(m,(uint8_t)(m->ax>>8),2,8);
#endif
goto L_193DF;
L_193DF: ; /* IDA 0x193DF; jnz     short locret_193E4 */
m->ip=0x93df;
rs_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_193E4;}
goto L_193E1;
L_193E1: ; /* IDA 0x193E1; call    sub_10241 */
m->ip=0x93e1;
rs_skip_first_enter=0;
{uint16_t target=577;ec_push(m,0x93e4);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193E4: ; /* IDA 0x193E4; retn */
m->ip=0x93e4;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_193E5: ; /* IDA 0x193E5; call    sub_102C2 */
m->ip=0x93e5;
rs_skip_first_enter=0;
{uint16_t target=706;ec_push(m,0x93e8);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_193E8: ; /* IDA 0x193E8; retn */
m->ip=0x93e8;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_1AB48: ; /* IDA 0x1AB48; jmp     short sub_1AB7C */
m->ip=0xab48;
/* Original entry: raw_entry_1AB48 */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0xab48,h->user);
goto L_1AB7C;
L_1AB7C: ; /* IDA 0x1AB7C; cmp     byte ptr [si+19h], 28h ; '(' */
m->ip=0xab7c;
/* Original entry: sub_1AB7C */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0xab7c,h->user);
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+25)),40,8);
goto L_1AB80;
L_1AB80: ; /* IDA 0x1AB80; jb      short loc_1AB86 */
m->ip=0xab80;
rs_skip_first_enter=0;
if(m->flags&1) {goto L_1AB86;}
goto L_1AB82;
L_1AB82: ; /* IDA 0x1AB82; mov     byte ptr [si+19h], 28h ; '(' */
m->ip=0xab82;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 7
dl_write8(m,m->ds,(uint16_t)(m->si+25),41);
#else
dl_write8(m,m->ds,(uint16_t)(m->si+25),40);
#endif
goto L_1AB86;
L_1AB86: ; /* IDA 0x1AB86; mov     bl, [si+4] */
m->ip=0xab86;
rs_skip_first_enter=0;
gw_bl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+4))));
goto L_1AB89;
L_1AB89: ; /* IDA 0x1AB89; cmp     bl, 24h ; '$' */
m->ip=0xab89;
rs_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->bx,36,8);
goto L_1AB8C;
L_1AB8C: ; /* IDA 0x1AB8C; jz      short loc_1AB94 */
m->ip=0xab8c;
rs_skip_first_enter=0;
if(m->flags&0x40) {goto L_1AB94;}
goto L_1AB8E;
L_1AB8E: ; /* IDA 0x1AB8E; ja      short loc_1AB98 */
m->ip=0xab8e;
rs_skip_first_enter=0;
if(!(m->flags&(1|0x40))) {goto L_1AB98;}
goto L_1AB90;
L_1AB90: ; /* IDA 0x1AB90; call    sub_1ABD2 */
m->ip=0xab90;
rs_skip_first_enter=0;
{uint16_t target=43986;ec_push(m,0xab93);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_1AB93: ; /* IDA 0x1AB93; retn */
m->ip=0xab93;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_1AB94: ; /* IDA 0x1AB94; call    sub_1ABFF */
m->ip=0xab94;
rs_skip_first_enter=0;
{uint16_t target=44031;ec_push(m,0xab97);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_1AB97: ; /* IDA 0x1AB97; retn */
m->ip=0xab97;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_1AB98: ; /* IDA 0x1AB98; call    sub_1AC55 */
m->ip=0xab98;
rs_skip_first_enter=0;
{uint16_t target=44117;ec_push(m,0xab9b);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_1AB9B: ; /* IDA 0x1AB9B; retn */
m->ip=0xab9b;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_1C300: ; /* IDA 0x1C300; adc     bl, dl */
m->ip=0xc300;
/* Original entry: sub_1C300 */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0xc300,h->user);
gw_bl(m,(uint8_t)(ec_math(m,(uint8_t)m->bx,(uint8_t)m->dx,8,m->flags&1,0)));
goto L_1C302;
L_1C302: ; /* IDA 0x1C302; and     [di+2], dh */
m->ip=0xc302;
rs_skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->di+2)))&((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->di+2),result);}
goto L_1C305;
L_1C305: ; /* IDA 0x1C305; xor     ah, ah */
m->ip=0xc305;
rs_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_1C307;
L_1C307: ; /* IDA 0x1C307; mov     al, 5 */
m->ip=0xc307;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 8
vg_al(m,(uint8_t)(6));
#else
vg_al(m,(uint8_t)(5));
#endif
goto L_1C309;
L_1C309: ; /* IDA 0x1C309; call    sub_1C6BF */
m->ip=0xc309;
rs_skip_first_enter=0;
{uint16_t target=50879;ec_push(m,0xc30c);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}
L_1C30C: ; /* IDA 0x1C30C; retn */
m->ip=0xc30c;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_1E567: ; /* IDA 0x1E567; mov     [si], di */
m->ip=0xe567;
/* Original entry: sub_1E567 */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0xe567,h->user);
dl_write16(m,m->ds,(uint16_t)(m->si),m->di);
goto L_1E569;
L_1E569: ; /* IDA 0x1E569; mov     [si+2], bp */
m->ip=0xe569;
rs_skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+2),m->bp);
goto L_1E56C;
L_1E56C: ; /* IDA 0x1E56C; mov     [si+4], cx */
m->ip=0xe56c;
rs_skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+4),m->cx);
goto L_1E56F;
L_1E56F: ; /* IDA 0x1E56F; mov     [si+6], ah */
m->ip=0xe56f;
rs_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+6),(uint8_t)(m->ax>>8));
goto L_1E572;
L_1E572: ; /* IDA 0x1E572; mov     [si+7], al */
m->ip=0xe572;
rs_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+7),(uint8_t)m->ax);
goto L_1E575;
L_1E575: ; /* IDA 0x1E575; mov     es, dx */
m->ip=0xe575;
rs_skip_first_enter=0;
m->es=(uint16_t)(m->dx);
goto L_1E577;
L_1E577: ; /* IDA 0x1E577; mov     al, cs:byte_1E47E */
m->ip=0xe577;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 9
vg_al(m,(uint8_t)(dl_read8(m,m->cs,58495)));
#else
vg_al(m,(uint8_t)(dl_read8(m,m->cs,58494)));
#endif
goto L_1E57B;
L_1E57B: ; /* IDA 0x1E57B; mov     es:[di], al */
m->ip=0xe57b;
rs_skip_first_enter=0;
dl_write8(m,m->es,(uint16_t)(m->di),(uint8_t)m->ax);
goto L_1E57E;
L_1E57E: ; /* IDA 0x1E57E; retn */
m->ip=0xe57e;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
L_1E6FF: ; /* IDA 0x1E6FF; mov     [si], di */
m->ip=0xe6ff;
/* Original entry: sub_1E6FF */
if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0xe6ff,h->user);
dl_write16(m,m->ds,(uint16_t)(m->si),m->di);
goto L_1E701;
L_1E701: ; /* IDA 0x1E701; mov     [si+2], bp */
m->ip=0xe701;
rs_skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+2),m->bp);
goto L_1E704;
L_1E704: ; /* IDA 0x1E704; mov     [si+4], cx */
m->ip=0xe704;
rs_skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+4),m->cx);
goto L_1E707;
L_1E707: ; /* IDA 0x1E707; mov     [si+6], ah */
m->ip=0xe707;
rs_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+6),(uint8_t)(m->ax>>8));
goto L_1E70A;
L_1E70A: ; /* IDA 0x1E70A; mov     [si+7], al */
m->ip=0xe70a;
rs_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+7),(uint8_t)m->ax);
goto L_1E70D;
L_1E70D: ; /* IDA 0x1E70D; mov     es, dx */
m->ip=0xe70d;
rs_skip_first_enter=0;
m->es=(uint16_t)(m->dx);
goto L_1E70F;
L_1E70F: ; /* IDA 0x1E70F; mov     al, cs:byte_1E47F */
m->ip=0xe70f;
rs_skip_first_enter=0;
#if KI_RESIDUAL_MUTATION == 10
vg_al(m,(uint8_t)(dl_read8(m,m->cs,58494)));
#else
vg_al(m,(uint8_t)(dl_read8(m,m->cs,58495)));
#endif
goto L_1E713;
L_1E713: ; /* IDA 0x1E713; mov     es:[di], al */
m->ip=0xe713;
rs_skip_first_enter=0;
dl_write8(m,m->es,(uint16_t)(m->di),(uint8_t)m->ax);
goto L_1E716;
L_1E716: ; /* IDA 0x1E716; retn */
m->ip=0xe716;
rs_skip_first_enter=0;
m->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;
}
int ki_residual_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
if(!rs_known(t))return 0;
uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,ss,sp);
rs_graph(m,t,h,1); /* Old invoke already emitted the entry hook. */
/* Old external ABI owns one final POP; reject abandoned native frames. */
assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
ec_push(m,ret);return 1;}
void ki_residual_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {if(rs_known(t))rs_graph(m,t,h,0);else ki_startup_invoke(m,t,h);}
void sub_145F8(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0x45f8,h,0);}
void sub_15AD1(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0x5ad1,h,0);}
void sub_15AFC(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0x5afc,h,0);}
void sub_15C14(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0x5c14,h,0);}
void sub_19377(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0x9377,h,0);}
void sub_193D7(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0x93d7,h,0);}
void sub_1AB7C(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0xab7c,h,0);}
void sub_1C300(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0xc300,h,0);}
void sub_1E567(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0xe567,h,0);}
void sub_1E6FF(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0xe6ff,h,0);}
void raw_entry_1AB48(KiMachine16 *m,const KiEconomyHooks *h) {rs_graph(m,0xab48,h,0);}
