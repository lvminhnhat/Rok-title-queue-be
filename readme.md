# ROK Helper Backend API

Backend API service cho ứng dụng ROK Helper, hỗ trợ quản lý title trong game Rise of Kingdoms.

## Tổng quan

API này cung cấp các endpoint để:
- Quản lý hàng đợi title (Duke, Architect, Scientist, Justice)
- Quản lý kingdom và người chơi
- Cấu hình thời gian cooldown cho mỗi loại title
- Quản lý bản đồ (Home Kingdom và Lost Kingdom)

## Công nghệ sử dụng

- [Go](https://golang.org/) - Ngôn ngữ lập trình
- [Fiber](https://gofiber.io/) - Web framework
- [MongoDB](https://www.mongodb.com/) - Cơ sở dữ liệu
- [Swagger](https://swagger.io/) - API documentation

## Cài đặt

### Yêu cầu

- Go 1.22 hoặc mới hơn
- MongoDB
- Git

### Các bước cài đặt

1. Clone repository:
```
git clone https://github.com/lvminhnhat/Rok-title-queue-be.git
cd rok-helper-backend
```

2. Cài đặt dependencies:

```
go mod download
```

3. Tạo file .env:

```
cp .env.example .env
```


4. Chỉnh sửa file .env với thông tin cấu hình của bạn

5. Chạy ứng dụng:

```
go run main.go
```


## API Documentation

Swagger UI có sẵn tại: `http://localhost:3000/swagger/`

### Endpoints chính:

#### Title Management
- `GET /api/title/{id}` - Lấy title tiếp theo từ hàng đợi
- `POST /api/title/{id}` - Thêm title mới vào hàng đợi
- `POST /api/title/finish/{id}` - Đánh dấu title đã hoàn thành (với cooldown)
- `POST /api/title/done/{id}` - Đánh dấu title đã hoàn thành ngay lập tức

#### Configuration
- `GET /api/config/{id}` - Lấy cấu hình title
- `PUT /api/config/{id}` - Cập nhật cấu hình title

#### Maps
- `PUT /api/maps/{id}` - Cập nhật bản đồ kingdom

#### Kingdom Management
- `GET /api/kingdom/{id}` - Lấy thông tin về kingdom
- `POST /api/kingdom` - Tạo kingdom mới
- `GET /api/kingdom` - Lấy danh sách tất cả kingdom

## Cơ chế hoạt động

### Hàng đợi Title
1. Title được thêm vào hàng đợi với `AddTitle()`
2. `GetTitleAssignment()` lấy title phù hợp tiếp theo dựa trên:
- Thời gian chờ lâu nhất (TimeAdd nhỏ nhất)
- Thời gian cooldown đã hết (currentTime > timeDone) 
3. Khi title hoàn thành, `Finish()` hoặc `Done()` được gọi
4. `CleanExpiredTitles()` xóa các title đã hết thời hạn

## Contributing

Vui lòng tạo issue hoặc pull request cho các cải tiến hoặc sửa lỗi.

## License

[MIT](LICENSE)

