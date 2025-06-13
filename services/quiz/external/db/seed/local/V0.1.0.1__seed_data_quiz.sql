-- AI generated
-- Quiz
INSERT INTO public.quiz (quiz_id, content)
VALUES
    ('quiz_id_01', 'Địa danh'),
    ('quiz_id_02', 'Toán học'),
    ('quiz_id_03', 'Kiến thức chung');

-- Question
INSERT INTO question (question_id, content, correct_answer, answers)
VALUES
-- Địa danh (q1–q4)
('q1', 'Thủ đô của Việt Nam là gì?', 'A',
 '[{"title": "A", "content": "Hà Nội"}, {"title": "B", "content": "TP. Hồ Chí Minh"}, {"title": "C", "content": "Đà Nẵng"}]'::jsonb),
('q2', 'Sapa thuộc tỉnh nào?', 'B',
 '[{"title": "A", "content": "Hà Giang"}, {"title": "B", "content": "Lào Cai"}, {"title": "C", "content": "Yên Bái"}]'::jsonb),
('q3', 'Vịnh Hạ Long nằm ở tỉnh nào?', 'C',
 '[{"title": "A", "content": "Thanh Hóa"}, {"title": "B", "content": "Nghệ An"}, {"title": "C", "content": "Quảng Ninh"}]'::jsonb),
('q4', 'Cố đô Huế thuộc miền nào?', 'B',
 '[{"title": "A", "content": "Miền Bắc"}, {"title": "B", "content": "Miền Trung"}, {"title": "C", "content": "Miền Nam"}]'::jsonb),

-- Toán học (q5–q7)
('q5', '2 + 3 bằng mấy?', 'C',
 '[{"title": "A", "content": "4"}, {"title": "B", "content": "6"}, {"title": "C", "content": "5"}]'::jsonb),
('q6', 'Bình phương của 4 là?', 'A',
 '[{"title": "A", "content": "16"}, {"title": "B", "content": "12"}, {"title": "C", "content": "14"}]'::jsonb),
('q7', 'Số nguyên tố đầu tiên là?', 'A',
 '[{"title": "A", "content": "2"}, {"title": "B", "content": "3"}, {"title": "C", "content": "1"}]'::jsonb),

-- Kiến thức chung riêng biệt (q8–q10)
('q8', 'Loài vật nào được mệnh danh là vua sơn lâm?', 'C',
 '[{"title": "A", "content": "Voi"}, {"title": "B", "content": "Báo"}, {"title": "C", "content": "Sư tử"}]'::jsonb),
('q9', 'Ngôn ngữ lập trình dùng phổ biến cho iOS?', 'A',
 '[{"title": "A", "content": "Swift"}, {"title": "B", "content": "Java"}, {"title": "C", "content": "Python"}]'::jsonb),
('q10', 'Trái Đất quay quanh gì?', 'B',
 '[{"title": "A", "content": "Mặt trăng"}, {"title": "B", "content": "Mặt trời"}, {"title": "C", "content": "Sao Hỏa"}]'::jsonb);

-- Gán các câu hỏi cho quiz
INSERT INTO quiz_question (quiz_id, question_id)
VALUES
-- Quiz Địa danh
('quiz_id_01', 'q1'),
('quiz_id_01', 'q2'),
('quiz_id_01', 'q3'),
('quiz_id_01', 'q4'),

-- Quiz Toán học
('quiz_id_02', 'q5'),
('quiz_id_02', 'q6'),
('quiz_id_02', 'q7'),

-- Quiz Kiến thức chung (bao gồm các câu đã có ở các quiz khác)
('quiz_id_03', 'q1'),  -- trùng với Địa danh
('quiz_id_03', 'q5'),  -- trùng với Toán học
('quiz_id_03', 'q8'),
('quiz_id_03', 'q9'),
('quiz_id_03', 'q10');
