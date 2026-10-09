#ifndef WOLONG_C_VERDICT_FIXTURE_H
#define WOLONG_C_VERDICT_FIXTURE_H
#include "resume_fixture.h"
typedef KiResourceFixture KiVerdictFixture;
static void verdict_external(KiMachine16 *m,uint16_t t,void *user) {
 KiEconomyHooks h={rs_snapshot,verdict_external,user};
 if(!ki_verdict_body(m,t,&h)) resume_external(m,t,user);
}
static void verdict_run(KiMachine16 *m,uint16_t t,KiVerdictFixture *f) {
 KiEconomyHooks h={rs_snapshot,verdict_external,f};ki_verdict_invoke(m,t,&h);
}
#endif
