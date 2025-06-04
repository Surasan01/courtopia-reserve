import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";
import { useAuth } from "@/hooks/useAuth";
import { createBooking } from "@/lib/api";

interface BookingDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  courtName: string;
  courtTime: string;
  courtNumber: number;
  bookingDate: string;
}

const BookingDialog = ({ 
  open, 
  onOpenChange, 
  courtName, 
  courtTime, 
  courtNumber, 
  bookingDate 
}: BookingDialogProps) => {
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuth();
  
  // Extract start time and end time from courtTime (format: "HH:MM - HH:MM")
  const [startTime, endTime] = courtTime.split(" - ").map(t => t.replace(" น.", ""));

  const handleBooking = async () => {
    if (!isAuthenticated) {
      toast({
        title: "กรุณาเข้าสู่ระบบ",
        description: "คุณจำเป็นต้องเข้าสู่ระบบก่อนจองคอร์ท",
        variant: "destructive",
      });
      navigate("/login");
      return;
    }

    const currentTime = new Date();
  const bookingStartTime = new Date(`${bookingDate}T${startTime}:00`);
  if (currentTime > bookingStartTime) {
    toast({
      title: "ไม่สามารถจองคอร์ทได้",
      description: "เวลาที่คุณเลือกเลยเวลาปัจจุบันแล้ว",
      variant: "destructive",
    });
    return;
  }
  
    setIsLoading(true);
    
    try {
      const bookingData = {
        courtNumber,
        bookingDate,
        startTime,
        endTime,
      };
      
      await createBooking(user!.token, bookingData);
      
      toast({
        title: "จองคอร์ทสำเร็จ",
        description: `คุณได้จอง ${courtName} เวลา ${courtTime}`,
      });
      
      // Navigate to bookings page
      navigate("/bookings");
    } catch (error) {
      toast({
        title: "ไม่สามารถจองคอร์ทได้",
        description: error instanceof Error ? error.message : "กรุณาลองใหม่อีกครั้ง",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
      onOpenChange(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>จองคอร์ทแบดมินตัน</DialogTitle>
          <DialogDescription>
            {courtName} เวลา {courtTime}
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-4 py-4">
          <div className="flex flex-col space-y-1.5">
            <Label htmlFor="name">ข้อมูลการจอง</Label>
            <div className="grid grid-cols-2 gap-2">
              <div>
                <p className="text-sm text-muted-foreground">วันที่:</p>
                <p>{bookingDate}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">เวลา:</p>
                <p>{courtTime}</p>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button 
            type="button" 
            className="bg-court-orange hover:bg-court-orange/90"
            onClick={handleBooking}
            disabled={isLoading}
          >
            {isLoading ? 'กำลังดำเนินการ...' : 'ยืนยันการจอง'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default BookingDialog;