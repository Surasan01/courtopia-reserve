import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useToast } from '@/hooks/use-toast';
import { register as registerUser } from '@/lib/api';

import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Eye, EyeOff } from 'lucide-react';

const registerSchema = z.object({
    studentId: z.string().min(1, 'กรุณากรอกรหัสนักศึกษา'),
    name: z.string().min(1, 'กรุณากรอกชื่อ-นามสกุล'),
    email: z.string().email('อีเมลไม่ถูกต้อง').optional().or(z.literal('')),
    password: z.string().min(6, 'รหัสผ่านต้องมีอย่างน้อย 6 ตัวอักษร'),
    confirmPassword: z.string().min(1, 'กรุณายืนยันรหัสผ่าน'),
}).refine((data) => data.password === data.confirmPassword, {
    message: "รหัสผ่านไม่ตรงกัน",
    path: ["confirmPassword"],
});

type RegisterFormValues = z.infer<typeof registerSchema>;

const Register = () => {
    const navigate = useNavigate();
    const { toast } = useToast();
    const [isLoading, setIsLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const form = useForm<RegisterFormValues>({
        resolver: zodResolver(registerSchema),
        defaultValues: {
            studentId: '',
            name: '',
            email: '',
            password: '',
            confirmPassword: '',
        },
    });

    const onSubmit = async (data: RegisterFormValues) => {
        setIsLoading(true);
        try {
            await registerUser(
                data.studentId,
                data.password,
                data.name,
                data.email || undefined
            );

            toast({
                title: 'ลงทะเบียนสำเร็จ',
                description: 'กรุณาเข้าสู่ระบบด้วยข้อมูลที่ลงทะเบียน',
            });

            navigate('/login');
        } catch (error) {
            toast({
                title: 'ลงทะเบียนไม่สำเร็จ',
                description: error instanceof Error ? error.message : 'เกิดข้อผิดพลาด กรุณาลองอีกครั้ง',
                variant: 'destructive',
            });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="container flex items-center justify-center min-h-screen py-8">
            <Card className="w-full max-w-md">
                <CardHeader className="text-center">
                    <div className="flex justify-center mb-2">
                        <img
                            src="/lovable-uploads/da3c37a1-a8c1-4fe0-b7ba-7b4033af8179.png"
                            alt="Shuttlecock Logo"
                            className="w-16 h-16 object-contain"
                        />
                    </div>
                    <CardTitle className="text-2xl font-bold text-court-orange">สมัครสมาชิก</CardTitle>
                    <CardDescription>กรอกข้อมูลเพื่อสมัครสมาชิกใช้งานระบบจองคอร์ท</CardDescription>
                </CardHeader>
                <CardContent>
                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                            <FormField
                                control={form.control}
                                name="studentId"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>รหัสนักศึกษา</FormLabel>
                                        <FormControl>
                                            <Input placeholder="กรอกรหัสนักศึกษา" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="name"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>ชื่อ-นามสกุล</FormLabel>
                                        <FormControl>
                                            <Input placeholder="กรอกชื่อและนามสกุล" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="email"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>อีเมล (ไม่จำเป็น)</FormLabel>
                                        <FormControl>
                                            <Input placeholder="กรอกอีเมล (ถ้ามี)" type="email" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="password"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>รหัสผ่าน</FormLabel>
                                        <FormControl>
                                            <div className="relative">
                                                <Input
                                                    type={showPassword ? 'text' : 'password'}
                                                    placeholder="กรอกรหัสผ่าน"
                                                    {...field}
                                                />
                                                <Button
                                                    type="button"
                                                    variant="ghost"
                                                    size="icon"
                                                    className="absolute right-0 top-0 h-full px-3"
                                                    onClick={() => setShowPassword(!showPassword)}
                                                >
                                                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                                                </Button>
                                            </div>
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="confirmPassword"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>ยืนยันรหัสผ่าน</FormLabel>
                                        <FormControl>
                                            <div className="relative">
                                                <Input
                                                    type={showConfirmPassword ? 'text' : 'password'}
                                                    placeholder="กรอกรหัสผ่านอีกครั้ง"
                                                    {...field}
                                                />
                                                <Button
                                                    type="button"
                                                    variant="ghost"
                                                    size="icon"
                                                    className="absolute right-0 top-0 h-full px-3"
                                                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                                >
                                                    {showConfirmPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                                                </Button>
                                            </div>
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <Button
                                type="submit"
                                className="w-full bg-court-orange hover:bg-court-orange/90"
                                disabled={isLoading}
                            >
                                {isLoading ? 'กำลังดำเนินการ...' : 'สมัครสมาชิก'}
                            </Button>
                        </form>
                    </Form>
                </CardContent>
                <CardFooter className="flex justify-center">
                    <p className="text-sm text-center">
                        มีบัญชีอยู่แล้ว?{' '}
                        <Link to="/login" className="text-court-orange hover:underline">
                            เข้าสู่ระบบ
                        </Link>
                    </p>
                </CardFooter>
            </Card>
        </div>
    );
};

export default Register;