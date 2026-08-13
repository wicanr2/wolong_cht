// IDA Pro 9.4 IDC evidence export for DOS/V event 10.
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

static PrintWindow(ea, radius)
{
  auto fstart, fend, p, i, start;
  fstart = GetFunctionAttr(ea, FUNCATTR_START);
  fend = GetFunctionAttr(ea, FUNCATTR_END);
  if (fstart == BADADDR || fend == BADADDR)
    return;
  start = fstart;
  for (i = 0; i < radius && start < ea; i = i + 1)
  {
    p = PrevHead(ea, fstart);
    if (p == BADADDR || p >= ea)
      break;
    ea = p;
  }
  printf("  WINDOW %08X %s\n", fstart, GetFunctionName(fstart));
  p = start;
  while (p != BADADDR && p < fend)
  {
    printf("    %08X  %s\n", p, GetDisasm(p));
    p = NextHead(p, fend);
  }
}

static PrintSymbolRefs(name)
{
  auto target, xref, dref;
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
    printf("  XREF %08X %s  %s\n", xref,
           GetFunctionName(GetFunctionAttr(xref, FUNCATTR_START)),
           GetDisasm(xref));
    xref = RnextB(target, xref);
  }
  dref = DfirstB(target);
  while (dref != BADADDR)
  {
    printf("  DREF %08X %s  %s\n", dref,
           GetFunctionName(GetFunctionAttr(dref, FUNCATTR_START)),
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

static PrintImmediateCandidates()
{
  auto seg, ea, end, mnem, op0, op1;
  printf("IMMEDIATE 0x0A CANDIDATES\n");
  seg = FirstSeg();
  while (seg != BADADDR)
  {
    ea = SegStart(seg);
    end = SegEnd(seg);
    while (ea != BADADDR && ea < end)
    {
      if (isCode(GetFlags(ea)))
      {
        mnem = GetMnem(ea);
        op0 = GetOpnd(ea, 0);
        op1 = GetOpnd(ea, 1);
        if ((mnem == "mov" || mnem == "add" || mnem == "or" ||
             mnem == "xor" || mnem == "and" || mnem == "cmp" ||
             mnem == "sub") &&
            (strstr(op0, "0Ah") != -1 || strstr(op1, "0Ah") != -1 ||
             strstr(op0, "0A00h") != -1 || strstr(op1, "0A00h") != -1))
        {
          printf("  %08X %s  %s\n", ea,
                 GetFunctionName(GetFunctionAttr(ea, FUNCATTR_START)),
                 GetDisasm(ea));
        }
      }
      ea = NextHead(ea, end);
    }
    seg = NextSeg(seg);
  }
}

static main()
{
  printf("IDA_VERSION 9.4 (container image ida-pro-9.4-ver2:uidfix-v1)\n");
  printf("EVIDENCE_LEVELS callers/xrefs/raw instructions=proven; semantic producer=review required\n");

  PrintCallers("sub_12FB1");
  PrintCallers("sub_12FBF");
  PrintCallers("sub_1301C");
  PrintCallers("sub_131AE");
  PrintCallers("sub_13496");

  PrintSymbolRefs("word_10D20");
  PrintSymbolRefs("word_10D52");
  PrintSymbolRefs("word_10D56");
  PrintSymbolRefs("byte_131AD");
  PrintSymbolRefs("funcs_131E8");
  PrintFunction("sub_12FB1");
  PrintFunction("sub_12FBF");
  PrintFunction("sub_1301C");
  PrintFunction("sub_1300E");
  PrintFunction("sub_131AE");
  PrintFunction("sub_13496");
  PrintFunction("sub_15715");
  PrintFunction("sub_1578F");
  PrintFunction("sub_157FE");
  PrintFunction("sub_15940");
  PrintFunction("sub_16623");
  PrintImmediateCandidates();
  Exit(0);
}
