# codelight.py
import pandas as pd
import numpy as np
from typing import List, Dict, Any, Optional, Tuple
from sklearn.linear_model import LinearRegression
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_absolute_error
from sklearn.preprocessing import StandardScaler
import pulp
import logging
from datetime import datetime

# Logger yapılandırması
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# ==================== VERİ YÜKLEME VE İŞLEME ====================
def load_data() -> Dict[str, pd.DataFrame]:
    """Veri setlerini yükler ve ön işleme yapar"""
    try:
        data = {
            "users": pd.read_csv("users.csv"),
            "coverage": pd.read_csv("coverage.csv"),
            "mobile_plans": pd.read_csv("mobile_plans.csv"),
            "home_plans": pd.read_csv("home_plans.csv"),
            "tv_plans": pd.read_csv("tv_plans.csv"),
            "bundling_rules": pd.read_csv("bundling_rules.csv"),
            "household": pd.read_csv("household.csv"),
            "current_services": pd.read_csv("current_services.csv", converters={'mobile_plan_ids': lambda x: x.split(';')}),
            "install_slots": pd.read_csv("install_slots.csv")
        }
        
        # applies_to sütunu için güvenli işleme
        if 'applies_to' in data['bundling_rules'].columns:
            data['bundling_rules']['applies_to'] = data['bundling_rules']['applies_to'].fillna('all').str.lower()
        else:
            data['bundling_rules']['applies_to'] = 'all'
            
        # Eksik sütunları ekle
        if 'rule_type' not in data['bundling_rules'].columns:
            data['bundling_rules']['rule_type'] = data['bundling_rules']['combo_key'].apply(
                lambda x: 'bundle' if '+' in x else 'extra_line'
            )
        
        if 'description' not in data['bundling_rules'].columns:
            data['bundling_rules']['description'] = data['bundling_rules']['combo_key'].apply(
                lambda x: {
                    'MOBILE+HOME': 'Mobil+Ev',
                    'MOBILE+HOME+TV': 'Mobil+Ev+TV',
                    'MOBILE': 'Mobil'
                }.get(x, x)
            )
        
        if 'discount_percent' not in data['bundling_rules'].columns:
            data['bundling_rules']['discount_percent'] = data['bundling_rules']['discount_pct'] * 100
            
        logger.info("Veri setleri başarıyla yüklendi")
        return data
    except Exception as e:
        logger.error(f"Veri yükleme hatası: {str(e)}")
        raise Exception(f"Veri yükleme hatası: {str(e)}")

# ==================== YAPAY ZEKA TAHMİN FONKSİYONLARI ====================
# Global model değişkenleri
usage_model = None
scaler = None
model_trained = False

def train_usage_model(data: Dict[str, pd.DataFrame]):
    """Kullanım tahmini için yapay zeka modelini eğitir"""
    global usage_model, scaler, model_trained
    
    try:
        # Model eğitimi için veri hazırlama
        household_data = data['household'].copy()
        user_data = data['users'].copy()
        
        # Kullanıcı verilerini birleştir
        merged = pd.merge(household_data, user_data, on='user_id')
        
        # Özellikler ve hedef değişkenler
        features = ['line_id', 'user_id', 'age', 'gender', 'city']
        targets = ['expected_gb', 'expected_min', 'tv_hd_hours']
        
        # Kategorik verileri sayısala çevir
        merged['line_id'] = merged['line_id'].str.replace('L-', '').astype(int)
        merged = pd.get_dummies(merged, columns=['gender', 'city'], drop_first=True)
        
        # Eksik sütunları kontrol et
        for col in features:
            if col not in merged.columns:
                merged[col] = 0
        
        X = merged[features]
        y = merged[targets]
        
        # Veriyi ölçeklendir
        scaler = StandardScaler()
        X_scaled = scaler.fit_transform(X)
        
        # Veri setini böl
        X_train, X_test, y_train, y_test = train_test_split(X_scaled, y, test_size=0.2, random_state=42)
        
        # Modeli eğit
        usage_model = LinearRegression()
        usage_model.fit(X_train, y_train)
        
        # Model performansını değerlendir
        y_pred = usage_model.predict(X_test)
        mae = mean_absolute_error(y_test, y_pred)
        logger.info(f"Model eğitildi. MAE: {mae:.2f}")
        
        model_trained = True
    except Exception as e:
        logger.error(f"Model eğitimi hatası: {str(e)}")
        raise Exception(f"Model eğitimi hatası: {str(e)}")

