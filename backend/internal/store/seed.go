package store

import (
	"database/sql"
	"fmt"
)

// SeedCurriculum inserts sample lessons and activities for Math & Coding for Builders age group (L1-L3).
func SeedCurriculum(db *sql.DB) error {
	var mathCatID, codingCatID int64
	err := db.QueryRow(`SELECT id FROM categories WHERE slug = 'math'`).Scan(&mathCatID)
	if err != nil {
		return fmt.Errorf("get math category id: %w", err)
	}

	err = db.QueryRow(`SELECT id FROM categories WHERE slug = 'coding'`).Scan(&codingCatID)
	if err != nil {
		return fmt.Errorf("get coding category id: %w", err)
	}

	// ==========================================
	// MATH LESSONS (BUILDERS AGE GROUP: L1 - L3)
	// ==========================================

	// 1. Math Level 1 Lesson 1: Penjumlahan 1–100
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (1, ?, 'builders', 1, 1, 'Penjumlahan 1–100', 'Belajar menjumlahkan dua angka hingga 100 dengan mudah', 
		 '{"intro_text":"Halo! Hari ini kita akan belajar menjumlahkan angka. Ingat, tambahkan angka satuan terlebih dahulu, lalu puluhan!","icon":"🧮"}', 10, 20)`,
		mathCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1, 1, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Berapakah hasil dari 15 + 12?",
			"options": ["25", "27", "30", "22"],
			"answer": "27",
			"hint": "Coba tambahkan 5 + 2 = 7, lalu 10 + 10 = 20. Hasilnya 20 + 7!",
			"explanation": "15 + 12 = (10 + 5) + (10 + 2) = 20 + 7 = 27."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2, 1, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Urutkan angka-angka berikut dari yang TERKECIL ke yang TERBESAR!",
			"items": ["45", "12", "89", "23"],
			"expected_order": ["12", "23", "45", "89"],
			"hint": "Cari angka puluhan terkecil dulu, yaitu 12!",
			"explanation": "Urutan dari kecil ke besar: 12, 23, 45, 89."
		}`,
	)

	// 2. Math Level 1 Lesson 2: Pengurangan 1–100
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (5, ?, 'builders', 1, 2, 'Pengurangan 1–100', 'Belajar pengurangan dua angka hingga 100', 
		 '{"intro_text":"Pengurangan adalah proses mengambil sejumlah nilai dari angka awal. Yuk latihan!","icon":"➖"}', 10, 20)`,
		mathCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (8, 5, 'fill_blank', 1, ?, 10)`,
		`{
			"prompt": "Berapakah hasil dari 50 - 15? Ketik jawabanmu di bawah!",
			"answer": "35",
			"hint": "Hitung 50 - 10 = 40, lalu kurangi 5 lagi!",
			"explanation": "50 - 15 = 35."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (9, 5, 'multiple_choice', 2, ?, 10)`,
		`{
			"prompt": "Ibu punya 25 apel dan diberikan ke tetangga 8 apel. Sisa apel ibu adalah?",
			"options": ["17", "15", "18", "20"],
			"answer": "17",
			"hint": "Kurangi 25 dengan 8!",
			"explanation": "25 - 8 = 17 sisa apel ibu."
		}`,
	)

	// 3. Math Level 2 Lesson 1: Perkalian Dasar (Tabel 1–5)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (2, ?, 'builders', 2, 1, 'Perkalian Dasar (Tabel 1–5)', 'Belajar perkalian sebagai penjumlahan berulang', 
		 '{"intro_text":"Perkalian adalah penjumlahan yang diulang-ulang. Contohnya: 3 x 4 artinya 4 + 4 + 4!","icon":"✖️"}', 12, 25)`,
		mathCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (3, 2, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Berapakah hasil dari 4 x 5?",
			"options": ["15", "20", "24", "18"],
			"answer": "20",
			"hint": "Hitung 5 + 5 + 5 + 5!",
			"explanation": "4 x 5 = 5 + 5 + 5 + 5 = 20."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (4, 2, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Susunlah urutan hasil kelipatan 3 dari yang terkecil!",
			"items": ["9", "3", "12", "6"],
			"expected_order": ["3", "6", "9", "12"],
			"hint": "Mulailah dari 3 x 1 = 3, lalu 3 x 2 = 6...",
			"explanation": "Kelipatan 3: 3, 6, 9, 12."
		}`,
	)

	// 4. Math Level 2 Lesson 2: Pembagian Dasar
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (6, ?, 'builders', 2, 2, 'Pembagian Dasar 1–20', 'Membagi benda menjadi beberapa bagian sama banyak', 
		 '{"intro_text":"Pembagian adalah kebalikan dari perkalian. Membagi artinya membagikan sama rata!","icon":"➗"}', 12, 25)`,
		mathCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (10, 6, 'fill_blank', 1, ?, 10)`,
		`{
			"prompt": "Budi mempunyai 12 permen dan dibagikan sama rata kepada 3 temannya. Setiap teman mendapat berapa permen?",
			"answer": "4",
			"hint": "Berapa dikali 3 hasilnya 12?",
			"explanation": "12 : 3 = 4 permen per anak."
		}`,
	)

	// 5. Math Level 3 Lesson 1: Soal Cerita Matematika Sehari-hari
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (7, ?, 'builders', 3, 1, 'Soal Cerita Matematika', 'Memecahkan masalah matematika dalam kehidupan sehari-hari', 
		 '{"intro_text":"Matematika ada di mana-mana! Mari memecahkan soal cerita seru di toko mainan dan sekolah.","icon":"📖"}', 15, 30)`,
		mathCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (11, 7, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Siti membeli 3 pensil seharga Rp 2.000 per pensil. Berapa total uang yang harus dibayar Siti?",
			"options": ["Rp 5.000", "Rp 6.000", "Rp 7.000", "Rp 8.000"],
			"answer": "Rp 6.000",
			"hint": "Hitung 3 x 2.000!",
			"explanation": "3 x Rp 2.000 = Rp 6.000."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (12, 7, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Sebuah bus membawa 30 penumpang. Di selter pertama turun 8 orang, dan naik 5 orang. Berapa jumlah penumpang sekarang?",
			"answer": "27",
			"hint": "30 - 8 + 5 = ?",
			"explanation": "30 - 8 = 22, kemudian 22 + 5 = 27 penumpang."
		}`,
	)

	// 6. Math Level 3 Lesson 2: Mengenal Bangun Datar & Simetri
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (8, ?, 'builders', 3, 2, 'Bangun Datar & Simetri', 'Mengenal persegi, segitiga, lingkaran, dan garis simetri', 
		 '{"intro_text":"Bangun datar memiliki sisi dan sudut. Yuk kenali sifat-sifatnya!","icon":"📐"}', 15, 30)`,
		mathCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (13, 8, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Bangun datar manakah yang memiliki 3 buah sisi dan 3 buah sudut?",
			"options": ["Segitiga", "Persegi", "Lingkaran", "Trapesium"],
			"answer": "Segitiga",
			"hint": "Sesuai namanya, 'segi-tiga'!",
			"explanation": "Segitiga memiliki 3 sisi dan 3 sudut."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (14, 8, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Urutkan bangun datar berikut berdasarkan JUMLAH SISINYA dari yang TERSEDIKIT!",
			"items": ["Persegi", "Lingkaran", "Segitiga", "Segilima"],
			"expected_order": ["Lingkaran", "Segitiga", "Persegi", "Segilima"],
			"hint": "Lingkaran memiliki 0/1 sisi lengkung, Segitiga 3, Persegi 4, Segilima 5!",
			"explanation": "Jumlah sisi: Lingkaran (0/1), Segitiga (3), Persegi (4), Segilima (5)."
		}`,
	)

	// ==========================================
	// CODING LESSONS (BUILDERS AGE GROUP: L1 - L2)
	// ==========================================

	// 7. Coding Lesson 1 (Level 1)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (3, ?, 'builders', 1, 1, 'Algoritma & urutan Langkah', 'Belajar menyusun instruksi komputer dengan urutan yang tepat', 
		 '{"intro_text":"Komputer itu sangat patuh tapi butuh petunjuk yang jelas! Urutan langkah petunjuk ini disebut Algoritma.","icon":"🤖"}', 10, 20)`,
		codingCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (5, 3, 'drag_drop', 1, ?, 10)`,
		`{
			"prompt": "Urutkan langkah membuang sampah dengan benar!",
			"items": ["Buka tempat sampah", "Masukkan sampah", "Ambil sampah", "Tutup tempat sampah"],
			"expected_order": ["Ambil sampah", "Buka tempat sampah", "Masukkan sampah", "Tutup tempat sampah"],
			"hint": "Langkah pertama adalah mengambil sampahnya dulu!",
			"explanation": "Algoritma yang benar: Ambil sampah -> Buka tempat sampah -> Masukkan sampah -> Tutup tempat sampah."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (6, 3, 'multiple_choice', 2, ?, 10)`,
		`{
			"prompt": "Apa langkah pertama yang benar sebelum menyeberang jalan?",
			"options": ["Tengok kanan dan kiri", "Langsung lari cepat", "Tutup mata", "Bermain ponsel"],
			"answer": "Tengok kanan dan kiri",
			"hint": "Kita harus memastikan jalan aman dari kendaraan!",
			"explanation": "Selalu tengok kanan dan kiri untuk memastikan jalan aman sebelum menyeberang."
		}`,
	)

	// 8. Coding Lesson 2 (Level 2)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (4, ?, 'builders', 2, 1, 'Pengulangan (Looping)', 'Belajar menggunakan perintah perulangan agar kode lebih hemat', 
		 '{"intro_text":"Daripada menulis instruksi yang sama berulang kali, kita bisa menggunakan Loop!","icon":"🔄"}', 12, 25)`,
		codingCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (7, 4, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Robot ingin berjalan 4 kali ke depan. Instruksi mana yang paling ringkas?",
			"options": ["Ulangi 4x: Maju 1 langkah", "Maju 10 langkah", "Mundur 4 langkah", "Diam saja"],
			"answer": "Ulangi 4x: Maju 1 langkah",
			"hint": "Gunakan instruksi perulangan (Loop)!",
			"explanation": "'Ulangi 4x: Maju 1 langkah' melakukan hal yang sama seperti menulis 'Maju 1 langkah' sebanyak 4 kali."
		}`,
	)

	fmt.Println("Seeded Math (L1-L3) & Coding curriculum content for Builders group")
	return nil
}
