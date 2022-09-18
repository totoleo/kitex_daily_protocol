
include "base.thrift"

namespace go item
struct GetItemsRequest {
    1: list<i64> ids;
    2: optional bool with_cdn; //默认false，true返回带域名的图片img_v2,cover_v2
    255: base.Base Base,
}
struct Resp {
    
}
service Item {
   Resp GetItems(1:GetItemsRequest req)
}