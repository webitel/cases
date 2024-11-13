  // simple filter query language [SFQL]
  // ( -> start of the filter or filter node
  // ) -> end of the filter or the filter node
  // last ( -> always start of the filter
  // [ -> start of the operator
  // ] -> end of the operator
  // operands - simple strings, first string should always be a valid field name, second should
  // operators:
  // [and] [or] -- only applied to the filter node, after them always should be a ( as filter or another filter node
  // [eq], [gte], [gt], [lte], [lt], [regex], [neq] -- only applied to the filter, value after them and to the ) symbol considered as string
  //...?size=10&page=1&filter=((file[eq]fghj)[and](name[eq]yehor))[or])
//  oneof filters {
//    FilterNode node = 7;
//    Filter filter = 8;
//  }