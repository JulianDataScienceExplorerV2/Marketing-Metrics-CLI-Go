import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

# Load the campaign data
df = pd.read_csv("C:/tmp/go-mkt/data/sample_campaigns.csv")

# Calculate metrics
df['CTR'] = (df['clicks'] / df['impressions']) * 100
df['CVR'] = (df['conversions'] / df['clicks']) * 100
df['CPA'] = df['spend'] / df['conversions']
df['ROAS'] = df['revenue'] / df['spend']

# Drop any infinite/nan values for plotting if conversions were 0
df.replace([np.inf, -np.inf], np.nan, inplace=True)
df.dropna(subset=['CPA', 'ROAS'], inplace=True)

plt.style.use('dark_background')
fig, ax = plt.subplots(figsize=(12, 8))

# Define colors based on ROAS performance (>3 is great, 1-3 is ok, <1 is bad)
colors = []
for roas in df['ROAS']:
    if roas >= 3.0:
        colors.append('#00ff9d') # Green
    elif roas >= 1.0:
        colors.append('#00e5ff') # Cyan/Blue
    else:
        colors.append('#ff3366') # Red

# Create bubble chart
scatter = ax.scatter(
    df['CPA'], 
    df['ROAS'], 
    s=df['spend'] * 0.5, # Bubble size based on spend
    c=colors, 
    alpha=0.7, 
    edgecolors='#ffffff', 
    linewidth=1.5
)

# Label the bubbles
for i, row in df.iterrows():
    ax.annotate(
        row['campaign_name'].replace(" - ", "\n"), 
        (row['CPA'], row['ROAS']),
        xytext=(0, 15), 
        textcoords='offset points',
        ha='center', 
        va='bottom', 
        fontsize=9,
        color='white',
        weight='bold'
    )

# Formatting
ax.set_title('Campaign Performance Matrix: CPA vs ROAS', fontsize=18, pad=20, weight='bold')
ax.set_xlabel('Cost Per Acquisition (CPA) - $ (Lower is Better)', fontsize=12)
ax.set_ylabel('Return on Ad Spend (ROAS) - Multiplier (Higher is Better)', fontsize=12)

# Add quadrant lines
ax.axhline(y=1.0, color='#666666', linestyle='--', alpha=0.5)
ax.axvline(x=df['CPA'].median(), color='#666666', linestyle='--', alpha=0.5)

# Add ROAS threshold area shading
ax.axhspan(1.0, 30.0, alpha=0.1, color='#00ff9d')
ax.axhspan(0.0, 1.0, alpha=0.1, color='#ff3366')

# Grid and limits
ax.grid(alpha=0.2, linestyle=':')
ax.set_xlim(left=0)
ax.set_ylim(bottom=0)

plt.tight_layout()
plt.savefig("C:/tmp/go-mkt/assets/campaign_performance.png", dpi=300, bbox_inches='tight', transparent=True)
print("Plot generated successfully!")
