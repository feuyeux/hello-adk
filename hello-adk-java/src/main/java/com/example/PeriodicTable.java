package com.example;

import java.util.HashMap;
import java.util.Map;

public class PeriodicTable {
    
    private static PeriodicTable instance;
    private final Map<String, Element> symbolMap;
    private final Map<String, String> chineseToSymbolMap;

    private PeriodicTable() {
        this.symbolMap = new HashMap<>();
        this.chineseToSymbolMap = new HashMap<>();
        initializePeriodicTable();
    }

    public static synchronized PeriodicTable getInstance() {
        if (instance == null) {
            instance = new PeriodicTable();
        }
        return instance;
    }

    private void initializePeriodicTable() {
        addElement("H", "Hydrogen", "氢", 1, 1.008);
        addElement("He", "Helium", "氦", 2, 4.0026);
        addElement("Li", "Lithium", "锂", 3, 6.94);
        addElement("Be", "Beryllium", "铍", 4, 9.0122);
        addElement("B", "Boron", "硼", 5, 10.81);
        addElement("C", "Carbon", "碳", 6, 12.011);
        addElement("N", "Nitrogen", "氮", 7, 14.007);
        addElement("O", "Oxygen", "氧", 8, 15.999);
        addElement("F", "Fluorine", "氟", 9, 18.998);
        addElement("Ne", "Neon", "氖", 10, 20.180);
        addElement("Na", "Sodium", "钠", 11, 22.990);
        addElement("Mg", "Magnesium", "镁", 12, 24.305);
        addElement("Al", "Aluminum", "铝", 13, 26.982);
        addElement("Si", "Silicon", "硅", 14, 28.085);
        addElement("P", "Phosphorus", "磷", 15, 30.974);
        addElement("S", "Sulfur", "硫", 16, 32.06);
        addElement("Cl", "Chlorine", "氯", 17, 35.45);
        addElement("Ar", "Argon", "氩", 18, 39.948);
        addElement("K", "Potassium", "钾", 19, 39.098);
        addElement("Ca", "Calcium", "钙", 20, 40.078);
        addElement("Sc", "Scandium", "钪", 21, 44.956);
        addElement("Ti", "Titanium", "钛", 22, 47.867);
        addElement("V", "Vanadium", "钒", 23, 50.942);
        addElement("Cr", "Chromium", "铬", 24, 51.996);
        addElement("Mn", "Manganese", "锰", 25, 54.938);
        addElement("Fe", "Iron", "铁", 26, 55.845);
        addElement("Co", "Cobalt", "钴", 27, 58.933);
        addElement("Ni", "Nickel", "镍", 28, 58.693);
        addElement("Cu", "Copper", "铜", 29, 63.546);
        addElement("Zn", "Zinc", "锌", 30, 65.38);
        addElement("Ga", "Gallium", "镓", 31, 69.723);
        addElement("Ge", "Germanium", "锗", 32, 72.630);
        addElement("As", "Arsenic", "砷", 33, 74.922);
        addElement("Se", "Selenium", "硒", 34, 78.971);
        addElement("Br", "Bromine", "溴", 35, 79.904);
        addElement("Kr", "Krypton", "氪", 36, 83.798);
        // Adding more common elements
        addElement("Ag", "Silver", "银", 47, 107.868);
        addElement("Au", "Gold", "金", 79, 196.967);
        addElement("Pt", "Platinum", "铂", 78, 195.084);
        addElement("Pb", "Lead", "铅", 82, 207.2);
        addElement("Hg", "Mercury", "汞", 80, 200.592);
    }

    private void addElement(String symbol, String name, String chineseName, int atomicNumber, double atomicWeight) {
        Element element = new Element(symbol, name, chineseName, atomicNumber, atomicWeight);
        symbolMap.put(symbol, element);
        chineseToSymbolMap.put(chineseName, symbol);
    }

    public Element getElement(String query) {
        // Try to find by symbol first
        if (symbolMap.containsKey(query)) {
            return symbolMap.get(query);
        }
        
        // Then try to find by Chinese name
        if (chineseToSymbolMap.containsKey(query)) {
            String symbol = chineseToSymbolMap.get(query);
            return symbolMap.get(symbol);
        }
        
        return null;
    }

    public static class Element {
        public String symbol;
        public String name;
        public String chineseName;
        public int atomicNumber;
        public double atomicWeight;

        public Element(String symbol, String name, String chineseName, int atomicNumber, double atomicWeight) {
            this.symbol = symbol;
            this.name = name;
            this.chineseName = chineseName;
            this.atomicNumber = atomicNumber;
            this.atomicWeight = atomicWeight;
        }
    }
}
