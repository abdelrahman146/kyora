package country

import (
	"sort"
	"sync"
)

type Country struct {
	Name           string `json:"name"`
	NameAr         string `json:"nameAr"`
	Code           string `json:"code"`
	IsoCode        string `json:"iso_code"`
	Flag           string `json:"flag"`
	PhonePrefix    string `json:"phonePrefix"`
	CurrencyCode   string `json:"currencyCode"`
	CurrencyLabel  string `json:"currencyLabel"`
	CurrencySymbol string `json:"currencySymbol"`
}

var countries = []Country{
	{Name: "United States", NameAr: "الولايات المتحدة", Code: "US", IsoCode: "USA", Flag: "🇺🇸", PhonePrefix: "+1", CurrencyCode: "USD", CurrencyLabel: "USD — $", CurrencySymbol: "$"},
	{Name: "United Kingdom", NameAr: "المملكة المتحدة", Code: "GB", IsoCode: "GBR", Flag: "🇬🇧", PhonePrefix: "+44", CurrencyCode: "GBP", CurrencyLabel: "GBP — £", CurrencySymbol: "£"},
	{Name: "Japan", NameAr: "اليابان", Code: "JP", IsoCode: "JPN", Flag: "🇯🇵", PhonePrefix: "+81", CurrencyCode: "JPY", CurrencyLabel: "JPY — ¥", CurrencySymbol: "¥"},
	{Name: "China", NameAr: "الصين", Code: "CN", IsoCode: "CHN", Flag: "🇨🇳", PhonePrefix: "+86", CurrencyCode: "CNY", CurrencyLabel: "CNY — ¥", CurrencySymbol: "¥"},
	{Name: "Australia", NameAr: "أستراليا", Code: "AU", IsoCode: "AUS", Flag: "🇦🇺", PhonePrefix: "+61", CurrencyCode: "AUD", CurrencyLabel: "AUD — A$", CurrencySymbol: "A$"},
	{Name: "Canada", NameAr: "كندا", Code: "CA", IsoCode: "CAN", Flag: "🇨🇦", PhonePrefix: "+1", CurrencyCode: "CAD", CurrencyLabel: "CAD — C$", CurrencySymbol: "C$"},
	{Name: "India", NameAr: "الهند", Code: "IN", IsoCode: "IND", Flag: "🇮🇳", PhonePrefix: "+91", CurrencyCode: "INR", CurrencyLabel: "INR — ₹", CurrencySymbol: "₹"},
	{Name: "United Arab Emirates", NameAr: "الإمارات العربية المتحدة", Code: "AE", IsoCode: "ARE", Flag: "🇦🇪", PhonePrefix: "+971", CurrencyCode: "AED", CurrencyLabel: "AED — د.إ", CurrencySymbol: "د.إ"},
	{Name: "Saudi Arabia", NameAr: "المملكة العربية السعودية", Code: "SA", IsoCode: "SAU", Flag: "🇸🇦", PhonePrefix: "+966", CurrencyCode: "SAR", CurrencyLabel: "SAR — ر.س", CurrencySymbol: "ر.س"},
	{Name: "Egypt", NameAr: "مصر", Code: "EG", IsoCode: "EGY", Flag: "🇪🇬", PhonePrefix: "+20", CurrencyCode: "EGP", CurrencyLabel: "EGP — E£", CurrencySymbol: "E£"},
	{Name: "Kuwait", NameAr: "الكويت", Code: "KW", IsoCode: "KWT", Flag: "🇰🇼", PhonePrefix: "+965", CurrencyCode: "KWD", CurrencyLabel: "KWD — د.ك", CurrencySymbol: "د.ك"},
	{Name: "Qatar", NameAr: "قطر", Code: "QA", IsoCode: "QAT", Flag: "🇶🇦", PhonePrefix: "+974", CurrencyCode: "QAR", CurrencyLabel: "QAR — ر.ق", CurrencySymbol: "ر.ق"},
	{Name: "Oman", NameAr: "سلطنة عمان", Code: "OM", IsoCode: "OMN", Flag: "🇴🇲", PhonePrefix: "+968", CurrencyCode: "OMR", CurrencyLabel: "OMR — ر.ع.", CurrencySymbol: "ر.ع."},
	{Name: "Algeria", NameAr: "الجزائر", Code: "DZ", IsoCode: "DZA", Flag: "🇩🇿", PhonePrefix: "+213", CurrencyCode: "DZD", CurrencyLabel: "DZD — د.ج", CurrencySymbol: "د.ج"},
	{Name: "Morocco", NameAr: "المغرب", Code: "MA", IsoCode: "MAR", Flag: "🇲🇦", PhonePrefix: "+212", CurrencyCode: "MAD", CurrencyLabel: "MAD — د.م.", CurrencySymbol: "د.م."},
	{Name: "Tunisia", NameAr: "تونس", Code: "TN", IsoCode: "TUN", Flag: "🇹🇳", PhonePrefix: "+216", CurrencyCode: "TND", CurrencyLabel: "TND — د.ت.", CurrencySymbol: "د.ت."},
	{Name: "Jordan", NameAr: "الأردن", Code: "JO", IsoCode: "JOR", Flag: "🇯🇴", PhonePrefix: "+962", CurrencyCode: "JOD", CurrencyLabel: "JOD — د.ا", CurrencySymbol: "د.ا"},
	{Name: "Bahrain", NameAr: "البحرين", Code: "BH", IsoCode: "BHR", Flag: "🇧🇭", PhonePrefix: "+973", CurrencyCode: "BHD", CurrencyLabel: "BHD — د.ب", CurrencySymbol: "د.ب"},
	{Name: "Libya", NameAr: "ليبيا", Code: "LY", IsoCode: "LBY", Flag: "🇱🇾", PhonePrefix: "+218", CurrencyCode: "LYD", CurrencyLabel: "LYD — ل.د", CurrencySymbol: "ل.د"},
	{Name: "Sudan", NameAr: "السودان", Code: "SD", IsoCode: "SDN", Flag: "🇸🇩", PhonePrefix: "+249", CurrencyCode: "SDG", CurrencyLabel: "SDG — ج.س.", CurrencySymbol: "ج.س."},
	{Name: "Yemen", NameAr: "اليمن", Code: "YE", IsoCode: "YEM", Flag: "🇾🇪", PhonePrefix: "+967", CurrencyCode: "YER", CurrencyLabel: "YER — ﷼", CurrencySymbol: "﷼"},
	{Name: "Syria", NameAr: "سوريا", Code: "SY", IsoCode: "SYR", Flag: "🇸🇾", PhonePrefix: "+963", CurrencyCode: "SYP", CurrencyLabel: "SYP — £S", CurrencySymbol: "£S"},
	{Name: "Iraq", NameAr: "العراق", Code: "IQ", IsoCode: "IRQ", Flag: "🇮🇶", PhonePrefix: "+964", CurrencyCode: "IQD", CurrencyLabel: "IQD — ع.د", CurrencySymbol: "ع.د"},
	{Name: "Palestine", NameAr: "فلسطين", Code: "PS", IsoCode: "PSE", Flag: "🇵🇸", PhonePrefix: "+970", CurrencyCode: "ILS", CurrencyLabel: "ILS — ₪", CurrencySymbol: "₪"},
	{Name: "Lebanon", NameAr: "لبنان", Code: "LB", IsoCode: "LBN", Flag: "🇱🇧", PhonePrefix: "+961", CurrencyCode: "LBP", CurrencyLabel: "LBP — ل.ل", CurrencySymbol: "ل.ل"},
	{Name: "Mauritania", NameAr: "موريتانيا", Code: "MR", IsoCode: "MRT", Flag: "🇲🇷", PhonePrefix: "+222", CurrencyCode: "MRU", CurrencyLabel: "MRU — UM", CurrencySymbol: "UM"},
	{Name: "Turkey", NameAr: "تركيا", Code: "TR", IsoCode: "TUR", Flag: "🇹🇷", PhonePrefix: "+90", CurrencyCode: "TRY", CurrencyLabel: "TRY — ₺", CurrencySymbol: "₺"},
	{Name: "Iran", NameAr: "إيران", Code: "IR", IsoCode: "IRN", Flag: "🇮🇷", PhonePrefix: "+98", CurrencyCode: "IRR", CurrencyLabel: "IRR — ﷼", CurrencySymbol: "﷼"},
	{Name: "South Korea", NameAr: "كوريا الجنوبية", Code: "KR", IsoCode: "KOR", Flag: "🇰🇷", PhonePrefix: "+82", CurrencyCode: "KRW", CurrencyLabel: "KRW — ₩", CurrencySymbol: "₩"},
	{Name: "Singapore", NameAr: "سنغافورة", Code: "SG", IsoCode: "SGP", Flag: "🇸🇬", PhonePrefix: "+65", CurrencyCode: "SGD", CurrencyLabel: "SGD — S$", CurrencySymbol: "S$"},
	{Name: "Hong Kong", NameAr: "هونغ كونغ", Code: "HK", IsoCode: "HKG", Flag: "🇭🇰", PhonePrefix: "+852", CurrencyCode: "HKD", CurrencyLabel: "HKD — HK$", CurrencySymbol: "HK$"},
	{Name: "Thailand", NameAr: "تايلاند", Code: "TH", IsoCode: "THA", Flag: "🇹🇭", PhonePrefix: "+66", CurrencyCode: "THB", CurrencyLabel: "THB — ฿", CurrencySymbol: "฿"},
	{Name: "Malaysia", NameAr: "ماليزيا", Code: "MY", IsoCode: "MYS", Flag: "🇲🇾", PhonePrefix: "+60", CurrencyCode: "MYR", CurrencyLabel: "MYR — RM", CurrencySymbol: "RM"},
	{Name: "Indonesia", NameAr: "إندونيسيا", Code: "ID", IsoCode: "IDN", Flag: "🇮🇩", PhonePrefix: "+62", CurrencyCode: "IDR", CurrencyLabel: "IDR — Rp", CurrencySymbol: "Rp"},
	{Name: "Philippines", NameAr: "الفلبين", Code: "PH", IsoCode: "PHL", Flag: "🇵🇭", PhonePrefix: "+63", CurrencyCode: "PHP", CurrencyLabel: "PHP — ₱", CurrencySymbol: "₱"},
	{Name: "Vietnam", NameAr: "فيتنام", Code: "VN", IsoCode: "VNM", Flag: "🇻🇳", PhonePrefix: "+84", CurrencyCode: "VND", CurrencyLabel: "VND — ₫", CurrencySymbol: "₫"},
	{Name: "Pakistan", NameAr: "باكستان", Code: "PK", IsoCode: "PAK", Flag: "🇵🇰", PhonePrefix: "+92", CurrencyCode: "PKR", CurrencyLabel: "PKR — ₨", CurrencySymbol: "₨"},
	{Name: "Bangladesh", NameAr: "بنغلاديش", Code: "BD", IsoCode: "BGD", Flag: "🇧🇩", PhonePrefix: "+880", CurrencyCode: "BDT", CurrencyLabel: "BDT — ৳", CurrencySymbol: "৳"},
	{Name: "Sri Lanka", NameAr: "سريلانكا", Code: "LK", IsoCode: "LKA", Flag: "🇱🇰", PhonePrefix: "+94", CurrencyCode: "LKR", CurrencyLabel: "LKR — Rs", CurrencySymbol: "Rs"},
	{Name: "Myanmar", NameAr: "ميانمار", Code: "MM", IsoCode: "MMR", Flag: "🇲🇲", PhonePrefix: "+95", CurrencyCode: "MMK", CurrencyLabel: "MMK — Ks", CurrencySymbol: "Ks"},
	{Name: "Cambodia", NameAr: "كمبوديا", Code: "KH", IsoCode: "KHM", Flag: "🇰🇭", PhonePrefix: "+855", CurrencyCode: "KHR", CurrencyLabel: "KHR — ៛", CurrencySymbol: "៛"},
	{Name: "Laos", NameAr: "لاوس", Code: "LA", IsoCode: "LAO", Flag: "🇱🇦", PhonePrefix: "+856", CurrencyCode: "LAK", CurrencyLabel: "LAK — ₭", CurrencySymbol: "₭"},
	{Name: "Mongolia", NameAr: "منغوليا", Code: "MN", IsoCode: "MNG", Flag: "🇲🇳", PhonePrefix: "+976", CurrencyCode: "MNT", CurrencyLabel: "MNT — ₮", CurrencySymbol: "₮"},
	{Name: "Kazakhstan", NameAr: "كازاخستان", Code: "KZ", IsoCode: "KAZ", Flag: "🇰🇿", PhonePrefix: "+7", CurrencyCode: "KZT", CurrencyLabel: "KZT — ₸", CurrencySymbol: "₸"},
	{Name: "Uzbekistan", NameAr: "أوزبكستان", Code: "UZ", IsoCode: "UZB", Flag: "🇺🇿", PhonePrefix: "+998", CurrencyCode: "UZS", CurrencyLabel: "UZS — so'm", CurrencySymbol: "so'm"},
	{Name: "Tajikistan", NameAr: "طاجيكستان", Code: "TJ", IsoCode: "TJK", Flag: "🇹🇯", PhonePrefix: "+992", CurrencyCode: "TJS", CurrencyLabel: "TJS — SM", CurrencySymbol: "SM"},
	{Name: "Kyrgyzstan", NameAr: "قرغيزستان", Code: "KG", IsoCode: "KGZ", Flag: "🇰🇬", PhonePrefix: "+996", CurrencyCode: "KGS", CurrencyLabel: "KGS — сом", CurrencySymbol: "сом"},
	{Name: "Afghanistan", NameAr: "أفغانستان", Code: "AF", IsoCode: "AFG", Flag: "🇦🇫", PhonePrefix: "+93", CurrencyCode: "AFN", CurrencyLabel: "AFN — ؋", CurrencySymbol: "؋"},
	{Name: "Nepal", NameAr: "نيبال", Code: "NP", IsoCode: "NPL", Flag: "🇳🇵", PhonePrefix: "+977", CurrencyCode: "NPR", CurrencyLabel: "NPR — ₨", CurrencySymbol: "₨"},
	{Name: "Germany", NameAr: "ألمانيا", Code: "DE", IsoCode: "DEU", Flag: "🇩🇪", PhonePrefix: "+49", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "France", NameAr: "فرنسا", Code: "FR", IsoCode: "FRA", Flag: "🇫🇷", PhonePrefix: "+33", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Italy", NameAr: "إيطاليا", Code: "IT", IsoCode: "ITA", Flag: "🇮🇹", PhonePrefix: "+39", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Spain", NameAr: "إسبانيا", Code: "ES", IsoCode: "ESP", Flag: "🇪🇸", PhonePrefix: "+34", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Netherlands", NameAr: "هولندا", Code: "NL", IsoCode: "NLD", Flag: "🇳🇱", PhonePrefix: "+31", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Belgium", NameAr: "بلجيكا", Code: "BE", IsoCode: "BEL", Flag: "🇧🇪", PhonePrefix: "+32", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Switzerland", NameAr: "سويسرا", Code: "CH", IsoCode: "CHE", Flag: "🇨🇭", PhonePrefix: "+41", CurrencyCode: "CHF", CurrencyLabel: "CHF — CHF", CurrencySymbol: "CHF"},
	{Name: "Austria", NameAr: "النمسا", Code: "AT", IsoCode: "AUT", Flag: "🇦🇹", PhonePrefix: "+43", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Sweden", NameAr: "السويد", Code: "SE", IsoCode: "SWE", Flag: "🇸🇪", PhonePrefix: "+46", CurrencyCode: "SEK", CurrencyLabel: "SEK — kr", CurrencySymbol: "kr"},
	{Name: "Norway", NameAr: "النرويج", Code: "NO", IsoCode: "NOR", Flag: "🇳🇴", PhonePrefix: "+47", CurrencyCode: "NOK", CurrencyLabel: "NOK — kr", CurrencySymbol: "kr"},
	{Name: "Denmark", NameAr: "الدنمارك", Code: "DK", IsoCode: "DNK", Flag: "🇩🇰", PhonePrefix: "+45", CurrencyCode: "DKK", CurrencyLabel: "DKK — kr", CurrencySymbol: "kr"},
	{Name: "Finland", NameAr: "فنلندا", Code: "FI", IsoCode: "FIN", Flag: "🇫🇮", PhonePrefix: "+358", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Poland", NameAr: "بولندا", Code: "PL", IsoCode: "POL", Flag: "🇵🇱", PhonePrefix: "+48", CurrencyCode: "PLN", CurrencyLabel: "PLN — zł", CurrencySymbol: "zł"},
	{Name: "Czech Republic", NameAr: "جمهورية التشيك", Code: "CZ", IsoCode: "CZE", Flag: "🇨🇿", PhonePrefix: "+420", CurrencyCode: "CZK", CurrencyLabel: "CZK — Kč", CurrencySymbol: "Kč"},
	{Name: "Hungary", NameAr: "المجر", Code: "HU", IsoCode: "HUN", Flag: "🇭🇺", PhonePrefix: "+36", CurrencyCode: "HUF", CurrencyLabel: "HUF — Ft", CurrencySymbol: "Ft"},
	{Name: "Portugal", NameAr: "البرتغال", Code: "PT", IsoCode: "PRT", Flag: "🇵🇹", PhonePrefix: "+351", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Greece", NameAr: "اليونان", Code: "GR", IsoCode: "GRC", Flag: "🇬🇷", PhonePrefix: "+30", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Ireland", NameAr: "أيرلندا", Code: "IE", IsoCode: "IRL", Flag: "🇮🇪", PhonePrefix: "+353", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Romania", NameAr: "رومانيا", Code: "RO", IsoCode: "ROU", Flag: "🇷🇴", PhonePrefix: "+40", CurrencyCode: "RON", CurrencyLabel: "RON — lei", CurrencySymbol: "lei"},
	{Name: "Bulgaria", NameAr: "بلغاريا", Code: "BG", IsoCode: "BGR", Flag: "🇧🇬", PhonePrefix: "+359", CurrencyCode: "BGN", CurrencyLabel: "BGN — лв", CurrencySymbol: "лв"},
	{Name: "Croatia", NameAr: "كرواتيا", Code: "HR", IsoCode: "HRV", Flag: "🇭🇷", PhonePrefix: "+385", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Slovakia", NameAr: "سلوفاكيا", Code: "SK", IsoCode: "SVK", Flag: "🇸🇰", PhonePrefix: "+421", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Slovenia", NameAr: "سلوفينيا", Code: "SI", IsoCode: "SVN", Flag: "🇸🇮", PhonePrefix: "+386", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Estonia", NameAr: "إستونيا", Code: "EE", IsoCode: "EST", Flag: "🇪🇪", PhonePrefix: "+372", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Latvia", NameAr: "لاتفيا", Code: "LV", IsoCode: "LVA", Flag: "🇱🇻", PhonePrefix: "+371", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Lithuania", NameAr: "ليتوانيا", Code: "LT", IsoCode: "LTU", Flag: "🇱🇹", PhonePrefix: "+370", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Iceland", NameAr: "أيسلندا", Code: "IS", IsoCode: "ISL", Flag: "🇮🇸", PhonePrefix: "+354", CurrencyCode: "ISK", CurrencyLabel: "ISK — kr", CurrencySymbol: "kr"},
	{Name: "Luxembourg", NameAr: "لوكسمبورغ", Code: "LU", IsoCode: "LUX", Flag: "🇱🇺", PhonePrefix: "+352", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Malta", NameAr: "مالطا", Code: "MT", IsoCode: "MLT", Flag: "🇲🇹", PhonePrefix: "+356", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
	{Name: "Cyprus", NameAr: "قبرص", Code: "CY", IsoCode: "CYP", Flag: "🇨🇾", PhonePrefix: "+357", CurrencyCode: "EUR", CurrencyLabel: "EUR — €", CurrencySymbol: "€"},
}

var countryByCodeMap map[string]Country
var countryByPhonePrefixMap map[string]Country
var countryByIsoCodeMap map[string]Country
var sortedCountries []Country

var initOnce sync.Once

func ensureInit() {
	initOnce.Do(func() {
		countryByCodeMap = make(map[string]Country, len(countries))
		countryByPhonePrefixMap = make(map[string]Country, len(countries))
		countryByIsoCodeMap = make(map[string]Country, len(countries))
		for _, c := range countries {
			countryByCodeMap[c.Code] = c
			countryByPhonePrefixMap[c.PhonePrefix] = c
			countryByIsoCodeMap[c.IsoCode] = c
		}

		sortedCountries = append([]Country(nil), countries...)
		sort.Slice(sortedCountries, func(i, j int) bool {
			return sortedCountries[i].Name < sortedCountries[j].Name
		})
	})
}

func FindByCode(code string) Country {
	ensureInit()
	if country, ok := countryByCodeMap[code]; ok {
		return country
	}
	return Country{}
}

func FindByPhonePrefix(prefix string) Country {
	ensureInit()
	if country, ok := countryByPhonePrefixMap[prefix]; ok {
		return country
	}
	return Country{}
}

func FindByIsoCode(isoCode string) Country {
	ensureInit()
	if country, ok := countryByIsoCodeMap[isoCode]; ok {
		return country
	}
	return Country{}
}

func Countries() []Country {
	ensureInit()
	return append([]Country(nil), sortedCountries...)
}
