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

	fmt.Println("Seeded Math, Coding, Toddler, Science, Language, and Logic curriculum content")
	return nil
}
