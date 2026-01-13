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
	addElement("Ag", "Silver", "银", 47, 107.87)
	addElement("Sn", "Tin", "锡", 50, 118.71)
	addElement("I", "Iodine", "碘", 53, 126.90)
	addElement("Xe", "Xenon", "氙", 54, 131.29)
	addElement("Cs", "Cesium", "铯", 55, 132.91)
	addElement("Ba", "Barium", "钡", 56, 137.33)
	addElement("Au", "Gold", "金", 79, 196.97)
	addElement("Hg", "Mercury", "汞", 80, 200.59)
	addElement("Pb", "Lead", "铅", 82, 207.2)
	addElement("U", "Uranium", "铀", 92, 238.03)
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
