# 这是供外部服务使用muxi审核的sdk-v2

### Description:

1. 相较于v1,v2属于破坏性的更新,故独立出v1版本单独更新为v2.

2. 增添了创建请求对象的方法，不暴露内部实现；字段封装为引用类型以区分三态。请始终使用内置的创建方法来创建对象，不建议直接构造对象，如果自己构造，请保证字段空值情况处理，避免panic或无效对象.

3. v1已经废弃，请尽快更新到v2.

   ```go
   hookUrl := "http://example.com/webhook"
   id := uint(1)
   contents := internal.NewContents(internal.WithTopicText("11", "11"))
   req, err := NewAuditReq(hookUrl, id, contents)
   if err!=nil {
       // ...
   }
   ```

   