def predict_usage(user_id: int, line_id: str, user_data: pd.DataFrame) -> Dict[str, float]:
    """Belirtilen kullanıcı ve hat için kullanım tahmini yapar"""
    if not model_trained:
        raise Exception("Model eğitilmemiş!")
    
    try:
        # Line_id'yi sayısala çevir
        line_num = int(line_id.replace('L-', ''))
        
        # Kullanıcı verilerini hazırla
        user = user_data[user_data['user_id'] == user_id].iloc[0]
        
        # Özellik vektörünü oluştur
        features = {
            'line_id': line_num,
            'user_id': user_id,
            'age': user.get('age', 30),
            'gender': user.get('gender', 'unknown'),
            'city': user.get('city', 'unknown')
        }
        
        # DataFrame oluştur
        df = pd.DataFrame([features])
        
        # Kategorik verileri sayısala çevir
        df = pd.get_dummies(df, columns=['gender', 'city'], drop_first=True)
        
        # Eksik sütunları ekle
        for col in scaler.feature_names_in_:
            if col not in df.columns:
                df[col] = 0
        
        # Sıralamayı düzelt
        df = df[scaler.feature_names_in_]
        
        # Veriyi ölçeklendir
        X_scaled = scaler.transform(df)
        
        # Tahmin yap
        prediction = usage_model.predict(X_scaled)[0]
        
        return {
            'expected_gb': max(0, prediction[0]),
            'expected_min': max(0, prediction[1]),
            'tv_hd_hours': max(0, prediction[2])
        }
    except Exception as e:
        logger.error(f"Tahmin hatası: {str(e)}")
        raise Exception(f"Tahmin hatası: {str(e)}")

# ==================== MALİYET HESAPLAMA FONKSİYONLARI ====================
def calculate_mobile_line_cost(plan: pd.Series, expected_gb: float, expected_min: float) -> float:
    """Tek bir mobil hattın maliyetini hesaplar."""
    over_gb = max(0.0, expected_gb - plan['quota_gb']) * plan['overage_gb']
    over_min = max(0.0, expected_min - plan['quota_min']) * plan['overage_min']
    return plan['monthly_price'] + over_gb + over_min

def optimize_mobile_plans(expected_gb: List[float], expected_min: List[float], 
                         plans: pd.DataFrame) -> Tuple[float, List[Dict]]:
    """
    Birden fazla hat için en uygun mobil plan atamasını yapar.
    ILP (Integer Linear Programming) kullanarak en optimal çözümü bulur.
    """
    n = len(expected_gb)
    m = len(plans)
    cost_matrix = np.zeros((n, m))
    
    # Maliyet matrisini oluştur
    for i in range(n):
        for j in range(m):
            cost_matrix[i][j] = calculate_mobile_line_cost(plans.iloc[j], expected_gb[i], expected_min[i])
    
    # ILP problemi oluştur
    prob = pulp.LpProblem("MobilePlanAssignment", pulp.LpMinimize)
    
    # Karar değişkenleri: x[i,j] = 1 ise i. hat j. plana atanır
    x = pulp.LpVariable.dicts("x", ((i, j) for i in range(n) for j in range(m)), cat='Binary')
    
    # Amaç fonksiyonu: Toplam maliyeti minimize et
    prob += pulp.lpSum(cost_matrix[i][j] * x[i, j] for i in range(n) for j in range(m))
    
    # Kısıtlamalar: Her hat tam olarak bir plana atanmalı
    for i in range(n):
        prob += pulp.lpSum(x[i, j] for j in range(m)) == 1
    
    # Problemi çöz
    prob.solve(pulp.PULP_CBC_CMD(msg=False))
    
    # Sonuçları al
    chosen = []
    for i in range(n):
        for j in range(m):
            if pulp.value(x[i, j]) == 1:
                chosen.append(j)
                break
    
    # Toplam maliyet ve detayları hazırla
    total_cost = sum(cost_matrix[i][chosen[i]] for i in range(n))
    details = []
    for i in range(n):
        plan = plans.iloc[chosen[i]]
        details.append({
            "line_id": f"L-{i+1}",
            "plan_id": int(plan['plan_id']),
            "plan_name": plan['plan_name'],
            "cost": float(cost_matrix[i][chosen[i]]),
            "quota_gb": plan['quota_gb'],
            "quota_min": plan['quota_min'],
            "overage_gb": plan['overage_gb'],
            "overage_min": plan['overage_min']
        })
    
    return total_cost, details

