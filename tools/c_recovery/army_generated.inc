/* Fixed original army C graph; re/132, spec/252. No guest instruction decoder. */
static int am_known(uint16_t t) {switch(t) {
case 0x25a3:
case 0x25a8:
case 0x25ad:
case 0x25b2:
case 0x25b5:
case 0x25b6:
case 0x25bb:
case 0x25bd:
case 0x25c1:
case 0x25c4:
case 0x25c6:
case 0x25c9:
case 0x25cc:
case 0x25cf:
case 0x25d2:
case 0x25d5:
case 0x25d8:
case 0x25dc:
case 0x25de:
case 0x25e3:
case 0x25e5:
case 0x25e8:
case 0x25e9:
case 0x25ec:
case 0x25ee:
case 0x25f2:
case 0x25f4:
case 0x25f6:
case 0x25fb:
case 0x25fd:
case 0x25ff:
case 0x2600:
case 0x2606:
case 0x2608:
case 0x2609:
case 0x260c:
case 0x2611:
case 0x2613:
case 0x2615:
case 0x2617:
case 0x2619:
case 0x261b:
case 0x261d:
case 0x2620:
case 0x2623:
case 0x2624:
case 0x2626:
case 0x2628:
case 0x2629:
case 0x262b:
case 0x262e:
case 0x2631:
case 0x2634:
case 0x2636:
case 0x2638:
case 0x263a:
case 0x263d:
case 0x2641:
case 0x2644:
case 0x2646:
case 0x2649:
case 0x264a:
case 0x264d:
case 0x264f:
case 0x2653:
case 0x2657:
case 0x2658:
case 0x265b:
case 0x265d:
case 0x2661:
case 0x2662:
case 0x2665:
case 0x2668:
case 0x266a:
case 0x266e:
case 0x2671:
case 0x2674:
case 0x2675:
case 0x2678:
case 0x267a:
case 0x267d:
case 0x2681:
case 0x2683:
case 0x2686:
case 0x268b:
case 0x268d:
case 0x2690:
case 0x2692:
case 0x2694:
case 0x2697:
case 0x2699:
case 0x269c:
case 0x269e:
case 0x26a0:
case 0x26a3:
case 0x26a5:
case 0x26a8:
case 0x26ab:
case 0x26ad:
case 0x26b0:
case 0x26b2:
case 0x26b4:
case 0x26b9:
case 0x26bb:
case 0x26be:
case 0x26c0:
case 0x26c3:
case 0x26c6:
case 0x26c8:
case 0x26cb:
case 0x26cf:
case 0x26d1:
case 0x26d3:
case 0x26d5:
case 0x26d7:
case 0x26da:
case 0x26dc:
case 0x26e0:
case 0x26e2:
case 0x26e4:
case 0x26e8:
case 0x26ea:
case 0x26ed:
case 0x26f0:
case 0x26f2:
case 0x26f5:
case 0x26f7:
case 0x26fa:
case 0x26fc:
case 0x26fe:
case 0x26ff:
case 0x2702:
case 0x2705:
case 0x2706:
case 0x2708:
case 0x270b:
case 0x270f:
case 0x2711:
case 0x2713:
case 0x2715:
case 0x2717:
case 0x2719:
case 0x271b:
case 0x271d:
case 0x271f:
case 0x2721:
case 0x2726:
case 0x2728:
case 0x272b:
case 0x272d:
case 0x2732:
case 0x2734:
case 0x2736:
case 0x2739:
case 0x273b:
case 0x273c:
case 0x273e:
case 0x2743:
case 0x2745:
case 0x2747:
case 0x274c:
case 0x274f:
case 0x2751:
case 0x2753:
case 0x2755:
case 0x2757:
case 0x2759:
case 0x275c:
case 0x275e:
case 0x275f:
case 0x2764:
case 0x2767:
case 0x276a:
case 0x276d:
case 0x2771:
case 0x2774:
case 0x2777:
case 0x2779:
case 0x277b:
case 0x277d:
case 0x277f:
case 0x2782:
case 0x2783:
case 0x2786:
case 0x2788:
case 0x278c:
case 0x278e:
case 0x2791:
case 0x2794:
case 0x2795:
case 0x2799:
case 0x279b:
case 0x279e:
case 0x27a1:
case 0x27a2:
case 0x27a5:
case 0x27a9:
case 0x27ab:
case 0x27af:
case 0x27b1:
case 0x27b5:
case 0x27b8:
case 0x27bc:
case 0x27be:
case 0x27c0:
case 0x27c2:
case 0x27c6:
case 0x27ca:
case 0x27cd:
case 0x27d0:
case 0x27d2:
case 0x27d4:
case 0x27d6:
case 0x27d8:
case 0x27da:
case 0x27dc:
case 0x27de:
case 0x27e3:
case 0x27e6:
case 0x27e9:
case 0x27ec:
case 0x27ef:
case 0x27f1:
case 0x27f5:
case 0x27f6:
case 0x27f9:
case 0x27fb:
case 0x27fc:
case 0x27ff:
case 0x2803:
case 0x2804:
case 0x2807:
case 0x2808:
case 0x280b:
case 0x280d:
case 0x2810:
case 0x2813:
case 0x2815:
case 0x2818:
case 0x281a:
case 0x281d:
case 0x281e:
case 0x2821:
case 0x2825:
case 0x2827:
case 0x2829:
case 0x282b:
case 0x282d:
case 0x2830:
case 0x2831:
case 0x2832:
case 0x2833:
case 0x2834:
case 0x2837:
case 0x283a:
case 0x283d:
case 0x2840:
case 0x2842:
case 0x2845:
case 0x2847:
case 0x284a:
case 0x284c:
case 0x284f:
case 0x2851:
case 0x2853:
case 0x2856:
case 0x2859:
case 0x285b:
case 0x285e:
case 0x2862:
case 0x2864:
case 0x2866:
case 0x286a:
case 0x286c:
case 0x286e:
case 0x2871:
case 0x2873:
case 0x2876:
case 0x2877:
case 0x2878:
case 0x2879:
case 0x287a:
case 0x287b:
case 0x287c:
case 0x287d:
case 0x287e:
case 0x287f:
case 0x2880:
case 0x2881:
case 0x2882:
case 0x2883:
case 0x2886:
case 0x288a:
case 0x288c:
case 0x2890:
case 0x2892:
case 0x2896:
case 0x2898:
case 0x289a:
case 0x289d:
case 0x28a1:
case 0x28a3:
case 0x28a6:
case 0x28aa:
case 0x28ac:
case 0x28ae:
case 0x28b2:
case 0x28b4:
case 0x28b6:
case 0x28b9:
case 0x28bb:
case 0x28bf:
case 0x28c2:
case 0x28c3:
case 0x28c4:
case 0x28c5:
case 0x28c6:
case 0x28c7:
case 0x28c8:
case 0x28c9:
case 0x28ca:
case 0x28cb:
case 0x28f4:
case 0x28f6:
case 0x28f8:
case 0x28fa:
case 0x28fd:
case 0x2901:
case 0x2903:
case 0x2906:
case 0x2908:
case 0x290a:
case 0x290c:
case 0x290e:
case 0x2910:
case 0x2912:
case 0x2915:
case 0x2916:
case 0x2917:
case 0x2918:
case 0x2a7e:
case 0x2a82:
case 0x2a84:
case 0x2a85:
case 0x2a8a:
case 0x2a8c:
case 0x2a8e:
case 0x2a93:
case 0x2a97:
case 0x2a99:
case 0x2a9b:
case 0x2a9d:
case 0x2aa0:
case 0x2aa2:
case 0x2aa7:
case 0x2aa8:
case 0x2aac:
case 0x2ab0:
case 0x2ab2:
case 0x2ab6:
case 0x2ab7:
case 0x2ab9:
case 0x2abc:
case 0x2abf:
case 0x2ac1:
case 0x2ac4:
case 0x2ac5:
case 0x2ac8:
case 0x2acb:
case 0x2ace:
case 0x2ad1:
case 0x42ab:
case 0x42ac:
case 0x42ad:
case 0x42af:
case 0x42b1:
case 0x42b5:
case 0x42b7:
case 0x42b8:
case 0x42b9:
case 0x42bd:
case 0x42bf:
case 0x42c1:
case 0x42c4:
case 0x42c8:
case 0x42ca:
case 0x42cc:
case 0x42ce:
case 0x42d0:
case 0x42d2:
case 0x42d4:
case 0x42d6:
case 0x42d8:
case 0x42da:
case 0x42dc:
case 0x42de:
case 0x42e2:
case 0x42e3:
case 0x42e5:
case 0x42e7:
case 0x42ea:
case 0x42ee:
case 0x42f1:
case 0x42f3:
case 0x42f5:
case 0x42f7:
case 0x42fa:
case 0x42fd:
case 0x42fe:
case 0x42ff:
case 0x4300:
case 0x4302:
case 0x4304:
case 0x4309:
case 0x430b:
case 0x430f:
case 0x4311:
case 0x4316:
case 0x4318:
case 0x431a:
case 0x431c:
case 0x431e:
case 0x4321:
case 0x4322:
case 0x4323:
case 0x4324:
case 0x562b:
case 0x562c:
case 0x562e:
case 0x5632:
case 0x5634:
case 0x5636:
case 0x5639:
case 0x563a:
return 1;default:return 0;}}
static void am_graph(KiMachine16 *m,uint16_t entry,const KiEconomyHooks *h,int skip_first_enter) {
const uint16_t core_cs=m->cs;int invoke_external=0;m->ip=entry;
AM_DISPATCH:
if(m->cs!=core_cs) {assert(!invoke_external);return;}
switch(m->ip) {
case 0x25a3:goto L_125A3;
case 0x25a8:goto L_125A8;
case 0x25ad:goto L_125AD;
case 0x25b2:goto L_125B2;
case 0x25b5:goto L_125B5;
case 0x25b6:goto L_125B6;
case 0x25bb:goto L_125BB;
case 0x25bd:goto L_125BD;
case 0x25c1:goto L_125C1;
case 0x25c4:goto L_125C4;
case 0x25c6:goto L_125C6;
case 0x25c9:goto L_125C9;
case 0x25cc:goto L_125CC;
case 0x25cf:goto L_125CF;
case 0x25d2:goto L_125D2;
case 0x25d5:goto L_125D5;
case 0x25d8:goto L_125D8;
case 0x25dc:goto L_125DC;
case 0x25de:goto L_125DE;
case 0x25e3:goto L_125E3;
case 0x25e5:goto L_125E5;
case 0x25e8:goto L_125E8;
case 0x25e9:goto L_125E9;
case 0x25ec:goto L_125EC;
case 0x25ee:goto L_125EE;
case 0x25f2:goto L_125F2;
case 0x25f4:goto L_125F4;
case 0x25f6:goto L_125F6;
case 0x25fb:goto L_125FB;
case 0x25fd:goto L_125FD;
case 0x25ff:goto L_125FF;
case 0x2600:goto L_12600;
case 0x2606:goto L_12606;
case 0x2608:goto L_12608;
case 0x2609:goto L_12609;
case 0x260c:goto L_1260C;
case 0x2611:goto L_12611;
case 0x2613:goto L_12613;
case 0x2615:goto L_12615;
case 0x2617:goto L_12617;
case 0x2619:goto L_12619;
case 0x261b:goto L_1261B;
case 0x261d:goto L_1261D;
case 0x2620:goto L_12620;
case 0x2623:goto L_12623;
case 0x2624:goto L_12624;
case 0x2626:goto L_12626;
case 0x2628:goto L_12628;
case 0x2629:goto L_12629;
case 0x262b:goto L_1262B;
case 0x262e:goto L_1262E;
case 0x2631:goto L_12631;
case 0x2634:goto L_12634;
case 0x2636:goto L_12636;
case 0x2638:goto L_12638;
case 0x263a:goto L_1263A;
case 0x263d:goto L_1263D;
case 0x2641:goto L_12641;
case 0x2644:goto L_12644;
case 0x2646:goto L_12646;
case 0x2649:goto L_12649;
case 0x264a:goto L_1264A;
case 0x264d:goto L_1264D;
case 0x264f:goto L_1264F;
case 0x2653:goto L_12653;
case 0x2657:goto L_12657;
case 0x2658:goto L_12658;
case 0x265b:goto L_1265B;
case 0x265d:goto L_1265D;
case 0x2661:goto L_12661;
case 0x2662:goto L_12662;
case 0x2665:goto L_12665;
case 0x2668:goto L_12668;
case 0x266a:goto L_1266A;
case 0x266e:goto L_1266E;
case 0x2671:goto L_12671;
case 0x2674:goto L_12674;
case 0x2675:goto L_12675;
case 0x2678:goto L_12678;
case 0x267a:goto L_1267A;
case 0x267d:goto L_1267D;
case 0x2681:goto L_12681;
case 0x2683:goto L_12683;
case 0x2686:goto L_12686;
case 0x268b:goto L_1268B;
case 0x268d:goto L_1268D;
case 0x2690:goto L_12690;
case 0x2692:goto L_12692;
case 0x2694:goto L_12694;
case 0x2697:goto L_12697;
case 0x2699:goto L_12699;
case 0x269c:goto L_1269C;
case 0x269e:goto L_1269E;
case 0x26a0:goto L_126A0;
case 0x26a3:goto L_126A3;
case 0x26a5:goto L_126A5;
case 0x26a8:goto L_126A8;
case 0x26ab:goto L_126AB;
case 0x26ad:goto L_126AD;
case 0x26b0:goto L_126B0;
case 0x26b2:goto L_126B2;
case 0x26b4:goto L_126B4;
case 0x26b9:goto L_126B9;
case 0x26bb:goto L_126BB;
case 0x26be:goto L_126BE;
case 0x26c0:goto L_126C0;
case 0x26c3:goto L_126C3;
case 0x26c6:goto L_126C6;
case 0x26c8:goto L_126C8;
case 0x26cb:goto L_126CB;
case 0x26cf:goto L_126CF;
case 0x26d1:goto L_126D1;
case 0x26d3:goto L_126D3;
case 0x26d5:goto L_126D5;
case 0x26d7:goto L_126D7;
case 0x26da:goto L_126DA;
case 0x26dc:goto L_126DC;
case 0x26e0:goto L_126E0;
case 0x26e2:goto L_126E2;
case 0x26e4:goto L_126E4;
case 0x26e8:goto L_126E8;
case 0x26ea:goto L_126EA;
case 0x26ed:goto L_126ED;
case 0x26f0:goto L_126F0;
case 0x26f2:goto L_126F2;
case 0x26f5:goto L_126F5;
case 0x26f7:goto L_126F7;
case 0x26fa:goto L_126FA;
case 0x26fc:goto L_126FC;
case 0x26fe:goto L_126FE;
case 0x26ff:goto L_126FF;
case 0x2702:goto L_12702;
case 0x2705:goto L_12705;
case 0x2706:goto L_12706;
case 0x2708:goto L_12708;
case 0x270b:goto L_1270B;
case 0x270f:goto L_1270F;
case 0x2711:goto L_12711;
case 0x2713:goto L_12713;
case 0x2715:goto L_12715;
case 0x2717:goto L_12717;
case 0x2719:goto L_12719;
case 0x271b:goto L_1271B;
case 0x271d:goto L_1271D;
case 0x271f:goto L_1271F;
case 0x2721:goto L_12721;
case 0x2726:goto L_12726;
case 0x2728:goto L_12728;
case 0x272b:goto L_1272B;
case 0x272d:goto L_1272D;
case 0x2732:goto L_12732;
case 0x2734:goto L_12734;
case 0x2736:goto L_12736;
case 0x2739:goto L_12739;
case 0x273b:goto L_1273B;
case 0x273c:goto L_1273C;
case 0x273e:goto L_1273E;
case 0x2743:goto L_12743;
case 0x2745:goto L_12745;
case 0x2747:goto L_12747;
case 0x274c:goto L_1274C;
case 0x274f:goto L_1274F;
case 0x2751:goto L_12751;
case 0x2753:goto L_12753;
case 0x2755:goto L_12755;
case 0x2757:goto L_12757;
case 0x2759:goto L_12759;
case 0x275c:goto L_1275C;
case 0x275e:goto L_1275E;
case 0x275f:goto L_1275F;
case 0x2764:goto L_12764;
case 0x2767:goto L_12767;
case 0x276a:goto L_1276A;
case 0x276d:goto L_1276D;
case 0x2771:goto L_12771;
case 0x2774:goto L_12774;
case 0x2777:goto L_12777;
case 0x2779:goto L_12779;
case 0x277b:goto L_1277B;
case 0x277d:goto L_1277D;
case 0x277f:goto L_1277F;
case 0x2782:goto L_12782;
case 0x2783:goto L_12783;
case 0x2786:goto L_12786;
case 0x2788:goto L_12788;
case 0x278c:goto L_1278C;
case 0x278e:goto L_1278E;
case 0x2791:goto L_12791;
case 0x2794:goto L_12794;
case 0x2795:goto L_12795;
case 0x2799:goto L_12799;
case 0x279b:goto L_1279B;
case 0x279e:goto L_1279E;
case 0x27a1:goto L_127A1;
case 0x27a2:goto L_127A2;
case 0x27a5:goto L_127A5;
case 0x27a9:goto L_127A9;
case 0x27ab:goto L_127AB;
case 0x27af:goto L_127AF;
case 0x27b1:goto L_127B1;
case 0x27b5:goto L_127B5;
case 0x27b8:goto L_127B8;
case 0x27bc:goto L_127BC;
case 0x27be:goto L_127BE;
case 0x27c0:goto L_127C0;
case 0x27c2:goto L_127C2;
case 0x27c6:goto L_127C6;
case 0x27ca:goto L_127CA;
case 0x27cd:goto L_127CD;
case 0x27d0:goto L_127D0;
case 0x27d2:goto L_127D2;
case 0x27d4:goto L_127D4;
case 0x27d6:goto L_127D6;
case 0x27d8:goto L_127D8;
case 0x27da:goto L_127DA;
case 0x27dc:goto L_127DC;
case 0x27de:goto L_127DE;
case 0x27e3:goto L_127E3;
case 0x27e6:goto L_127E6;
case 0x27e9:goto L_127E9;
case 0x27ec:goto L_127EC;
case 0x27ef:goto L_127EF;
case 0x27f1:goto L_127F1;
case 0x27f5:goto L_127F5;
case 0x27f6:goto L_127F6;
case 0x27f9:goto L_127F9;
case 0x27fb:goto L_127FB;
case 0x27fc:goto L_127FC;
case 0x27ff:goto L_127FF;
case 0x2803:goto L_12803;
case 0x2804:goto L_12804;
case 0x2807:goto L_12807;
case 0x2808:goto L_12808;
case 0x280b:goto L_1280B;
case 0x280d:goto L_1280D;
case 0x2810:goto L_12810;
case 0x2813:goto L_12813;
case 0x2815:goto L_12815;
case 0x2818:goto L_12818;
case 0x281a:goto L_1281A;
case 0x281d:goto L_1281D;
case 0x281e:goto L_1281E;
case 0x2821:goto L_12821;
case 0x2825:goto L_12825;
case 0x2827:goto L_12827;
case 0x2829:goto L_12829;
case 0x282b:goto L_1282B;
case 0x282d:goto L_1282D;
case 0x2830:goto L_12830;
case 0x2831:goto L_12831;
case 0x2832:goto L_12832;
case 0x2833:goto L_12833;
case 0x2834:goto L_12834;
case 0x2837:goto L_12837;
case 0x283a:goto L_1283A;
case 0x283d:goto L_1283D;
case 0x2840:goto L_12840;
case 0x2842:goto L_12842;
case 0x2845:goto L_12845;
case 0x2847:goto L_12847;
case 0x284a:goto L_1284A;
case 0x284c:goto L_1284C;
case 0x284f:goto L_1284F;
case 0x2851:goto L_12851;
case 0x2853:goto L_12853;
case 0x2856:goto L_12856;
case 0x2859:goto L_12859;
case 0x285b:goto L_1285B;
case 0x285e:goto L_1285E;
case 0x2862:goto L_12862;
case 0x2864:goto L_12864;
case 0x2866:goto L_12866;
case 0x286a:goto L_1286A;
case 0x286c:goto L_1286C;
case 0x286e:goto L_1286E;
case 0x2871:goto L_12871;
case 0x2873:goto L_12873;
case 0x2876:goto L_12876;
case 0x2877:goto L_12877;
case 0x2878:goto L_12878;
case 0x2879:goto L_12879;
case 0x287a:goto L_1287A;
case 0x287b:goto L_1287B;
case 0x287c:goto L_1287C;
case 0x287d:goto L_1287D;
case 0x287e:goto L_1287E;
case 0x287f:goto L_1287F;
case 0x2880:goto L_12880;
case 0x2881:goto L_12881;
case 0x2882:goto L_12882;
case 0x2883:goto L_12883;
case 0x2886:goto L_12886;
case 0x288a:goto L_1288A;
case 0x288c:goto L_1288C;
case 0x2890:goto L_12890;
case 0x2892:goto L_12892;
case 0x2896:goto L_12896;
case 0x2898:goto L_12898;
case 0x289a:goto L_1289A;
case 0x289d:goto L_1289D;
case 0x28a1:goto L_128A1;
case 0x28a3:goto L_128A3;
case 0x28a6:goto L_128A6;
case 0x28aa:goto L_128AA;
case 0x28ac:goto L_128AC;
case 0x28ae:goto L_128AE;
case 0x28b2:goto L_128B2;
case 0x28b4:goto L_128B4;
case 0x28b6:goto L_128B6;
case 0x28b9:goto L_128B9;
case 0x28bb:goto L_128BB;
case 0x28bf:goto L_128BF;
case 0x28c2:goto L_128C2;
case 0x28c3:goto L_128C3;
case 0x28c4:goto L_128C4;
case 0x28c5:goto L_128C5;
case 0x28c6:goto L_128C6;
case 0x28c7:goto L_128C7;
case 0x28c8:goto L_128C8;
case 0x28c9:goto L_128C9;
case 0x28ca:goto L_128CA;
case 0x28cb:goto L_128CB;
case 0x28f4:goto L_128F4;
case 0x28f6:goto L_128F6;
case 0x28f8:goto L_128F8;
case 0x28fa:goto L_128FA;
case 0x28fd:goto L_128FD;
case 0x2901:goto L_12901;
case 0x2903:goto L_12903;
case 0x2906:goto L_12906;
case 0x2908:goto L_12908;
case 0x290a:goto L_1290A;
case 0x290c:goto L_1290C;
case 0x290e:goto L_1290E;
case 0x2910:goto L_12910;
case 0x2912:goto L_12912;
case 0x2915:goto L_12915;
case 0x2916:goto L_12916;
case 0x2917:goto L_12917;
case 0x2918:goto L_12918;
case 0x2a7e:goto L_12A7E;
case 0x2a82:goto L_12A82;
case 0x2a84:goto L_12A84;
case 0x2a85:goto L_12A85;
case 0x2a8a:goto L_12A8A;
case 0x2a8c:goto L_12A8C;
case 0x2a8e:goto L_12A8E;
case 0x2a93:goto L_12A93;
case 0x2a97:goto L_12A97;
case 0x2a99:goto L_12A99;
case 0x2a9b:goto L_12A9B;
case 0x2a9d:goto L_12A9D;
case 0x2aa0:goto L_12AA0;
case 0x2aa2:goto L_12AA2;
case 0x2aa7:goto L_12AA7;
case 0x2aa8:goto L_12AA8;
case 0x2aac:goto L_12AAC;
case 0x2ab0:goto L_12AB0;
case 0x2ab2:goto L_12AB2;
case 0x2ab6:goto L_12AB6;
case 0x2ab7:goto L_12AB7;
case 0x2ab9:goto L_12AB9;
case 0x2abc:goto L_12ABC;
case 0x2abf:goto L_12ABF;
case 0x2ac1:goto L_12AC1;
case 0x2ac4:goto L_12AC4;
case 0x2ac5:goto L_12AC5;
case 0x2ac8:goto L_12AC8;
case 0x2acb:goto L_12ACB;
case 0x2ace:goto L_12ACE;
case 0x2ad1:goto L_12AD1;
case 0x42ab:goto L_142AB;
case 0x42ac:goto L_142AC;
case 0x42ad:goto L_142AD;
case 0x42af:goto L_142AF;
case 0x42b1:goto L_142B1;
case 0x42b5:goto L_142B5;
case 0x42b7:goto L_142B7;
case 0x42b8:goto L_142B8;
case 0x42b9:goto L_142B9;
case 0x42bd:goto L_142BD;
case 0x42bf:goto L_142BF;
case 0x42c1:goto L_142C1;
case 0x42c4:goto L_142C4;
case 0x42c8:goto L_142C8;
case 0x42ca:goto L_142CA;
case 0x42cc:goto L_142CC;
case 0x42ce:goto L_142CE;
case 0x42d0:goto L_142D0;
case 0x42d2:goto L_142D2;
case 0x42d4:goto L_142D4;
case 0x42d6:goto L_142D6;
case 0x42d8:goto L_142D8;
case 0x42da:goto L_142DA;
case 0x42dc:goto L_142DC;
case 0x42de:goto L_142DE;
case 0x42e2:goto L_142E2;
case 0x42e3:goto L_142E3;
case 0x42e5:goto L_142E5;
case 0x42e7:goto L_142E7;
case 0x42ea:goto L_142EA;
case 0x42ee:goto L_142EE;
case 0x42f1:goto L_142F1;
case 0x42f3:goto L_142F3;
case 0x42f5:goto L_142F5;
case 0x42f7:goto L_142F7;
case 0x42fa:goto L_142FA;
case 0x42fd:goto L_142FD;
case 0x42fe:goto L_142FE;
case 0x42ff:goto L_142FF;
case 0x4300:goto L_14300;
case 0x4302:goto L_14302;
case 0x4304:goto L_14304;
case 0x4309:goto L_14309;
case 0x430b:goto L_1430B;
case 0x430f:goto L_1430F;
case 0x4311:goto L_14311;
case 0x4316:goto L_14316;
case 0x4318:goto L_14318;
case 0x431a:goto L_1431A;
case 0x431c:goto L_1431C;
case 0x431e:goto L_1431E;
case 0x4321:goto L_14321;
case 0x4322:goto L_14322;
case 0x4323:goto L_14323;
case 0x4324:goto L_14324;
case 0x562b:goto L_1562B;
case 0x562c:goto L_1562C;
case 0x562e:goto L_1562E;
case 0x5632:goto L_15632;
case 0x5634:goto L_15634;
case 0x5636:goto L_15636;
case 0x5639:goto L_15639;
case 0x563a:goto L_1563A;
default:break;}
if(!invoke_external)return;
ki_engagement_invoke(m,m->ip,h);invoke_external=0;goto AM_DISPATCH;
L_125A3: ; /* IDA 0x125A3; mov     ds, cs:word_10D52 */
m->ip=0x25a3;
/* Original entry: sub_125A3 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x25a3,h->user);
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_125A8;
L_125A8: ; /* IDA 0x125A8; mov     es, cs:word_19874 */
m->ip=0x25a8;
skip_first_enter=0;
m->es=(uint16_t)(dl_read16(m,m->cs,39028));
goto L_125AD;
L_125AD: ; /* IDA 0x125AD; mov     si, cs:word_10D18 */
m->ip=0x25ad;
skip_first_enter=0;
m->si=(uint16_t)(dl_read16(m,m->cs,3352));
goto L_125B2;
L_125B2: ; /* IDA 0x125B2; mov     cx, 10h */
m->ip=0x25b2;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 1
m->cx=15;
#else
m->cx=(uint16_t)(16);
#endif
goto L_125B5;
L_125B5: ; /* IDA 0x125B5; push    cx */
m->ip=0x25b5;
skip_first_enter=0;
ec_push(m,m->cx);
goto L_125B6;
L_125B6: ; /* IDA 0x125B6; cmp     byte ptr [si+2240h], 80h */
m->ip=0x25b6;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+8768)),128,8);
goto L_125BB;
L_125BB: ; /* IDA 0x125BB; jb      short loc_125DE */
m->ip=0x25bb;
skip_first_enter=0;
if(m->flags&1) {goto L_125DE;}
goto L_125BD;
L_125BD: ; /* IDA 0x125BD; add     si, 2240h */
m->ip=0x25bd;
skip_first_enter=0;
m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,0));
goto L_125C1;
L_125C1: ; /* IDA 0x125C1; dec     byte ptr [si+0Bh] */
m->ip=0x25c1;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+11),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->si+11)),8));
goto L_125C4;
L_125C4: ; /* IDA 0x125C4; jnz     short loc_125D2 */
m->ip=0x25c4;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_125D2;}
goto L_125C6;
L_125C6: ; /* IDA 0x125C6; mov     al, [si+1Eh] */
m->ip=0x25c6;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+30))));
goto L_125C9;
L_125C9: ; /* IDA 0x125C9; mov     [si+0Bh], al */
m->ip=0x25c9;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 2
; /* omitted original timer reload */
#else
dl_write8(m,m->ds,(uint16_t)(m->si+11),(uint8_t)m->ax);
#endif
goto L_125CC;
L_125CC: ; /* IDA 0x125CC; and     byte ptr [si], 0DFh */
m->ip=0x25cc;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(223))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
goto L_125CF;
L_125CF: ; /* IDA 0x125CF; call    sub_12662 */
m->ip=0x25cf;
skip_first_enter=0;
{uint16_t target=9826;ec_push(m,0x25d2);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_125D2: ; /* IDA 0x125D2; call    sub_12600 */
m->ip=0x25d2;
skip_first_enter=0;
{uint16_t target=9728;ec_push(m,0x25d5);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_125D5: ; /* IDA 0x125D5; call    sub_1264A */
m->ip=0x25d5;
skip_first_enter=0;
{uint16_t target=9802;ec_push(m,0x25d8);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_125D8: ; /* IDA 0x125D8; sub     si, 2240h */
m->ip=0x25d8;
skip_first_enter=0;
m->si=(uint16_t)(ec_math(m,m->si,8768,16,0,1));
goto L_125DC;
L_125DC: ; /* IDA 0x125DC; jmp     short loc_125E8 */
m->ip=0x25dc;
skip_first_enter=0;
goto L_125E8;
L_125DE: ; /* IDA 0x125DE; test    byte ptr [si+2240h], 8 */
m->ip=0x25de;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si+8768)))&(8))&255;ec_logic(m,result,8);}
goto L_125E3;
L_125E3: ; /* IDA 0x125E3; jz      short loc_125E8 */
m->ip=0x25e3;
skip_first_enter=0;
if(m->flags&0x40) {goto L_125E8;}
goto L_125E5;
L_125E5: ; /* IDA 0x125E5; call    sub_12A7E */
m->ip=0x25e5;
skip_first_enter=0;
{uint16_t target=10878;ec_push(m,0x25e8);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_125E8: ; /* IDA 0x125E8; pop     cx */
m->ip=0x25e8;
skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_125E9;
L_125E9: ; /* IDA 0x125E9; add     si, 40h ; '@' */
m->ip=0x25e9;
skip_first_enter=0;
m->si=(uint16_t)(ec_math(m,m->si,64,16,0,0));
goto L_125EC;
L_125EC: ; /* IDA 0x125EC; loop    loc_125B5 */
m->ip=0x25ec;
skip_first_enter=0;
if(--m->cx) {goto L_125B5;}
goto L_125EE;
L_125EE: ; /* IDA 0x125EE; cmp     si, 1FC0h */
m->ip=0x25ee;
skip_first_enter=0;
ec_cmp(m,m->si,8128,16);
goto L_125F2;
L_125F2: ; /* IDA 0x125F2; jb      short loc_125F6 */
m->ip=0x25f2;
skip_first_enter=0;
if(m->flags&1) {goto L_125F6;}
goto L_125F4;
L_125F4: ; /* IDA 0x125F4; xor     si, si */
m->ip=0x25f4;
skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_125F6;
L_125F6: ; /* IDA 0x125F6; mov     cs:word_10D18, si */
m->ip=0x25f6;
skip_first_enter=0;
dl_write16(m,m->cs,3352,m->si);
goto L_125FB;
L_125FB: ; /* IDA 0x125FB; mov     ax, cs */
m->ip=0x25fb;
skip_first_enter=0;
m->ax=(uint16_t)(m->cs);
goto L_125FD;
L_125FD: ; /* IDA 0x125FD; mov     ds, ax */
m->ip=0x25fd;
skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_125FF;
L_125FF: ; /* IDA 0x125FF; retn */
m->ip=0x25ff;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12600: ; /* IDA 0x12600; cmp     cs:byte_10CF3, 1 */
m->ip=0x2600;
/* Original entry: sub_12600 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x2600,h->user);
ec_cmp(m,dl_read8(m,m->cs,3315),1,8);
goto L_12606;
L_12606: ; /* IDA 0x12606; jz      short loc_12609 */
m->ip=0x2606;
skip_first_enter=0;
if(m->flags&0x40) {goto L_12609;}
goto L_12608;
L_12608: ; /* IDA 0x12608; retn */
m->ip=0x2608;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12609: ; /* IDA 0x12609; mov     ax, [si+4] */
m->ip=0x2609;
skip_first_enter=0;
m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+4)));
goto L_1260C;
L_1260C: ; /* IDA 0x1260C; cmp     word ptr [si+0Eh], 800h */
m->ip=0x260c;
skip_first_enter=0;
ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+14)),2048,16);
goto L_12611;
L_12611: ; /* IDA 0x12611; jb      short loc_12624 */
m->ip=0x2611;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 3
if(!(m->flags&1)) {goto L_12624;}
#else
if(m->flags&1) {goto L_12624;}
#endif
goto L_12613;
L_12613: ; /* IDA 0x12613; shr     ax, 1 */
m->ip=0x2613;
skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
goto L_12615;
L_12615: ; /* IDA 0x12615; mov     dx, ax */
m->ip=0x2615;
skip_first_enter=0;
m->dx=(uint16_t)(m->ax);
goto L_12617;
L_12617: ; /* IDA 0x12617; shr     ax, 1 */
m->ip=0x2617;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 4
; /* omitted original quarter shift */
#else
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
#endif
goto L_12619;
L_12619: ; /* IDA 0x12619; add     ax, dx */
m->ip=0x2619;
skip_first_enter=0;
m->ax=(uint16_t)(ec_math(m,m->ax,m->dx,16,0,0));
goto L_1261B;
L_1261B: ; /* IDA 0x1261B; xor     dl, dl */
m->ip=0x261b;
skip_first_enter=0;
{uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
goto L_1261D;
L_1261D: ; /* IDA 0x1261D; mov     dh, [si+1] */
m->ip=0x261d;
skip_first_enter=0;
vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_12620;
L_12620: ; /* IDA 0x12620; call    sub_1562B */
m->ip=0x2620;
skip_first_enter=0;
{uint16_t target=22059;ec_push(m,0x2623);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12623: ; /* IDA 0x12623; retn */
m->ip=0x2623;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12624: ; /* IDA 0x12624; mov     cl, 5 */
m->ip=0x2624;
skip_first_enter=0;
vg_cl(m,(uint8_t)(5));
goto L_12626;
L_12626: ; /* IDA 0x12626; shr     ax, cl */
m->ip=0x2626;
skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,(uint8_t)m->cx,16,5));
goto L_12628;
L_12628: ; /* IDA 0x12628; inc     ax */
m->ip=0x2628;
skip_first_enter=0;
m->ax=(uint16_t)(ec_inc(m,m->ax,16));
goto L_12629;
L_12629: ; /* IDA 0x12629; xor     dl, dl */
m->ip=0x2629;
skip_first_enter=0;
{uint16_t result=(((uint8_t)m->dx)^((uint8_t)m->dx))&255;ec_logic(m,result,8);vg_dl(m,(uint8_t)(result));}
goto L_1262B;
L_1262B: ; /* IDA 0x1262B; mov     dh, [si+1] */
m->ip=0x262b;
skip_first_enter=0;
vg_dh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_1262E;
L_1262E: ; /* IDA 0x1262E; call    sub_1562B */
m->ip=0x262e;
skip_first_enter=0;
{uint16_t target=22059;ec_push(m,0x2631);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12631: ; /* IDA 0x12631; mov     bh, [si+1] */
m->ip=0x2631;
skip_first_enter=0;
gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_12634;
L_12634: ; /* IDA 0x12634; xor     bl, bl */
m->ip=0x2634;
skip_first_enter=0;
{uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
goto L_12636;
L_12636: ; /* IDA 0x12636; shr     bx, 1 */
m->ip=0x2636;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_12638;
L_12638: ; /* IDA 0x12638; shr     bx, 1 */
m->ip=0x2638;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_1263A;
L_1263A: ; /* IDA 0x1263A; mov     al, [bx+1Dh] */
m->ip=0x263a;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+29))));
goto L_1263D;
L_1263D: ; /* IDA 0x1263D; add     byte ptr [si+6], 0Ah */
m->ip=0x263d;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 5
dl_write8(m,m->ds,(uint16_t)(m->si+6),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->si+6)),9,8,0,0));
#else
dl_write8(m,m->ds,(uint16_t)(m->si+6),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->si+6)),10,8,0,0));
#endif
goto L_12641;
L_12641: ; /* IDA 0x12641; cmp     [si+6], al */
m->ip=0x2641;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+6)),(uint8_t)m->ax,8);
goto L_12644;
L_12644: ; /* IDA 0x12644; jb      short locret_12649 */
m->ip=0x2644;
skip_first_enter=0;
if(m->flags&1) {goto L_12649;}
goto L_12646;
L_12646: ; /* IDA 0x12646; mov     [si+6], al */
m->ip=0x2646;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+6),(uint8_t)m->ax);
goto L_12649;
L_12649: ; /* IDA 0x12649; retn */
m->ip=0x2649;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_1264A: ; /* IDA 0x1264A; test    byte ptr [si], 20h */
m->ip=0x264a;
/* Original entry: sub_1264A */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x264a,h->user);
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(32))&255;ec_logic(m,result,8);}
goto L_1264D;
L_1264D: ; /* IDA 0x1264D; jnz     short loc_12658 */
m->ip=0x264d;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_12658;}
goto L_1264F;
L_1264F: ; /* IDA 0x1264F; mov     byte ptr [si+3], 0 */
m->ip=0x264f;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+3),0);
goto L_12653;
L_12653: ; /* IDA 0x12653; mov     byte ptr [si+21h], 0 */
m->ip=0x2653;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+33),0);
goto L_12657;
L_12657: ; /* IDA 0x12657; retn */
m->ip=0x2657;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12658: ; /* IDA 0x12658; dec     byte ptr [si+3] */
m->ip=0x2658;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+3),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->si+3)),8));
goto L_1265B;
L_1265B: ; /* IDA 0x1265B; jnz     short locret_12661 */
m->ip=0x265b;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_12661;}
goto L_1265D;
L_1265D: ; /* IDA 0x1265D; mov     byte ptr [si+3], 1 */
m->ip=0x265d;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+3),1);
goto L_12661;
L_12661: ; /* IDA 0x12661; retn */
m->ip=0x2661;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12662: ; /* IDA 0x12662; mov     bx, [si+0Eh] */
m->ip=0x2662;
/* Original entry: sub_12662 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x2662,h->user);
m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
goto L_12665;
L_12665: ; /* IDA 0x12665; cmp     bx, [si+14h] */
m->ip=0x2665;
skip_first_enter=0;
ec_cmp(m,m->bx,dl_read16(m,m->ds,(uint16_t)(m->si+20)),16);
goto L_12668;
L_12668: ; /* IDA 0x12668; jnz     short loc_12675 */
m->ip=0x2668;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_12675;}
goto L_1266A;
L_1266A: ; /* IDA 0x1266A; mov     byte ptr [si+8], 4 */
m->ip=0x266a;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+8),4);
goto L_1266E;
L_1266E: ; /* IDA 0x1266E; call    sub_128F4 */
m->ip=0x266e;
skip_first_enter=0;
{uint16_t target=10484;ec_push(m,0x2671);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12671: ; /* IDA 0x12671; call    sub_14325 */
m->ip=0x2671;
skip_first_enter=0;
{uint16_t target=17189;ec_push(m,0x2674);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12674: ; /* IDA 0x12674; retn */
m->ip=0x2674;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12675: ; /* IDA 0x12675; test    byte ptr [si], 10h */
m->ip=0x2675;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(16))&255;ec_logic(m,result,8);}
goto L_12678;
L_12678: ; /* IDA 0x12678; jz      short loc_1267D */
m->ip=0x2678;
skip_first_enter=0;
if(m->flags&0x40) {goto L_1267D;}
goto L_1267A;
L_1267A: ; /* IDA 0x1267A; call    sub_12BA8 */
m->ip=0x267a;
skip_first_enter=0;
{uint16_t target=11176;ec_push(m,0x267d);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_1267D: ; /* IDA 0x1267D; cmp     bx, 800h */
m->ip=0x267d;
skip_first_enter=0;
ec_cmp(m,m->bx,2048,16);
goto L_12681;
L_12681: ; /* IDA 0x12681; jnb     short loc_12694 */
m->ip=0x2681;
skip_first_enter=0;
if(!(m->flags&1)) {goto L_12694;}
goto L_12683;
L_12683: ; /* IDA 0x12683; mov     al, [si+1] */
m->ip=0x2683;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_12686;
L_12686: ; /* IDA 0x12686; cmp     al, cs:byte_10CFF */
m->ip=0x2686;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->cs,3327),8);
goto L_1268B;
L_1268B: ; /* IDA 0x1268B; jz      short loc_12697 */
m->ip=0x268b;
skip_first_enter=0;
if(m->flags&0x40) {goto L_12697;}
goto L_1268D;
L_1268D: ; /* IDA 0x1268D; call    sub_14300 */
m->ip=0x268d;
skip_first_enter=0;
{uint16_t target=17152;ec_push(m,0x2690);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12690: ; /* IDA 0x12690; jb      short loc_1266A */
m->ip=0x2690;
skip_first_enter=0;
if(m->flags&1) {goto L_1266A;}
goto L_12692;
L_12692: ; /* IDA 0x12692; jmp     short loc_12697 */
m->ip=0x2692;
skip_first_enter=0;
goto L_12697;
L_12694: ; /* IDA 0x12694; call    sub_142AB */
m->ip=0x2694;
skip_first_enter=0;
{uint16_t target=17067;ec_push(m,0x2697);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12697: ; /* IDA 0x12697; mov     dx, ds */
m->ip=0x2697;
skip_first_enter=0;
m->dx=(uint16_t)(m->ds);
goto L_12699;
L_12699: ; /* IDA 0x12699; lds     di, [si+1Ah] */
m->ip=0x2699;
skip_first_enter=0;
{uint16_t off=(uint16_t)(m->si+26),seg=m->ds;uint16_t value=dl_read16(m,seg,off),bank=dl_read16(m,seg,(uint16_t)(off+2));m->di=(uint16_t)(value);m->ds=bank;}
goto L_1269C;
L_1269C: ; /* IDA 0x1269C; dec     byte ptr [di] */
m->ip=0x269c;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 7
; /* omitted original occupancy decrement */
#else
dl_write8(m,m->ds,(uint16_t)(m->di),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->di)),8));
#endif
goto L_1269E;
L_1269E: ; /* IDA 0x1269E; mov     ds, dx */
m->ip=0x269e;
skip_first_enter=0;
m->ds=(uint16_t)(m->dx);
goto L_126A0;
L_126A0: ; /* IDA 0x126A0; test    byte ptr [si], 2 */
m->ip=0x26a0;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(2))&255;ec_logic(m,result,8);}
goto L_126A3;
L_126A3: ; /* IDA 0x126A3; jz      short loc_126B4 */
m->ip=0x26a3;
skip_first_enter=0;
if(m->flags&0x40) {goto L_126B4;}
goto L_126A5;
L_126A5: ; /* IDA 0x126A5; and     byte ptr [si], 0FDh */
m->ip=0x26a5;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(253))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
goto L_126A8;
L_126A8: ; /* IDA 0x126A8; call    sub_147BB */
m->ip=0x26a8;
skip_first_enter=0;
{uint16_t target=18363;ec_push(m,0x26ab);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_126AB: ; /* IDA 0x126AB; jb      short loc_126F5 */
m->ip=0x26ab;
skip_first_enter=0;
if(m->flags&1) {goto L_126F5;}
goto L_126AD;
L_126AD: ; /* IDA 0x126AD; test    byte ptr [si], 1 */
m->ip=0x26ad;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(1))&255;ec_logic(m,result,8);}
goto L_126B0;
L_126B0: ; /* IDA 0x126B0; jz      short loc_126C0 */
m->ip=0x26b0;
skip_first_enter=0;
if(m->flags&0x40) {goto L_126C0;}
goto L_126B2;
L_126B2: ; /* IDA 0x126B2; jmp     short loc_126C8 */
m->ip=0x26b2;
skip_first_enter=0;
goto L_126C8;
L_126B4: ; /* IDA 0x126B4; cmp     word ptr [si+0Eh], 800h */
m->ip=0x26b4;
skip_first_enter=0;
ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+14)),2048,16);
goto L_126B9;
L_126B9: ; /* IDA 0x126B9; jnb     short loc_126C8 */
m->ip=0x26b9;
skip_first_enter=0;
if(!(m->flags&1)) {goto L_126C8;}
goto L_126BB;
L_126BB: ; /* IDA 0x126BB; call    sub_147BB */
m->ip=0x26bb;
skip_first_enter=0;
{uint16_t target=18363;ec_push(m,0x26be);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_126BE: ; /* IDA 0x126BE; jb      short loc_126F5 */
m->ip=0x26be;
skip_first_enter=0;
if(m->flags&1) {goto L_126F5;}
goto L_126C0;
L_126C0: ; /* IDA 0x126C0; mov     bx, [si+0Ch] */
m->ip=0x26c0;
skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+12)));
goto L_126C3;
L_126C3: ; /* IDA 0x126C3; call    sub_12708 */
m->ip=0x26c3;
skip_first_enter=0;
{uint16_t target=9992;ec_push(m,0x26c6);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_126C6: ; /* IDA 0x126C6; jmp     short loc_126F5 */
m->ip=0x26c6;
skip_first_enter=0;
goto L_126F5;
L_126C8: ; /* IDA 0x126C8; mov     bx, [si+0Ch] */
m->ip=0x26c8;
skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+12)));
goto L_126CB;
L_126CB: ; /* IDA 0x126CB; mov     al, es:[bx+3] */
m->ip=0x26cb;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->es,(uint16_t)(m->bx+3))));
goto L_126CF;
L_126CF: ; /* IDA 0x126CF; mov     ah, al */
m->ip=0x26cf;
skip_first_enter=0;
vg_ah(m,(uint8_t)((uint8_t)m->ax));
goto L_126D1;
L_126D1: ; /* IDA 0x126D1; and     al, 7 */
m->ip=0x26d1;
skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)&(7))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_126D3;
L_126D3: ; /* IDA 0x126D3; cmp     al, 2 */
m->ip=0x26d3;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,2,8);
goto L_126D5;
L_126D5: ; /* IDA 0x126D5; jb      short loc_126F2 */
m->ip=0x26d5;
skip_first_enter=0;
if(m->flags&1) {goto L_126F2;}
goto L_126D7;
L_126D7: ; /* IDA 0x126D7; test    ah, 40h */
m->ip=0x26d7;
skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))&(64))&255;ec_logic(m,result,8);}
goto L_126DA;
L_126DA: ; /* IDA 0x126DA; jz      short loc_126E4 */
m->ip=0x26da;
skip_first_enter=0;
if(m->flags&0x40) {goto L_126E4;}
goto L_126DC;
L_126DC: ; /* IDA 0x126DC; cmp     byte ptr [si+0Ah], 0 */
m->ip=0x26dc;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+10)),0,8);
goto L_126E0;
L_126E0: ; /* IDA 0x126E0; jg      short loc_126F2 */
m->ip=0x26e0;
skip_first_enter=0;
if(ec_greater(m)) {goto L_126F2;}
goto L_126E2;
L_126E2: ; /* IDA 0x126E2; jmp     short loc_126EA */
m->ip=0x26e2;
skip_first_enter=0;
goto L_126EA;
L_126E4: ; /* IDA 0x126E4; cmp     byte ptr [si+0Ah], 0 */
m->ip=0x26e4;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+10)),0,8);
goto L_126E8;
L_126E8: ; /* IDA 0x126E8; jl      short loc_126F2 */
m->ip=0x26e8;
skip_first_enter=0;
if(ec_less(m)) {goto L_126F2;}
goto L_126EA;
L_126EA: ; /* IDA 0x126EA; call    sub_127F6 */
m->ip=0x26ea;
skip_first_enter=0;
{uint16_t target=10230;ec_push(m,0x26ed);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_126ED: ; /* IDA 0x126ED; call    sub_127A2 */
m->ip=0x26ed;
skip_first_enter=0;
{uint16_t target=10146;ec_push(m,0x26f0);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_126F0: ; /* IDA 0x126F0; jmp     short loc_126F5 */
m->ip=0x26f0;
skip_first_enter=0;
goto L_126F5;
L_126F2: ; /* IDA 0x126F2; call    sub_126FF */
m->ip=0x26f2;
skip_first_enter=0;
{uint16_t target=9983;ec_push(m,0x26f5);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_126F5: ; /* IDA 0x126F5; mov     dx, ds */
m->ip=0x26f5;
skip_first_enter=0;
m->dx=(uint16_t)(m->ds);
goto L_126F7;
L_126F7: ; /* IDA 0x126F7; lds     di, [si+1Ah] */
m->ip=0x26f7;
skip_first_enter=0;
{uint16_t off=(uint16_t)(m->si+26),seg=m->ds;uint16_t value=dl_read16(m,seg,off),bank=dl_read16(m,seg,(uint16_t)(off+2));m->di=(uint16_t)(value);m->ds=bank;}
goto L_126FA;
L_126FA: ; /* IDA 0x126FA; inc     byte ptr [di] */
m->ip=0x26fa;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->di),ec_inc(m,dl_read8(m,m->ds,(uint16_t)(m->di)),8));
goto L_126FC;
L_126FC: ; /* IDA 0x126FC; mov     ds, dx */
m->ip=0x26fc;
skip_first_enter=0;
m->ds=(uint16_t)(m->dx);
goto L_126FE;
L_126FE: ; /* IDA 0x126FE; retn */
m->ip=0x26fe;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_126FF: ; /* IDA 0x126FF; or      byte ptr [si], 1 */
m->ip=0x26ff;
/* Original entry: sub_126FF */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x26ff,h->user);
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(1))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
goto L_12702;
L_12702: ; /* IDA 0x12702; mov     al, [si+0Ah] */
m->ip=0x2702;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+10))));
goto L_12705;
L_12705: ; /* IDA 0x12705; cbw */
m->ip=0x2705;
skip_first_enter=0;
m->ax=(uint16_t)(int16_t)(int8_t)m->ax;
goto L_12706;
L_12706: ; /* IDA 0x12706; add     bx, ax */
m->ip=0x2706;
skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,m->ax,16,0,0));
goto L_12708;
L_12708: ; /* IDA 0x12708; mov     di, es:[bx] */
m->ip=0x2708;
/* Original entry: sub_12708 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x2708,h->user);
m->di=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx)));
goto L_1270B;
L_1270B: ; /* IDA 0x1270B; mov     ax, es:[bx+2] */
m->ip=0x270b;
skip_first_enter=0;
m->ax=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+2)));
goto L_1270F;
L_1270F: ; /* IDA 0x1270F; mov     cl, al */
m->ip=0x270f;
skip_first_enter=0;
vg_cl(m,(uint8_t)((uint8_t)m->ax));
goto L_12711;
L_12711: ; /* IDA 0x12711; xor     ch, ch */
m->ip=0x2711;
skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->cx>>8))^((uint8_t)(m->cx>>8)))&255;ec_logic(m,result,8);vg_ch(m,(uint8_t)(result));}
goto L_12713;
L_12713: ; /* IDA 0x12713; shl     cx, 1 */
m->ip=0x2713;
skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
goto L_12715;
L_12715: ; /* IDA 0x12715; shl     cx, 1 */
m->ip=0x2715;
skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
goto L_12717;
L_12717: ; /* IDA 0x12717; shl     cx, 1 */
m->ip=0x2717;
skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
goto L_12719;
L_12719: ; /* IDA 0x12719; mov     dx, cx */
m->ip=0x2719;
skip_first_enter=0;
m->dx=(uint16_t)(m->cx);
goto L_1271B;
L_1271B: ; /* IDA 0x1271B; shl     cx, 1 */
m->ip=0x271b;
skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
goto L_1271D;
L_1271D: ; /* IDA 0x1271D; add     cx, dx */
m->ip=0x271d;
skip_first_enter=0;
m->cx=(uint16_t)(ec_math(m,m->cx,m->dx,16,0,0));
goto L_1271F;
L_1271F: ; /* IDA 0x1271F; mov     dx, cx */
m->ip=0x271f;
skip_first_enter=0;
m->dx=(uint16_t)(m->cx);
goto L_12721;
L_12721: ; /* IDA 0x12721; add     dx, cs:word_19872 */
m->ip=0x2721;
skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->cs,39026),16,0,0));
goto L_12726;
L_12726: ; /* IDA 0x12726; mov     ds, dx */
m->ip=0x2726;
skip_first_enter=0;
m->ds=(uint16_t)(m->dx);
goto L_12728;
L_12728: ; /* IDA 0x12728; cmp     byte ptr [di], 0 */
m->ip=0x2728;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di)),0,8);
goto L_1272B;
L_1272B: ; /* IDA 0x1272B; jz      short loc_1273C */
m->ip=0x272b;
skip_first_enter=0;
if(m->flags&0x40) {goto L_1273C;}
goto L_1272D;
L_1272D: ; /* IDA 0x1272D; mov     ds, cs:word_10D52 */
m->ip=0x272d;
skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_12732;
L_12732: ; /* IDA 0x12732; mov     dx, di */
m->ip=0x2732;
skip_first_enter=0;
m->dx=(uint16_t)(m->di);
goto L_12734;
L_12734: ; /* IDA 0x12734; xor     ah, ah */
m->ip=0x2734;
skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_12736;
L_12736: ; /* IDA 0x12736; call    sub_12831 */
m->ip=0x2736;
skip_first_enter=0;
{uint16_t target=10289;ec_push(m,0x2739);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12739: ; /* IDA 0x12739; jb      short loc_1273C */
m->ip=0x2739;
skip_first_enter=0;
if(m->flags&1) {goto L_1273C;}
goto L_1273B;
L_1273B: ; /* IDA 0x1273B; retn */
m->ip=0x273b;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_1273C: ; /* IDA 0x1273C; mov     ax, cx */
m->ip=0x273c;
skip_first_enter=0;
m->ax=(uint16_t)(m->cx);
goto L_1273E;
L_1273E: ; /* IDA 0x1273E; add     ax, cs:word_10D44 */
m->ip=0x273e;
skip_first_enter=0;
m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->cs,3396),16,0,0));
goto L_12743;
L_12743: ; /* IDA 0x12743; mov     ds, ax */
m->ip=0x2743;
skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_12745;
L_12745: ; /* IDA 0x12745; mov     al, [di] */
m->ip=0x2745;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di))));
goto L_12747;
L_12747: ; /* IDA 0x12747; mov     ds, cs:word_10D52 */
m->ip=0x2747;
skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_1274C;
L_1274C: ; /* IDA 0x1274C; test    byte ptr [si], 1 */
m->ip=0x274c;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))&(1))&255;ec_logic(m,result,8);}
goto L_1274F;
L_1274F: ; /* IDA 0x1274F; jz      short loc_1275F */
m->ip=0x274f;
skip_first_enter=0;
if(m->flags&0x40) {goto L_1275F;}
goto L_12751;
L_12751: ; /* IDA 0x12751; cmp     al, 0CEh */
m->ip=0x2751;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,206,8);
goto L_12753;
L_12753: ; /* IDA 0x12753; jb      short loc_1275F */
m->ip=0x2753;
skip_first_enter=0;
if(m->flags&1) {goto L_1275F;}
goto L_12755;
L_12755: ; /* IDA 0x12755; cmp     al, 0DDh */
m->ip=0x2755;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,221,8);
goto L_12757;
L_12757: ; /* IDA 0x12757; ja      short loc_1275F */
m->ip=0x2757;
skip_first_enter=0;
if(!(m->flags&(1|0x40))) {goto L_1275F;}
goto L_12759;
L_12759: ; /* IDA 0x12759; call    sub_12880 */
m->ip=0x2759;
skip_first_enter=0;
{uint16_t target=10368;ec_push(m,0x275c);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_1275C: ; /* IDA 0x1275C; jb      short loc_1275F */
m->ip=0x275c;
skip_first_enter=0;
if(m->flags&1) {goto L_1275F;}
goto L_1275E;
L_1275E: ; /* IDA 0x1275E; retn */
m->ip=0x275e;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_1275F: ; /* IDA 0x1275F; add     cx, cs:word_19872 */
m->ip=0x275f;
skip_first_enter=0;
m->cx=(uint16_t)(ec_math(m,m->cx,dl_read16(m,m->cs,39026),16,0,0));
goto L_12764;
L_12764: ; /* IDA 0x12764; mov     [si+1Ch], cx */
m->ip=0x2764;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+28),m->cx);
goto L_12767;
L_12767: ; /* IDA 0x12767; mov     [si+1Ah], di */
m->ip=0x2767;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+26),m->di);
goto L_1276A;
L_1276A: ; /* IDA 0x1276A; mov     [si+0Ch], bx */
m->ip=0x276a;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+12),m->bx);
goto L_1276D;
L_1276D: ; /* IDA 0x1276D; mov     ax, es:[bx+2] */
m->ip=0x276d;
skip_first_enter=0;
m->ax=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+2)));
goto L_12771;
L_12771: ; /* IDA 0x12771; mov     [si+10h], di */
m->ip=0x2771;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+16),m->di);
goto L_12774;
L_12774: ; /* IDA 0x12774; mov     [si+12h], al */
m->ip=0x2774;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+18),(uint8_t)m->ax);
goto L_12777;
L_12777: ; /* IDA 0x12777; mov     al, ah */
m->ip=0x2777;
skip_first_enter=0;
vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
goto L_12779;
L_12779: ; /* IDA 0x12779; and     al, 7 */
m->ip=0x2779;
skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)&(7))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_1277B;
L_1277B: ; /* IDA 0x1277B; cmp     al, 2 */
m->ip=0x277b;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,2,8);
goto L_1277D;
L_1277D: ; /* IDA 0x1277D; jnb     short loc_12783 */
m->ip=0x277d;
skip_first_enter=0;
if(!(m->flags&1)) {goto L_12783;}
goto L_1277F;
L_1277F: ; /* IDA 0x1277F; call    sub_12804 */
m->ip=0x277f;
skip_first_enter=0;
{uint16_t target=10244;ec_push(m,0x2782);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12782: ; /* IDA 0x12782; retn */
m->ip=0x2782;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12783: ; /* IDA 0x12783; test    ah, 40h */
m->ip=0x2783;
skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))&(64))&255;ec_logic(m,result,8);}
goto L_12786;
L_12786: ; /* IDA 0x12786; jz      short loc_12795 */
m->ip=0x2786;
skip_first_enter=0;
if(m->flags&0x40) {goto L_12795;}
goto L_12788;
L_12788: ; /* IDA 0x12788; cmp     byte ptr [si+0Ah], 0 */
m->ip=0x2788;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+10)),0,8);
goto L_1278C;
L_1278C: ; /* IDA 0x1278C; jg      short loc_1277F */
m->ip=0x278c;
skip_first_enter=0;
if(ec_greater(m)) {goto L_1277F;}
goto L_1278E;
L_1278E: ; /* IDA 0x1278E; call    sub_127F6 */
m->ip=0x278e;
skip_first_enter=0;
{uint16_t target=10230;ec_push(m,0x2791);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12791: ; /* IDA 0x12791; call    sub_127A2 */
m->ip=0x2791;
skip_first_enter=0;
{uint16_t target=10146;ec_push(m,0x2794);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12794: ; /* IDA 0x12794; retn */
m->ip=0x2794;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12795: ; /* IDA 0x12795; cmp     byte ptr [si+0Ah], 0 */
m->ip=0x2795;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+10)),0,8);
goto L_12799;
L_12799: ; /* IDA 0x12799; jl      short loc_1277F */
m->ip=0x2799;
skip_first_enter=0;
if(ec_less(m)) {goto L_1277F;}
goto L_1279B;
L_1279B: ; /* IDA 0x1279B; call    sub_127F6 */
m->ip=0x279b;
skip_first_enter=0;
{uint16_t target=10230;ec_push(m,0x279e);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_1279E: ; /* IDA 0x1279E; call    sub_127A2 */
m->ip=0x279e;
skip_first_enter=0;
{uint16_t target=10146;ec_push(m,0x27a1);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_127A1: ; /* IDA 0x127A1; retn */
m->ip=0x27a1;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_127A2: ; /* IDA 0x127A2; mov     bx, [si+0Eh] */
m->ip=0x27a2;
/* Original entry: sub_127A2 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x27a2,h->user);
m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
goto L_127A5;
L_127A5: ; /* IDA 0x127A5; cmp     byte ptr [si+0Ah], 4 */
m->ip=0x27a5;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+10)),4,8);
goto L_127A9;
L_127A9: ; /* IDA 0x127A9; jnz     short loc_127B1 */
m->ip=0x27a9;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_127B1;}
goto L_127AB;
L_127AB: ; /* IDA 0x127AB; mov     bx, es:[bx+8] */
m->ip=0x27ab;
skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+8)));
goto L_127AF;
L_127AF: ; /* IDA 0x127AF; jmp     short loc_127B5 */
m->ip=0x27af;
skip_first_enter=0;
goto L_127B5;
L_127B1: ; /* IDA 0x127B1; mov     bx, es:[bx+6] */
m->ip=0x27b1;
skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+6)));
goto L_127B5;
L_127B5: ; /* IDA 0x127B5; mov     [si+0Eh], bx */
m->ip=0x27b5;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+14),m->bx);
goto L_127B8;
L_127B8: ; /* IDA 0x127B8; cmp     bx, 600h */
m->ip=0x27b8;
skip_first_enter=0;
ec_cmp(m,m->bx,1536,16);
goto L_127BC;
L_127BC: ; /* IDA 0x127BC; jnb     short loc_127E9 */
m->ip=0x27bc;
skip_first_enter=0;
if(!(m->flags&1)) {goto L_127E9;}
goto L_127BE;
L_127BE: ; /* IDA 0x127BE; shl     bx, 1 */
m->ip=0x27be;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_127C0;
L_127C0: ; /* IDA 0x127C0; shl     bx, 1 */
m->ip=0x27c0;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_127C2;
L_127C2: ; /* IDA 0x127C2; mov     ax, [bx+848h] */
m->ip=0x27c2;
skip_first_enter=0;
m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+2120)));
goto L_127C6;
L_127C6: ; /* IDA 0x127C6; mov     bx, [bx+84Ah] */
m->ip=0x27c6;
skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+2122)));
goto L_127CA;
L_127CA: ; /* IDA 0x127CA; mov     [si+10h], ax */
m->ip=0x27ca;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+16),m->ax);
goto L_127CD;
L_127CD: ; /* IDA 0x127CD; mov     [si+12h], bx */
m->ip=0x27cd;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+18),m->bx);
goto L_127D0;
L_127D0: ; /* IDA 0x127D0; xor     bh, bh */
m->ip=0x27d0;
skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
goto L_127D2;
L_127D2: ; /* IDA 0x127D2; shl     bx, 1 */
m->ip=0x27d2;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_127D4;
L_127D4: ; /* IDA 0x127D4; shl     bx, 1 */
m->ip=0x27d4;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_127D6;
L_127D6: ; /* IDA 0x127D6; shl     bx, 1 */
m->ip=0x27d6;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_127D8;
L_127D8: ; /* IDA 0x127D8; mov     dx, bx */
m->ip=0x27d8;
skip_first_enter=0;
m->dx=(uint16_t)(m->bx);
goto L_127DA;
L_127DA: ; /* IDA 0x127DA; shl     bx, 1 */
m->ip=0x27da;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_127DC;
L_127DC: ; /* IDA 0x127DC; add     bx, dx */
m->ip=0x27dc;
skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
goto L_127DE;
L_127DE: ; /* IDA 0x127DE; add     bx, cs:word_19872 */
m->ip=0x27de;
skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,dl_read16(m,m->cs,39026),16,0,0));
goto L_127E3;
L_127E3: ; /* IDA 0x127E3; mov     [si+1Ah], ax */
m->ip=0x27e3;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+26),m->ax);
goto L_127E6;
L_127E6: ; /* IDA 0x127E6; mov     [si+1Ch], bx */
m->ip=0x27e6;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+28),m->bx);
goto L_127E9;
L_127E9: ; /* IDA 0x127E9; mov     ax, [si+0Eh] */
m->ip=0x27e9;
skip_first_enter=0;
m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
goto L_127EC;
L_127EC: ; /* IDA 0x127EC; cmp     ax, [si+14h] */
m->ip=0x27ec;
skip_first_enter=0;
ec_cmp(m,m->ax,dl_read16(m,m->ds,(uint16_t)(m->si+20)),16);
goto L_127EF;
L_127EF: ; /* IDA 0x127EF; jnz     short locret_127F5 */
m->ip=0x27ef;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_127F5;}
goto L_127F1;
L_127F1: ; /* IDA 0x127F1; mov     byte ptr [si+8], 4 */
m->ip=0x27f1;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+8),4);
goto L_127F5;
L_127F5: ; /* IDA 0x127F5; retn */
m->ip=0x27f5;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_127F6: ; /* IDA 0x127F6; mov     al, [si+0Ah] */
m->ip=0x27f6;
/* Original entry: sub_127F6 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x27f6,h->user);
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+10))));
goto L_127F9;
L_127F9: ; /* IDA 0x127F9; neg     al */
m->ip=0x27f9;
skip_first_enter=0;
vg_al(m,(uint8_t)(ec_math(m,0,(uint8_t)m->ax,8,0,1)));
goto L_127FB;
L_127FB: ; /* IDA 0x127FB; cbw */
m->ip=0x27fb;
skip_first_enter=0;
m->ax=(uint16_t)(int16_t)(int8_t)m->ax;
goto L_127FC;
L_127FC: ; /* IDA 0x127FC; call    sub_12808 */
m->ip=0x27fc;
skip_first_enter=0;
{uint16_t target=10248;ec_push(m,0x27ff);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_127FF: ; /* IDA 0x127FF; xor     byte ptr [si+8], 1 */
m->ip=0x27ff;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si+8)))^(1))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si+8),result);}
goto L_12803;
L_12803: ; /* IDA 0x12803; retn */
m->ip=0x2803;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12804: ; /* IDA 0x12804; mov     al, [si+0Ah] */
m->ip=0x2804;
/* Original entry: sub_12804 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x2804,h->user);
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+10))));
goto L_12807;
L_12807: ; /* IDA 0x12807; cbw */
m->ip=0x2807;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 6
m->ax=(uint8_t)m->ax;
#else
m->ax=(uint16_t)(int16_t)(int8_t)m->ax;
#endif
goto L_12808;
L_12808: ; /* IDA 0x12808; add     ax, [si+0Ch] */
m->ip=0x2808;
/* Original entry: sub_12808 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x2808,h->user);
m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->ds,(uint16_t)(m->si+12)),16,0,0));
goto L_1280B;
L_1280B: ; /* IDA 0x1280B; mov     bx, ax */
m->ip=0x280b;
skip_first_enter=0;
m->bx=(uint16_t)(m->ax);
goto L_1280D;
L_1280D: ; /* IDA 0x1280D; mov     ax, [si+10h] */
m->ip=0x280d;
skip_first_enter=0;
m->ax=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+16)));
goto L_12810;
L_12810: ; /* IDA 0x12810; sub     ax, es:[bx] */
m->ip=0x2810;
skip_first_enter=0;
m->ax=(uint16_t)(ec_math(m,m->ax,dl_read16(m,m->es,(uint16_t)(m->bx)),16,0,1));
goto L_12813;
L_12813: ; /* IDA 0x12813; jz      short loc_1281E */
m->ip=0x2813;
skip_first_enter=0;
if(m->flags&0x40) {goto L_1281E;}
goto L_12815;
L_12815: ; /* IDA 0x12815; and     ax, 8000h */
m->ip=0x2815;
skip_first_enter=0;
{uint16_t result=((m->ax)&(32768))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
goto L_12818;
L_12818: ; /* IDA 0x12818; rol     ax, 1 */
m->ip=0x2818;
skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,0));
goto L_1281A;
L_1281A: ; /* IDA 0x1281A; mov     [si+8], al */
m->ip=0x281a;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+8),(uint8_t)m->ax);
goto L_1281D;
L_1281D: ; /* IDA 0x1281D; retn */
m->ip=0x281d;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_1281E: ; /* IDA 0x1281E; mov     al, [si+12h] */
m->ip=0x281e;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+18))));
goto L_12821;
L_12821: ; /* IDA 0x12821; sub     al, es:[bx+2] */
m->ip=0x2821;
skip_first_enter=0;
vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,dl_read8(m,m->es,(uint16_t)(m->bx+2)),8,0,1)));
goto L_12825;
L_12825: ; /* IDA 0x12825; jz      short locret_12830 */
m->ip=0x2825;
skip_first_enter=0;
if(m->flags&0x40) {goto L_12830;}
goto L_12827;
L_12827: ; /* IDA 0x12827; and     al, 80h */
m->ip=0x2827;
skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)&(128))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_12829;
L_12829: ; /* IDA 0x12829; rol     al, 1 */
m->ip=0x2829;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_shift(m,(uint8_t)m->ax,1,8,0)));
goto L_1282B;
L_1282B: ; /* IDA 0x1282B; add     al, 2 */
m->ip=0x282b;
skip_first_enter=0;
vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,2,8,0,0)));
goto L_1282D;
L_1282D: ; /* IDA 0x1282D; mov     [si+8], al */
m->ip=0x282d;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+8),(uint8_t)m->ax);
goto L_12830;
L_12830: ; /* IDA 0x12830; retn */
m->ip=0x2830;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12831: ; /* IDA 0x12831; push    bx */
m->ip=0x2831;
/* Original entry: sub_12831 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x2831,h->user);
ec_push(m,m->bx);
goto L_12832;
L_12832: ; /* IDA 0x12832; push    cx */
m->ip=0x2832;
skip_first_enter=0;
ec_push(m,m->cx);
goto L_12833;
L_12833: ; /* IDA 0x12833; push    di */
m->ip=0x2833;
skip_first_enter=0;
ec_push(m,m->di);
goto L_12834;
L_12834: ; /* IDA 0x12834; mov     bx, [si+0Eh] */
m->ip=0x2834;
skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
goto L_12837;
L_12837: ; /* IDA 0x12837; mov     di, 2240h */
m->ip=0x2837;
skip_first_enter=0;
m->di=(uint16_t)(8768);
goto L_1283A;
L_1283A: ; /* IDA 0x1283A; mov     cx, 7Fh */
m->ip=0x283a;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 8
m->cx=126;
#else
m->cx=(uint16_t)(127);
#endif
goto L_1283D;
L_1283D: ; /* IDA 0x1283D; cmp     byte ptr [di], 80h */
m->ip=0x283d;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di)),128,8);
goto L_12840;
L_12840: ; /* IDA 0x12840; jb      short loc_1284C */
m->ip=0x2840;
skip_first_enter=0;
if(m->flags&1) {goto L_1284C;}
goto L_12842;
L_12842: ; /* IDA 0x12842; cmp     ax, [di+12h] */
m->ip=0x2842;
skip_first_enter=0;
ec_cmp(m,m->ax,dl_read16(m,m->ds,(uint16_t)(m->di+18)),16);
goto L_12845;
L_12845: ; /* IDA 0x12845; jnz     short loc_1284C */
m->ip=0x2845;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_1284C;}
goto L_12847;
L_12847: ; /* IDA 0x12847; cmp     dx, [di+10h] */
m->ip=0x2847;
skip_first_enter=0;
ec_cmp(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->di+16)),16);
goto L_1284A;
L_1284A: ; /* IDA 0x1284A; jz      short loc_12853 */
m->ip=0x284a;
skip_first_enter=0;
if(m->flags&0x40) {goto L_12853;}
goto L_1284C;
L_1284C: ; /* IDA 0x1284C; add     di, 40h ; '@' */
m->ip=0x284c;
skip_first_enter=0;
m->di=(uint16_t)(ec_math(m,m->di,64,16,0,0));
goto L_1284F;
L_1284F: ; /* IDA 0x1284F; loop    loc_1283D */
m->ip=0x284f;
skip_first_enter=0;
if(--m->cx) {goto L_1283D;}
goto L_12851;
L_12851: ; /* IDA 0x12851; jmp     short loc_1287B */
m->ip=0x2851;
skip_first_enter=0;
goto L_1287B;
L_12853: ; /* IDA 0x12853; mov     cl, [si+1] */
m->ip=0x2853;
skip_first_enter=0;
vg_cl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_12856;
L_12856: ; /* IDA 0x12856; cmp     cl, [di+1] */
m->ip=0x2856;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->cx,dl_read8(m,m->ds,(uint16_t)(m->di+1)),8);
goto L_12859;
L_12859: ; /* IDA 0x12859; jz      short loc_1287B */
m->ip=0x2859;
skip_first_enter=0;
if(m->flags&0x40) {goto L_1287B;}
goto L_1285B;
L_1285B: ; /* IDA 0x1285B; or      byte ptr [si], 20h */
m->ip=0x285b;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(32))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
goto L_1285E;
L_1285E: ; /* IDA 0x1285E; cmp     byte ptr [si+3], 1 */
m->ip=0x285e;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+3)),1,8);
goto L_12862;
L_12862: ; /* IDA 0x12862; ja      short loc_1286C */
m->ip=0x2862;
skip_first_enter=0;
if(!(m->flags&(1|0x40))) {goto L_1286C;}
goto L_12864;
L_12864: ; /* IDA 0x12864; jz      short loc_12873 */
m->ip=0x2864;
skip_first_enter=0;
if(m->flags&0x40) {goto L_12873;}
goto L_12866;
L_12866: ; /* IDA 0x12866; mov     byte ptr [si+3], 0Ch */
m->ip=0x2866;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 9
dl_write8(m,m->ds,(uint16_t)(m->si+3),11);
#else
dl_write8(m,m->ds,(uint16_t)(m->si+3),12);
#endif
goto L_1286A;
L_1286A: ; /* IDA 0x1286A; jmp     short loc_12876 */
m->ip=0x286a;
skip_first_enter=0;
goto L_12876;
L_1286C: ; /* IDA 0x1286C; mov     al, 3 */
m->ip=0x286c;
skip_first_enter=0;
vg_al(m,(uint8_t)(3));
goto L_1286E;
L_1286E: ; /* IDA 0x1286E; call    sub_102F5 */
m->ip=0x286e;
skip_first_enter=0;
{uint16_t target=757;ec_push(m,0x2871);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12871: ; /* IDA 0x12871; jmp     short loc_12876 */
m->ip=0x2871;
skip_first_enter=0;
goto L_12876;
L_12873: ; /* IDA 0x12873; call    sub_14A7B */
m->ip=0x2873;
skip_first_enter=0;
{uint16_t target=19067;ec_push(m,0x2876);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12876: ; /* IDA 0x12876; clc */
m->ip=0x2876;
skip_first_enter=0;
m->flags&=(uint16_t)~1u;
goto L_12877;
L_12877: ; /* IDA 0x12877; pop     di */
m->ip=0x2877;
skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_12878;
L_12878: ; /* IDA 0x12878; pop     cx */
m->ip=0x2878;
skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_12879;
L_12879: ; /* IDA 0x12879; pop     bx */
m->ip=0x2879;
skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_1287A;
L_1287A: ; /* IDA 0x1287A; retn */
m->ip=0x287a;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_1287B: ; /* IDA 0x1287B; stc */
m->ip=0x287b;
skip_first_enter=0;
m->flags|=1;
goto L_1287C;
L_1287C: ; /* IDA 0x1287C; pop     di */
m->ip=0x287c;
skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_1287D;
L_1287D: ; /* IDA 0x1287D; pop     cx */
m->ip=0x287d;
skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_1287E;
L_1287E: ; /* IDA 0x1287E; pop     bx */
m->ip=0x287e;
skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_1287F;
L_1287F: ; /* IDA 0x1287F; retn */
m->ip=0x287f;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12880: ; /* IDA 0x12880; push    es */
m->ip=0x2880;
/* Original entry: sub_12880 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x2880,h->user);
ec_push(m,m->es);
goto L_12881;
L_12881: ; /* IDA 0x12881; push    ax */
m->ip=0x2881;
skip_first_enter=0;
ec_push(m,m->ax);
goto L_12882;
L_12882: ; /* IDA 0x12882; push    di */
m->ip=0x2882;
skip_first_enter=0;
ec_push(m,m->di);
goto L_12883;
L_12883: ; /* IDA 0x12883; mov     di, [si+0Eh] */
m->ip=0x2883;
skip_first_enter=0;
m->di=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->si+14)));
goto L_12886;
L_12886: ; /* IDA 0x12886; cmp     byte ptr [si+0Ah], 4 */
m->ip=0x2886;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+10)),4,8);
goto L_1288A;
L_1288A: ; /* IDA 0x1288A; jnz     short loc_12892 */
m->ip=0x288a;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_12892;}
goto L_1288C;
L_1288C: ; /* IDA 0x1288C; mov     di, es:[di+8] */
m->ip=0x288c;
skip_first_enter=0;
m->di=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->di+8)));
goto L_12890;
L_12890: ; /* IDA 0x12890; jmp     short loc_12896 */
m->ip=0x2890;
skip_first_enter=0;
goto L_12896;
L_12892: ; /* IDA 0x12892; mov     di, es:[di+6] */
m->ip=0x2892;
skip_first_enter=0;
m->di=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->di+6)));
goto L_12896;
L_12896: ; /* IDA 0x12896; shl     di, 1 */
m->ip=0x2896;
skip_first_enter=0;
m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
goto L_12898;
L_12898: ; /* IDA 0x12898; shl     di, 1 */
m->ip=0x2898;
skip_first_enter=0;
m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
goto L_1289A;
L_1289A: ; /* IDA 0x1289A; mov     al, [si+1] */
m->ip=0x289a;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_1289D;
L_1289D: ; /* IDA 0x1289D; cmp     [di+841h], al */
m->ip=0x289d;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+2113)),(uint8_t)m->ax,8);
goto L_128A1;
L_128A1: ; /* IDA 0x128A1; jz      short loc_128C7 */
m->ip=0x28a1;
skip_first_enter=0;
if(m->flags&0x40) {goto L_128C7;}
goto L_128A3;
L_128A3: ; /* IDA 0x128A3; or      byte ptr [si], 20h */
m->ip=0x28a3;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(32))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
goto L_128A6;
L_128A6: ; /* IDA 0x128A6; cmp     byte ptr [si+3], 1 */
m->ip=0x28a6;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+3)),1,8);
goto L_128AA;
L_128AA: ; /* IDA 0x128AA; ja      short loc_128B4 */
m->ip=0x28aa;
skip_first_enter=0;
if(!(m->flags&(1|0x40))) {goto L_128B4;}
goto L_128AC;
L_128AC: ; /* IDA 0x128AC; jz      short loc_128BB */
m->ip=0x28ac;
skip_first_enter=0;
if(m->flags&0x40) {goto L_128BB;}
goto L_128AE;
L_128AE: ; /* IDA 0x128AE; mov     byte ptr [si+3], 0Ch */
m->ip=0x28ae;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+3),12);
goto L_128B2;
L_128B2: ; /* IDA 0x128B2; jmp     short loc_128C2 */
m->ip=0x28b2;
skip_first_enter=0;
goto L_128C2;
L_128B4: ; /* IDA 0x128B4; mov     al, 3 */
m->ip=0x28b4;
skip_first_enter=0;
vg_al(m,(uint8_t)(3));
goto L_128B6;
L_128B6: ; /* IDA 0x128B6; call    sub_102F5 */
m->ip=0x28b6;
skip_first_enter=0;
{uint16_t target=757;ec_push(m,0x28b9);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_128B9: ; /* IDA 0x128B9; jmp     short loc_128C2 */
m->ip=0x28b9;
skip_first_enter=0;
goto L_128C2;
L_128BB: ; /* IDA 0x128BB; add     di, 840h */
m->ip=0x28bb;
skip_first_enter=0;
m->di=(uint16_t)(ec_math(m,m->di,2112,16,0,0));
goto L_128BF;
L_128BF: ; /* IDA 0x128BF; call    sub_14ADE */
m->ip=0x28bf;
skip_first_enter=0;
{uint16_t target=19166;ec_push(m,0x28c2);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_128C2: ; /* IDA 0x128C2; clc */
m->ip=0x28c2;
skip_first_enter=0;
m->flags&=(uint16_t)~1u;
goto L_128C3;
L_128C3: ; /* IDA 0x128C3; pop     di */
m->ip=0x28c3;
skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_128C4;
L_128C4: ; /* IDA 0x128C4; pop     ax */
m->ip=0x28c4;
skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_128C5;
L_128C5: ; /* IDA 0x128C5; pop     es */
m->ip=0x28c5;
skip_first_enter=0;
m->es=(uint16_t)(ec_pop(m));
goto L_128C6;
L_128C6: ; /* IDA 0x128C6; retn */
m->ip=0x28c6;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_128C7: ; /* IDA 0x128C7; stc */
m->ip=0x28c7;
skip_first_enter=0;
m->flags|=1;
goto L_128C8;
L_128C8: ; /* IDA 0x128C8; pop     di */
m->ip=0x28c8;
skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_128C9;
L_128C9: ; /* IDA 0x128C9; pop     ax */
m->ip=0x28c9;
skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_128CA;
L_128CA: ; /* IDA 0x128CA; pop     es */
m->ip=0x28ca;
skip_first_enter=0;
m->es=(uint16_t)(ec_pop(m));
goto L_128CB;
L_128CB: ; /* IDA 0x128CB; retn */
m->ip=0x28cb;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_128F4: ; /* IDA 0x128F4; mov     di, bx */
m->ip=0x28f4;
/* Original entry: sub_128F4 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x28f4,h->user);
m->di=(uint16_t)(m->bx);
goto L_128F6;
L_128F6: ; /* IDA 0x128F6; shl     di, 1 */
m->ip=0x28f6;
skip_first_enter=0;
m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
goto L_128F8;
L_128F8: ; /* IDA 0x128F8; shl     di, 1 */
m->ip=0x28f8;
skip_first_enter=0;
m->di=(uint16_t)(dl_shift(m,m->di,1,16,4));
goto L_128FA;
L_128FA: ; /* IDA 0x128FA; mov     al, [si+1] */
m->ip=0x28fa;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_128FD;
L_128FD: ; /* IDA 0x128FD; cmp     al, [di+841h] */
m->ip=0x28fd;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->di+2113)),8);
goto L_12901;
L_12901: ; /* IDA 0x12901; jz      short loc_12917 */
m->ip=0x2901;
skip_first_enter=0;
if(m->flags&0x40) {goto L_12917;}
goto L_12903;
L_12903: ; /* IDA 0x12903; mov     al, [si+20h] */
m->ip=0x2903;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+32))));
goto L_12906;
L_12906: ; /* IDA 0x12906; xor     ah, ah */
m->ip=0x2906;
skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_12908;
L_12908: ; /* IDA 0x12908; shl     ax, 1 */
m->ip=0x2908;
skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
goto L_1290A;
L_1290A: ; /* IDA 0x1290A; shl     ax, 1 */
m->ip=0x290a;
skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
goto L_1290C;
L_1290C: ; /* IDA 0x1290C; shl     ax, 1 */
m->ip=0x290c;
skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
goto L_1290E;
L_1290E: ; /* IDA 0x1290E; cmp     ax, bx */
m->ip=0x290e;
skip_first_enter=0;
ec_cmp(m,m->ax,m->bx,16);
goto L_12910;
L_12910: ; /* IDA 0x12910; jnz     short loc_12917 */
m->ip=0x2910;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_12917;}
goto L_12912;
L_12912: ; /* IDA 0x12912; call    sub_1291A */
m->ip=0x2912;
skip_first_enter=0;
{uint16_t target=10522;ec_push(m,0x2915);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12915: ; /* IDA 0x12915; stc */
m->ip=0x2915;
skip_first_enter=0;
m->flags|=1;
goto L_12916;
L_12916: ; /* IDA 0x12916; retn */
m->ip=0x2916;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12917: ; /* IDA 0x12917; clc */
m->ip=0x2917;
skip_first_enter=0;
m->flags&=(uint16_t)~1u;
goto L_12918;
L_12918: ; /* IDA 0x12918; retn */
m->ip=0x2918;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12A7E: ; /* IDA 0x12A7E; dec     byte ptr [si+2243h] */
m->ip=0x2a7e;
/* Original entry: sub_12A7E */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x2a7e,h->user);
dl_write8(m,m->ds,(uint16_t)(m->si+8771),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->si+8771)),8));
goto L_12A82;
L_12A82: ; /* IDA 0x12A82; jz      short loc_12A85 */
m->ip=0x2a82;
skip_first_enter=0;
if(m->flags&0x40) {goto L_12A85;}
goto L_12A84;
L_12A84: ; /* IDA 0x12A84; retn */
m->ip=0x2a84;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12A85: ; /* IDA 0x12A85; mov     byte ptr [si+2240h], 0 */
m->ip=0x2a85;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+8768),0);
goto L_12A8A;
L_12A8A: ; /* IDA 0x12A8A; mov     di, si */
m->ip=0x2a8a;
skip_first_enter=0;
m->di=(uint16_t)(m->si);
goto L_12A8C;
L_12A8C: ; /* IDA 0x12A8C; shr     di, 1 */
m->ip=0x2a8c;
skip_first_enter=0;
m->di=(uint16_t)(dl_shift(m,m->di,1,16,5));
goto L_12A8E;
L_12A8E: ; /* IDA 0x12A8E; mov     byte ptr [di+4257h], 0 */
m->ip=0x2a8e;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 11
; /* omitted original general duty clear */
#else
dl_write8(m,m->ds,(uint16_t)(m->di+16983),0);
#endif
goto L_12A93;
L_12A93: ; /* IDA 0x12A93; mov     bh, [si+2241h] */
m->ip=0x2a93;
skip_first_enter=0;
gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+8769))));
goto L_12A97;
L_12A97: ; /* IDA 0x12A97; xor     bl, bl */
m->ip=0x2a97;
skip_first_enter=0;
{uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
goto L_12A99;
L_12A99: ; /* IDA 0x12A99; shr     bx, 1 */
m->ip=0x2a99;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_12A9B;
L_12A9B: ; /* IDA 0x12A9B; shr     bx, 1 */
m->ip=0x2a9b;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_12A9D;
L_12A9D: ; /* IDA 0x12A9D; cmp     byte ptr [bx], 80h */
m->ip=0x2a9d;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx)),128,8);
goto L_12AA0;
L_12AA0: ; /* IDA 0x12AA0; jnb     short loc_12AA8 */
m->ip=0x2aa0;
skip_first_enter=0;
if(!(m->flags&1)) {goto L_12AA8;}
goto L_12AA2;
L_12AA2: ; /* IDA 0x12AA2; mov     byte ptr [di+425Ch], 0FFh */
m->ip=0x2aa2;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->di+16988),255);
goto L_12AA7;
L_12AA7: ; /* IDA 0x12AA7; retn */
m->ip=0x2aa7;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_12AA8: ; /* IDA 0x12AA8; mov     al, cs:byte_10CFF */
m->ip=0x2aa8;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->cs,3327)));
goto L_12AAC;
L_12AAC: ; /* IDA 0x12AAC; cmp     al, [si+2241h] */
m->ip=0x2aac;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+8769)),8);
goto L_12AB0;
L_12AB0: ; /* IDA 0x12AB0; jnz     short locret_12AD1 */
m->ip=0x2ab0;
skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_12AD1;}
goto L_12AB2;
L_12AB2: ; /* IDA 0x12AB2; add     di, 4240h */
m->ip=0x2ab2;
skip_first_enter=0;
m->di=(uint16_t)(ec_math(m,m->di,16960,16,0,0));
goto L_12AB6;
L_12AB6: ; /* IDA 0x12AB6; push    di */
m->ip=0x2ab6;
skip_first_enter=0;
ec_push(m,m->di);
goto L_12AB7;
L_12AB7: ; /* IDA 0x12AB7; mov     di, sp */
m->ip=0x2ab7;
skip_first_enter=0;
m->di=(uint16_t)(m->sp);
goto L_12AB9;
L_12AB9: ; /* IDA 0x12AB9; call    sub_10CDE */
m->ip=0x2ab9;
skip_first_enter=0;
{uint16_t target=3294;ec_push(m,0x2abc);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12ABC: ; /* IDA 0x12ABC; mov     cx, 23h ; '#' */
m->ip=0x2abc;
skip_first_enter=0;
m->cx=(uint16_t)(35);
goto L_12ABF;
L_12ABF: ; /* IDA 0x12ABF; mov     al, 93h */
m->ip=0x2abf;
skip_first_enter=0;
vg_al(m,(uint8_t)(147));
goto L_12AC1;
L_12AC1: ; /* IDA 0x12AC1; call    sub_18810 */
m->ip=0x2ac1;
skip_first_enter=0;
{uint16_t target=34832;ec_push(m,0x2ac4);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12AC4: ; /* IDA 0x12AC4; pop     bx */
m->ip=0x2ac4;
skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_12AC5;
L_12AC5: ; /* IDA 0x12AC5; mov     cx, 198h */
m->ip=0x2ac5;
skip_first_enter=0;
m->cx=(uint16_t)(408);
goto L_12AC8;
L_12AC8: ; /* IDA 0x12AC8; mov     ah, [bx+1Eh] */
m->ip=0x2ac8;
skip_first_enter=0;
vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+30))));
goto L_12ACB;
L_12ACB: ; /* IDA 0x12ACB; mov     al, [bx+1] */
m->ip=0x2acb;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+1))));
goto L_12ACE;
L_12ACE: ; /* IDA 0x12ACE; call    sub_18810 */
m->ip=0x2ace;
skip_first_enter=0;
{uint16_t target=34832;ec_push(m,0x2ad1);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_12AD1: ; /* IDA 0x12AD1; retn */
m->ip=0x2ad1;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_142AB: ; /* IDA 0x142AB; push    di */
m->ip=0x42ab;
/* Original entry: sub_142AB */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x42ab,h->user);
ec_push(m,m->di);
goto L_142AC;
L_142AC: ; /* IDA 0x142AC; push    bp */
m->ip=0x42ac;
skip_first_enter=0;
ec_push(m,m->bp);
goto L_142AD;
L_142AD: ; /* IDA 0x142AD; mov     bp, bx */
m->ip=0x42ad;
skip_first_enter=0;
m->bp=(uint16_t)(m->bx);
goto L_142AF;
L_142AF: ; /* IDA 0x142AF; xor     di, di */
m->ip=0x42af;
skip_first_enter=0;
{uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
goto L_142B1;
L_142B1: ; /* IDA 0x142B1; cmp     byte ptr [si+0Ah], 0FCh */
m->ip=0x42b1;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+10)),252,8);
goto L_142B5;
L_142B5: ; /* IDA 0x142B5; jz      short loc_142B9 */
m->ip=0x42b5;
skip_first_enter=0;
if(m->flags&0x40) {goto L_142B9;}
goto L_142B7;
L_142B7: ; /* IDA 0x142B7; inc     di */
m->ip=0x42b7;
skip_first_enter=0;
m->di=(uint16_t)(ec_inc(m,m->di,16));
goto L_142B8;
L_142B8: ; /* IDA 0x142B8; inc     di */
m->ip=0x42b8;
skip_first_enter=0;
m->di=(uint16_t)(ec_inc(m,m->di,16));
goto L_142B9;
L_142B9: ; /* IDA 0x142B9; mov     bx, es:[bx+di+6] */
m->ip=0x42b9;
skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bx+m->di+6)));
goto L_142BD;
L_142BD: ; /* IDA 0x142BD; shl     bx, 1 */
m->ip=0x42bd;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_142BF;
L_142BF: ; /* IDA 0x142BF; shl     bx, 1 */
m->ip=0x42bf;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_142C1;
L_142C1: ; /* IDA 0x142C1; mov     dl, [si+1] */
m->ip=0x42c1;
skip_first_enter=0;
vg_dl(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_142C4;
L_142C4: ; /* IDA 0x142C4; mov     al, [bx+841h] */
m->ip=0x42c4;
skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+2113))));
goto L_142C8;
L_142C8: ; /* IDA 0x142C8; cmp     al, dl */
m->ip=0x42c8;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,(uint8_t)m->dx,8);
goto L_142CA;
L_142CA: ; /* IDA 0x142CA; jz      short loc_142FD */
m->ip=0x42ca;
skip_first_enter=0;
if(m->flags&0x40) {goto L_142FD;}
goto L_142CC;
L_142CC: ; /* IDA 0x142CC; cmp     al, 18h */
m->ip=0x42cc;
skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,24,8);
goto L_142CE;
L_142CE: ; /* IDA 0x142CE; jz      short loc_142FD */
m->ip=0x42ce;
skip_first_enter=0;
if(m->flags&0x40) {goto L_142FD;}
goto L_142D0;
L_142D0: ; /* IDA 0x142D0; xor     dh, dh */
m->ip=0x42d0;
skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->dx>>8))^((uint8_t)(m->dx>>8)))&255;ec_logic(m,result,8);vg_dh(m,(uint8_t)(result));}
goto L_142D2;
L_142D2: ; /* IDA 0x142D2; shl     dx, 1 */
m->ip=0x42d2;
skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
goto L_142D4;
L_142D4: ; /* IDA 0x142D4; shl     dx, 1 */
m->ip=0x42d4;
skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
goto L_142D6;
L_142D6: ; /* IDA 0x142D6; shl     dx, 1 */
m->ip=0x42d6;
skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
goto L_142D8;
L_142D8: ; /* IDA 0x142D8; mov     bx, dx */
m->ip=0x42d8;
skip_first_enter=0;
m->bx=(uint16_t)(m->dx);
goto L_142DA;
L_142DA: ; /* IDA 0x142DA; shl     dx, 1 */
m->ip=0x42da;
skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,4));
goto L_142DC;
L_142DC: ; /* IDA 0x142DC; add     bx, dx */
m->ip=0x42dc;
skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,m->dx,16,0,0));
goto L_142DE;
L_142DE: ; /* IDA 0x142DE; add     bx, 600h */
m->ip=0x42de;
skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,1536,16,0,0));
goto L_142E2;
L_142E2: ; /* IDA 0x142E2; xlat */
m->ip=0x42e2;
skip_first_enter=0;
vg_al(m,dl_read8(m,m->ds,(uint16_t)(m->bx+(uint8_t)m->ax)));
goto L_142E3;
L_142E3: ; /* IDA 0x142E3; cmp     al, 80h */
m->ip=0x42e3;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 10
ec_cmp(m,(uint8_t)m->ax,0x81,8);
#else
ec_cmp(m,(uint8_t)m->ax,128,8);
#endif
goto L_142E5;
L_142E5: ; /* IDA 0x142E5; jb      short loc_142FD */
m->ip=0x42e5;
skip_first_enter=0;
if(m->flags&1) {goto L_142FD;}
goto L_142E7;
L_142E7: ; /* IDA 0x142E7; xor     di, 2 */
m->ip=0x42e7;
skip_first_enter=0;
{uint16_t result=((m->di)^(2))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
goto L_142EA;
L_142EA: ; /* IDA 0x142EA; mov     bx, es:[bp+di+6] */
m->ip=0x42ea;
skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->bp+m->di+6)));
goto L_142EE;
L_142EE: ; /* IDA 0x142EE; mov     [si+14h], bx */
m->ip=0x42ee;
skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+20),m->bx);
goto L_142F1;
L_142F1: ; /* IDA 0x142F1; shr     bx, 1 */
m->ip=0x42f1;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_142F3;
L_142F3: ; /* IDA 0x142F3; shr     bx, 1 */
m->ip=0x42f3;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_142F5;
L_142F5: ; /* IDA 0x142F5; shr     bx, 1 */
m->ip=0x42f5;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_142F7;
L_142F7: ; /* IDA 0x142F7; mov     [si+20h], bl */
m->ip=0x42f7;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)m->bx);
goto L_142FA;
L_142FA: ; /* IDA 0x142FA; or      byte ptr [si], 2 */
m->ip=0x42fa;
skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,(uint16_t)(m->si)))|(2))&255;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si),result);}
goto L_142FD;
L_142FD: ; /* IDA 0x142FD; pop     bp */
m->ip=0x42fd;
skip_first_enter=0;
m->bp=(uint16_t)(ec_pop(m));
goto L_142FE;
L_142FE: ; /* IDA 0x142FE; pop     di */
m->ip=0x42fe;
skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_142FF;
L_142FF: ; /* IDA 0x142FF; retn */
m->ip=0x42ff;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_14300: ; /* IDA 0x14300; shl     bx, 1 */
m->ip=0x4300;
/* Original entry: sub_14300 */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x4300,h->user);
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_14302;
L_14302: ; /* IDA 0x14302; shl     bx, 1 */
m->ip=0x4302;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_14304;
L_14304: ; /* IDA 0x14304; cmp     byte ptr [bx+858h], 1 */
m->ip=0x4304;
skip_first_enter=0;
#if KI_ARMY_MUTATION == 12
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+0x858)),0,8);
#else
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+2136)),1,8);
#endif
goto L_14309;
L_14309: ; /* IDA 0x14309; ja      short loc_14323 */
m->ip=0x4309;
skip_first_enter=0;
if(!(m->flags&(1|0x40))) {goto L_14323;}
goto L_1430B;
L_1430B: ; /* IDA 0x1430B; cmp     byte ptr [si+8], 4 */
m->ip=0x430b;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+8)),4,8);
goto L_1430F;
L_1430F: ; /* IDA 0x1430F; jz      short loc_14323 */
m->ip=0x430f;
skip_first_enter=0;
if(m->flags&0x40) {goto L_14323;}
goto L_14311;
L_14311: ; /* IDA 0x14311; cmp     byte ptr [bx+840h], 80h */
m->ip=0x4311;
skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+2112)),128,8);
goto L_14316;
L_14316: ; /* IDA 0x14316; jb      short loc_14323 */
m->ip=0x4316;
skip_first_enter=0;
if(m->flags&1) {goto L_14323;}
goto L_14318;
L_14318: ; /* IDA 0x14318; shl     bx, 1 */
m->ip=0x4318;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_1431A;
L_1431A: ; /* IDA 0x1431A; shl     bx, 1 */
m->ip=0x431a;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_1431C;
L_1431C: ; /* IDA 0x1431C; shl     bx, 1 */
m->ip=0x431c;
skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_1431E;
L_1431E: ; /* IDA 0x1431E; mov     [si+20h], bh */
m->ip=0x431e;
skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+32),(uint8_t)(m->bx>>8));
goto L_14321;
L_14321: ; /* IDA 0x14321; stc */
m->ip=0x4321;
skip_first_enter=0;
m->flags|=1;
goto L_14322;
L_14322: ; /* IDA 0x14322; retn */
m->ip=0x4322;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_14323: ; /* IDA 0x14323; clc */
m->ip=0x4323;
skip_first_enter=0;
m->flags&=(uint16_t)~1u;
goto L_14324;
L_14324: ; /* IDA 0x14324; retn */
m->ip=0x4324;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
L_1562B: ; /* IDA 0x1562B; push    si */
m->ip=0x562b;
/* Original entry: sub_1562B */
if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x562b,h->user);
ec_push(m,m->si);
goto L_1562C;
L_1562C: ; /* IDA 0x1562C; mov     si, dx */
m->ip=0x562c;
skip_first_enter=0;
m->si=(uint16_t)(m->dx);
goto L_1562E;
L_1562E: ; /* IDA 0x1562E; and     si, 0FF00h */
m->ip=0x562e;
skip_first_enter=0;
{uint16_t result=((m->si)&(65280))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_15632;
L_15632: ; /* IDA 0x15632; shr     si, 1 */
m->ip=0x5632;
skip_first_enter=0;
m->si=(uint16_t)(dl_shift(m,m->si,1,16,5));
goto L_15634;
L_15634: ; /* IDA 0x15634; shr     si, 1 */
m->ip=0x5634;
skip_first_enter=0;
m->si=(uint16_t)(dl_shift(m,m->si,1,16,5));
goto L_15636;
L_15636: ; /* IDA 0x15636; call    sub_1563B */
m->ip=0x5636;
skip_first_enter=0;
{uint16_t target=22075;ec_push(m,0x5639);m->ip=target;invoke_external=1;goto AM_DISPATCH;}
L_15639: ; /* IDA 0x15639; pop     si */
m->ip=0x5639;
skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_1563A;
L_1563A: ; /* IDA 0x1563A; retn */
m->ip=0x563a;
skip_first_enter=0;
m->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;
}
int ki_army_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
if(!am_known(t))return 0;
uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,ss,sp);
am_graph(m,t,h,1); /* External resolver already emitted the entry hook. */
assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
ec_push(m,ret);return 1;}
void ki_army_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {if(am_known(t))am_graph(m,t,h,0);else ki_engagement_invoke(m,t,h);}
void sub_125A3(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x25a3,h,0);}
void sub_12600(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x2600,h,0);}
void sub_1264A(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x264a,h,0);}
void sub_12662(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x2662,h,0);}
void sub_126FF(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x26ff,h,0);}
void sub_12708(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x2708,h,0);}
void sub_127A2(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x27a2,h,0);}
void sub_127F6(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x27f6,h,0);}
void sub_12804(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x2804,h,0);}
void sub_12808(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x2808,h,0);}
void sub_12831(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x2831,h,0);}
void sub_12880(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x2880,h,0);}
void sub_128F4(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x28f4,h,0);}
void sub_12A7E(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x2a7e,h,0);}
void sub_142AB(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x42ab,h,0);}
void sub_14300(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x4300,h,0);}
void sub_1562B(KiMachine16 *m,const KiEconomyHooks *h) {am_graph(m,0x562b,h,0);}
