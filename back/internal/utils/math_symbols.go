package utils

import (
	"regexp"
	"strings"

	"docx-converter-demo/internal/types"
)

// ----------------- TeX -> Unicode -----------------
// สัญลักษณ์
var symPairs = [][2]string{
	{`\geq`, "≥"}, {`\ge`, "≥"},
	{`\leq`, "≤"}, {`\le`, "≤"},
	{`\neq`, "≠"},
	{`\pm`, "±"},
	{`\times`, "×"},
	{`\cdot`, "·"},
	{`\infty`, "∞"},
	{`\approx`, "≈"},
	{`\propto`, "∝"},
	{`\sum`, "∑"},
	{`\prod`, "∏"},
	{`\int`, "∫"},
}

// กรีก
var greekPairs = [][2]string{
	{`\Gamma`, `Γ`}, {`\Delta`, `Δ`}, {`\Theta`, `Θ`}, {`\Lambda`, `Λ`},
	{`\Xi`, `Ξ`}, {`\Pi`, `Π`}, {`\Sigma`, `Σ`}, {`\Upsilon`, `Υ`},
	{`\Phi`, `Φ`}, {`\Psi`, `Ψ`}, {`\Omega`, `Ω`},

	{`\alpha`, `α`}, {`\beta`, `β`}, {`\gamma`, `γ`}, {`\delta`, `δ`},
	{`\epsilon`, `ε`}, {`\zeta`, `ζ`}, {`\eta`, `η`}, {`\theta`, `θ`},
	{`\iota`, `ι`}, {`\kappa`, `κ`}, {`\lambda`, `λ`}, {`\mu`, `μ`},
	{`\nu`, `ν`}, {`\xi`, `ξ`}, {`\pi`, `π`}, {`\rho`, `ρ`},
	{`\sigma`, `σ`}, {`\tau`, `τ`}, {`\upsilon`, `υ`}, {`\phi`, `φ`},
	{`\chi`, `χ`}, {`\psi`, `ψ`}, {`\omega`, `ω`},
}

// combining overline U+0305
const overline = "\u0305"

var (
	reSpace1    = regexp.MustCompile(`\\,`)                        // "\," → space
	reSpace2    = regexp.MustCompile(`\\\s`)                       // "\ " → space
	reBarLetter = regexp.MustCompile(`\\bar\s*{\s*([A-Za-z])\s*}`) // \bar{x} (ตัวเดียว)
	// \overline{a} หรือ \overline{\sigma} (รองรับตัวอักษรเดี่ยวหรือตัวสั่งกรีกหนึ่งตัว)
	reOverAny    = regexp.MustCompile(`\\overline\s*{\s*([A-Za-z]|\\[A-Za-z]+)\s*}`)
	reBfOne      = regexp.MustCompile(`\\mathbf\s*{\s*([A-Za-z])\s*}`)       // \mathbf{s} -> s
	reBfMulti    = regexp.MustCompile(`\\mathbf\s*{\s*([A-Za-z0-9\s]+)\s*}`) // \mathbf{ABC 123} -> ABC 123
	reMultiSpace = regexp.MustCompile(`\s+`)
)

// แทนที่ตามลำดับที่กำหนด (ยาวก่อนสั้น)
func replaceAllOrdered(s string, pairs [][2]string) string {
	for _, p := range pairs {
		s = strings.ReplaceAll(s, p[0], p[1])
	}
	return s
}