def calculate_home_cost(address_id: str, coverage: pd.DataFrame, 
                       home_plans: pd.DataFrame, prefer_techs: List[str]) -> Tuple[Optional[float], Dict, str]:
    """Adrese uygun en iyi ev internet planını bulur."""
    try:
        # Kapsama kontrolü
        cov = coverage[coverage['address_id'] == address_id].iloc[0]
        available_techs = [tech for tech in prefer_techs if cov.get(tech, 0) == 1]
        
        if not available_techs:
            return None, {}, "None"
        
        # Teknoloji önceliğine göre plan ara
        for tech in prefer_techs:
            if tech in available_techs:
                tech_plans = home_plans[home_plans['tech'] == tech].copy()
                
                # VDSL için hız tahmini
                if tech == 'vdsl' and 'distance_km' in cov:
                    speed = vdsl_speed_from_distance(cov['distance_km'])
                    tech_plans = tech_plans[tech_plans['down_mbps'] <= speed]
                
                if not tech_plans.empty:
                    best_plan = tech_plans.sort_values('monthly_price').iloc[0]
                    return float(best_plan['monthly_price']), best_plan.to_dict(), tech
        
        return None, {}, "None"
    except Exception as e:
        logger.error(f"Ev interneti maliyet hesaplama hatası: {str(e)}")
        return None, {}, "None"

def calculate_tv_cost(expected_hours: float, tv_plans: pd.DataFrame) -> Tuple[float, Dict]:
    """TV ihtiyacına uygun en iyi planı bulur."""
    if expected_hours <= 0:
        return 0.0, {}
    
    filtered = tv_plans[tv_plans['hd_hours_included'] >= expected_hours]
    if filtered.empty:
        best = tv_plans.sort_values('hd_hours_included', ascending=False).iloc[0]
    else:
        best = filtered.sort_values('monthly_price').iloc[0]
    
    return float(best['monthly_price']), best.to_dict()

# ==================== İNDİRİM HESAPLAMA ====================
def apply_extra_line_discount(mobile_total: float, num_lines: int, rules: pd.DataFrame) -> float:
    """Ek hat indirimi uygular (sadece mobil bileşene)."""
    if num_lines >= 3:
        rule = rules[(rules['rule_type'] == 'extra_line') & 
                    (rules['description'].str.contains("3\+ Hat", na=False))]
    elif num_lines == 2:
        rule = rules[(rules['rule_type'] == 'extra_line') & 
                    (rules['description'].str.contains("2. Hat", na=False))]
    else:
        return mobile_total
    
    if not rule.empty:
        discount = rule.iloc[0]['discount_percent'] / 100
        return mobile_total * (1 - discount)
    return mobile_total

