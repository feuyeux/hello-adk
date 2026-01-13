package main

import "sync"

// Element holds basic periodic table data.
type Element struct {
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	ChineseName  string  `json:"chinese_name"`
	AtomicNumber int     `json:"atomic_number"`
	AtomicWeight float64 `json:"atomic_weight"`
}

var (
	tableOnce    sync.Once
	elementByKey map[string]Element
)

func lookupElement(key string) (Element, bool) {
	tableOnce.Do(loadTable)
	el, ok := elementByKey[key]
	return el, ok
}

func loadTable() {
	elementByKey = make(map[string]Element)
	addElement("H", "Hydrogen", "氢", 1, 1.008)
	addElement("He", "Helium", "氦", 2, 4.0026)
	addElement("Li", "Lithium", "锂", 3, 6.94)
	addElement("Be", "Beryllium", "铍", 4, 9.0122)
	addElement("B", "Boron", "硼", 5, 10.81)
	addElement("C", "Carbon", "碳", 6, 12.011)
	addElement("N", "Nitrogen", "氮", 7, 14.007)
	addElement("O", "Oxygen", "氧", 8, 15.999)
	addElement("F", "Fluorine", "氟", 9, 18.998)
	addElement("Ne", "Neon", "氖", 10, 20.180)
	addElement("Na", "Sodium", "钠", 11, 22.990)
	addElement("Mg", "Magnesium", "镁", 12, 24.305)
	addElement("Al", "Aluminum", "铝", 13, 26.982)
	addElement("Si", "Silicon", "硅", 14, 28.085)
	addElement("P", "Phosphorus", "磷", 15, 30.974)
	addElement("S", "Sulfur", "硫", 16, 32.06)
	addElement("Cl", "Chlorine", "氯", 17, 35.45)
	addElement("Ar", "Argon", "氩", 18, 39.948)
	addElement("K", "Potassium", "钾", 19, 39.098)
	addElement("Ca", "Calcium", "钙", 20, 40.078)
	addElement("Sc", "Scandium", "钪", 21, 44.956)
	addElement("Ti", "Titanium", "钛", 22, 47.867)
	addElement("V", "Vanadium", "钒", 23, 50.942)
	addElement("Cr", "Chromium", "铬", 24, 51.996)
	addElement("Mn", "Manganese", "锰", 25, 54.938)
	addElement("Fe", "Iron", "铁", 26, 55.845)
	addElement("Co", "Cobalt", "钴", 27, 58.933)
	addElement("Ni", "Nickel", "镍", 28, 58.693)
	addElement("Cu", "Copper", "铜", 29, 63.546)
	addElement("Zn", "Zinc", "锌", 30, 65.38)
	addElement("Ga", "Gallium", "镓", 31, 69.723)
	addElement("Ge", "Germanium", "锗", 32, 72.630)
	addElement("As", "Arsenic", "砷", 33, 74.922)
	addElement("Se", "Selenium", "硒", 34, 78.971)
	addElement("Br", "Bromine", "溴", 35, 79.904)
	addElement("Kr", "Krypton", "氪", 36, 83.798)
	addElement("Rb", "Rubidium", "铷", 37, 85.468)
	addElement("Sr", "Strontium", "锶", 38, 87.62)
	addElement("Y", "Yttrium", "钇", 39, 88.906)
	addElement("Zr", "Zirconium", "锆", 40, 91.224)
	addElement("Nb", "Niobium", "铌", 41, 92.906)
	addElement("Mo", "Molybdenum", "钼", 42, 95.95)
	addElement("Tc", "Technetium", "锝", 43, 98)
	addElement("Ru", "Ruthenium", "钌", 44, 101.07)
	addElement("Rh", "Rhodium", "铑", 45, 102.91)
	addElement("Pd", "Palladium", "钯", 46, 106.42)
	addElement("Ag", "Silver", "银", 47, 107.87)
	addElement("Cd", "Cadmium", "镉", 48, 112.41)
	addElement("In", "Indium", "铟", 49, 114.82)
	addElement("Sn", "Tin", "锡", 50, 118.71)
	addElement("Sb", "Antimony", "锑", 51, 121.76)
	addElement("Te", "Tellurium", "碲", 52, 127.60)
	addElement("I", "Iodine", "碘", 53, 126.90)
	addElement("Xe", "Xenon", "氙", 54, 131.29)
	addElement("Cs", "Cesium", "铯", 55, 132.91)
	addElement("Ba", "Barium", "钡", 56, 137.33)
	addElement("La", "Lanthanum", "镧", 57, 138.91)
	addElement("Ce", "Cerium", "铈", 58, 140.12)
	addElement("Pr", "Praseodymium", "镨", 59, 140.91)
	addElement("Nd", "Neodymium", "钕", 60, 144.24)
	addElement("Pm", "Promethium", "钷", 61, 145)
	addElement("Sm", "Samarium", "钐", 62, 150.36)
	addElement("Eu", "Europium", "铕", 63, 151.96)
	addElement("Gd", "Gadolinium", "钆", 64, 157.25)
	addElement("Tb", "Terbium", "铽", 65, 158.93)
	addElement("Dy", "Dysprosium", "镝", 66, 162.50)
	addElement("Ho", "Holmium", "钬", 67, 164.93)
	addElement("Er", "Erbium", "铒", 68, 167.26)
	addElement("Tm", "Thulium", "铥", 69, 168.93)
	addElement("Yb", "Ytterbium", "镱", 70, 173.05)
	addElement("Lu", "Lutetium", "镥", 71, 174.97)
	addElement("Hf", "Hafnium", "铪", 72, 178.49)
	addElement("Ta", "Tantalum", "钽", 73, 180.95)
	addElement("W", "Tungsten", "钨", 74, 183.84)
	addElement("Re", "Rhenium", "铼", 75, 186.21)
	addElement("Os", "Osmium", "锇", 76, 190.23)
	addElement("Ir", "Iridium", "铱", 77, 192.22)
	addElement("Pt", "Platinum", "铂", 78, 195.08)
	addElement("Au", "Gold", "金", 79, 196.97)
	addElement("Hg", "Mercury", "汞", 80, 200.59)
	addElement("Tl", "Thallium", "铊", 81, 204.38)
	addElement("Pb", "Lead", "铅", 82, 207.2)
	addElement("Bi", "Bismuth", "铋", 83, 208.98)
	addElement("Po", "Polonium", "钋", 84, 209)
	addElement("At", "Astatine", "砹", 85, 210)
	addElement("Rn", "Radon", "氡", 86, 222)
	addElement("Fr", "Francium", "钫", 87, 223)
	addElement("Ra", "Radium", "镭", 88, 226)
	addElement("Ac", "Actinium", "锕", 89, 227)
	addElement("Th", "Thorium", "钍", 90, 232.04)
	addElement("Pa", "Protactinium", "镤", 91, 231.04)
	addElement("U", "Uranium", "铀", 92, 238.03)
	addElement("Np", "Neptunium", "镎", 93, 237)
	addElement("Pu", "Plutonium", "钚", 94, 244)
	addElement("Am", "Americium", "镅", 95, 243)
	addElement("Cm", "Curium", "锔", 96, 247)
	addElement("Bk", "Berkelium", "锫", 97, 247)
	addElement("Cf", "Californium", "锎", 98, 251)
	addElement("Es", "Einsteinium", "锿", 99, 252)
	addElement("Fm", "Fermium", "镄", 100, 257)
	addElement("Md", "Mendelevium", "钔", 101, 258)
	addElement("No", "Nobelium", "锘", 102, 259)
	addElement("Lr", "Lawrencium", "铹", 103, 266)
	addElement("Rf", "Rutherfordium", "𬬻", 104, 267)
	addElement("Db", "Dubnium", "𬭊", 105, 270)
	addElement("Sg", "Seaborgium", "𬭳", 106, 271)
	addElement("Bh", "Bohrium", "𬭛", 107, 270)
	addElement("Hs", "Hassium", "𬭶", 108, 277)
	addElement("Mt", "Meitnerium", "鿏", 109, 278)
	addElement("Ds", "Darmstadtium", "𫟼", 110, 281)
	addElement("Rg", "Roentgenium", "𬬭", 111, 282)
	addElement("Cn", "Copernicium", "鿔", 112, 285)
	addElement("Fl", "Flerovium", "𫓧", 114, 289)
	addElement("Lv", "Livermorium", "𫟷", 116, 293)
	addElement("Ts", "Tennessine", "鿬", 117, 294)
	addElement("Og", "Oganesson", "鿫", 118, 294)
}

func addElement(symbol, name, chinese string, number int, weight float64) {
	el := Element{
		Symbol:       symbol,
		Name:         name,
		ChineseName:  chinese,
		AtomicNumber: number,
		AtomicWeight: weight,
	}
	elementByKey[symbol] = el
	elementByKey[chinese] = el
}
