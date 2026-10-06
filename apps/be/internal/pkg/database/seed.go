package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func intPtr(i int) *int {
	return &i
}

// SeedTickets seeds initial categories and event tickets
func SeedTickets(ctx context.Context, db *gorm.DB) error {
	defaultCategories := []domain.Category{
		{
			ID:          uuid.MustParse("019468e1-0000-7000-8000-000000000001"),
			Name:        "Âm nhạc",
			Slug:        "am-nhac",
			Description: "Các sự kiện âm nhạc, đại nhạc hội và liveshow đỉnh cao",
		},
		{
			ID:          uuid.MustParse("019468e1-0000-7000-8000-000000000002"),
			Name:        "Công nghệ",
			Slug:        "cong-nghe",
			Description: "Hội thảo và sự kiện công nghệ, lập trình và trí tuệ nhân tạo",
		},
		{
			ID:          uuid.MustParse("019468e1-0000-7000-8000-000000000003"),
			Name:        "Workshop",
			Slug:        "workshop",
			Description: "Các buổi đào tạo chuyên sâu và chia sẻ kỹ năng thực chiến",
		},
		{
			ID:          uuid.MustParse("019468e1-0000-7000-8000-000000000004"),
			Name:        "Thể thao",
			Slug:        "the-thao",
			Description: "Giải chạy, thi đấu thể thao và hoạt động vận động ngoài trời",
		},
		{
			ID:          uuid.MustParse("019468e1-0000-7000-8000-000000000005"),
			Name:        "Triển lãm",
			Slug:        "trien-lam",
			Description: "Triển lãm nghệ thuật thị giác và không gian sáng tạo",
		},
		{
			ID:          uuid.MustParse("019468e1-0000-7000-8000-000000000006"),
			Name:        "Hài kịch",
			Slug:        "hai-kich",
			Description: "Các đêm diễn kịch nghệ và hài độc thoại",
		},
		{
			ID:          uuid.MustParse("019468e1-0000-7000-8000-000000000007"),
			Name:        "Điện ảnh",
			Slug:        "dien-anh",
			Description: "Công chiếu phim điện ảnh và sự kiện thảm đỏ đặc biệt",
		},
	}

	categoryMap := make(map[string]uuid.UUID)

	for _, cat := range defaultCategories {
		var existing domain.Category
		err := db.WithContext(ctx).Where("slug = ? OR name = ?", cat.Slug, cat.Name).First(&existing).Error
		if err == nil {
			existing.Name = cat.Name
			existing.Slug = cat.Slug
			existing.Description = cat.Description
			existing.UpdatedAt = time.Now()
			if err := db.WithContext(ctx).Save(&existing).Error; err != nil {
				return fmt.Errorf("failed to update category %s: %w", cat.Slug, err)
			}
			categoryMap[cat.Slug] = existing.ID
			slog.Info("Synced category", "id", existing.ID, "name", existing.Name)
		} else {
			c := cat
			c.CreatedAt = time.Now()
			c.UpdatedAt = time.Now()
			if err := db.WithContext(ctx).Create(&c).Error; err != nil {
				return fmt.Errorf("failed to create category %s: %w", c.Slug, err)
			}
			categoryMap[cat.Slug] = c.ID
			slog.Info("Created category", "id", c.ID, "name", c.Name)
		}
	}

	getCatID := func(slug string) *uuid.UUID {
		if id, ok := categoryMap[slug]; ok {
			return &id
		}
		return nil
	}

	tickets := []domain.Ticket{
		{
			ID:                uuid.MustParse("019468e2-63b7-7bb2-8069-b570e3037f01"),
			Name:              "CyberSound Arena: Electro Symphony 2026",
			Description:       "Đại tiệc âm nhạc điện tử kết hợp dàn nhạc giao hưởng 60 nhạc công với sân khấu visual Hologram 360 độ hoành tráng bậc nhất Đông Nam Á.",
			Price:             850000,
			TotalQuantity:     100,
			AvailableStock:    45,
			MaxBookingPerUser: intPtr(4),
			Status:            domain.TicketStatusActive,
			CategoryID:        getCatID("am-nhac"),
			Venue:             "Sân vận động Quân khu 7, TP. Hồ Chí Minh",
			Date:              "15/11/2026",
			Time:              "19:30",
			ImageURL:          "https://images.unsplash.com/photo-1470225620780-dba8ba36b745?auto=format&fit=crop&w=1200&q=80",
			Tags:              pq.StringArray{"Hot", "Music", "Hologram"},
			Featured:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.MustParse("019468e2-63b7-7bb2-8069-b570e3037f02"),
			Name:              "Vietnam AI & Cloud Summit 2026",
			Description:       "Diễn đàn công nghệ hàng đầu hội tụ hơn 2.500 kỹ sư, nhà sáng lập và chuyên gia AI từ Google, OpenAI, AWS chia sẻ kiến trúc Agentic AI & High-Load Backend.",
			Price:             1200000,
			TotalQuantity:     50,
			AvailableStock:    12,
			MaxBookingPerUser: intPtr(2),
			Status:            domain.TicketStatusActive,
			CategoryID:        getCatID("cong-nghe"),
			Venue:             "Trung tâm Hội nghị Quốc gia, Hà Nội",
			Date:              "28/11/2026",
			Time:              "08:30",
			ImageURL:          "https://images.unsplash.com/photo-1540575467063-178a50c2df87?auto=format&fit=crop&w=1200&q=80",
			Tags:              pq.StringArray{"Tech", "VIP", "Limited"},
			Featured:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.MustParse("019468e2-63b7-7bb2-8069-b570e3037f03"),
			Name:              "Indie Acoustic Night: Thanh Âm Mùa Thu",
			Description:       "Đêm nhạc acoustic ấm cúng với những nghệ sĩ indie được yêu mến. Trải nghiệm không gian mộc mạc, gần gũi với đồ uống thủ công tặng kèm.",
			Price:             350000,
			TotalQuantity:     120,
			AvailableStock:    80,
			MaxBookingPerUser: intPtr(5),
			Status:            domain.TicketStatusActive,
			CategoryID:        getCatID("am-nhac"),
			Venue:             "Soul Live Project Complex, TP. Hồ Chí Minh",
			Date:              "05/12/2026",
			Time:              "20:00",
			ImageURL:          "https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?auto=format&fit=crop&w=1200&q=80",
			Tags:              pq.StringArray{"Acoustic", "Chill", "Indie"},
			Featured:          false,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.MustParse("019468e2-63b7-7bb2-8069-b570e3037f04"),
			Name:              "Masterclass: High-Performance Go Microservices",
			Description:       "Khóa huấn luyện chuyên sâu 1 ngày về tối ưu hóa concurrency, zero-allocation buffers, Redis caching và xử lý hàng triệu giao dịch vé mỗi giây.",
			Price:             2500000,
			TotalQuantity:     20,
			AvailableStock:    8,
			MaxBookingPerUser: intPtr(1),
			Status:            domain.TicketStatusActive,
			CategoryID:        getCatID("workshop"),
			Venue:             "Dreamplex Tech Hub, Quận 1, TP. Hồ Chí Minh",
			Date:              "12/12/2026",
			Time:              "09:00",
			ImageURL:          "https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?auto=format&fit=crop&w=1200&q=80",
			Tags:              pq.StringArray{"Golang", "Masterclass", "Sắp hết"},
			Featured:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.MustParse("019468e2-63b7-7bb2-8069-b570e3037f05"),
			Name:              "Hanoi International Marathon: Đêm Di Sản 2026",
			Description:       "Cung đường chạy đêm đi qua các di tích lịch sử Hồ Gươm, Cầu Long Biên, Hoàng thành Thăng Long với huy chương mạ vàng độc bản.",
			Price:             650000,
			TotalQuantity:     200,
			AvailableStock:    120,
			MaxBookingPerUser: intPtr(4),
			Status:            domain.TicketStatusActive,
			CategoryID:        getCatID("the-thao"),
			Venue:             "Quảng trường Đông Kinh Nghĩa Thục, Hà Nội",
			Date:              "20/12/2026",
			Time:              "23:00",
			ImageURL:          "https://images.unsplash.com/photo-1530549387789-4c1017266635?auto=format&fit=crop&w=1200&q=80",
			Tags:              pq.StringArray{"Sports", "Night Run", "Medal"},
			Featured:          false,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.MustParse("019468e2-63b7-7bb2-8069-b570e3037f06"),
			Name:              "Art & Light Exhibition: Lạc Vào Hư Ảo",
			Description:       "Không gian triển lãm nghệ thuật thị giác kỹ thuật số tương tác với hơn 15 phòng trải nghiệm ánh sáng đa chiều và âm thanh vòm binaural.",
			Price:             290000,
			TotalQuantity:     150,
			AvailableStock:    0,
			MaxBookingPerUser: intPtr(6),
			Status:            domain.TicketStatusActive,
			CategoryID:        getCatID("trien-lam"),
			Venue:             "Trung tâm Nghệ thuật Đương đại VCCA, Hà Nội",
			Date:              "Hằng ngày (Đến 31/12/2026)",
			Time:              "10:00 - 21:00",
			ImageURL:          "https://images.unsplash.com/photo-1508997449629-303059a039c0?auto=format&fit=crop&w=1200&q=80",
			Tags:              pq.StringArray{"Art", "Light", "Sold Out"},
			Featured:          false,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.MustParse("019468e2-63b7-7bb2-8069-b570e3037f07"),
			Name:              "Saigon Comedy Club: Cười Xuyên Đêm",
			Description:       "Đêm hài độc thoại đỉnh cao cùng 6 diễn viên hài tài năng hàng đầu Việt Nam. Tặng kèm 01 phần thức uống và bắp rang bơ cao cấp.",
			Price:             320000,
			TotalQuantity:     60,
			AvailableStock:    28,
			MaxBookingPerUser: intPtr(4),
			Status:            domain.TicketStatusActive,
			CategoryID:        getCatID("hai-kich"),
			Venue:             "Nhà hát Kịch TP.HCM, Quận 1, TP. Hồ Chí Minh",
			Date:              "25/12/2026",
			Time:              "20:00",
			ImageURL:          "https://images.unsplash.com/photo-1585699324551-f6c309eedeca?auto=format&fit=crop&w=1200&q=80",
			Tags:              pq.StringArray{"Comedy", "Weekend", "Fun"},
			Featured:          false,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.MustParse("019468e2-63b7-7bb2-8069-b570e3037f08"),
			Name:              "Cinema Gala Premiere: Vũ Trụ Điện Ảnh 2026",
			Description:       "Suất chiếu đặc biệt thảm đỏ đầu tiên tại Việt Nam trên màn hình IMAX Laser thế hệ mới, giao lưu cùng đạo diễn và dàn diễn viên chính.",
			Price:             450000,
			TotalQuantity:     40,
			AvailableStock:    18,
			MaxBookingPerUser: intPtr(2),
			Status:            domain.TicketStatusActive,
			CategoryID:        getCatID("dien-anh"),
			Venue:             "Cụm rạp IMAX Landmark 81, TP. Hồ Chí Minh",
			Date:              "30/12/2026",
			Time:              "18:45",
			ImageURL:          "https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?auto=format&fit=crop&w=1200&q=80",
			Tags:              pq.StringArray{"Cinema", "IMAX", "Red Carpet"},
			Featured:          true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}

	for _, ticket := range tickets {
		t := ticket
		err := db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"name", "description", "price", "total_quantity", "available_stock",
					"max_booking_per_user", "status", "category_id", "venue", "date", "time",
					"image_url", "tags", "featured", "updated_at",
				}),
			}).
			Create(&t).Error
		if err != nil {
			return fmt.Errorf("failed to upsert ticket %s: %w", t.ID, err)
		}
		slog.Info("Seeded ticket", "id", t.ID, "name", t.Name)
	}

	return nil
}