def apply_bundle_discount(total_cost: float, bundle_type: str, rules: pd.DataFrame) -> float:
    """Kombinasyon indirimlerini uygular."""
    if bundle_type == "MOBILE+HOME":
        rule = rules[(rules['rule_type'] == 'bundle') & 
                    (rules['description'].str.contains("Mobil\+Ev", na=False))]
    elif bundle_type == "MOBILE+HOME+TV":
        rule = rules[(rules['rule_type'] == 'bundle') & 
                    (rules['description'].str.contains("Mobil\+Ev\+TV", na=False))]
    else:
        return total_cost
    
    if not rule.empty:
        discount = rule.iloc[0]['discount_percent'] / 100
        return total_cost * (1 - discount)
    return total_cost

# ==================== VDSL HIZ TAHMİNİ ====================
def vdsl_speed_from_distance(distance_km: float) -> float:
    """Mesafeye göre VDSL hızını tahmin eder (mock)."""
    if pd.isna(distance_km):
        return 35.0
    return max(8.0, 100.0 * np.exp(-0.6 * distance_km))

# ==================== DÖVİZ NORMALİZASYONU ====================
def normalize_currency(amount: float, from_currency: str = 'TL', to_currency: str = 'TL') -> float:
    """
    Para birimini normalize eder. Şimdilik sadece TL'yi destekliyor.
    """
    if from_currency == to_currency:
        return amount
    
    # Gerçek bir uygulamada burada döviz kurları çekilir
    # Şimdilik sabit kurlar kullanalım
    exchange_rates = {
        'USD': 0.03,  # 1 TL = 0.03 USD
        'EUR': 0.028, # 1 TL = 0.028 EUR
    }
    
    if to_currency in exchange_rates:
        # İki ondalık basamağa yuvarla
        return round(amount * exchange_rates[to_currency], 2)
    
    return amount

