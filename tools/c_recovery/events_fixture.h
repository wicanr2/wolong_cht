#ifndef WOLONG_C_EVENTS_FIXTURE_H
#define WOLONG_C_EVENTS_FIXTURE_H
#include "hourly_fixture.h"
static void ev_external(KiMachine16 *m,uint16_t t,void *user) {
    KiEconomyHooks h={hr_snapshot,ev_external,user};
    if(t==0x3c3d) {gw_al(m,ec_byte(m,m->cs,0x7000));return;}
    if(!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&!ki_politics_body(m,t,&h)&&!ki_world_body(m,t,&h)) (void)ki_settlement_body(m,t,&h);
}
static void events_run(KiMachine16 *m,uint16_t t,KiHourlyFixture *f) {
    KiEconomyHooks h={hr_snapshot,ev_external,f};ki_events_invoke(m,t,&h);
}
#endif
