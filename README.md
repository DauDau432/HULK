# Hulk DoS tool

Công cụ HULK DoS được chuyển sang ngôn ngữ Go từ Python.
Tiện ích Python gốc của Barry Shteiman http://www.sectorix.com/2012/05/17/hulk-web-server-dos-tool/
Tôi vừa chuyển mã vì nó nhanh và bẩn. Tên các hàm gốc được giữ lại và phần lớn logic gốc cũng được giữ lại.

Sự khác biệt chính so với phiên bản Python nằm trong kiến ​​trúc Golang cho đồng thời: goroutines. Hulk.py chạy
một luồng mới cho mỗi kết nối trong nhóm kết nối để nó sử dụng hàng trăm và hàng nghìn luồng.
hulk.go chỉ sử dụng các goroutines nhẹ chỉ sử dụng hàng chục luồng (thời gian chạy golang thường bắt đầu một luồng cho
Lõi CPU + một số luồng dịch vụ). Kiến trúc này cho phép phiên bản golang tiêu thụ tài nguyên tốt hơn và cao hơn nhiều
nhóm kết nối trên cùng một phần cứng hơn phiên bản Python có thể.

Công cụ này được nhắm mục tiêu để kiểm tra căng thẳng và có thể thực sự làm hỏng máy chủ được định cấu hình kém hoặc ứng dụng được tạo ra kém. Hãy sử dụng nó một cách cẩn thận.

Môi trường hữu ích vars:

* GOMAXPROCS
   Đặt nó thành số CPU của bạn hoặc cao hơn (không thực tế hơn đối với các phiên bản golang mới nhất).
* HULKMAXPROCS
   Giới hạn nhóm kết nối (1024 theo mặc định).

Thêm chi tiết: http://old.siberian.laika.name/node/7 

Cập nhật: tốt, tôi đã tạo tiện ích này cho tác vụ một lần khi tôi chỉ chơi một chút với golang. Thật ngạc nhiên khi tôi thấy rằng
tiện ích này được những người khác sử dụng, được đánh giá cao trên github và thậm chí còn được đưa vào [BlackArch Linux distro] (http://blackarch.org/dos.html). Vì vậy, tôi đã dọn dẹp mã một chút.

### Cài Đặt
- Yêu cầu python
```
apt install python
```
- Hoặc golang
```
apt install golang
```
- Tải xuống kho lưu trữ
```
git clone https://github.com/DauDau432/hulk.git
```
- Di chuyển vào thư mục hulk
```
cd hulk
```
### Cách sử dụng
- Với python
```
python hulk.py <url>
```
- Với golang
```
go run hulk.go -site <url>
```
### Giấy phép
Tôi nghĩ nó có thể là miền công cộng vì nó chỉ là một đoạn mã ngắn và đơn giản nhưng vì lý do gì tôi không nhớ nữa
Tôi đã chọn GPL cho nó. Okey. Vì vậy, hãy sử dụng phiên bản HULK được cấp phép theo GPLv3. Xem LICENSE.

Tôi không liên quan đến tiện ích HULK gốc bằng Python. Tiện ích HULK gốc là quyền của Barry Shteiman (http://sectorix.com). Không có bất kỳ tham chiếu nào đến giấy phép trong nguồn gốc thì nó không thuộc GPL. Hãy hỏi tác giả của tiện ích ban đầu về giấy phép.