# ==================== ÖNERİ MOTORU ====================
def recommend_bundles(user_id: int, data: Dict[str, pd.DataFrame], 
                     prefer_techs: List[str] = None, 
                     use_ai_prediction: bool = False) -> Tuple[List[Dict], Dict]:
    """En iyi 3 paket kombinasyonunu önerir (AI tahmini seçeneği ile)"""
    if prefer_techs is None:
        prefer_techs = ['fiber', 'vdsl', 'fwa']
    
    try:
        # Kullanıcı verilerini al
        user = data['users'][data['users']['user_id'] == user_id].iloc[0]
        household = data['household'][data['household']['user_id'] == user_id]
        current = data['current_services'][data['current_services']['user_id'] == user_id].iloc[0]
        
        # AI tahmini kullanılacaksa verileri güncelle
        if use_ai_prediction:
            predicted_data = []
            for _, row in household.iterrows():
                pred = predict_usage(user_id, row['line_id'], data['users'])
                predicted_data.append({
                    'line_id': row['line_id'],
                    'expected_gb': pred['expected_gb'],
                    'expected_min': pred['expected_min'],
                    'tv_hd_hours': pred['tv_hd_hours']
                })
            household = pd.DataFrame(predicted_data)
        
        # Mevcut maliyeti hesapla
        current_mobile_cost = sum(
            data['mobile_plans'][data['mobile_plans']['plan_id'].astype(str) == pid]['monthly_price'].iloc[0]
            for pid in current['mobile_plan_ids']
        )

        current_total = current_mobile_cost
        if current['has_home']:
            current_total += data['home_plans'][
                (data['home_plans']['tech'] == current['home_tech']) & 
                (data['home_plans']['down_mbps'] == current['home_speed'])
            ]['monthly_price'].iloc[0]

        if current['has_tv']:
            current_total += data['tv_plans']['monthly_price'].min()
        
        # Mobil hatlar için optimizasyon yap
        mobile_cost, mobile_details = optimize_mobile_plans(
            household['expected_gb'].tolist(),
            household['expected_min'].tolist(),
            data['mobile_plans']
        )
        
        # Ek hat indirimi uygula (sadece mobil bileşene)
        mobile_discounted = apply_extra_line_discount(
            mobile_cost, len(household), data['bundling_rules']
        )
        
        # Ev interneti öner
        home_cost, home_info, tech_label = calculate_home_cost(
            user['address_id'], data['coverage'], data['home_plans'], prefer_techs
        )
        
        # TV öner
        tv_hours = household['tv_hd_hours'].sum()
        tv_cost, tv_info = calculate_tv_cost(tv_hours, data['tv_plans'])
        
        # Kombinasyonları oluştur
        combos = []
        
        # 1. Sadece Mobil
        combos.append({
            "combo_label": "Sadece Mobil Paket",
            "items": {"mobile": mobile_details},
            "monthly_total": mobile_discounted,
            "tech_label": "None",
            "reasoning": "Sadece mobil hatlar için optimize edilmiş paket."
        })
        
        # 2. Mobil + Ev
        if home_cost is not None:
            total = mobile_discounted + home_cost
            disc_total = apply_bundle_discount(total, "MOBILE+HOME", data['bundling_rules'])
            
            combos.append({
                "combo_label": "Mobil + Ev İnterneti",
                "items": {
                    "mobile": mobile_details,
                    "home": home_info
                },
                "monthly_total": disc_total,
                "tech_label": tech_label,
                "reasoning": f"Mobil ve ev interneti birleşimi. {tech_label.upper()} teknolojisi."
            })
            
            # 3. Mobil + Ev + TV
            if tv_cost > 0:
                total = mobile_discounted + home_cost + tv_cost
                disc_total = apply_bundle_discount(total, "MOBILE+HOME+TV", data['bundling_rules'])
                
                combos.append({
                    "combo_label": "Mobil + Ev + TV Paketi",
                    "items": {
                        "mobile": mobile_details,
                        "home": home_info,
                        "tv": tv_info
                    },
                    "monthly_total": disc_total,
                    "tech_label": tech_label,
                    "reasoning": "Tam paket: Mobil, ev interneti ve TV."
                })
        
        # Sırala ve en iyi 3'ü seç
        combos.sort(key=lambda x: x['monthly_total'])
        top3 = combos[:3]
        
        # Tasarrufu hesapla
        for combo in top3:
            combo['savings'] = current_total - combo['monthly_total']
        
        # Tahmin verileri
        pred_data = {
            "total_gb": sum(household['expected_gb']),
            "total_min": sum(household['expected_min']),
            "total_tv_hours": sum(household['tv_hd_hours'])
        }
        
        logger.info(f"Kullanıcı {user_id} için öneriler oluşturuldu")
        return top3, pred_data
    
    except Exception as e:
        logger.error(f"Öneri oluşturulamadı: {str(e)}")
        raise Exception(f"Öneri oluşturulamadı: {str(e)}")

# ==================== RANDEVU İŞLEMLERİ ====================
def get_available_slots(address_id: str, tech: str, data: Dict[str, pd.DataFrame]) -> List[Dict]:
    """Uygun kurulum slotlarını getirir."""
    try:
        slots = data['install_slots']
        available = slots[(slots['address_id'] == address_id) & 
                         (slots['tech'] == tech)].copy()
        
        if available.empty:
            return []
        
        available['slot_start'] = pd.to_datetime(available['slot_start'])
        available['slot_end'] = pd.to_datetime(available['slot_end'])
        available = available.sort_values('slot_start')
        
        return available.to_dict('records')
    except Exception as e:
        logger.error(f"Slot getirme hatası: {str(e)}")
        return []

def book_installation(user_id: int, slot_id: str, data: Dict[str, pd.DataFrame]) -> bool:
    """Randevu kaydı yapar (mock)."""
    try:
        # Gerçek uygulamada burada veritabanı güncellemesi yapılır
        logger.info(f"Mock: {user_id} kullanıcısı için {slot_id} slotu rezerve edildi")
        return True
    except Exception as e:
        logger.error(f"Randevu kaydı hatası: {str(e)}")
        return False