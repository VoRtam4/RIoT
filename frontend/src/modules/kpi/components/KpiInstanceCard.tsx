/**
 * @file KpiInstanceCard.tsx
 * @brief Karta instance navázané na vybranou KPI definici.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";

type Props = {
  instance: {
    label?: string | null;
    uid?: string | null;
  };
  selected: boolean;
  onClick: () => void;
};

export default function KpiInstanceCard({
  instance,
  selected,
  onClick,
}: Props) {
  return (
    <Card
      sx={{
        flexShrink: 0,
        height: "100%",
        overflow: "hidden",
        backgroundColor: selected ? "var(--border-input)" : "var(--bg-card)",
        "&:hover": {
          backgroundColor: selected ? "var(--border-input)" : "var(--border)",
        },
      }}
    >
      <CardActionArea
        onClick={onClick}
        sx={{
          height: "100%",
        }}
      >
        <CardContent
          sx={{
            py: 1.5,
            height: "100%",
            display: "flex",
            alignItems: "center",
            overflow: "hidden",
          }}
        >
          <Typography
            variant="body1"
            sx={{
              display: "-webkit-box",
              overflow: "hidden",
              textOverflow: "ellipsis",
              WebkitBoxOrient: "vertical",
              WebkitLineClamp: 2,
            }}
          >
            {instance.label || instance.uid || "\u00A0"}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
