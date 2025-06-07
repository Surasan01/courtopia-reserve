import { useState, useEffect } from "react";
import { format } from "date-fns";
import Header from "@/components/Header";
import Hero from "@/components/Hero";
import { DatePicker } from "@/components/DatePicker";
import { TimePicker } from "@/components/TimePicker";
import CourtCard from "@/components/CourtCard";
import BookingDialog from "@/components/BookingDialog";
import { fetchAvailableCourts } from "@/lib/api";
import { useToast } from "@/hooks/use-toast";

interface Court {
  courtNumber: number;
  name: string;
  time: string;
  isAvailable: boolean;
}

const Index = () => {
  const [selectedDate, setSelectedDate] = useState<Date | undefined>(new Date());
  const [selectedTime, setSelectedTime] = useState<string>("17:30 - 18:30");
  const [courts, setCourts] = useState<Court[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedCourt, setSelectedCourt] = useState<Court | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    if (!selectedDate || !selectedTime) return;

    const loadAvailableCourts = async () => {
      setIsLoading(true);
      try {
        const [startTime, endTime] = selectedTime.split(" - ");
        
        const formattedDate = format(selectedDate, "yyyy-MM-dd");
        
        const response = await fetchAvailableCourts(formattedDate, startTime, endTime);
        
        const transformedCourts = response.courts.map((court: any) => ({
          courtNumber: court.courtNumber,
          name: `คอร์ท ${court.courtNumber}`,
          time: `${response.startTime} - ${response.endTime}`,
          isAvailable: court.isAvailable
        }));
        
        setCourts(transformedCourts);
      } catch (error) {
        toast({
          title: 'ไม่สามารถโหลดข้อมูลคอร์ทได้',
          description: error instanceof Error ? error.message : 'กรุณาลองใหม่อีกครั้ง',
          variant: 'destructive',
        });
        
        setCourts([
          { courtNumber: 1, name: "คอร์ท 1", time: selectedTime, isAvailable: true },
          { courtNumber: 2, name: "คอร์ท 2", time: selectedTime, isAvailable: false },
          { courtNumber: 3, name: "คอร์ท 3", time: selectedTime, isAvailable: true },
          { courtNumber: 4, name: "คอร์ท 4", time: selectedTime, isAvailable: true },
          { courtNumber: 5, name: "คอร์ท 5", time: selectedTime, isAvailable: false },
          { courtNumber: 6, name: "คอร์ท 6", time: selectedTime, isAvailable: true }
        ]);
      } finally {
        setIsLoading(false);
      }
    };

    loadAvailableCourts();
  }, [selectedDate, selectedTime, toast]);

  const handleCourtClick = (court: Court) => {
    if (court.isAvailable) {
      setSelectedCourt(court);
      setDialogOpen(true);
    }
  };

  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <Hero />
      
      <main className="flex-1 bg-court-yellow/50">
        <div className="container mx-auto px-4 py-8">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-w-lg mx-auto mb-8">
            <DatePicker date={selectedDate} setDate={setSelectedDate} />
            <TimePicker selectedTime={selectedTime} setSelectedTime={setSelectedTime} />
          </div>
          
          {isLoading ? (
            <div className="flex justify-center items-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-court-orange"></div>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {courts.map((court) => (
                <CourtCard
                  key={court.courtNumber}
                  id={court.courtNumber}
                  name={court.name}
                  time={court.time}
                  isAvailable={court.isAvailable}
                  onClick={() => handleCourtClick(court)}
                />
              ))}
            </div>
          )}
        </div>
      </main>
      
      {selectedCourt && (
        <BookingDialog
          open={dialogOpen}
          onOpenChange={setDialogOpen}
          courtName={selectedCourt.name}
          courtTime={selectedCourt.time}
          courtNumber={selectedCourt.courtNumber}
          bookingDate={selectedDate ? format(selectedDate, "yyyy-MM-dd") : ""}
        />
      )}
      
      <footer className="bg-court-orange text-white p-4 text-center">
        <p>© {new Date().getFullYear()} Courtminton - King Mongkut's University of Technology Thonburi</p>
      </footer>
    </div>
  );
};

export default Index;