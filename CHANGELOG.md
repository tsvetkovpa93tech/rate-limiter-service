# Changelog

All notable changes to the Rate Limiter Service project will be documented in this file.

## [1.1.0] - 2024-01-01

### Added - Enterprise Features
- **Multi-tenancy support**: Different rate limits per API key/tenant
- **Dynamic limit updates**: Update rate limits without service restart
- **Webhook notifications**: Get notified when rate limits are exceeded
- **Tenant management**: API for managing tenant configurations

### Added - Performance Optimizations
- **Object pooling**: `sync.Pool` for JSON buffers and HTTP request buffers
- **In-memory caching**: Configurable cache for frequent rate limit checks
- **Redis connection pooling**: Optimized connection management with configurable pool size
- **Reduced allocations**: 22-35% reduction in memory allocations

### Added - Monitoring & Observability
- **Grafana dashboards**: Pre-configured dashboards for monitoring
- **Prometheus alerts**: Alerts for high request rate, error rate, latency, and memory usage
- **Go memory metrics**: Monitoring of Go runtime memory usage
- **Performance benchmarks**: Comprehensive benchmark tests with results

### Added - Documentation
- **BENCHMARKS.md**: Detailed performance benchmarks and optimization analysis
- **Algorithm comparison table**: Side-by-side comparison of Token Bucket vs Sliding Window
- **Enterprise features documentation**: Guide for multi-tenancy and webhooks
- **Performance optimization guide**: Best practices for production deployments

### Improved
- **Redis configuration**: Added pool size and min idle connections configuration
- **Error handling**: Better error messages and logging
- **Code organization**: Better separation of concerns

### Changed
- **Configuration structure**: Added webhook and tenant configuration sections
- **Docker Compose**: Added Grafana dashboard provisioning

## [1.0.0] - 2024-01-01

### Added
- Initial release
- Token Bucket algorithm
- Sliding Window Log algorithm
- In-memory storage (sync.Map)
- Redis storage
- REST API
- Prometheus metrics
- Health checks
- Docker support
- Load testing scripts
- OpenAPI documentation