// texToUnicode: รับเฉพาะ “เนื้อในสูตร” (ไม่รวม $ หรือ $$) แล้วคืนยูนิโค้ด
func texToUnicode(t string) string {
	// เว้นวรรคที่เจอบ่อย
	t = reSpace1.ReplaceAllString(t, " ")
	t = reSpace2.ReplaceAllString(t, " ")

	// สัญลักษณ์พื้นฐาน + ตัวกรีก (ordered)
	t = replaceAllOrdered(t, symPairs)
	t = replaceAllOrdered(t, greekPairs)

	// \overline{…} — รองรับ {a} หรือ {\sigma}
	t = reOverAny.ReplaceAllStringFunc(t, func(m string) string {
		sub := reOverAny.FindStringSubmatch(m)
		if len(sub) != 2 {
			return m
		}
		x := sub[1]
		// ถ้าเป็นคำสั่งกรีก ให้แปลงก่อน
		if strings.HasPrefix(x, `\`) {
			x = replaceAllOrdered(x, greekPairs)
		}
		// ใส่ combining overline เฉพาะกรณีตัวอักษรเดี่ยว/สัญลักษณ์เดี่ยว
		if len([]rune(x)) == 1 {
			return x + overline
		}
		// ถ้ายาวกว่าหนึ่งตัว ปล่อยเป็น plain (เลี่ยงการวาง overline ทับหลายตัว)
		return x
	})

	// \bar{…} — รองรับเฉพาะตัวอักษรเดี่ยว
	t = reBarLetter.ReplaceAllStringFunc(t, func(m string) string {
		sub := reBarLetter.FindStringSubmatch(m)
		if len(sub) == 2 {
			return sub[1] + overline
		}
		return m
	})

	// \mathbf{…} — คงข้อความด้านในเป็น plain (frontend จะจัด style เองถ้าต้องการ)
	t = reBfOne.ReplaceAllString(t, `$1`)
	t = reBfMulti.ReplaceAllString(t, `$1`)

	// เก็บกวาดเว้นวรรค
	t = reMultiSpace.ReplaceAllString(strings.TrimSpace(t), " ")
	return t
}

// ----------------- จับสูตรจากสตริงทั้งก้อน -----------------
// รองรับทั้ง plain text และ HTML fragment (tinyMCE, pandoc)
var (
	// <span class="math …"> … </span> (non-greedy, case-insensitive, dotall)
	reSpanMath = regexp.MustCompile(`(?is)<span[^>]*class="[^"]*math[^"]*"[^>]*>(.*?)</span>`)

	// $$…$$ ก่อน $…$ (กันซ้อน)
	reDisplay = regexp.MustCompile(`\$\$([^$]+?)\$\$`)
	reInline  = regexp.MustCompile(`\$(.+?)\$`)
)

// ConvertMathInString: แปลงสูตรที่ฝังอยู่ในสตริงเป็นยูนิโค้ด โดยคง HTML อื่น ๆ ไว้
func ConvertMathInString(s string) string {
	if s == "" {
		return s
	}

	// 1) <span class="math …">…</span>
	s = reSpanMath.ReplaceAllStringFunc(s, func(m string) string {
		sub := reSpanMath.FindStringSubmatch(m)
		if len(sub) == 2 {
			inner := strings.Trim(sub[1], " \t\r\n$")
			return texToUnicode(inner)
		}
		return m
	})

	// 2) $$ … $$
	s = reDisplay.ReplaceAllStringFunc(s, func(m string) string {
		sub := reDisplay.FindStringSubmatch(m)
		if len(sub) == 2 {
			return texToUnicode(strings.TrimSpace(sub[1]))
		}
		return m
	})

	// 3) $ … $
	s = reInline.ReplaceAllStringFunc(s, func(m string) string {
		sub := reInline.FindStringSubmatch(m)
		if len(sub) == 2 {
			return texToUnicode(strings.TrimSpace(sub[1]))
		}
		return m
	})

	return s
}

// ConvertBlocks: map ทั้ง []string
func ConvertBlocks(items []string) []string {
	if len(items) == 0 {
		return items
	}
	out := make([]string, len(items))
	for i, v := range items {
		out[i] = ConvertMathInString(v)
	}
	return out
}

// MathifyOutput: โพสต์โปรเซสทุกฟิลด์ใน Output ที่เป็น []string ให้แปลงสูตรเป็นยูนิโค้ด
func MathifyOutput(o *types.Output) {
	o.Rationale = ConvertBlocks(o.Rationale)
	o.Objective = ConvertBlocks(o.Objective)
	o.Content = ConvertBlocks(o.Content)
	o.Evaluation = ConvertBlocks(o.Evaluation)
	o.Keywords = ConvertBlocks(o.Keywords)
	o.Overview = ConvertBlocks(o.Overview)
}
