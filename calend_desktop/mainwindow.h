#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include "auth.h"
#include "profile.h"
#include "cal_30.h"
#include "event.h"
#include "eventData.h"
#include "cal_1.h"
#include "cal_7.h"

QT_BEGIN_NAMESPACE
namespace Ui {
class MainWindow;
}
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

private slots:


    void on_prof_button_2_clicked();

    void on_mycal_clicked();

    void on_Event_push_clicked();

private:
    Ui::MainWindow *ui;
    QString uid;
    Event *event;
};
#endif // MAINWINDOW_H
