// IDA Pro 9.4 IDC evidence export for DOS/V idle-clock / event-10 hypothesis.
// Non-destructive: preserves original names, addresses, and operands.
#include <idc.idc>
#define printf Message

static PrintCallers(name)
{
  auto target, xref, fstart;
  target = LocByName(name);
  if (target == BADADDR)
  {
    printf("TARGET %s EA UNKNOWN\n", name);
    return;
  }
  printf("TARGET %s EA %08X\n", name, target);
  xref = RfirstB(target);
  while (xref != BADADDR)
  {
    fstart = GetFunctionAttr(xref, FUNCATTR_START);
    printf("  CALLER %08X %s  %s\n", xref,
           fstart == BADADDR ? "<no-function>" : GetFunctionName(fstart),
           GetDisasm(xref));
    xref = RnextB(target, xref);
  }
}

static PrintRefs(name)
{
  auto target, xref, dref, fstart;
  target = LocByName(name);
  if (target == BADADDR)
  {
    printf("SYMBOL %s EA UNKNOWN\n", name);
    return;
  }
  printf("SYMBOL %s EA %08X\n", name, target);
  xref = RfirstB(target);
  while (xref != BADADDR)
  {
    fstart = GetFunctionAttr(xref, FUNCATTR_START);
    printf("  XREF %08X %s  %s\n", xref,
           fstart == BADADDR ? "<no-function>" : GetFunctionName(fstart),
           GetDisasm(xref));
    xref = RnextB(target, xref);
  }
  dref = DfirstB(target);
  while (dref != BADADDR)
  {
    fstart = GetFunctionAttr(dref, FUNCATTR_START);
    printf("  DREF %08X %s  %s\n", dref,
           fstart == BADADDR ? "<no-function>" : GetFunctionName(fstart),
           GetDisasm(dref));
    dref = DnextB(target, dref);
  }
}

static PrintFunction(name)
{
  auto target, fstart, fend, ea;
  target = LocByName(name);
  if (target == BADADDR)
  {
    printf("FUNCTION %s EA UNKNOWN\n", name);
    return;
  }
  fstart = GetFunctionAttr(target, FUNCATTR_START);
  fend = GetFunctionAttr(target, FUNCATTR_END);
  printf("FUNCTION %s %08X-%08X\n", name, fstart, fend);
  ea = fstart;
  while (ea != BADADDR && ea < fend)
  {
    printf("  %08X  %s\n", ea, GetDisasm(ea));
    ea = NextHead(ea, fend);
  }
}

static main()
{
  printf("IDA_VERSION 9.4 (container image ida-pro-9.4-ver2:uidfix-v1)\n");
  printf("EVIDENCE_LEVELS raw xrefs/functions/instructions=proven; idle semantics=review required\n");

  PrintCallers("sub_11BE0");
  PrintCallers("sub_11F7F");
  PrintCallers("sub_11CD0");
  PrintCallers("sub_11D8E");
  PrintCallers("sub_13E11");
  PrintCallers("sub_125A3");
  PrintCallers("sub_12459");
  PrintCallers("sub_13EFD");
  PrintCallers("sub_131AE");
  PrintCallers("sub_1304E");

  PrintRefs("byte_198A3");
  PrintRefs("word_19886");
  PrintRefs("word_19888");
  PrintRefs("word_10D18");
  PrintRefs("word_10D1C");
  PrintRefs("word_10D1E");
  PrintRefs("byte_10D2A");
  PrintRefs("byte_10CF2");
  PrintRefs("byte_10CF3");
  PrintRefs("byte_10CF4");
  PrintRefs("word_10CF6");
  PrintRefs("byte_131AD");

  PrintFunction("sub_11BE0");
  PrintFunction("sub_11F7F");
  PrintFunction("sub_11CD0");
  PrintFunction("sub_11D8E");
  PrintFunction("sub_13E11");
  PrintFunction("sub_125A3");
  PrintFunction("sub_12662");
  PrintFunction("sub_127A2");
  PrintFunction("sub_13EFD");
  PrintFunction("sub_131AE");
  PrintFunction("sub_1304E");
  PrintFunction("sub_122DB");
  PrintFunction("sub_12286");
  PrintFunction("sub_15358");
  Exit(0);
}
