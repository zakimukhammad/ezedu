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
	// CODING LESSONS (BUILDERS AGE GROUP: L1 - L3)
	// ==========================================

	// 7. Coding Level 1 Lesson 1: Algoritma & Urutan Langkah
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (3, ?, 'builders', 1, 1, 'Algoritma & Urutan Langkah', 'Belajar menyusun instruksi komputer dengan urutan yang tepat', 
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

	// 8. Coding Level 1 Lesson 2: Navigasi Arah Robot (Block Coding)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (9, ?, 'builders', 1, 2, 'Navigasi Arah Robot', 'Menyusun blok kode arah untuk menggerakkan robot ke bendera', 
		 '{"intro_text":"Robot berada di posisi start. Susun blok kode arah untuk mengantarkan robot ke garis finish!","icon":"🧩"}', 10, 20)`,
		codingCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (15, 9, 'block_code', 1, ?, 10)`,
		`{
			"prompt": "Susun blok kode agar robot berjalan: Maju 1 langkah, Belok Kanan, lalu Maju 1 langkah!",
			"available_blocks": ["Maju ⬆️", "Belok Kiri ⬅️", "Belok Kanan ➡️", "Ulangi 3x 🔄"],
			"expected_order": ["Maju ⬆️", "Belok Kanan ➡️", "Maju ⬆️"],
			"hint": "Robot butuh maju, belok kanan, lalu maju lagi!",
			"explanation": "Urutan blok kode yang benar: Maju ⬆️ -> Belok Kanan ➡️ -> Maju ⬆️."
		}`,
	)

	// 9. Coding Level 2 Lesson 1: Pengulangan (Looping)
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

	// 10. Coding Level 2 Lesson 2: Menghemat Kode dengan Loop (Block Coding)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (10, ?, 'builders', 2, 2, 'Menghemat Kode dengan Loop', 'Gunakan blok Ulangi 3x untuk memutar pola berjalan', 
		 '{"intro_text":"Loop menghemat jumlah baris kode yang kamu tulis!","icon":"⚡"}', 12, 25)`,
		codingCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (16, 10, 'block_code', 1, ?, 10)`,
		`{
			"prompt": "Susun blok kode ringkas menggunakan Ulangi 3x diikuti Maju 1 langkah!",
			"available_blocks": ["Maju ⬆️", "Belok Kiri ⬅️", "Ulangi 3x 🔄", "Tutup Loop 🔚"],
			"expected_order": ["Ulangi 3x 🔄", "Maju ⬆️"],
			"hint": "Gunakan blok Ulangi 3x 🔄 terlebih dahulu!",
			"explanation": "Blok perulangan menginstruksikan robot untuk mengulang perintah Maju sebanyak 3 kali."
		}`,
	)

	// 11. Coding Level 3 Lesson 1: Kondisi & Pengandaian (Jika - Maka)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (11, ?, 'builders', 3, 1, 'Kondisi & Pengandaian (If - Else)', 'Belajar membuat keputusan otomatis berdasarkan kondisi', 
		 '{"intro_text":"Komputer bisa mengambil keputusan! Jika ada rintangan di depan, maka robot harus berbelok.","icon":"🔀"}', 15, 30)`,
		codingCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (17, 11, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Manakah contoh aturan logika 'Jika - Maka' (If - Then) dalam kehidupan sehari-hari?",
			"options": ["Jika hujan, maka pakailah payung", "Lari tanpa melihat jalan", "Tidur siang saat belajar", "Makan es krim saat mandi"],
			"answer": "Jika hujan, maka pakailah payung",
			"hint": "Kondisinya adalah Hujan, keputusannya adalah Pakai Payung!",
			"explanation": "'Jika hujan' adalah kondisi, dan 'pakailah payung' adalah tindakan yang diambil."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (18, 11, 'block_code', 2, ?, 10)`,
		`{
			"prompt": "Susun blok kode kondisi: Jika Ada Rintangan ⚠️ -> Belok Kiri ⬅️ -> Maju ⬆️!",
			"available_blocks": ["Maju ⬆️", "Belok Kiri ⬅️", "Jika Ada Rintangan ⚠️", "Ulangi 3x 🔄"],
			"expected_order": ["Jika Ada Rintangan ⚠️", "Belok Kiri ⬅️", "Maju ⬆️"],
			"hint": "Mulailah dengan blok pengandaian 'Jika Ada Rintangan ⚠️'!",
			"explanation": "Blok kondisi memeriksa adanya rintangan terlebih dahulu sebelum berbelok dan maju."
		}`,
	)

	// 12. Coding Level 3 Lesson 2: Debugging & Mencari Kesalahan Kode
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (12, ?, 'builders', 3, 2, 'Debugging & Perbaikan Kode', 'Menemukan dan memperbaiki error/bug dalam urutan kode', 
		 '{"intro_text":"Programmer hebat adalah yang jago menemukan Bug (kesalahan kode) dan memperbaikinya!","icon":"🐛"}', 15, 30)`,
		codingCatID,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (19, 12, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Istilah untuk kesalahan atau masalah pada kode komputer disebut apa?",
			"options": ["Bug", "Cat", "Fish", "Bird"],
			"answer": "Bug",
			"hint": "Sesuai nama serangga kecil dalam bahasa Inggris!",
			"explanation": "Kesalahan atau cacat pada program komputer secara historis disebut 'Bug'."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (20, 12, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Urutkan proses Debugging yang benar!",
			"items": ["Perbaiki kesalahan", "Temukan posisi bug", "Jalankan ulang kode", "Amati masalah"],
			"expected_order": ["Amati masalah", "Temukan posisi bug", "Perbaiki kesalahan", "Jalankan ulang kode"],
			"hint": "Amati dulu masalahnya sebelum mencari posisi bug!",
			"explanation": "Urutan debugging: Amati masalah -> Temukan posisi bug -> Perbaiki kesalahan -> Jalankan ulang kode."
		}`,
	)

	// ==========================================
	// TODDLER LESSONS ("MENGENAL DUNIA": L1)
	// ==========================================
	var toddlerCatID int64
	err = db.QueryRow(`SELECT id FROM categories WHERE slug = 'toddlers'`).Scan(&toddlerCatID)
	if err == nil && toddlerCatID > 0 {
		// 1. Toddler Lesson 1: Bentuk & Warna Dasar
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (21, ?, 'toddlers', 1, 1, 'Bentuk & Warna Dasar', 'Mengenal lingkaran, persegi, dan warna-warni', 
			 '{"intro_text":"Lihat bentuk dan warna yang indah ini! Sentuh gambar untuk mendengar namanya!","icon":"🟡"}', 5, 10)`,
			toddlerCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (21, 21, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Mana gambar Lingkaran Kuning 🟡?",
				"options": ["Lingkaran Kuning 🟡", "Persegi Biru 🟦", "Segitiga Merah 🔺"],
				"answer": "Lingkaran Kuning 🟡",
				"hint": "Cari yang berbentuk bulat dan berwarna kuning cerah!",
				"explanation": "Pintar! Ini adalah Lingkaran Kuning 🟡."
			}`,
		)

		// 2. Toddler Lesson 2: Suara Hewan Ceria
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (22, ?, 'toddlers', 1, 2, 'Suara Hewan Ceria', 'Mengenal suara sapi, kucing, dan ayam', 
			 '{"intro_text":"Dengarkan suara hewan lucu di sekitarmu!","icon":"🐮"}', 5, 10)`,
			toddlerCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (22, 22, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Hewan manakah yang bersuara 'Muuu... Muuu...' 🐮?",
				"options": ["Sapi 🐮", "Kucing 🐱", "Bebek 🦆"],
				"answer": "Sapi 🐮",
				"hint": "Sapi penghasil susu yang bersuara Muuu!",
				"explanation": "Hebat! Sapi 🐮 bersuara Muuu!"
			}`,
		)

		// 3. Toddler Lesson 3: Benda & Kendaraan
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (23, ?, 'toddlers', 1, 3, 'Benda & Kendaraan', 'Mengenal mobil, sepeda, dan bola', 
			 '{"intro_text":"Ayo kenali benda-benda dan kendaraan favoritmu!","icon":"🚗"}', 5, 10)`,
			toddlerCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (23, 23, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Mana kendaraan Mobil Merah 🚗?",
				"options": ["Mobil Merah 🚗", "Sepeda 🚲", "Bola ⚽"],
				"answer": "Mobil Merah 🚗",
				"hint": "Mobil memiliki roda dan bersuara Brumm!",
				"explanation": "Luar biasa! Ini adalah Mobil Merah 🚗."
			}`,
		)

		// 4. Toddler Lesson 4: Buah & Tanaman Segar
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (24, ?, 'toddlers', 1, 4, 'Buah & Tanaman Segar', 'Mengenal buah apel, pisang, dan jeruk', 
			 '{"intro_text":"Buah-buahan sangat sehat dan rasanya lezat!","icon":"🍎"}', 5, 10)`,
			toddlerCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (24, 24, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Mana buah Apel Merah 🍎?",
				"options": ["Apel Merah 🍎", "Pisang Kuning 🍌", "Jeruk 🍊"],
				"answer": "Apel Merah 🍎",
				"hint": "Apel berwarna merah manis!",
				"explanation": "Yum! Apel Merah 🍎 sangat manis dan sehat."
			}`,
		)

		// 5. Toddler Lesson 5: Angka & Huruf Pertama
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (25, ?, 'toddlers', 1, 5, 'Angka & Huruf Pertama', 'Belajar angka 1–5 dan huruf vokal A-I-U-E-O', 
			 '{"intro_text":"Ayo mengenal angka awal dan bunyi huruf pertama!","icon":"🅰️"}', 5, 10)`,
			toddlerCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (25, 25, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Manakah Huruf Vokal A 🅰️?",
				"options": ["Huruf A 🅰️", "Angka 1 1️⃣", "Huruf O ⭕"],
				"answer": "Huruf A 🅰️",
				"hint": "Huruf pertama dalam abjad: A seperti Apel!",
				"explanation": "Hebat sekali! Huruf A 🅰️ adalah awal kata Apel."
			}`,
		)
	}

	// ==========================================
	// SCIENCE LESSONS (BUILDERS: L1 - L3)
	// ==========================================
	var scienceCatID int64
	err = db.QueryRow(`SELECT id FROM categories WHERE slug = 'science'`).Scan(&scienceCatID)
	if err == nil && scienceCatID > 0 {
		// Lesson 26: Hewan & Tempat Tinggalnya
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (26, ?, 'builders', 1, 1, 'Hewan & Tempat Tinggalnya', 'Mengenal habitat hewan darat, air, dan udara', 
			 '{"intro_text":"Setiap hewan memiliki tempat tinggal alamiah yang disebut habitat. Yuk pelajari!","icon":"🌿"}', 10, 20)`,
			scienceCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (26, 26, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Di manakah tempat tinggal (habitat) alami ikan 🐟?",
				"options": ["Air 🌊", "Darat 🏜️", "Udara ☁️"],
				"answer": "Air 🌊",
				"hint": "Ikan bernapas menggunakan insang di dalam air!",
				"explanation": "Ikan hidup di air laut atau sungai dan bernapas dengan insang."
			}`,
		)

		// Lesson 27: Wujud Benda (Padat, Cair, Gas)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (27, ?, 'builders', 1, 2, 'Wujud Benda (Padat, Cair, Gas)', 'Membedakan tiga wujud benda di sekitar kita', 
			 '{"intro_text":"Benda di sekitar kita dibagi menjadi 3 wujud: Padat, Cair, dan Gas!","icon":"🧊"}', 10, 20)`,
			scienceCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (27, 27, 'drag_drop', 1, ?, 10)`,
			`{
				"prompt": "Urutkan contoh benda berikut berdasarkan wujudnya: [Padat, Cair, Gas]!",
				"items": ["Air Minum 💧", "Batu 🪨", "Uap Air 💨"],
				"expected_order": ["Batu 🪨", "Air Minum 💧", "Uap Air 💨"],
				"hint": "Batu adalah benda Padat, Air Minum adalah Cair, Uap Air adalah Gas!",
				"explanation": "Wujud benda: Batu (Padat) -> Air Minum (Cair) -> Uap Air (Gas)."
			}`,
		)

		// Lesson 28: Daur Air & Terjadinya Hujan
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (28, ?, 'builders', 2, 1, 'Daur Air & Terjadinya Hujan', 'Bagaimana air menguap dan turun menjadi hujan', 
			 '{"intro_text":"Air di bumi terus berputar dalam siklus daur air yang menakjubkan!","icon":"🌧️"}', 12, 25)`,
			scienceCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (28, 28, 'sequencing', 1, ?, 10)`,
			`{
				"prompt": "Urutkan proses Daur Air dari awal!",
				"items": ["Air menguap karea panas", "Pengembunan menjadi awan", "Hujan turun ke bumi"],
				"expected_order": ["Air menguap karea panas", "Pengembunan menjadi awan", "Hujan turun ke bumi"],
				"hint": "Mulai dari air laut/sungai yang dipanaskan matahari!",
				"explanation": "Daur air: Penguapan -> Pembentukan awan -> Hujan."
			}`,
		)

		// Lesson 29: Tata Surya & Planet
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (29, ?, 'builders', 3, 1, 'Tata Surya & Planet', 'Mengenal planet-planet yang mengelilingi matahari', 
			 '{"intro_text":"Matahari adalah pusat tata surya yang dikelilingi oleh 8 planet hebat!","icon":"🪐"}', 15, 30)`,
			scienceCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (29, 29, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Planet manakah yang sering disebut sebagai 'Planet Merah' 🔴?",
				"options": ["Mars 🔴", "Bumi 🌍", "Jupiter 🪐"],
				"answer": "Mars 🔴",
				"hint": "Planet ini tampak kemerahan karena permukaan besi oksida!",
				"explanation": "Mars disebut Planet Merah karena kandungan besi oksida yang melimpah di permukaannya."
			}`,
		)
	}

	// ==========================================
	// LANGUAGE LESSONS (BUILDERS: L1 - L3)
	// ==========================================
	var langCatID int64
	err = db.QueryRow(`SELECT id FROM categories WHERE slug = 'language'`).Scan(&langCatID)
	if err == nil && langCatID > 0 {
		// Lesson 30: Abjad & Kosakata Dasar
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (30, ?, 'builders', 1, 1, 'Abjad & Kosakata Dasar', 'Mengenal huruf dan kata-kata benda awal', 
			 '{"intro_text":"Bahasa adalah kunci komunikasi. Yuk tambah kosakatamu!","icon":"🔤"}', 10, 20)`,
			langCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (30, 30, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Kata manakah yang diawali dengan huruf B 🅱️?",
				"options": ["Buku 📖", "Apel 🍎", "Gajah 🐘"],
				"answer": "Buku 📖",
				"hint": "Buku diawali dengan huruf B!",
				"explanation": "B - U - K - U diawali huruf B."
			}`,
		)

		// Lesson 31: Tata Bahasa Dasar (Kata Kerja & Benda)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (31, ?, 'builders', 2, 1, 'Kata Kerja & Benda', 'Membedakan tindakan (kata kerja) dan nama benda', 
			 '{"intro_text":"Kata kerja menyatakan tindakan, sedangkan kata benda menunjukkan barang atau orang!","icon":"✍️"}', 12, 25)`,
			langCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (31, 31, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Manakah di bawah ini yang termasuk KATA KERJA (Tindakan)?",
				"options": ["Membaca 📖", "Meja 🪑", "Pensil ✏️"],
				"answer": "Membaca 📖",
				"hint": "Membaca adalah kegiatan/tindakan yang dilakukan seseorang!",
				"explanation": "'Membaca' adalah kata kerja karena menunjukkan suatu aktivitas."
			}`,
		)

		// Lesson 32: Membaca Paham & Ringkasan Cerita
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (32, ?, 'builders', 3, 1, 'Membaca Paham & Cerita', 'Memahami isi bacaan pendek dan menjawab pertanyaan', 
			 '{"intro_text":"Bacalah cerita pendek dengan cermat untuk menemukan jawabannya!","icon":"📚"}', 15, 30)`,
			langCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (32, 32, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Cerita: 'Budi membawa payung merah karena langit gelap mendung.' Mengapa Budi membawa payung?",
				"options": ["Karena cuaca mendung 🌧️", "Karena cuaca panas ☀️", "Karena ingin bermain ⚽"],
				"answer": "Karena cuaca mendung 🌧️",
				"hint": "Budi mengamati langit gelap mendung!",
				"explanation": "Budi membawa payung karena langit gelap mendung dan bersiap hujan."
			}`,
		)
	}

	// ==========================================
	// LOGIC LESSONS (BUILDERS: L1 - L3)
	// ==========================================
	var logicCatID int64
	err = db.QueryRow(`SELECT id FROM categories WHERE slug = 'logic'`).Scan(&logicCatID)
	if err == nil && logicCatID > 0 {
		// Lesson 33: Pola Gambar & Perbedaan
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (33, ?, 'builders', 1, 1, 'Pola Gambar & Perbedaan', 'Melatih ketelitian dan pengenalan pola gambar', 
			 '{"intro_text":"Amati bentuk dan warna pola dengan teliti!","icon":"🧩"}', 10, 20)`,
			logicCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (33, 33, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Lengkapi pola berikut: 🔴 - 🔵 - 🔴 - 🔵 - __ ?",
				"options": ["🔴", "🔵", "🟡"],
				"answer": "🔴",
				"hint": "Polanya selang-seling antara Merah 🔴 dan Biru 🔵!",
				"explanation": "Setelah Biru 🔵, pola berulang kembali ke Merah 🔴."
			}`,
		)

		// Lesson 34: Labirin Logika & Jalur
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (34, ?, 'builders', 2, 1, 'Labirin Logika & Jalur', 'Menemukan rute tercepat dan menyelesaikan teka-teki', 
			 '{"intro_text":"Gunakan logika untuk memilih rute terbaik menuju tujuan!","icon":"🧭"}', 12, 25)`,
			logicCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (34, 34, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Jika jalur A ada rintangan ⚠️ dan jalur B aman 🟢, jalur mana yang harus dipilih?",
				"options": ["Jalur B 🟢", "Jalur A ⚠️", "Kembali ke awal 🏁"],
				"answer": "Jalur B 🟢",
				"hint": "Pilihlah jalur yang aman tanpa rintangan!",
				"explanation": "Jalur B aman (🟢) sehingga menjadi pilihan logis."
			}`,
		)

		// Lesson 35: Pola Logika Angka (Sudoku Lite)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (35, ?, 'builders', 3, 1, 'Pola Logika Angka', 'Menemukan angka yang hilang dalam barisan pola', 
			 '{"intro_text":"Temukan kelipatan dan rahasia di balik barisan angka!","icon":"🔢"}', 15, 30)`,
			logicCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (35, 35, 'fill_blank', 1, ?, 10)`,
			`{
				"prompt": "Isi angka yang hilang dalam pola berikut: 2, 4, 6, __, 10 !",
				"answer": "8",
				"hint": "Pola ini bertambah 2 setiap langkahnya!",
				"explanation": "6 + 2 = 8. Pola kelipatan 2: 2, 4, 6, 8, 10."
			}`,
		)
	}

	// ==========================================
	// ART LESSONS (BUILDERS: L1 - L3)
	// ==========================================
	var artCatID int64
	err = db.QueryRow(`SELECT id FROM categories WHERE slug = 'art'`).Scan(&artCatID)
	if err == nil && artCatID > 0 {
		// Lesson 36: Warna Dasar & Lukisan Kreasi
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (36, ?, 'builders', 1, 1, 'Warna Dasar & Lukisan Kreasi', 'Mengenal warna dasar dan melukis pada kanvas digital', 
			 '{"intro_text":"Gunakan imajinasi dan kuas warnamu untuk membuat karya indah!","icon":"🎨"}', 10, 20)`,
			artCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (36, 36, 'drawing', 1, ?, 10)`,
			`{
				"prompt": "Gunakan kuas dan stempel stiker untuk melukis pemandangan indahmu! 🎨",
				"hint": "Pilih warna dan stempel bintang atau pelangi!",
				"explanation": "Karya seni melukis yang luar biasa!"
			}`,
		)

		// Lesson 37: Seni Pixel & Pola Gambar
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (37, ?, 'builders', 2, 1, 'Seni Pixel & Pola Gambar', 'Membuat gambar retro 8x8 menggunakan seni pixel', 
			 '{"intro_text":"Seni pixel dibentuk dari kotak-kotak kecil berwarna yang tersusun rapi!","icon":"👾"}', 12, 25)`,
			artCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (37, 37, 'pixel_art', 1, ?, 10)`,
			`{
				"prompt": "Warnai grid pixel di bawah untuk membuat emoji atau karakter retro favoritmu! 👾",
				"hint": "Sentuh setiap kotak pixel untuk memberi warna!",
				"explanation": "Sangat kreatif! Karakter pixel art buatanmu sungguh menarik."
			}`,
		)

		// Lesson 38: Warna Komplementer & Roda Warna
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (38, ?, 'builders', 3, 1, 'Campuran & Roda Warna', 'Pelajari hasil pencampuran dua warna dasar', 
			 '{"intro_text":"Mencampurkan dua warna primer akan menghasilkan warna sekunder yang baru!","icon":"🌈"}', 15, 30)`,
			artCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (38, 38, 'multiple_choice', 1, ?, 10)`,
			`{
				"prompt": "Pencampuran warna Merah 🔴 dan Kuning 🟡 akan menghasilkan warna apa?",
				"options": ["Oranye 🟧", "Hijau 🟩", "Ungu 🟪"],
				"answer": "Oranye 🟧",
				"hint": "Pikirkan warna buah Jeruk!",
				"explanation": "Merah + Kuning = Oranye 🟧."
			}`,
		)
	}

	// ==========================================
	// MINI-GAME LESSONS (CROSS-CATEGORY)
	// ==========================================

	// Lesson 39: Math Racer (Math category, timed mini-game)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (39, ?, 'builders', 4, 1, 'Math Racer ⏱️', 'Pecahkan soal aritmatika sebanyak mungkin dalam 60 detik!', 
		 '{"intro_text":"Seberapa cepat kamu bisa menghitung? Ayo buktikan!","icon":"⏱️","game_type":"math_racer"}', 5, 30)`,
		mathCatID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (39, 39, 'timed', 1, ?, 100)`,
		`{
			"prompt": "Pecahkan soal aritmatika sebanyak mungkin dalam 60 detik! ⏱️🚀",
			"time_limit": 60,
			"operations": ["add", "subtract"],
			"max_number": 50,
			"game_type": "math_racer"
		}`,
	)

	// Lesson 40: Susun Kata (Language category, word builder mini-game)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'language'`).Scan(&langCatID)
	if langCatID > 0 {
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (40, ?, 'builders', 4, 1, 'Susun Kata 📝', 'Susun huruf-huruf acak menjadi kata yang benar', 
			 '{"intro_text":"Dapatkah kamu menyusun huruf acak menjadi kata yang benar? Ayo coba!","icon":"📝","game_type":"word_builder"}', 10, 25)`,
			langCatID,
		)
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (40, 40, 'fill_blank', 1, ?, 50)`,
			`{
				"prompt": "Susun huruf acak menjadi kata yang benar! 📝",
				"game_type": "word_builder",
				"words": [
					{"word": "KUCING", "clue": "🐱 Hewan berkaki empat yang suka ikan"},
					{"word": "BUNGA", "clue": "🌸 Tumbuhan indah yang harum"},
					{"word": "MATAHARI", "clue": "☀️ Benda langit yang bersinar terang di siang hari"},
					{"word": "SEKOLAH", "clue": "🏫 Tempat untuk belajar bersama teman-teman"},
					{"word": "PELANGI", "clue": "🌈 Lengkungan warna-warni di langit setelah hujan"}
				]
			}`,
		)
	}

	// Lesson 41: Labirin Logika (Logic category, maze mini-game)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'logic'`).Scan(&logicCatID)
	if logicCatID > 0 {
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
			 VALUES (41, ?, 'builders', 4, 1, 'Labirin Logika 🧩', 'Arahkan karakter melewati labirin menuju bintang!', 
			 '{"intro_text":"Gunakan perintah arah untuk memandu kucing melewati labirin!","icon":"🧩","game_type":"maze_logic"}', 10, 25)`,
			logicCatID,
		)
		const maze41JSON = `{
			"prompt": "Arahkan 🐱 melewati labirin menuju ⭐! Gunakan tombol arah.",
			"game_type": "maze_logic",
			"maze_data": {
				"width": 6,
				"height": 6,
				"start": [0, 0],
				"goal": [5, 5],
				"walls": [
					[0,1],[1,1],[2,1],[3,1],[4,1],
					[1,3],[2,3],[3,3],[4,3],[5,3],
					[0,5],[1,5],[2,5],[3,5],[4,5]
				]
			},
			"hint": "Ikuti jalan berliku: Ke Kanan (➡️), Turun (⬇️), Ke Kiri (⬅️), Turun (⬇️), Ke Kanan (➡️)!"
		}`
		_, _ = db.Exec(
			`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
			 VALUES (41, 41, 'sequencing', 1, ?, 10)`,
			maze41JSON,
		)
		_, _ = db.Exec(`UPDATE activities SET question_json = ? WHERE id = 41`, maze41JSON)
	}

	// =========================================================================
	// EXPLORERS LESSONS (4-6 YEARS: L1 - L3) & CHALLENGERS (10-12 YEARS: L1 - L3)
	// =========================================================================
	var mathID, codingID, scienceID, languageID, logicID, artID int64
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'math'`).Scan(&mathID)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'coding'`).Scan(&codingID)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'science'`).Scan(&scienceID)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'language'`).Scan(&languageID)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'logic'`).Scan(&logicID)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'art'`).Scan(&artID)

	// --- Explorers Math ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (50, ?, 'explorers', 1, 1, 'Belajar Berhitung 1–10 🧮', 'Belajar menghitung benda dari 1 sampai 10', 
		 '{"intro_text":"Halo adik pintar! Ayo belajar menghitung benda-benda lucu bersama! 🍎","icon":"🧮"}', 5, 10)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (50, 50, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Berapa jumlah buah apel di gambar ini? 🍎🍎🍎",
			"options": ["3 🍎", "5 🍎", "2 🍎"],
			"answer": "3 🍎",
			"hint": "Coba tunjuk dan hitung satu-satu: satu, dua, tiga!",
			"explanation": "Benar! Ada 3 buah apel."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (51, ?, 'explorers', 2, 1, 'Mengenal Bentuk Dasar 🟡', 'Mengenal bentuk lingkaran, persegi, dan segitiga', 
		 '{"intro_text":"Bentuk ada di mana-mana! Mari mengenali bentuk mainanmu! 🧸","icon":"🟡"}', 5, 10)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (51, 51, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Bentuk apakah roda sepeda itu? 🚲",
			"options": ["Lingkaran 🟡", "Persegi 🟦", "Segitiga 🔺"],
			"answer": "Lingkaran 🟡",
			"hint": "Roda sepeda berbentuk bulat dan bisa menggelinding!",
			"explanation": "Benar! Roda berbentuk bulat atau lingkaran."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (52, ?, 'explorers', 3, 1, 'Pola Warna Bergantian 🔴', 'Belajar mengidentifikasi pola warna sederhana', 
		 '{"intro_text":"Pola adalah urutan yang berulang. Ayo melengkapi pola warna! 🎨","icon":"🔴"}', 5, 10)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (52, 52, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Lengkapi pola warna ini: 🟥 - 🟦 - 🟥 - 🟦 - __ ?",
			"options": ["🟥", "🟦", "🟢"],
			"answer": "🟥",
			"hint": "Polanya selang-seling: merah, biru, merah, biru...",
			"explanation": "Benar! Setelah biru kembali ke merah."
		}`,
	)

	// --- Explorers Science ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (53, ?, 'explorers', 1, 1, 'Hewan Darat & Air 🦁', 'Membedakan tempat hidup hewan darat dan air', 
		 '{"intro_text":"Ada hewan yang tinggal di darat, dan ada yang di air. Yuk kita kelompokkan! 🐠","icon":"🦁"}', 5, 10)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (53, 53, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Hewan manakah yang hidup di dalam air laut? 🌊",
			"options": ["Ikan Lumba-lumba 🐬", "Singa 🦁", "Burung Elang 🦅"],
			"answer": "Ikan Lumba-lumba 🐬",
			"hint": "Cari hewan yang berenang menggunakan sirip!",
			"explanation": "Benar! Lumba-lumba tinggal dan berenang di dalam air."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (54, ?, 'explorers', 2, 1, 'Matahari & Bulan ☀️', 'Mengenal perbedaan siang dan malam hari', 
		 '{"intro_text":"Tuhan menciptakan matahari untuk siang hari, dan bulan untuk malam hari! 🌙","icon":"☀️"}', 5, 10)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (54, 54, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Benda langit apa yang bersinar terang and membuat hari menjadi hangat di siang hari? ☀️",
			"options": ["Matahari ☀️", "Bulan 🌙", "Bintang ✨"],
			"answer": "Matahari ☀️",
			"hint": "Benda ini berwarna kuning oranye cerah dan muncul saat kita sekolah!",
			"explanation": "Hebat! Matahari menyinari bumi di siang hari."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (55, ?, 'explorers', 3, 1, 'Bagian Tubuh Kita 👀', 'Mengenal panca indra dan fungsinya', 
		 '{"intro_text":"Tubuh kita luar biasa! Ayo belajar tentang panca indra kita! 👂","icon":"👀"}', 5, 10)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (55, 55, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Bagian tubuh mana yang kita gunakan untuk mendengarkan lagu dan cerita? 🎵",
			"options": ["Telinga 👂", "Mata 👀", "Hidung 👃"],
			"answer": "Telinga 👂",
			"hint": "Letaknya ada di sebelah kanan dan kiri kepala kita!",
			"explanation": "Tepat! Kita mendengarkan suara menggunakan telinga."
		}`,
	)

	// --- Explorers Coding ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (56, ?, 'explorers', 1, 1, 'Urutan Langkah Pagi 🌅', 'Belajar mengurutkan aktivitas sehari-hari', 
		 '{"intro_text":"Koding adalah tentang urutan langkah yang benar! Mari urutkan kegiatan pagimu.","icon":"🌅"}', 5, 10)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (56, 56, 'drag_drop', 1, ?, 10)`,
		`{
			"prompt": "Urutkan kegiatan bangun pagi yang benar!",
			"items": ["Mandi 🚿", "Bangun Tidur 🛌", "Sarapan Pagi 🍳"],
			"expected_order": ["Bangun Tidur 🛌", "Mandi 🚿", "Sarapan Pagi 🍳"],
			"hint": "Kita harus bangun dari tempat tidur terlebih dahulu!",
			"explanation": "Pintar! Urutannya: Bangun tidur -> Mandi -> Sarapan."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (57, ?, 'explorers', 2, 1, 'Perintah Arah Kucing 🐱', 'Belajar arah maju, kiri, dan kanan sederhana', 
		 '{"intro_text":"Ayo bantu kucing berjalan ke tujuannya dengan blok arah! ➡️","icon":"🐱"}', 5, 10)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (57, 57, 'block_code', 1, ?, 10)`,
		`{
			"prompt": "Susun perintah agar kucing melangkah: Maju, lalu Belok Kanan!",
			"available_blocks": ["Maju ⬆️", "Belok Kanan ➡️", "Belok Kiri ⬅️"],
			"expected_order": ["Maju ⬆️", "Belok Kanan ➡️"],
			"hint": "Susun blok 'Maju ⬆️' terlebih dahulu, baru belok!",
			"explanation": "Hebat! Kucing berhasil melangkah maju lalu belok kanan."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (58, ?, 'explorers', 3, 1, 'Ulang Langkah Seru 🔄', 'Mengenal konsep perulangan sederhana', 
		 '{"intro_text":"Jika ingin melakukan hal yang sama berkali-kali, kita bisa menggunakan Loop/Pengulangan!","icon":"🔄"}', 5, 10)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (58, 58, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Jika kita ingin berjalan ke depan sebanyak 3 kali, mana perintah yang paling singkat?",
			"options": ["Ulangi 3x: Maju", "Maju, Mundur, Maju", "Maju 10x"],
			"answer": "Ulangi 3x: Maju",
			"hint": "Gunakan perintah 'Ulangi' untuk menghemat tenaga!",
			"explanation": "Benar! Ulangi 3x membuat perintah menjadi sangat singkat."
		}`,
	)

	// --- Explorers Language ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (59, ?, 'explorers', 1, 1, 'Mengenal Huruf Vokal 🅰️', 'Belajar huruf vokal A, I, U, E, O', 
		 '{"intro_text":"Huruf vokal sangat penting dalam kata. Yuk kenali bunyinya! 🍎","icon":"🅰️"}', 5, 10)`,
		languageID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (59, 59, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Kata 'APEL' diawali dengan huruf vokal apa? 🍎",
			"options": ["A 🅰️", "I ℹ️", "U 🇺"],
			"answer": "A 🅰️",
			"hint": "Bunyikan kata A-A-Apel!",
			"explanation": "Hebat! Apel diawali dengan huruf A."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (60, ?, 'explorers', 2, 1, 'Mengeja Kata Hewan 🐶', 'Mengeja nama-nama hewan sederhana', 
		 '{"intro_text":"Ayo mengeja nama teman-teman hewan kita yang lucu! 🐈","icon":"🐶"}', 5, 10)`,
		languageID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (60, 60, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Bagaimana ejaan yang benar untuk nama hewan ini: 🐈?",
			"options": ["K-U-C-I-N-G", "K-U-J-I-N-G", "K-U-S-I-N-G"],
			"answer": "K-U-C-I-N-G",
			"hint": "Huruf tengahnya adalah C seperti Cacing!",
			"explanation": "Pintar! Ejaannya adalah K-U-C-I-N-G (Kucing)."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (61, ?, 'explorers', 3, 1, 'Mencocokkan Kata Benda 🚗', 'Menghubungkan gambar dengan kata benda yang tepat', 
		 '{"intro_text":"Lihat benda di sekitarmu dan sebutkan namanya! 🧸","icon":"🚗"}', 5, 10)`,
		languageID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (61, 61, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Kendaraan roda empat yang berbunyi 'Brumm' disebut?",
			"options": ["Mobil 🚗", "Sepeda 🚲", "Kereta 🚂"],
			"answer": "Mobil 🚗",
			"hint": "Mobil memiliki pintu dan roda empat!",
			"explanation": "Luar biasa! Itu adalah Mobil."
		}`,
	)

	// --- Explorers Logic ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (62, ?, 'explorers', 1, 1, 'Mencari Gambar Berbeda 🔎', 'Melatih ketelitian dengan mencari benda yang berbeda', 
		 '{"intro_text":"Amati gambar-gambar berikut dengan sangat teliti! Ada satu yang berbeda.","icon":"🔎"}', 5, 10)`,
		logicID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (62, 62, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Manakah gambar yang paling BERBEDA di antara pilihan berikut?",
			"options": ["Pisang 🍌 (Buah)", "Apel 🍎 (Buah)", "Mobil 🚗 (Kendaraan)"],
			"answer": "Mobil 🚗 (Kendaraan)",
			"hint": "Dua pilihan adalah buah yang bisa dimakan, satu adalah kendaraan!",
			"explanation": "Hebat! Mobil bukan buah, jadi ia berbeda."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (63, ?, 'explorers', 2, 1, 'Mengingat Kartu Hewan 🃏', 'Permainan mencocokkan ingatan visual', 
		 '{"intro_text":"Gunakan ingatanmu untuk mencocokkan sepasang hewan yang sama! 🐵","icon":"🃏"}', 5, 10)`,
		logicID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (63, 63, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Jika kartu pertama adalah Kucing 🐱, pasangannya yang tepat adalah kartu bergambar apa?",
			"options": ["Kucing 🐱", "Anjing 🐶", "Kelinci 🐰"],
			"answer": "Kucing 🐱",
			"hint": "Cari hewan yang sama persis bentuk dan suaranya!",
			"explanation": "Tepat sekali! Pasangan Kucing adalah Kucing 🐱."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (64, ?, 'explorers', 3, 1, 'Labirin Keju Tikus 🧀', 'Menemukan rute sederhana untuk tikus kecil', 
		 '{"intro_text":"Bantu tikus kecil menemukan jalannya menuju keju lezat! 🐭","icon":"🧀"}', 5, 10)`,
		logicID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (64, 64, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Jalan ke atas tertutup dinding 🧱, jalan ke kanan terbuka bebas 🟢. Ke mana tikus harus berjalan?",
			"options": ["Ke Kanan ➡️", "Ke Atas ⬆️", "Kembali ke Belakang ⬅️"],
			"answer": "Ke Kanan ➡️",
			"hint": "Pilihlah jalan yang tidak tertutup dinding!",
			"explanation": "Benar! Jalan ke kanan aman dan terbuka."
		}`,
	)

	// --- Explorers Art ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (65, ?, 'explorers', 1, 1, 'Mewarnai Bebas 🎨', 'Melukis dan mewarnai di kanvas dengan kuas', 
		 '{"intro_text":"Gunakan kuas warnamu untuk membuat lukisan yang indah! 🖌️","icon":"🎨"}', 5, 10)`,
		artID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (65, 65, 'drawing', 1, ?, 10)`,
		`{
			"prompt": "Gambarlah pemandangan atau bunga kesukaanmu! 🎨",
			"hint": "Gunakan stempel stiker matahari ☀️ atau pelangi 🌈 agar lebih indah!",
			"explanation": "Sangat bagus sekali gambarmu! Sungguh indah."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (66, ?, 'explorers', 2, 1, 'Membuat Pola Warna-Warni 🌈', 'Mengenal keindahan susunan warna pelangi', 
		 '{"intro_text":"Pelangi memiliki susunan warna yang indah. Mari susun bersama! 🌈","icon":"🌈"}', 5, 10)`,
		artID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (66, 66, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Warna pertama pada pelangi yang indah adalah?",
			"options": ["Merah 🔴", "Hijau 🟢", "Biru 🔵"],
			"answer": "Merah 🔴",
			"hint": "Ingat lagu: 'Me-Ji-Ku-Hi-Bi-Ni-Yu'!",
			"explanation": "Benar! Merah adalah warna teratas pelangi."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (67, ?, 'explorers', 3, 1, 'Mencampur Warna Dasar 🧪', 'Belajar hasil pencampuran warna primer', 
		 '{"intro_text":"Jika dua warna dicampur, akan lahir warna baru yang ajaib! 🎨","icon":"🧪"}', 5, 10)`,
		artID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (67, 67, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Jika warna Biru 🔵 dicampur dengan Kuning 🟡, warna apa yang akan terbentuk?",
			"options": ["Hijau 🟩", "Oranye 🟧", "Ungu 🟪"],
			"answer": "Hijau 🟩",
			"hint": "Pikirkan warna daun pohon!",
			"explanation": "Campuran Biru and Kuning menghasilkan warna Hijau 🟩."
		}`,
	)

	// --- Challengers Math ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (80, ?, 'challengers', 1, 1, 'Penjumlahan & Pengurangan Besar 🧮', 'Latihan berhitung angka ratusan dan ribuan secara cepat', 
		 '{"intro_text":"Selamat datang penantang hebat! Mari mengasah kemampuan berhitung angka besar!","icon":"🧮"}', 10, 20)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (80, 80, 'fill_blank', 1, ?, 10)`,
		`{
			"prompt": "Berapakah hasil dari 1450 + 2350? Ketik jawabanmu!",
			"answer": "3800",
			"hint": "Gunakan penjumlahan susun: 1400 + 2300 = 3700, lalu tambahkan 50 + 50 = 100!",
			"explanation": "Benar! 1450 + 2350 = 3800."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (81, ?, 'challengers', 2, 1, 'Perkalian & Pembagian Cepat ⚡', 'Perkalian dua digit dan pembagian bersusun', 
		 '{"intro_text":"Perkalian cepat adalah teknik melatih konsentrasi dan pemecahan masalah numerik.","icon":"⚡"}', 10, 20)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (81, 81, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Berapakah hasil dari 25 × 12?",
			"options": ["300", "280", "320", "250"],
			"answer": "300",
			"hint": "25 × 10 = 250. Tambahkan dengan 25 × 2 = 50!",
			"explanation": "Tepat! 25 × 12 = 250 + 50 = 300."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (82, ?, 'challengers', 3, 1, 'Mengenal Pecahan & Desimal 🍰', 'Belajar konsep bagian dari keseluruhan dan desimal', 
		 '{"intro_text":"Pecahan membantu kita membagi kue atau pizza sama rata dengan teman! 🍕","icon":"🍰"}', 10, 20)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (82, 82, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Manakah bentuk desimal dari pecahan 3/4?",
			"options": ["0.75", "0.50", "0.34", "0.80"],
			"answer": "0.75",
			"hint": "Pikirkan 3 perempat bagian dari 100!",
			"explanation": "Benar! 3/4 sama dengan 75/100, yaitu 0.75."
		}`,
	)

	// --- Challengers Science ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (83, ?, 'challengers', 1, 1, 'Sistem Organ Manusia 🫁', 'Mengenal organ tubuh pernapasan dan pencernaan', 
		 '{"intro_text":"Di dalam tubuh kita, organ bekerja sama tanpa henti untuk menjaga kita tetap hidup!","icon":"🫁"}', 10, 20)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (83, 83, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Organ tubuh manakah yang berfungsi memompa darah ke seluruh tubuh kita? 🩸",
			"options": ["Jantung ❤️", "Paru-paru 🫁", "Lambung 🥣"],
			"answer": "Jantung ❤️",
			"hint": "Organ ini berdetak di sebelah kiri rongga dada kita!",
			"explanation": "Tepat sekali! Jantung memompa darah kaya oksigen ke seluruh tubuh."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (84, ?, 'challengers', 2, 1, 'Rantai Makanan & Ekosistem 🌾', 'Mengenal produsen, konsumen, dan pengurai di alam', 
		 '{"intro_text":"Semua makhluk hidup saling bergantung untuk makanan. Urutan makan-memakan ini disebut rantai makanan!","icon":"🌾"}', 10, 20)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (84, 84, 'sequencing', 1, ?, 10)`,
		`{
			"prompt": "Urutkan rantai makanan di sawah dari produsen hingga konsumen puncak!",
			"items": ["Ular 🐍", "Padi 🌾", "Belalang 🦗", "Katak 🐸"],
			"expected_order": ["Padi 🌾", "Belalang 🦗", "Katak 🐸", "Ular 🐍"],
			"hint": "Mulai dari padi sebagai tumbuhan hijau (produsen)!",
			"explanation": "Sempurna! Padi dimakan belalang, belalang dimakan katak, katak dimakan ular."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (85, ?, 'challengers', 3, 1, 'Rangkaian Listrik Logika 🔌', 'Mengenal rangkaian seri, paralel, dan penghantar listrik', 
		 '{"intro_text":"Arus listrik mengalir dalam jalur tertutup. Jalur ini disebut rangkaian listrik!","icon":"🔌"}', 10, 20)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (85, 85, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Jika salah satu lampu pada rangkaian listrik padam dan semua lampu lain ikut padam, rangkaian tersebut adalah rangkaian?",
			"options": ["Seri 🔗", "Paralel 🔀", "Campuran 🔄"],
			"answer": "Seri 🔗",
			"hint": "Semua komponen dihubungkan berurutan dalam satu jalur aliran tunggal.",
			"explanation": "Benar! Pada rangkaian seri, aliran listrik hanya satu jalur sehingga jika terputus semua mati."
		}`,
	)

	// --- Challengers Coding ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (86, ?, 'challengers', 1, 1, 'Variabel & Menyimpan Data 📦', 'Memahami konsep variabel sebagai kotak penyimpanan informasi', 
		 '{"intro_text":"Variabel digunakan oleh komputer untuk menyimpan angka atau teks dengan nama tertentu!","icon":"📦"}', 10, 20)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (86, 86, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Jika kita menulis kode: skor = 10 dan kemudian skor = skor + 5, berapa isi variabel skor sekarang?",
			"options": ["15", "10", "5", "skor5"],
			"answer": "15",
			"hint": "Ambil nilai awal (10) lalu tambahkan dengan 5!",
			"explanation": "Hebat! Variabel menyimpan nilai baru: 10 + 5 = 15."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (87, ?, 'challengers', 2, 1, 'Logika Kondisi Bercabang 🔀', 'Menyusun keputusan rumit menggunakan If - Else', 
		 '{"intro_text":"Komputer menggunakan logika If-Else bercabang untuk menangani kondisi yang kompleks!","icon":"🔀"}', 10, 20)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (87, 87, 'block_code', 1, ?, 10)`,
		`{
			"prompt": "Susun blok kode: Jika Nilai >= 70 -> Tampilkan 'Lulus' -> Jika Tidak -> Tampilkan 'Ulangi'!",
			"available_blocks": ["Jika Nilai >= 70 🎓", "Tampilkan 'Lulus' ✅", "Jika Tidak ❌", "Tampilkan 'Ulangi' 🔄"],
			"expected_order": ["Jika Nilai >= 70 🎓", "Tampilkan 'Lulus' ✅", "Jika Tidak ❌", "Tampilkan 'Ulangi' 🔄"],
			"hint": "Susun sesuai urutan kondisi, tindakan benar, jika tidak, dan tindakan salah!",
			"explanation": "Sempurna! Logika kondisi bercabang berhasil disusun dengan tepat."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (88, ?, 'challengers', 3, 1, 'Nested Loop (Loop Bersarang) 🔄', 'Mengenal konsep perulangan di dalam perulangan', 
		 '{"intro_text":"Nested Loop terjadi ketika ada loop di dalam loop lain. Berguna membuat pola baris-kolom!","icon":"🔄"}', 10, 20)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (88, 88, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Jika Loop Luar berjalan 3 kali, dan di dalamnya ada Loop Dalam yang berjalan 2 kali, berapa total tindakan yang dieksekusi?",
			"options": ["6 kali", "5 kali", "3 kali", "2 kali"],
			"answer": "6 kali",
			"hint": "Kalikan jumlah perulangan luar dengan perulangan dalam (3 × 2)!",
			"explanation": "Hebat! Loop dalam berjalan 2 kali untuk setiap langkah loop luar, totalnya 3 × 2 = 6 kali."
		}`,
	)

	// --- Challengers Language ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (89, ?, 'challengers', 1, 1, 'Kosakata Bahasa Indonesia 📚', 'Mengenal sinonim, antonim, dan padanan kata tingkat lanjut', 
		 '{"intro_text":"Mari perbanyak kosakata Bahasa Indonesiamu agar menulis dan membaca semakin asyik!","icon":"📚"}', 10, 20)`,
		languageID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (89, 89, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Manakah padanan kata (sinonim) yang paling tepat untuk kata 'EVALUASI'?",
			"options": ["Penilaian", "Permulaan", "Peralatan", "Penyebaran"],
			"answer": "Penilaian",
			"hint": "Mengevaluasi pekerjaan berarti memberikan nilai pada hasil kerja.",
			"explanation": "Benar! Sinonim evaluasi adalah penilaian atau penaksiran."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (90, ?, 'challengers', 2, 1, 'Menyusun Kalimat Efektif ✍️', 'Menyusun kata acak menjadi kalimat sesuai kaidah SPOK', 
		 '{"intro_text":"Kalimat efektif mempermudah orang lain memahami maksud tulisan kita!","icon":"✍️"}', 10, 20)`,
		languageID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (90, 90, 'sequencing', 1, ?, 10)`,
		`{
			"prompt": "Susun kalimat berikut agar mengikuti kaidah Subjek, Predikat, Objek, Keterangan (SPOK)!",
			"items": ["di perpustakaan", "membaca", "Siti", "buku cerita"],
			"expected_order": ["Siti", "membaca", "buku cerita", "di perpustakaan"],
			"hint": "Mulai dengan pelaku/nama orang (Subjek) diikuti kata kerja (Predikat)!",
			"explanation": "Bagus sekali! Kalimat efektif: Siti (S) membaca (P) buku cerita (O) di perpustakaan (K)."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (91, ?, 'challengers', 3, 1, 'Pemahaman Bacaan Kritis 📖', 'Membaca teks panjang dan menarik kesimpulan kritis', 
		 '{"intro_text":"Bacalah paragraf pendek di bawah dan tentukan ide pokok ceritanya!","icon":"📖"}', 10, 20)`,
		languageID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (91, 91, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Paragraf: 'Hutan bakau memiliki peran penting bagi ekosistem pantai. Akar bakau mencegah abrasi tanah oleh ombak laut. Selain itu, hutan bakau menjadi rumah bagi kepiting dan ikan kecil.' Apakah ide pokok paragraf tersebut?",
			"options": ["Manfaat hutan bakau bagi pantai 🌊", "Cara kepiting mencari makan 🦀", "Proses terjadinya ombak laut 🌊"],
			"answer": "Manfaat hutan bakau bagi pantai 🌊",
			"hint": "Perhatikan kalimat pertama dan kalimat pendukung setelahnya!",
			"explanation": "Benar! Seluruh paragraf membahas peran penting (manfaat) hutan bakau."
		}`,
	)

	// --- Challengers Logic ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (92, ?, 'challengers', 1, 1, 'Teka-teki Pola Deret 🔢', 'Menemukan pola angka deret matematika logika', 
		 '{"intro_text":"Teka-teki deret angka melatih otak menemukan hubungan logis antar elemen!","icon":"🔢"}', 10, 20)`,
		logicID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (92, 92, 'fill_blank', 1, ?, 10)`,
		`{
			"prompt": "Lengkapi angka berikutnya dari deret ini: 1, 3, 6, 10, 15, __ ?",
			"answer": "21",
			"hint": "Selisih angkanya bertambah: +2, +3, +4, +5... Langkah berikutnya adalah +6!",
			"explanation": "Tepat! 15 + 6 = 21. Pola pertambahan berurutan."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (93, ?, 'challengers', 2, 1, 'Logika Sudoku Grid 🧩', 'Memahami aturan peletakan angka tanpa duplikasi', 
		 '{"intro_text":"Sudoku adalah teka-teki penempatan angka berdasarkan baris, kolom, dan kotak tanpa duplikat!","icon":"🧩"}', 10, 20)`,
		logicID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (93, 93, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Pada baris Sudoku 4x4: [1, 2, 4, __ ]. Angka berapakah yang harus mengisi tempat kosong?",
			"options": ["3", "4", "5", "1"],
			"answer": "3",
			"hint": "Gunakan angka 1 sampai 4. Angka mana yang belum ada di baris tersebut?",
			"explanation": "Benar! Angka 3 melengkapi baris agar berisi angka 1-4 tanpa berulang."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (94, ?, 'challengers', 3, 1, 'Kriptografi Caesar Cipher 🔑', 'Belajar menyandikan dan memecahkan pesan rahasia', 
		 '{"intro_text":"Caesar Cipher adalah metode sandi geser tertua. Setiap huruf digeser beberapa langkah!","icon":"🔑"}', 10, 20)`,
		logicID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (94, 94, 'fill_blank', 1, ?, 10)`,
		`{
			"prompt": "Jika kunci geser adalah +1 (A menjadi B, B menjadi C), pecahkan pesan rahasia berikut: 'BUB'!",
			"answer": "CVC",
			"hint": "Ganti huruf B dengan C, U dengan V, dan B dengan C!",
			"explanation": "Hebat! Pesan rahasianya adalah CVC (B+1=C, U+1=V, B+1=C)."
		}`,
	)

	// --- Challengers Art ---
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (95, ?, 'challengers', 1, 1, 'Roda Warna Lanjutan 🌈', 'Mengenal harmoni warna analog, monokrom, dan komplementer', 
		 '{"intro_text":"Seni digital yang baik menggunakan skema harmoni warna agar enak dipandang!","icon":"🌈"}', 10, 20)`,
		artID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (95, 95, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Skema warna yang menggunakan tingkat kecerahan berbeda dari satu warna tunggal disebut?",
			"options": ["Monokromatik", "Komplementer", "Analog"],
			"answer": "Monokromatik",
			"hint": "Mono berarti satu, kromatik berarti warna.",
			"explanation": "Benar! Monokromatik menggunakan variasi terang-gelap dari satu warna dasar."
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (96, ?, 'challengers', 2, 1, 'Seni Pixel Grid 12x12 👾', 'Membuat lukisan pixel art berukuran 12x12', 
		 '{"intro_text":"Seni pixel 12x12 memberimu ruang lebih besar untuk detail gambar yang keren!","icon":"👾"}', 10, 20)`,
		artID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (96, 96, 'pixel_art', 1, ?, 10)`,
		`{
			"prompt": "Buatlah karakter pixel art favoritmu dalam grid 12x12! 👾",
			"hint": "Gunakan palet warna retro untuk mewarnai grid!",
			"explanation": "Karya seni pixel 12x12 buatanmu sungguh menakjubkan!"
		}`,
	)

	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (97, ?, 'challengers', 3, 1, 'Digital Canvas Master 🎨', 'Menggambar sketsa kreatif dengan kombinasi kuas dan stempel', 
		 '{"intro_text":"Kanvas digital adalah tempat bebas berekspresi. Mari buktikan kreativitasmu!","icon":"🎨"}', 10, 20)`,
		artID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (97, 97, 'drawing', 1, ?, 10)`,
		`{
			"prompt": "Lukislah poster bertema alam atau luar angkasa impianmu! 🚀🌲",
			"hint": "Gunakan kuas tebal untuk latar belakang dan kuas tipis untuk detail!",
			"explanation": "Luar biasa! Poster buatanmu siap dipajang!"
		}`,
	)

	fmt.Println("Seeded Math, Coding, Toddler, Science, Language, Logic, Art, Explorers, and Challengers content")
	return nil
}

// SeedL4L5Content inserts Level 4 and Level 5 lessons and activities for Math, Science, and Coding.
func SeedL4L5Content(db *sql.DB) error {
	var mathID, codingID, scienceID int64
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'math'`).Scan(&mathID)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'coding'`).Scan(&codingID)
	_ = db.QueryRow(`SELECT id FROM categories WHERE slug = 'science'`).Scan(&scienceID)

	// =========================================================================
	// MATH LEVEL 4 & 5 (EXPLORERS, BUILDERS, CHALLENGERS)
	// =========================================================================

	// 100. Explorers Math L4: Pola & Urutan Gambar 🔢
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (100, ?, 'explorers', 4, 1, 'Pola & Urutan Gambar 🔢', 'Belajar menemukan pola urutan bentuk dan gambar yang berulang',
		 '{"intro_text":"Lihat urutan gambarnya ya! Pola apa yang muncul berikutnya?","icon":"🔢"}', 8, 15)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1000, 100, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Pola: 🔴 🔵 🔴 🔵 🔴 ... Gambar berikutnya adalah?",
			"options": ["🔵", "🔴", "🟡", "🟢"],
			"answer": "🔵",
			"hint": "Perhatikan warnanya bergantian: merah, biru, merah, biru...",
			"explanation": "Hebat! Setelah merah (🔴) selalu diikuti biru (🔵)."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1001, 100, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Urutkan bentuk dari yang TERKECIL ke yang TERBESAR!",
			"items": ["Lingkaran Besar", "Lingkaran Kecil", "Lingkaran Sedang"],
			"expected_order": ["Lingkaran Kecil", "Lingkaran Sedang", "Lingkaran Besar"],
			"hint": "Mulai dari yang paling kecil!",
			"explanation": "Pintar! Urutannya: Kecil -> Sedang -> Besar."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1002, 100, 'matching', 3, ?, 10)`,
		`{
			"prompt": "Jodohkan angka dengan jumlah benda yang sesuai!",
			"pairs": [
				{"left": "3 🍎", "right": "Tiga Apel"},
				{"left": "5 🌟", "right": "Lima Bintang"},
				{"left": "2 🐱", "right": "Dua Kucing"}
			],
			"hint": "Hitung jumlah benda di sebelah kiri!",
			"explanation": "Luar biasa! Semua pasangan sudah tepat."
		}`,
	)

	// 101. Explorers Math L5: Ukuran: Besar/Kecil, Panjang/Pendek 📏
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (101, ?, 'explorers', 5, 1, 'Ukuran: Besar/Kecil, Panjang/Pendek 📏', 'Membandingkan ukuran benda di sekitar kita',
		 '{"intro_text":"Setiap benda punya ukuran berbeda. Ada yang besar, kecil, panjang, atau pendek!","icon":"📏"}', 8, 15)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1003, 101, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Manakah di antara berikut yang ukurannya PALING BESAR?",
			"options": ["Gajah 🐘", "Kucing 🐱", "Semut 🐜", "Burung 🐦"],
			"answer": "Gajah 🐘",
			"hint": "Hewan mana yang badannya paling tinggi dan besar?",
			"explanation": "Benar! Gajah adalah hewan paling besar di antara pilihan tersebut."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1004, 101, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Urutkan pensil berikut dari yang PALING PENDEK ke PALING PANJANG!",
			"items": ["Pensil Panjang (15 cm)", "Pensil Pendek (5 cm)", "Pensil Sedang (10 cm)"],
			"expected_order": ["Pensil Pendek (5 cm)", "Pensil Sedang (10 cm)", "Pensil Panjang (15 cm)"],
			"hint": "Bandingkan angka sentimeternya!",
			"explanation": "Hebat! Urutannya: 5 cm -> 10 cm -> 15 cm."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1005, 101, 'matching', 3, ?, 10)`,
		`{
			"prompt": "Jodohkan benda dengan kata ukuran yang tepat!",
			"pairs": [
				{"left": "Jerapah 🦒", "right": "Tinggi"},
				{"left": "Kura-kura 🐢", "right": "Pendek"},
				{"left": "Paus 🐋", "right": "Besar"}
			],
			"hint": "Bayangkan tinggi dan besarnya hewan-hewan tersebut!",
			"explanation": "Tepat sekali! Kamu memahami perbandingan ukuran."
		}`,
	)

	// 102. Builders Math L4: Mengenal Pecahan ½, ¼, ¾ 🍕
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (102, ?, 'builders', 4, 1, 'Mengenal Pecahan ½, ¼, ¾ 🍕', 'Memahami bagian dari satu kesatuan utuh',
		 '{"intro_text":"Membagi pizza sama rata mengajarkan kita tentang pecahan!","icon":"🍕"}', 10, 20)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1006, 102, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Jika 1 martabak dipotong menjadi 4 bagian sama besar, dan kamu makan 1 potong, berapa bagian yang kamu makan?",
			"options": ["1/4", "1/2", "3/4", "1/3"],
			"answer": "1/4",
			"hint": "Pembilang adalah potong yang dimakan (1), penyebut adalah total potong (4).",
			"explanation": "Tepat! 1 dari 4 potong bernilai 1/4 (seperempat)."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1007, 102, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Setengah dari 20 kue adalah ... kue. (Ketik angka saja)",
			"answer": "10",
			"hint": "Setengah artinya 20 dibagi 2.",
			"explanation": "Hebat! 20 ÷ 2 = 10 kue."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1008, 102, 'drag_drop', 3, ?, 10)`,
		`{
			"prompt": "Urutkan nilai pecahan berikut dari yang TERKECIL ke TERBESAR!",
			"items": ["3/4", "1/4", "1/2"],
			"expected_order": ["1/4", "1/2", "3/4"],
			"hint": "Ingat: 1/4 = 0.25, 1/2 = 0.5, 3/4 = 0.75",
			"explanation": "Benar! 1/4 < 1/2 < 3/4."
		}`,
	)

	// 103. Builders Math L5: Bangun Datar & Simetri Lanjutan 📐
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (103, ?, 'builders', 5, 1, 'Bangun Datar & Simetri Lanjutan 📐', 'Menghitung sisi, sudut, dan simetri lipat bangun datar',
		 '{"intro_text":"Bangun datar memiliki sisi dan sudut. Mari pelajari ciri khas masing-masing!","icon":"📐"}', 10, 20)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1009, 103, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Bangun datar yang memiliki 3 sisi dan 3 sudut adalah?",
			"options": ["Segitiga", "Persegi", "Jajar Genjang", "Trapesium"],
			"answer": "Segitiga",
			"hint": "Namanya diawali kata 'segi' diikuti angka jumlah sisinya.",
			"explanation": "Benar! Segitiga selalu memiliki 3 sisi dan 3 sudut."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1010, 103, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Berapa banyak sumbu simetri lipat yang dimiliki oleh sebuah persegi?",
			"answer": "4",
			"hint": "Persegi bisa dilipat vertikal, horizontal, dan 2 kali secara diagonal.",
			"explanation": "Tepat! Persegi memiliki 4 sumbu simetri lipat."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1011, 103, 'matching', 3, ?, 10)`,
		`{
			"prompt": "Jodohkan bangun datar dengan jumlah sisinya!",
			"pairs": [
				{"left": "Segilima", "right": "5 Sisi"},
				{"left": "Segienam", "right": "6 Sisi"},
				{"left": "Persegi Panjang", "right": "4 Sisi"}
			],
			"hint": "Perhatikan nama bangun datarnya!",
			"explanation": "Luar biasa! Semua pasangan tepat."
		}`,
	)

	// 104. Challengers Math L4: Geometri: Luas, Keliling & Sudut 📐
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (104, ?, 'challengers', 4, 1, 'Geometri: Luas, Keliling & Sudut 📐', 'Menghitung luas dan keliling bangun datar serta jenis-jenis sudut',
		 '{"intro_text":"Luas adalah daerah di dalam bangun, keliling adalah panjang garis tepinya!","icon":"📐"}', 12, 25)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1012, 104, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Berapakah luas persegi panjang dengan panjang 8 cm dan lebar 5 cm?",
			"options": ["40 cm²", "26 cm²", "13 cm²", "35 cm²"],
			"answer": "40 cm²",
			"hint": "Rumus luas persegi panjang = panjang × lebar.",
			"explanation": "Benar! Luas = 8 cm × 5 cm = 40 cm²."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1013, 104, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Berapakah keliling persegi yang memiliki panjang sisi 6 cm?",
			"answer": "24",
			"hint": "Rumus keliling persegi = 4 × sisi.",
			"explanation": "Tepat! Keliling = 4 × 6 cm = 24 cm."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1014, 104, 'drag_drop', 3, ?, 10)`,
		`{
			"prompt": "Urutkan jenis sudut dari yang terkecil hingga terbesar ukurannya!",
			"items": ["Sudut Tumpul (>90°)", "Sudut Lancip (<90°)", "Sudut Siku-siku (90°)"],
			"expected_order": ["Sudut Lancip (<90°)", "Sudut Siku-siku (90°)", "Sudut Tumpul (>90°)"],
			"hint": "Lancip itu kurang dari 90 derajat, siku-siku tepat 90 derajat.",
			"explanation": "Hebat! Urutan sudut: Lancip (<90°) -> Siku-siku (90°) -> Tumpul (>90°)."
		}`,
	)

	// 105. Challengers Math L5: Data & Statistik: Mean, Median, Modus 📊
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (105, ?, 'challengers', 5, 1, 'Data & Statistik: Mean, Median, Modus 📊', 'Menganalisis data, menghitung rata-rata dan nilai tengah',
		 '{"intro_text":"Statistik membantu kita mengolah kumpulan angka menjadi informasi yang berguna!","icon":"📊"}', 12, 25)`,
		mathID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1015, 105, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Nilai ujian 5 siswa: 6, 7, 8, 8, 9. Berapakah modus (nilai yang paling sering muncul)?",
			"options": ["8", "7", "6", "7.6"],
			"answer": "8",
			"hint": "Modus adalah angka yang frekuensi kemunculannya terbanyak.",
			"explanation": "Benar! Angka 8 muncul 2 kali, lebih sering dari angka lainnya."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1016, 105, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Berapakah rata-rata (mean) dari data angka berikut: 4, 6, 8, 10, 12?",
			"answer": "8",
			"hint": "Jumlahkan semua angka (4+6+8+10+12 = 40), lalu bagi dengan jumlah datanya (5).",
			"explanation": "Tepat! Mean = 40 ÷ 5 = 8."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1017, 105, 'drag_drop', 3, ?, 10)`,
		`{
			"prompt": "Urutkan langkah menghitung Median (nilai tengah) dari data acak!",
			"items": ["Hitung jumlah data", "Urutkan data dari terkecil ke terbesar", "Pilih nilai yang berada tepat di tengah"],
			"expected_order": ["Urutkan data dari terkecil ke terbesar", "Hitung jumlah data", "Pilih nilai yang berada tepat di tengah"],
			"hint": "Langkah pertama selalu mengurutkan data terlebih dahulu!",
			"explanation": "Luar biasa! Urutan statistik median sudah tepat."
		}`,
	)

	// =========================================================================
	// SCIENCE LEVEL 4 & 5 (EXPLORERS, BUILDERS, CHALLENGERS)
	// =========================================================================

	// 106. Explorers Science L4: Tubuhku & Panca Indra 👁️
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (106, ?, 'explorers', 4, 1, 'Tubuhku & Panca Indra 👁️', 'Mengenal 5 indra manusia dan kegunaannya',
		 '{"intro_text":"Tubuh kita luar biasa! Kita punya 5 indra untuk merasakan dunia.","icon":"👁️"}', 8, 15)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1018, 106, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Indra manakah yang kita gunakan untuk mendengarkan musik?",
			"options": ["Telinga 👂", "Mata 👁️", "Hidung 👃", "Lidah 👅"],
			"answer": "Telinga 👂",
			"hint": "Alat tubuh di samping kepala yang mendengar suara.",
			"explanation": "Pintar! Telinga adalah indra pendengaran."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1019, 106, 'matching', 2, ?, 10)`,
		`{
			"prompt": "Pasangkan organ indra dengan fungsinya!",
			"pairs": [
				{"left": "Mata 👁️", "right": "Melihat"},
				{"left": "Hidung 👃", "right": "Mencium Bau"},
				{"left": "Lidah 👅", "right": "Merasakan Makanan"}
			],
			"hint": "Pikirkan apa yang kamu lakukan dengan organ tersebut!",
			"explanation": "Hebat! Semua panca indra terpasang sempurna."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1020, 106, 'drag_drop', 3, ?, 10)`,
		`{
			"prompt": "Kelompokkan bagian tubuh berikut sesuai letaknya!",
			"items": ["Mata (Kepala)", "Kaki (Bawah)", "Tangan (Tengah)"],
			"expected_order": ["Mata (Kepala)", "Tangan (Tengah)", "Kaki (Bawah)"],
			"hint": "Urutkan dari atas tubuh ke bawah tubuh!",
			"explanation": "Bagus sekali! Urutan tubuh dari atas ke bawah."
		}`,
	)

	// 107. Explorers Science L5: Siang & Malam, Matahari & Bulan 🌙
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (107, ?, 'explorers', 5, 1, 'Siang & Malam, Matahari & Bulan 🌙', 'Memahami siklus siang dan malam serta benda langit',
		 '{"intro_text":"Matahari menerangi siang hari, dan Bulan menghiasi langit malam!","icon":"🌙"}', 8, 15)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1021, 107, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Benda langit apa yang terbit di pagi hari dan membuat langit menjadi terang?",
			"options": ["Matahari ☀️", "Bulan 🌙", "Bintang 🌟", "Awan ☁️"],
			"answer": "Matahari ☀️",
			"hint": "Benda langit ini sangat panas dan bersinar terang di siang hari.",
			"explanation": "Tepat! Matahari menerangi bumi di siang hari."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1022, 107, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Urutkan kegiatan sehari-hari dari pagi hingga malam!",
			"items": ["Tidur Malam", "Bangun Pagi", "Bermain Siang"],
			"expected_order": ["Bangun Pagi", "Bermain Siang", "Tidur Malam"],
			"hint": "Dimulai saat matahari terbit!",
			"explanation": "Pintar! Urutannya: Bangun Pagi -> Bermain Siang -> Tidur Malam."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1023, 107, 'matching', 3, ?, 10)`,
		`{
			"prompt": "Jodohkan benda langit dengan suasananya!",
			"pairs": [
				{"left": "Matahari ☀️", "right": "Siang Terang"},
				{"left": "Bulan 🌙", "right": "Malam Gelap"},
				{"left": "Pelangi 🌈", "right": "Setelah Hujan"}
			],
			"hint": "Kapan kita biasanya melihat hal-hal tersebut?",
			"explanation": "Hebat! Pilihanmu benar semua."
		}`,
	)

	// 108. Builders Science L4: Rantai Makanan & Jaring Kehidupan 🦁
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (108, ?, 'builders', 4, 1, 'Rantai Makanan & Jaring Kehidupan 🦁', 'Mengenal produsen, konsumen, dan pengurai dalam ekosistem',
		 '{"intro_text":"Di alam, makhluk hidup saling membutuhkan untuk mendapatkan energi makanan!","icon":"🦁"}', 10, 20)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1024, 108, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Dalam rantai makanan, tumbuhan berperan sebagai apa karena bisa membuat makanan sendiri?",
			"options": ["Produsen", "Konsumen I", "Konsumen II", "Pengurai"],
			"answer": "Produsen",
			"hint": "Tumbuhan memproduksi makanan sendiri melalui fotosintesis.",
			"explanation": "Benar! Tumbuhan adalah Produsen dalam rantai makanan."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1025, 108, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Urutkan rantai makanan di padang rumput dari bawah ke atas!",
			"items": ["Singa (Predator)", "Rumput (Produsen)", "Rusa (Herbivora)"],
			"expected_order": ["Rumput (Produsen)", "Rusa (Herbivora)", "Singa (Predator)"],
			"hint": "Rumput dimakan rusa, rusa dimakan singa!",
			"explanation": "Tepat! Rumput -> Rusa -> Singa."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1026, 108, 'sequencing', 3, ?, 10)`,
		`{
			"prompt": "Susun tingkatan rantai makanan secara berurutan!",
			"steps": ["Produsen (Tumbuhan)", "Konsumen Primer (Herbivora)", "Konsumen Sekunder (Karnivora)", "Pengurai (Jamur/Bakteri)"],
			"expected_order": ["Produsen (Tumbuhan)", "Konsumen Primer (Herbivora)", "Konsumen Sekunder (Karnivora)", "Pengurai (Jamur/Bakteri)"],
			"hint": "Mulai dari tumbuhan penghasil makanan!",
			"explanation": "Luar biasa! Urutan aliran energi dalam rantai makanan sudah sempurna."
		}`,
	)

	// 109. Builders Science L5: Pesawat Sederhana & Gaya 🔧
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (109, ?, 'builders', 5, 1, 'Pesawat Sederhana & Gaya 🔧', 'Belajar tentang tuas, katrol, bidang miring, dan gaya',
		 '{"intro_text":"Pesawat sederhana membantu manusia memindahkan beban berat dengan lebih mudah!","icon":"🔧"}', 10, 20)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1027, 109, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Alat manakah yang memanfaatkan prinsip bidang miring untuk memudahkan memindahkan benda ke tempat tinggi?",
			"options": ["Papan Seluncur / Tangga", "Gunting", "Jungkit-jungkit", "Roda"],
			"answer": "Papan Seluncur / Tangga",
			"hint": "Permukaan miring melandaikan lintasan naik.",
			"explanation": "Benar! Papan miring dan tangga adalah contoh bidang miring."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1028, 109, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Gaya yang menyebabkan buah apel jatuh dari pohon ke tanah adalah gaya ...",
			"answer": "gravitasi",
			"hint": "Gaya tarik bumi.",
			"explanation": "Tepat! Gaya gravitasi bumi menarik semua benda ke bawah."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1029, 109, 'matching', 3, ?, 10)`,
		`{
			"prompt": "Jodohkan pesawat sederhana dengan contoh kegunaannya!",
			"pairs": [
				{"left": "Katrol ⚙️", "right": "Menimba Air Sumur"},
				{"left": "Tuas / Pengungkit 🔨", "right": "Membuka Tutup Botol"},
				{"left": "Roda Berporos 🛞", "right": "Sepeda & Mobil"}
			],
			"hint": "Fikirkan alat sehari-hari yang menggunakan prinsip ini!",
			"explanation": "Luar biasa! Pengetahuan mekanikmu hebat."
		}`,
	)

	// 110. Challengers Science L4: Ilmu Bumi: Lempeng Tektonik & Batuan 🌍
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (110, ?, 'challengers', 4, 1, 'Ilmu Bumi: Lempeng Tektonik & Batuan 🌍', 'Memahami struktur bumi, gempa, dan jenis-jenis batuan',
		 '{"intro_text":"Bumi kita terus bergerak! Lempeng tektonik menyebabkan terbentuknya gunung dan gempa bumi.","icon":"🌍"}', 12, 25)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1030, 110, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Batuan yang terbentuk dari pembekuan magma atau lava gunung berapi disebut?",
			"options": ["Batuan Beku", "Batuan Sedimen", "Batuan Metamorf", "Batuan Kapur"],
			"answer": "Batuan Beku",
			"hint": "Kata puncaknya adalah pembekuan cairan lava panas.",
			"explanation": "Benar! Magma dingin membeku menjadi Batuan Beku (seperti batu apung dan batu obsidian)."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1031, 110, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Alat yang digunakan untuk mengukur dan mencatat kekuatan gempa bumi adalah ...",
			"answer": "seismograf",
			"hint": "Berawalan kata 'seismo'.",
			"explanation": "Tepat! Seismograf mencatat gelombang seismik gempa."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1032, 110, 'drag_drop', 3, ?, 10)`,
		`{
			"prompt": "Urutkan lapisan bumi dari yang PALING LUAR ke yang PALING DALAM!",
			"items": ["Inti Bumi (Core)", "Kerak Bumi (Crust)", "Mantel Bumi (Mantle)"],
			"expected_order": ["Kerak Bumi (Crust)", "Mantel Bumi (Mantle)", "Inti Bumi (Core)"],
			"hint": "Kita tinggal di permukaan paling luar (Kerak Bumi)!",
			"explanation": "Hebat! Urutannya: Kerak Bumi -> Mantel Bumi -> Inti Bumi."
		}`,
	)

	// 111. Challengers Science L5: Metode Ilmiah & Eksperimen 🧪
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (111, ?, 'challengers', 5, 1, 'Metode Ilmiah & Eksperimen 🧪', 'Langkah-langkah penelitian ilmiah dan perumusan hipotesis',
		 '{"intro_text":"Ilmuwan bekerja secara sistematis dengan metode ilmiah untuk menemukan kebenaran!","icon":"🧪"}', 12, 25)`,
		scienceID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1033, 111, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Dugaan awal atau jawaban sementara terhadap rumusan masalah dalam penelitian ilmiah disebut?",
			"options": ["Hipotesis", "Kesimpulan", "Observasi", "Variabel"],
			"answer": "Hipotesis",
			"hint": "Dugaan teruji yang belum dibuktikan eksperimen.",
			"explanation": "Benar! Hipotesis adalah dugaan ilmiah sementara."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1034, 111, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Faktor yang diubah-ubah dalam sebuah eksperimen ilmiah disebut variabel ...",
			"answer": "bebas",
			"hint": "Lawan dari variabel terikat.",
			"explanation": "Tepat! Variabel bebas adalah faktor yang sengaja dimanipulasi oleh peneliti."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1035, 111, 'sequencing', 3, ?, 10)`,
		`{
			"prompt": "Susun langkah-langkah Metode Ilmiah dari awal sampai akhir!",
			"steps": ["Merumuskan Masalah", "Membuat Hipotesis", "Melakukan Eksperimen", "Menganalisis Data & Tarik Kesimpulan"],
			"expected_order": ["Merumuskan Masalah", "Membuat Hipotesis", "Melakukan Eksperimen", "Menganalisis Data & Tarik Kesimpulan"],
			"hint": "Mulai dari pertanyaan atau masalah dulu!",
			"explanation": "Luar biasa! Urutan metode ilmiah sudah sangat ilmiah dan terstruktur."
		}`,
	)

	// =========================================================================
	// CODING LEVEL 4 & 5 (BUILDERS, CHALLENGERS)
	// =========================================================================

	// 112. Builders Coding L4: Fungsi & Blok Perintah Bernama 📦
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (112, ?, 'builders', 4, 1, 'Fungsi & Blok Perintah Bernama 📦', 'Membungkus serangkaian perintah ke dalam satu blok fungsi',
		 '{"intro_text":"Fungsi seperti kotak resep! Sekali dibuat, bisa dipanggil berkali-kali tanpa tulis ulang!","icon":"📦"}', 10, 20)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1036, 112, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Apa keuntungan utama menggunakan Fungsi (Function) dalam pemograman?",
			"options": ["Menghindari penulisan kode berulang", "Membuat komputer berjalan lebih lambat", "Menghapus semua gambar", "Mengubah warna layar"],
			"answer": "Menghindari penulisan kode berulang",
			"hint": "Kode jadi lebih rapi dan hemat tempat.",
			"explanation": "Tepat! Fungsi memungkinkan kode dipakai ulang (reusable)."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1037, 112, 'sequencing', 2, ?, 10)`,
		`{
			"prompt": "Urutkan cara membuat dan menggunakan fungsi 'LompatDuaKali'!",
			"steps": ["Definisikan Fungsi LompatDuaKali", "Isi instruksi: Lompat, Lompat", "Panggil nama fungsi: LompatDuaKali()"],
			"expected_order": ["Definisikan Fungsi LompatDuaKali", "Isi instruksi: Lompat, Lompat", "Panggil nama fungsi: LompatDuaKali()"],
			"hint": "Buat resepnya dulu baru panggil resepnya!",
			"explanation": "Hebat! Kamu paham alur pembuatan dan pemanggilan fungsi."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1038, 112, 'block_code', 3, ?, 10)`,
		`{
			"prompt": "Gunakan blok fungsi untuk menggerakkan karakter mencapai sasaran!",
			"available_blocks": ["Maju", "Lompat", "Ulangi 2x"],
			"target_sequence": ["Maju", "Lompat", "Maju", "Lompat"],
			"hint": "Kelompokkan 'Maju, Lompat' lalu ulangi 2x!",
			"explanation": "Luar biasa! Blok kode fungsi berhasil diselesaikan."
		}`,
	)

	// 113. Builders Coding L5: Debugging: Temukan Kesalahan! 🐛
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (113, ?, 'builders', 5, 1, 'Debugging: Temukan Kesalahan! 🐛', 'Menganalisis kode yang rusak dan memperbaikinya',
		 '{"intro_text":"Programmer sejati pandai berburu bug! Mari temukan dan perbaiki kesalahan dalam kode.","icon":"🐛"}', 10, 20)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1039, 113, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Kode untuk mengambil air: [Jalan, Jalan, AmbilAir]. Tapi robot berhenti SEBELUM air. Di mana bug-nya?",
			"options": ["Kurang langkah 'Jalan'", "Kelebihan 'AmbilAir'", "Tidak ada bug", "Arah robot salah"],
			"answer": "Kurang langkah 'Jalan'",
			"hint": "Hitung jarak kotak menuju ember air!",
			"explanation": "Benar! Robot butuh 3 langkah 'Jalan', bukan 2."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1040, 113, 'sequencing', 2, ?, 10)`,
		`{
			"prompt": "Susun ulang kode yang berantakan berikut agar robot bisa menyalakan lampu!",
			"steps": ["Maju ke saklar", "Tekan Tombol Lampu", "Pastikan Saklar Aktif"],
			"expected_order": ["Maju ke saklar", "Tekan Tombol Lampu", "Pastikan Saklar Aktif"],
			"hint": "Robot harus mendekat dulu sebelum menekan tombol!",
			"explanation": "Hebat! Bug urutan berhasil dibetulkan."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1041, 113, 'fill_blank', 3, ?, 10)`,
		`{
			"prompt": "Istilah dalam pemrograman untuk mencari dan membasmi kesalahan kode adalah ...",
			"answer": "debugging",
			"hint": "Berasal dari kata 'bug' (serangga).",
			"explanation": "Tepat! Debugging adalah proses perbaikan bug."
		}`,
	)

	// 114. Challengers Coding L4: Fungsi Sederhana, Input & Output ⚙️
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (114, ?, 'challengers', 4, 1, 'Fungsi Sederhana, Input & Output ⚙️', 'Memahami parameter input dan nilai kembalian (return) fungsi',
		 '{"intro_text":"Fungsi bisa menerima data masuk (Input/Parameter) dan mengembalikan hasil (Output/Return)!","icon":"⚙️"}', 12, 25)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1042, 114, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Diberikan fungsi: hitungLuas(p, l) { return p * l }. Jika p=5 dan l=4, berapa output return-nya?",
			"options": ["20", "9", "1", "54"],
			"answer": "20",
			"hint": "Fungsi mengalikan p dengan l (5 × 4).",
			"explanation": "Benar! 5 × 4 = 20."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1043, 114, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Fungsi sapa(nama) menampilkan 'Halo ' + nama. Jika dimasukkan input 'Budi', kata apakah yang muncul?",
			"answer": "Halo Budi",
			"hint": "Gabungkan kata 'Halo ' dengan input 'Budi'.",
			"explanation": "Tepat! Output-nya adalah 'Halo Budi'."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1044, 114, 'sequencing', 3, ?, 10)`,
		`{
			"prompt": "Urutkan alur kerja fungsi pemrosesan data!",
			"steps": ["Terima Parameter Input", "Proses Logika di Dalam Fungsi", "Kembalikan Nilai Return Output"],
			"expected_order": ["Terima Parameter Input", "Proses Logika di Dalam Fungsi", "Kembalikan Nilai Return Output"],
			"hint": "Input -> Proses -> Output!",
			"explanation": "Luar biasa! Alur I/O fungsi sudah kamu kuasai."
		}`,
	)

	// 115. Challengers Coding L5: Proyek Mini: Kalkulator Sederhana 🧮
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO lessons (id, category_id, age_group, level, sort_order, title, description, content_json, estimated_minutes, xp_reward)
		 VALUES (115, ?, 'challengers', 5, 1, 'Proyek Mini: Kalkulator Sederhana 🧮', 'Membangun logika program kalkulator dengan fungsi dan pengandaian',
		 '{"intro_text":"Mari gabungkan variabel, fungsi, dan kondisi If-Else untuk membuat aplikasi kalkulator!","icon":"🧮"}', 15, 30)`,
		codingID,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1045, 115, 'multiple_choice', 1, ?, 10)`,
		`{
			"prompt": "Dalam program kalkulator, jika operator yang dipilih adalah '+', fungsi mana yang harus dipanggil?",
			"options": ["tambah(a, b)", "kurang(a, b)", "kali(a, b)", "bagi(a, b)"],
			"answer": "tambah(a, b)",
			"hint": "Simbol '+' menandakan operasi penjumlahan.",
			"explanation": "Tepat! Operator '+' mengarah pada fungsi tambah."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1046, 115, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Dalam kalkulator, pembagian dengan angka 0 akan menyebabkan error. Kondisi penanganan ini disebut periksa ... Zero.",
			"answer": "Divide by",
			"hint": "Istilah bahasa Inggris: Divide by Zero.",
			"explanation": "Benar! Divide by Zero (pembagian nol) harus dicegah dengan If-Else."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (1047, 115, 'block_code', 3, ?, 10)`,
		`{
			"prompt": "Susun blok algoritma kalkulator sederhana!",
			"available_blocks": ["Input Angka A dan B", "Pilih Operator (+/-/*)", "Hitung & Tampilkan Hasil"],
			"target_sequence": ["Input Angka A dan B", "Pilih Operator (+/-/*)", "Hitung & Tampilkan Hasil"],
			"hint": "Minta angka dari pengguna, pilih operasi, lalu cetak hasilnya!",
			"explanation": "Selamat! Kamu berhasil merancang logika program kalkulator utuh!"
		}`,
	)

	fmt.Println("Seeded Level 4 and Level 5 curriculum for Math, Science, and Coding!")
	return nil
}

// SeedExtraActivities adds extra activities to existing L1-L3 lessons to enrich practice content.
func SeedExtraActivities(db *sql.DB) error {
	// Science Builders Extra (Lessons 26, 27, 28, 29)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2000, 26, 'matching', 2, ?, 10)`,
		`{
			"prompt": "Jodohkan hewan dengan kelompok jenisnya!",
			"pairs": [
				{"left": "Ikan Mas 🐟", "right": "Hewan Air"},
				{"left": "Burung Elang 🦅", "right": "Hewan Udara"},
				{"left": "Kelinci 🐇", "right": "Hewan Darat"}
			],
			"hint": "Perhatikan tempat bergerak hewan tersebut!",
			"explanation": "Hebat! Semua tempat tinggal hewan cocok."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2001, 27, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Kelompokkan benda berikut berdasarkan wujudnya!",
			"items": ["Batu (Padat)", "Air Es (Cair)", "Asap (Gas)"],
			"expected_order": ["Batu (Padat)", "Air Es (Cair)", "Asap (Gas)"],
			"hint": "Batu keras, air mengalir, asap melayang!",
			"explanation": "Bagus! Padat -> Cair -> Gas."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2002, 28, 'sequencing', 2, ?, 10)`,
		`{
			"prompt": "Urutkan proses Daur Air secara alami!",
			"steps": ["Air Laut Menguap (Evaporasi)", "Terbentuk Awan (Kondensasi)", "Turun Hujan (Presipitasi)"],
			"expected_order": ["Air Laut Menguap (Evaporasi)", "Terbentuk Awan (Kondensasi)", "Turun Hujan (Presipitasi)"],
			"hint": "Panas matahari membuat air naik ke atas terlebih dahulu!",
			"explanation": "Luar biasa! Urutan daur air sudah benar."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2003, 29, 'multiple_choice', 2, ?, 10)`,
		`{
			"prompt": "Planet yang paling dekat dengan Matahari dalam Tata Surya adalah?",
			"options": ["Merkurius", "Venus", "Bumi", "Mars"],
			"answer": "Merkurius",
			"hint": "Planet kecil pertama terdekat dengan pusat tata surya.",
			"explanation": "Tepat! Merkurius adalah planet pertama dari Matahari."
		}`,
	)

	// Language Builders Extra (Lessons 30, 31, 32)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2004, 30, 'matching', 2, ?, 10)`,
		`{
			"prompt": "Pasangkan kata lawan kata (antonim) berikut!",
			"pairs": [
				{"left": "Besar", "right": "Kecil"},
				{"left": "Tinggi", "right": "Rendah"},
				{"left": "Terang", "right": "Gelap"}
			],
			"hint": "Cari kata yang artinya berlawanan!",
			"explanation": "Pintar! Pasangan antonim kamu sudah tepat."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2005, 31, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Dalam kalimat 'Adik MEMBACA buku', kata MEMBACA termasuk jenis kata ...",
			"answer": "kerja",
			"hint": "Kata yang menunjukkan suatu tindakan/aktivitas.",
			"explanation": "Benar! Membaca adalah kata kerja (verba)."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2006, 32, 'multiple_choice', 2, ?, 10)`,
		`{
			"prompt": "Teks: 'Budi rajin menabung di bank setiap minggu.' Apa pesan utama dari kalimat tersebut?",
			"options": ["Budi hemat dan suka menabung", "Budi suka berbelanja", "Budi bekerja di bank", "Budi tidak punya uang"],
			"answer": "Budi hemat dan suka menabung",
			"hint": "Perhatikan kata 'rajin menabung'.",
			"explanation": "Tepat! Kalimat tersebut menceritakan kebiasaan hemat Budi."
		}`,
	)

	// Logic Builders Extra (Lessons 33, 34, 35)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2007, 33, 'multiple_choice', 2, ?, 10)`,
		`{
			"prompt": "Barisan angka: 2, 4, 6, 8, ... Angka berikutnya adalah?",
			"options": ["10", "9", "12", "11"],
			"answer": "10",
			"hint": "Pola ini bertambah +2 setiap langkahnya.",
			"explanation": "Benar! 8 + 2 = 10 (barisan bilangan genap)."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2008, 34, 'drag_drop', 2, ?, 10)`,
		`{
			"prompt": "Urutkan petunjuk arah labirin!",
			"items": ["Belok Kanan", "Lurus 3 Langkah", "Sampai Tujuan"],
			"expected_order": ["Lurus 3 Langkah", "Belok Kanan", "Sampai Tujuan"],
			"hint": "Jalan lurus dulu sebelum berbelok!",
			"explanation": "Tepat! Rute labirin berhasil dilewati."
		}`,
	)
	_, _ = db.Exec(
		`INSERT OR IGNORE INTO activities (id, lesson_id, type, sort_order, question_json, max_score)
		 VALUES (2009, 35, 'fill_blank', 2, ?, 10)`,
		`{
			"prompt": "Pola geometri: Segitiga, Persegi, Segilima, ... Bangun datar selanjutnya punya berapa sisi?",
			"answer": "6",
			"hint": "Jumlah sisi bertambah 1: 3, 4, 5, ...",
			"explanation": "Hebat! Segienam punya 6 sisi."
		}`,
	)

	fmt.Println("Seeded extra activities for existing lessons!")
	return nil
}

