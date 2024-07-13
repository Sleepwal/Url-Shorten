# Url-Shorten

## 项目架构

```mermaid
graph TB
    A[User] --> B[Golang Fiber] -->|存储原始和缩短的url| C[Redis]
```

```mermaid
graph LR
    API --> Database & Helpers & Routes & Dockerfile
```



## 测试结果

![](assets/2024-07-13-23-49-09-image.png)